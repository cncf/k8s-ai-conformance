# Evidence: DRA Support

---

## Overview

This test validates that the Kubernetes platform supports Dynamic Resource Allocation (DRA) APIs (resource.k8s.io/v1) for GPU resource management. DRA provides a more flexible mechanism for requesting and allocating hardware resources compared to traditional device plugins.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.36 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DRA (Dynamic Resource Allocation) driver support enabled is deployed through Blueprint as part of cluster provisioning. The GPU Operator will install all necessary components including the NVIDIA device plugin, container toolkit, and DRA driver required for this conformance test.

---

## Step 1: Verify GPU Operator and DRA Driver Installation

First, verify that the NVIDIA GPU Operator with DRA driver support is installed and running:

```
$ kubectl get po -A

NAMESPACE        NAME                                                              READY   STATUS      RESTARTS      AGE
ai-conformance   gpu-feature-discovery-lxdqn                                       1/1     Running     0             47h
ai-conformance   gpu-operator-56b8dfbbf7-sfxp9                                     1/1     Running     0             47h
ai-conformance   gpu-operator-cncf-node-feature-discovery-gc-696577d54d-fjx6z      1/1     Running     0             47h
ai-conformance   gpu-operator-cncf-node-feature-discovery-master-5fc6577d99splbq   1/1     Running     0             47h
ai-conformance   gpu-operator-cncf-node-feature-discovery-worker-n95d7             1/1     Running     0             47h
ai-conformance   nvidia-container-toolkit-daemonset-j7zfm                          1/1     Running     0             47h
ai-conformance   nvidia-cuda-validator-x67f4                                       0/1     Completed   0             6m41s
ai-conformance   nvidia-dcgm-exporter-dv8tv                                        1/1     Running     0             47h
ai-conformance   nvidia-dra-driver-gpu-controller-76864c66cb-kq2pw                 1/1     Running     0             3m32s
ai-conformance   nvidia-dra-driver-gpu-kubelet-plugin-m6qqx                        2/2     Running     0             3m32s
ai-conformance   nvidia-driver-daemonset-khrvr                                     1/1     Running     0             47h
ai-conformance   nvidia-operator-validator-zgghl                                   1/1     Running     0             6m43s
```

> **Note:** ✅ GPU Operator pods are running including the DRA driver components.

---

## Step 2: Create ResourceClaimTemplate

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

Verify the ResourceClaimTemplate was created:

```
$ kubectl get resourceclaimtemplates -A

NAMESPACE   NAME                 AGE
default     dra-test-gpu-claim   9s

```

> **Note:** ✅ ResourceClaimTemplate created successfully.

---

## Step 3: Create Pod with DRA ResourceClaim

Create a Pod that uses the ResourceClaimTemplate to request GPU access via DRA:

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

The Pod references the ResourceClaimTemplate via `spec.resourceClaims[].resourceClaimTemplateName`, and the container requests the claim via `spec.containers[].resources.claims`.

---

## Step 4: Verify ResourceClaim Allocation

Check that a ResourceClaim was automatically created from the template and allocated to the Pod:

```
$ kubectl get resourceclaims -A

NAMESPACE   NAME                               STATE                AGE
default     dra-test-gpu-pod-gpu-claim-g5f4m   allocated,reserved   44s
```

> **Note:** ✅ ResourceClaim was automatically created from the template and shows state `allocated,reserved` - the GPU has been allocated to the Pod.

---

## Step 5: Verify GPU Access in Pod

Check the Pod logs to confirm the container can access the GPU via nvidia-smi:

```
$ kubectl logs dra-test-gpu-pod
GPU 0: Tesla T4 (UUID: GPU-9856cd85-70af-d094-7dff-ca2f5b677dc8)
```

> **Note:** ✅ Pod successfully accessed the GPU via DRA ResourceClaim. The nvidia-smi command detected the Tesla T4 GPU.

---
