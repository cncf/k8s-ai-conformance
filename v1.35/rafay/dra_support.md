# Evidence: DRA Support

---

## Overview

This test validates that the Kubernetes platform supports Dynamic Resource Allocation (DRA) APIs (resource.k8s.io/v1) for GPU resource management. DRA provides a more flexible mechanism for requesting and allocating hardware resources compared to traditional device plugins.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DRA (Dynamic Resource Allocation) driver support enabled is deployed through Blueprint as part of cluster provisioning. The GPU Operator will install all necessary components including the NVIDIA device plugin, container toolkit, and DRA driver required for this conformance test.

---

## Step 1: Verify GPU Operator and DRA Driver Installation

First, verify that the NVIDIA GPU Operator with DRA driver support is installed and running:

```
$ kubectl get po -n ai-conformance

NAME                                                              READY   STATUS      RESTARTS      AGE
ai-conformance-gpu-operator-node-feature-discovery-gc-7bff8sb7l   1/1     Running     0             49m
ai-conformance-gpu-operator-node-feature-discovery-master-2rxkg   1/1     Running     0             49m
ai-conformance-gpu-operator-node-feature-discovery-worker-7fmvm   1/1     Running     0             49m
gpu-feature-discovery-bpsf4                                       1/1     Running     0             48m
gpu-operator-7f7dfb9975-9j6ll                                     1/1     Running     0             49m
nvidia-container-toolkit-daemonset-phbnh                          1/1     Running     0             48m
nvidia-cuda-validator-mrtwg                                       0/1     Completed   0             45m
nvidia-dcgm-4zbs7                                                 1/1     Running     1 (45m ago)   48m
nvidia-dcgm-exporter-bhqpl                                        1/1     Running     1 (45m ago)   48m
nvidia-dra-driver-gpu-controller-6fd47d97cf-wnqqv                 1/1     Running     0             43m
nvidia-dra-driver-gpu-kubelet-plugin-khpgq                        2/2     Running     0             25m
nvidia-driver-daemonset-mwzs7                                     1/1     Running     0             48m
nvidia-operator-validator-kmfqq                                   1/1     Running     0             48m
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
test        dra-test-gpu-claim   5s

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
    image: nvcr.io/nvidia/cuda:12.6.0-base-ubuntu22.04
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
default     dra-test-gpu-pod-gpu-claim-nn4tg   allocated,reserved   2s
```

> **Note:** ✅ ResourceClaim was automatically created from the template and shows state `allocated,reserved` - the GPU has been allocated to the Pod.

---

## Step 5: Verify GPU Access in Pod

Check the Pod logs to confirm the container can access the GPU via nvidia-smi:

```
$ kubectl logs -f dra-test-gpu-pod
GPU 0: Tesla T4 (UUID: GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea)
```

> **Note:** ✅ Pod successfully accessed the GPU via DRA ResourceClaim. The nvidia-smi command detected the Tesla T4 GPU.

---
