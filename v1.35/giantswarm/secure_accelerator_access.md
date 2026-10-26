# Secure Accelerator Access Tests

**MUST**: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Prerequisites

- A Kubernetes 1.35 Giant Swarm cluster.
- A GPU node pool whose instance type has **at least two GPUs**. The isolation test in
  Step 3 allocates one device per Pod, and DRA never double-allocates a device, so on a
  single-GPU instance the second Pod stays `Pending`. The runs below use
  `g4dn.12xlarge` (4 × NVIDIA Tesla T4). `p4d.24xlarge` is not a practical choice: it
  rarely gets capacity in `eu-north-1`, and 48 vCPU already consumes most of a default
  64 vCPU G/VT quota.
- [The NVIDIA GPU Operator and the NVIDIA DRA driver](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/dynamic-resource-allocation/),
  installed as described in the reproduction README. The GPU Operator is a hard
  prerequisite, not an optional extra: it brings Node Feature Discovery, whose
  `feature.node.kubernetes.io/pci-10de.present` label is what lets the DRA kubelet
  plugin schedule onto a GPU node at all, and it creates the `nvidia` `RuntimeClass`
  the Pods below reference.
- NVIDIA driver **570 or newer** on the GPU nodes. Images built by
  [capi-image-builder](https://github.com/giantswarm/capi-image-builder) default to
  570 since [#1322](https://github.com/giantswarm/capi-image-builder/pull/1322);
  Step 1 verifies it, and the troubleshooting note at the end covers older images.

All Pods below carry an explicit `securityContext`. Giant Swarm enforces the Pod
Security Standards *restricted* profile through Kyverno on every cluster, so a Pod
without one is rejected at admission and never reaches a node.

## Tests

### Test 1: Verify Isolated GPU Access via DRA

**Step 1**: With the DRA driver running, a `DeviceClass` named `gpu.nvidia.com` and one
`ResourceSlice` per GPU node are created automatically.

```nohighlight
$ kubectl get deviceclass
NAME                  AGE
gpu.nvidia.com        20s
mig.nvidia.com        20s
vfio.gpu.nvidia.com   20s

$ kubectl get resourceslices
NAME                                      NODE                                          DRIVER           AGE
00000-gpu.nvidia.com-ip-10-0-207-75...    ip-10-0-207-75.eu-north-1.compute.internal    gpu.nvidia.com   26s
```

The slice is the platform's own inventory of the node's accelerators. Read it back
before running anything: it names every device the scheduler can allocate and reports
the driver and CUDA runtime versions in use.

```nohighlight
$ kubectl get resourceslice -o jsonpath='{.items[0].spec.devices[*].name}'
gpu-0 gpu-1 gpu-2 gpu-3

$ kubectl get resourceslice -o yaml | yq '.items[0].spec.devices[0]'
attributes:
  architecture:          {string: Turing}
  brand:                 {string: Nvidia}
  cudaComputeCapability: {version: 7.5.0}
  cudaDriverVersion:     {version: 12.8.0}
  driverVersion:         {version: 570.195.3}
  productName:           {string: Tesla T4}
  type:                  {string: gpu}
  uuid:                  {string: GPU-5c0de657-e29f-91f0-987e-ed508ce78d5f}
capacity:
  memory: {value: 15Gi}
name: gpu-0
```

**Expected Result**: four devices are advertised, and `driverVersion` reports 570 or
newer. A `ResourceSlice` that never appears means the kubelet plugin is not running. 
Check the troubleshooting note if it is not the case.

**Step 2 [Accessible]**: Create a `ResourceClaimTemplate` that requests any GPU, then
deploy a Pod that references it.

```bash
$ kubectl apply -f - <<EOF
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: single-gpu
  namespace: default
spec:
  spec:
    devices:
      requests:
      - name: gpu
        exactly:
          deviceClassName: gpu.nvidia.com
---
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test-accessible
  namespace: default
spec:
  restartPolicy: Never
  runtimeClassName: nvidia
  tolerations:
  - key: nvidia.com/gpu
    operator: Exists
    effect: NoSchedule
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    seccompProfile:
      type: RuntimeDefault
  resourceClaims:
  - name: gpu
    resourceClaimTemplateName: single-gpu
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop: ["ALL"]
    resources:
      claims:
      - name: gpu
EOF

$ kubectl wait --for=condition=Ready pod/gpu-test-accessible --timeout=300s
```

Confirm which device the scheduler allocated, that the container can see it, and that a
real CUDA kernel runs on it:

```nohighlight
$ kubectl get resourceclaims -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.allocation.devices.results[0].device}{"\n"}{end}'
gpu-test-accessible-gpu-rlrz2	gpu-0

$ kubectl exec gpu-test-accessible -- sh -c 'ls /dev/ | grep "^nvidia"'
nvidia-modeset
nvidia-uvm
nvidia-uvm-tools
nvidia0
nvidiactl

$ kubectl exec gpu-test-accessible -- /cuda-samples/vectorAdd
[Vector addition of 50000 elements]
Copy input data from the host memory to the CUDA device
CUDA kernel launch with 196 blocks of 256 threads
Copy output data from the CUDA device to the host memory
Test PASSED
Done
```

**Expected Result**: the claim is bound to a single device, the container holds the
matching device nodes, and a CUDA workload completes on it. The DRA kubelet plugin
injected the device through CDI.

> **Note**: do not verify this with `nvidia-smi`. On Flatcar Container Linux the NVIDIA
> utilities live in `/opt/bin`, outside the container's filesystem and off its `PATH`.
> CDI injects the device nodes and the driver libraries a CUDA runtime needs, not the
> command-line tools, so `nvidia-smi` reports `executable file not found in $PATH`
> inside a container that has a perfectly good GPU. Running the workload proves more
> than querying the device would.

**Step 3 [Isolation]**: Deploy two further Pods on the same node, each with its own
claim from the same template. The manifests are identical to Step 2's Pod apart from
the name; create `gpu-test-pod1` and `gpu-test-pod2`.

```nohighlight
$ kubectl wait --for=condition=Ready pod/gpu-test-pod1 --timeout=300s
$ kubectl wait --for=condition=Ready pod/gpu-test-pod2 --timeout=300s

# Each Pod gets its own claim, allocated to a distinct device.
$ kubectl get resourceclaims -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.allocation.devices.results[0].device}{"\n"}{end}'
gpu-test-accessible-gpu-rlrz2	gpu-0
gpu-test-pod1-gpu-zqrxs	gpu-1
gpu-test-pod2-gpu-29cdm	gpu-2

# Each container sees only the device node it was allocated.
$ kubectl exec gpu-test-pod1 -- sh -c 'ls /dev/ | grep "^nvidia[0-9]"'
nvidia1

$ kubectl exec gpu-test-pod2 -- sh -c 'ls /dev/ | grep "^nvidia[0-9]"'
nvidia2

# Both run their own workload concurrently.
$ kubectl exec gpu-test-pod1 -- /cuda-samples/vectorAdd | grep Test
Test PASSED

$ kubectl exec gpu-test-pod2 -- /cuda-samples/vectorAdd | grep Test
Test PASSED
```

**Expected Result**: the three claims are bound to three different devices, and each
container holds exactly one device node, matching its own allocation and no other. DRA
enforces the separation at allocation time and CDI injection scopes device visibility
per container.

### Test 2: Verify Unauthorized Access Prevention

**Step 1**: Deploy a Pod that references no `ResourceClaim` and carries no
`runtimeClassName`, and verify it cannot reach a GPU.

```bash
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test-unauthorized
  namespace: default
spec:
  restartPolicy: Never
  tolerations:
  - key: nvidia.com/gpu
    operator: Exists
    effect: NoSchedule
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop: ["ALL"]
    # Note: no resourceClaims field, no runtimeClassName: nvidia
EOF

$ kubectl wait --for=condition=Ready pod/gpu-test-unauthorized --timeout=300s
```

The Pod is scheduled on the GPU node and runs, but the accelerator is not reachable
from it:

```nohighlight
$ kubectl exec gpu-test-unauthorized -- /cuda-samples/vectorAdd
[Vector addition of 50000 elements]
Failed to allocate device vector A (error code CUDA driver version is insufficient for CUDA runtime version)!
command terminated with exit code 1
```

**Expected Result**: the CUDA runtime cannot acquire a device. Without a
`ResourceClaim` and the matching `resources.claims` entry, the DRA kubelet plugin
performs no CDI injection for this Pod, so the container has no device nodes and no
usable driver.

**Step 2**: Confirm the container cannot bypass DRA by reaching the device files
directly.

```nohighlight
$ kubectl exec gpu-test-unauthorized -- ls -la /dev/nvidia*
ls: cannot access '/dev/nvidia*': No such file or directory
command terminated with exit code 2

$ kubectl exec gpu-test-unauthorized -- sh -c 'ls /dev/ | grep -c nvidia'
0
command terminated with exit code 1
```

**Expected Result**: no GPU device nodes (`/dev/nvidia0`, `/dev/nvidiactl`,
`/dev/nvidia-uvm`, ...) are present. CDI injection is the only path that exposes them,
and it runs only for containers whose Pod owns an allocated `ResourceClaim`.

Host-level escape is closed off independently: the platform's Kyverno policies reject
any Pod that sets `hostPath`, `hostPID` or `hostIPC`, so a workload cannot mount the
node's device tree to reach a GPU it was not allocated.

## Troubleshooting

**No `ResourceSlice` appears.** Check the kubelet plugin on the GPU node. If its init
container is stuck, read its log: `Option --version is not recognized` from
`nvidia-smi` means the node is on NVIDIA driver 535, which the DRA driver does not
support. Images built since
[capi-image-builder#1322](https://github.com/giantswarm/capi-image-builder/pull/1322)
default to 570; on an older image, write `NVIDIA_DRIVER_VERSION=570.195.03` to
`/etc/flatcar/nvidia-metadata` and replace the node. `setup-nvidia` sources that file
after Flatcar's own defaults, so it selects the newer driver without touching the
read-only `/usr`.

**The kubelet plugin DaemonSet reports 0 desired Pods.** Its node affinity requires an
NFD or GFD label (`feature.node.kubernetes.io/pci-10de.present` or
`nvidia.com/gpu.present`). Neither exists until the GPU Operator is running.

**The DRA driver chart refuses to render**, reporting that
`resources.gpus.enabled=true` is not supported alongside the standard device plugin.
That guard is deliberate, pending KEP 5004. To serve GPUs through DRA, disable the
standard device plugin and set both `resources.gpus.enabled=true` and
`gpuResourcesEnabledOverride=true`.
