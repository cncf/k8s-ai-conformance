# Evidence: Secure Accelerator Access

---

## Overview

This test validates that GPU access is properly isolated between containers. Pods that request GPU resources via DRA ResourceClaims should have access to GPU devices, while pods that do not request GPU resources should NOT have access to any GPU devices. This ensures proper security isolation of accelerator resources.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.36 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DRA (Dynamic Resource Allocation) driver support enabled is deployed through Blueprint as part of cluster provisioning. The GPU Operator will install all necessary components including the container runtime hooks that enforce GPU isolation.

---

## Step 1: Create ResourceClaimTemplate for GPU Access

Create a ResourceClaimTemplate that requests a GPU using the DRA DeviceClass:

```
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: dra-test-gpu-claim
  namespace: default
spec:
  spec:
    devices:
      requests:
      - name: gpu
        exactly:
          deviceClassName: gpu.nvidia.com
```

---

## Step 2: Create Pod WITH GPU ResourceClaim

Create a Pod that uses the ResourceClaimTemplate to request GPU access:

```
apiVersion: v1
kind: Pod
metadata:
  name: dra-test-gpu-pod
  namespace: default
spec:
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/cuda:13.4.3-base-ubuntu26.04
    command: ['sh', '-c', 'nvidia-smi -L && sleep 300']
    resources:
      claims:
      - name: gpu-claim
  resourceClaims:
  - name: gpu-claim
    resourceClaimTemplateName: dra-test-gpu-claim
  restartPolicy: Never
```

---

## Step 3: Verify GPU Access in Authorized Pod

Check that the GPU pod can access the GPU via nvidia-smi:

```
$ kubectl logs -f dra-test-gpu-pod
GPU 0: Tesla T4 (UUID: GPU-9856cd85-70af-d094-7dff-ca2f5b677dc8)
```

> **Note:** ✅ Pod with GPU ResourceClaim successfully accessed the GPU.

---

## Step 4: Verify NVIDIA Devices Present in GPU Pod

List NVIDIA devices in the authorized GPU pod:

```
$ kubectl exec -it dra-test-gpu-pod -- bash

root@dra-test-gpu-pod:/# ls -la /dev/nvidia*
crw-rw-rw- 1 root root 195, 254 Oct 29 21:12 /dev/nvidia-modeset
crw-rw-rw- 1 root root 508,   0 Oct 29 21:12 /dev/nvidia-uvm
crw-rw-rw- 1 root root 508,   1 Oct 29 21:12 /dev/nvidia-uvm-tools
crw-rw-rw- 1 root root 195,   0 Oct 29 21:12 /dev/nvidia0
crw-rw-rw- 1 root root 195, 255 Oct 29 21:12 /dev/nvidiactl

/dev/nvidia-caps:
total 0
drwxr-xr-x 2 root root     80 Oct 29 21:12 .
drwxr-xr-x 6 root root    480 Oct 29 21:12 ..
cr-------- 1 root root 511, 1 Oct 29 21:12 nvidia-cap1
cr--r--r-- 1 root root 511, 2 Oct 29 21:12 nvidia-cap2
```

> **Note:** ✅ GPU pod has access to NVIDIA device files (/dev/nvidia0, /dev/nvidiactl, etc.).

---

## Step 5: Create Pod WITHOUT GPU ResourceClaim (Isolation Test)

Create a Pod that does NOT request GPU resources (resources.claims is commented out):

```
apiVersion: v1
kind: Pod
metadata:
  name: dra-test-gpu-pod
  namespace: default
spec:
  containers:
  - name: cuda-container
    image: nvcr.io/nvidia/cuda:13.4.3-base-ubuntu26.04
    command: ['sh', '-c']
    args:
      - |
        while true;
        do
          nvidia-smi -L
          sleep 300
        done
      #resources:
      #claims:
      #- name: gpu-claim
  resourceClaims:
  - name: gpu-claim
    resourceClaimTemplateName: dra-test-gpu-claim
  restartPolicy: Never
```

