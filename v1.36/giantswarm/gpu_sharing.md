# GPU Sharing Tests

**SHOULD**: For accelerators that support static GPU sharing, provide well-defined mechanisms for at least one GPU sharing strategy to improve utilization for workloads that do not require a full dedicated GPU.

Giant Swarm exposes GPU sharing through the NVIDIA DRA driver. This test shows
time-slicing: two Pods share one physical GPU through a single `ResourceClaim`, and both
run CUDA workloads on it at the same time.

## Prerequisites

- The cluster, GPU Operator and DRA driver from the reproduction README.
- The DRA driver's `TimeSlicingSettings` feature gate. Time-slicing is an alpha feature of
  the upstream driver and is off by default. Without the gate the driver rejects any claim
  that carries a `sharing` config with `invalid GPU sharing settings`. Add it to
  `dra-values.yaml` and install or upgrade the app:

  ```yaml
  featureGates:
    TimeSlicingSettings: true
  ```

Note that time-slicing doesn't isolate GPU memory between the workloads that share a
device. Workloads that need memory isolation should claim a whole GPU, as in
[`secure_accelerator_access.md`](secure_accelerator_access.md).

## Test: Two Pods share one time-sliced GPU

**Step 1**: Create a `ResourceClaim` for one GPU with a time-slicing config, and two Pods
that both reference it. The Pods carry the same `securityContext` as the other tests,
because the platform enforces the Pod Security Standards *restricted* profile.

```bash
$ kubectl apply -f - <<EOF
apiVersion: resource.k8s.io/v1
kind: ResourceClaim
metadata:
  name: shared-gpu
  namespace: default
spec:
  devices:
    requests:
    - name: gpu
      exactly:
        deviceClassName: gpu.nvidia.com
    config:
    - requests: ["gpu"]
      opaque:
        driver: gpu.nvidia.com
        parameters:
          apiVersion: resource.nvidia.com/v1beta1
          kind: GpuConfig
          sharing:
            strategy: TimeSlicing
            timeSlicingConfig:
              interval: Long
---
apiVersion: v1
kind: Pod
metadata:
  name: gpu-share-a
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
    resourceClaimName: shared-gpu
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
```

Create `gpu-share-b` from the same Pod manifest with only the name changed, then wait for
both:

```bash
$ kubectl wait --for=condition=Ready pod/gpu-share-a pod/gpu-share-b --timeout=300s
```

**Step 2**: Confirm both Pods hold the same claim and the same device.

```nohighlight
$ kubectl get resourceclaim shared-gpu -o jsonpath='{.status.allocation.devices.results[0].device}{"  reservedFor="}{range .status.reservedFor[*]}{.name} {end}{"\n"}'
gpu-0  reservedFor=gpu-share-a gpu-share-b

$ kubectl exec gpu-share-a -- sh -c 'ls /dev/ | grep "^nvidia[0-9]"'
nvidia0

$ kubectl exec gpu-share-b -- sh -c 'ls /dev/ | grep "^nvidia[0-9]"'
nvidia0
```

**Step 3**: Run a CUDA workload in both Pods at the same time.

```nohighlight
$ kubectl exec gpu-share-a -- /cuda-samples/vectorAdd | grep Test &
$ kubectl exec gpu-share-b -- /cuda-samples/vectorAdd | grep Test
Test PASSED
Test PASSED
```

**Expected Result**: one device is allocated once and reserved for both Pods, and both
CUDA workloads complete on it. The GPU is oversubscribed, with the driver time-slicing it
between the two Pods.
