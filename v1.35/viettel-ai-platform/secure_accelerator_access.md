# Secure Accelerator Access Validation

## Overview

This guide validates that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework on Viettel AI Platform. The NVIDIA Device Plugin (`nvidia.com/gpu` extended resource) ensures that only pods with an explicit GPU resource request can access GPU devices.

## Step 1: Verify NVIDIA Device Plugin and GPU Operator

```bash
$ kubectl get pods -n gpu-operator | grep -E 'device-plugin|cuda-validator'
nvidia-cuda-validator-7hghr                                   0/1     Completed   0   27h
nvidia-cuda-validator-x72v5                                   0/1     Completed   0   5h7m
nvidia-device-plugin-daemonset-mz74d                          1/1     Running     0   5h8m
nvidia-device-plugin-daemonset-wcwpp                          1/1     Running     0   28h
```

CUDA validators completed successfully on both GPU nodes confirming driver/runtime correctness.

## Tests

### Test 1: Verify Pod Cannot Access Unallocated GPU Resources

#### 1. Create a pod WITHOUT GPU resource request

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: gpu-access-test
---
apiVersion: v1
kind: Pod
metadata:
  name: no-gpu-pod
  namespace: gpu-access-test
spec:
  containers:
  - name: test
    image: ubuntu:22.04
    command: ["sh", "-c"]
    args:
    - |
      echo '=== Test: Pod WITHOUT GPU request ==='
      if ls /dev/nvidia* 2>/dev/null; then
        echo 'FAIL: GPU devices are accessible'
      else
        echo 'PASS: No GPU devices visible (expected)'
      fi
      echo 'nvidia-smi check:'
      nvidia-smi 2>&1 || echo 'PASS: nvidia-smi not available (expected)'
  restartPolicy: Never
EOF
```

#### 2. Check logs — GPU must NOT be accessible

```bash
$ kubectl logs no-gpu-pod -n gpu-access-test
=== Test: Pod WITHOUT GPU request ===
PASS: No GPU devices visible (expected)
nvidia-smi check:
sh: 8: nvidia-smi: not found
PASS: nvidia-smi not available (expected)
```

No NVIDIA devices are mounted and `nvidia-smi` is unavailable in a container that did not request a GPU.

### Test 2: Verify Pod CAN Access GPU When Requested

#### 1. Create a pod WITH GPU resource request

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: gpu-pod
  namespace: gpu-access-test
spec:
  containers:
  - name: test
    image: nvidia/cuda:12.4.1-base-ubuntu22.04
    command: ["sh", "-c"]
    args:
    - |
      echo '=== Test: Pod WITH GPU request ==='
      nvidia-smi -L
      echo '---'
      ls /dev/nvidia*
    resources:
      limits:
        nvidia.com/gpu: "1"
  restartPolicy: Never
EOF
```

#### 2. Check logs — GPU MUST be accessible

```bash
$ kubectl logs gpu-pod -n gpu-access-test
=== Test: Pod WITH GPU request ===
GPU 0: NVIDIA L40 (UUID: GPU-c3c7096e-db9d-322f-efd9-510971f3d33c)
---
/dev/nvidia-modeset
/dev/nvidia-uvm
/dev/nvidia-uvm-tools
/dev/nvidia1
/dev/nvidiactl

/dev/nvidia-caps:
nvidia-cap1
nvidia-cap2
```

The GPU device (`nvidia1`) is exclusively mounted in the container and `nvidia-smi` correctly identifies the NVIDIA L40 GPU.

### Test 3: GPU Isolation Between Pods

Two pods each requesting `nvidia.com/gpu: 1` on a 2-GPU node will receive distinct exclusive GPU device allocations. The device plugin enforces that no two simultaneously-running pods share the same physical GPU device node (e.g., `/dev/nvidia0` vs `/dev/nvidia1`), preventing unauthorized cross-pod access.

## Cleanup

```bash
kubectl delete namespace gpu-access-test
```