---

## Step 6: Verify nvidia-smi Fails in Non-GPU Pod

Check that the non-GPU pod cannot run nvidia-smi (command not found because no GPU runtime injected):

```
$ kubectl logs -f dra-test-gpu-pod
sh: 3: nvidia-smi: not found
```

> **Note:** ✅ Pod without GPU claim cannot run nvidia-smi - isolation working correctly.

---

## Step 7: Verify NO NVIDIA Devices in Non-GPU Pod

List all devices in the non-GPU pod to confirm no NVIDIA devices are present:

```
kubectl exec -it dra-test-gpu-pod -- ls -la /dev/

total 4
drwxr-xr-x 5 root root  360 Oct 29 21:16 .
drwxr-xr-x 1 root root 4096 Oct 29 21:16 ..
lrwxrwxrwx 1 root root   11 Oct 29 21:16 core -> /proc/kcore
lrwxrwxrwx 1 root root   13 Oct 29 21:16 fd -> /proc/self/fd
crw-rw-rw- 1 root root 1, 7 Oct 29 21:16 full
drwxrwxrwt 2 root root   40 Oct 29 21:16 mqueue
crw-rw-rw- 1 root root 1, 3 Oct 29 21:16 null
lrwxrwxrwx 1 root root    8 Oct 29 21:16 ptmx -> pts/ptmx
drwxr-xr-x 2 root root    0 Oct 29 21:16 pts
crw-rw-rw- 1 root root 1, 8 Oct 29 21:16 random
drwxrwxrwt 2 root root   40 Oct 29 21:16 shm
lrwxrwxrwx 1 root root   15 Oct 29 21:16 stderr -> /proc/self/fd/2
lrwxrwxrwx 1 root root   15 Oct 29 21:16 stdin -> /proc/self/fd/0
lrwxrwxrwx 1 root root   15 Oct 29 21:16 stdout -> /proc/self/fd/1
-rw-rw-rw- 1 root root    0 Oct 29 21:16 termination-log
crw-rw-rw- 1 root root 5, 0 Oct 29 21:16 tty
crw-rw-rw- 1 root root 1, 9 Oct 29 21:16 urandom
crw-rw-rw- 1 root root 1, 5 Oct 29 21:16 zero

```

> **Note:** ✅ Non-GPU pod has NO NVIDIA devices in /dev/ - complete isolation confirmed. No /dev/nvidia* files present.

---

## Step 8: Run Kubernetes E2E DRA Test

Run the official Kubernetes e2e test for DRA device mapping to verify conformance:

```
$ go run github.com/onsi/ginkgo/v2/ginkgo -focus='must map configs and devices to the right containers' ./test/e2e


  I1029 21:34:40.468489  667251 test_context.go:528] The --provider flag is not set. Continuing as if --provider=skeleton had been used.
  I1029 21:34:40.468623  667251 e2e.go:109] Starting e2e run "5d4d636b-75c1-47d0-90f9-2c6ce68c3105" on Ginkgo node 1
Running Suite: Kubernetes e2e suite - /root/kubernetes/test/e2e
===============================================================
Random Seed: 1793309666 - will randomize all specs

Will run 1 of 3229 specs
SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS•SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS

Ran 1 of 3229 Specs in 12.582 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 3228 Skipped
PASS

Ginkgo ran 1 suite in 26.621720887s
Test Suite Passed
```

> **Note:** ✅ Official Kubernetes e2e test "must map configs and devices to the right containers" PASSED.

---


**Summary:** GPU access is properly isolated through DRA and the container runtime. Pods that request GPU resources via ResourceClaims receive access to NVIDIA device files (/dev/nvidia0, /dev/nvidiactl, etc.) and can run nvidia-smi. Pods that do not request GPU resources are completely isolated - they have no access to any NVIDIA devices and cannot interact with GPU hardware. This prevents unauthorized access or interference between workloads.
