# Secure Accelerator Access Tests

**MUST**: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Tests

### Test 1: Verify Isolated GPU Access via DRA

**Step 1**: Prepare the test environment, including:

- Creating a Kubernetes 1.34 cluster
- [Adding a GPU node pool and installing the NVIDIA DRA driver](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/dynamic-resource-allocation/)

With the DRA driver running, a `DeviceClass` named `gpu.nvidia.com` and a `ResourceSlice` per GPU node are created automatically — verify with `kubectl get deviceclass gpu.nvidia.com` and `kubectl get resourceslices`.

**Step 2 [Accessible]**: Create a `ResourceClaimTemplate` that requests any GPU, then deploy a Pod that references it. Inside the running container, execute a command to detect the accelerator device. The command should succeed and output the model of the accelerator device currently used by the container.

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
  resourceClaims:
  - name: gpu
    resourceClaimTemplateName: single-gpu
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    resources:
      claims:
      - name: gpu
EOF

$ kubectl wait --for=condition=Ready pod/gpu-test-accessible --timeout=300s

$ kubectl exec gpu-test-accessible -- nvidia-smi --query-gpu=name --format=csv,noheader
Tesla T4
```

**Expected Result**: The command should successfully return the GPU model name, confirming that the container has proper access to the GPU through the DRA framework. The kubelet plugin has injected the matching GPU into the container via CDI.

**Step 3 [Isolation]**: Deploy two Pods on the same node, each with its own `ResourceClaim` referencing the `gpu.nvidia.com` device class. The scheduler must allocate distinct devices to each Pod and the kubelet plugin must expose only the allocated device to each container.

> **Prerequisite**: this step requires the GPU node pool to use an instance type with at least two GPUs. The cluster configuration in the README uses `p4d.24xlarge` (8 × NVIDIA A100), which satisfies this. Single-GPU instances such as `g4dn.xlarge` would leave the second Pod `Pending` because DRA does not double-allocate a device.

```shell
# Deploy two Pods, each with an independent ResourceClaim
$ kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test-pod1
  namespace: default
spec:
  restartPolicy: Never
  runtimeClassName: nvidia
  tolerations:
  - key: nvidia.com/gpu
    operator: Exists
    effect: NoSchedule
  resourceClaims:
  - name: gpu
    resourceClaimTemplateName: single-gpu
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    resources:
      claims:
      - name: gpu
---
apiVersion: v1
kind: Pod
metadata:
  name: gpu-test-pod2
  namespace: default
spec:
  restartPolicy: Never
  runtimeClassName: nvidia
  tolerations:
  - key: nvidia.com/gpu
    operator: Exists
    effect: NoSchedule
  resourceClaims:
  - name: gpu
    resourceClaimTemplateName: single-gpu
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    resources:
      claims:
      - name: gpu
EOF

$ kubectl wait --for=condition=Ready pod/gpu-test-pod1 --timeout=300s
$ kubectl wait --for=condition=Ready pod/gpu-test-pod2 --timeout=300s

# Inspect the ResourceClaims that were created from the template - each
# Pod gets its own claim, allocated to a distinct GPU UUID.
$ kubectl get resourceclaims -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.allocation.devices.results[0].device}{"\n"}{end}'
gpu-test-pod1-gpu-<hash>    gpu-0
gpu-test-pod2-gpu-<hash>    gpu-1

# Each Pod should only see its allocated GPU
$ kubectl exec gpu-test-pod1 -- nvidia-smi -L
GPU 0: Tesla T4 (UUID: GPU-dabc57c1-250b-2979-2b6a-7fd7d9574143)

$ kubectl exec gpu-test-pod2 -- nvidia-smi -L
GPU 0: Tesla T4 (UUID: GPU-18705848-fd64-920c-22c5-e2f1a3d5a7c1)
```

**Expected Result**: The two `ResourceClaims` are bound to distinct underlying devices. Each Pod sees exactly one GPU, with different UUIDs across the two Pods, demonstrating that DRA enforces isolation at allocation time and the kubelet plugin's CDI injection scopes device visibility per container.

### Test 2: Verify Unauthorized Access Prevention

**Step 1**: Deploy a Pod that does **not** reference any `ResourceClaim` and verify that it cannot access GPU devices.

```shell
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
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda12.5.0
    command: ["sleep", "3600"]
    # Note: no resourceClaims field, no runtimeClassName: nvidia
EOF

$ kubectl wait --for=condition=Ready pod/gpu-test-unauthorized --timeout=300s

# Attempt to access GPU - this should fail
$ kubectl exec gpu-test-unauthorized -- nvidia-smi
OCI runtime exec failed: exec failed: unable to start container process: exec: "nvidia-smi": executable file not found in $PATH: unknown.
```

**Expected Result**: Without a `ResourceClaim` (and the matching `resources.claims` entry on the container), the DRA kubelet plugin performs no CDI injection for this Pod. The container starts with no NVIDIA libraries, binaries, or device nodes mounted, so `nvidia-smi` is not available.

**Step 2**: Verify that the container cannot bypass DRA by reaching the GPU device files directly.

```shell
# Check if GPU device files are accessible
$ kubectl exec gpu-test-unauthorized -- ls -la /dev/nvidia*
ls: cannot access '/dev/nvidia*': No such file or directory

# Verify that the container runtime has not mounted GPU devices
$ kubectl exec gpu-test-unauthorized -- ls -la /dev/ | grep nvidia
# Should return empty
```

**Expected Result**: GPU device nodes (`/dev/nvidia0`, `/dev/nvidiactl`, `/dev/nvidia-uvm`, …) must not be present inside containers that have not been allocated a GPU through DRA. CDI injection is the only path that exposes these devices, and it only runs for containers whose Pod owns an allocated `ResourceClaim`.
