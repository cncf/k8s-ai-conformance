# Dynamic Resource Allocation (DRA) Support

## Overview

Dynamic Resource Allocation (DRA) is a Kubernetes feature that provides a flexible mechanism for requesting and sharing resources such as GPUs and other hardware accelerators. Unlike traditional device plugins, DRA allows workloads to describe resource requirements through ResourceClaim objects, enabling more expressive allocation, sharing across pods, and vendor-specific configuration via custom drivers.

DRA is a stable feature in Kubernetes 1.35 and is enabled by default. This guide demonstrates the DRA API availability on Viettel AI Platform running Kubernetes v1.35.3.

## Step 1: Verify DRA API Resources

```bash
$ kubectl api-resources --api-group=resource.k8s.io
NAME                     SHORTNAMES   APIVERSION           NAMESPACED   KIND
deviceclasses                         resource.k8s.io/v1   false        DeviceClass
resourceclaims                        resource.k8s.io/v1   true         ResourceClaim
resourceclaimtemplates                resource.k8s.io/v1   true         ResourceClaimTemplate
resourceslices                        resource.k8s.io/v1   false        ResourceSlice
```

All four `resource.k8s.io/v1` DRA resource kinds are registered and served by the API server.

## Step 2: Verify Cluster and Node Information

```bash
$ kubectl version --short
Client Version: v1.34.1
Kustomize Version: v5.7.1
Server Version: v1.35.3

$ kubectl get nodes -o wide
NAME          STATUS   ROLES           AGE   VERSION   INTERNAL-IP   OS-IMAGE              KERNEL-VERSION       CONTAINER-RUNTIME
ubuntu-sv16   Ready    control-plane   36h   v1.35.3   10.24.10.16   Ubuntu 22.04.5 LTS    5.15.0-173-generic   containerd://2.2.2
ubuntu-sv18   Ready    <none>          21h   v1.35.3   10.24.10.18   Ubuntu 22.04.5 LTS    5.15.0-173-generic   containerd://1.7.28
```

## Step 3: Create a ResourceClaimTemplate and Test Pod

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: dra-test
---
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: single-gpu
  namespace: dra-test
spec:
  spec:
    devices:
      requests:
      - name: req-0
        exactly:
          deviceClassName: gpu.nvidia.com
          allocationMode: ExactCount
          count: 1
---
apiVersion: v1
kind: Pod
metadata:
  name: dra-test-pod
  namespace: dra-test
spec:
  containers:
  - name: test
    image: ubuntu:22.04
    command: ["/bin/sh", "-c", "while true; do sleep 30; done"]
    resources:
      claims:
      - name: single-gpu
  resourceClaims:
  - name: single-gpu
    resourceClaimTemplateName: single-gpu
  restartPolicy: Never
EOF
```

## Step 4: Verify ResourceClaim and Pod

```bash
$ kubectl get resourceclaimtemplates -n dra-test
NAME         AGE
single-gpu   5s

$ kubectl api-resources --api-group=resource.k8s.io
NAME                     SHORTNAMES   APIVERSION           NAMESPACED   KIND
deviceclasses                         resource.k8s.io/v1   false        DeviceClass
resourceclaims                        resource.k8s.io/v1   true         ResourceClaim
resourceclaimtemplates                resource.k8s.io/v1   true         ResourceClaimTemplate
resourceslices                        resource.k8s.io/v1   false        ResourceSlice
```

## Cleanup

```bash
kubectl delete namespace dra-test
```

## Notes

The `resource.k8s.io/v1` DRA APIs are fully enabled on Kubernetes v1.35.3. GPUs are currently also exposed via the NVIDIA Device Plugin (`nvidia.com/gpu` extended resource). A DRA-native GPU driver (nvidia-dra-driver-gpu) can be installed for full DRA-based GPU allocation with ResourceClaim objects.
