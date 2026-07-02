# DRA Support — Conformance Evidence

**Requirement (`dra_support`, MUST):** Support Dynamic Resource Allocation (DRA) APIs
to enable more flexible and fine-grained resource requests beyond simple counts.

**Cluster:** OpenShift Container Platform 4.22.0 (GCP, `n1-standard-8`, 2x NVIDIA Tesla T4)
**Kubernetes version:** v1.35.5
**Date:** 2026-06-09

## Summary

OpenShift 4.22 ships the Dynamic Resource Allocation APIs (`resource.k8s.io/v1`) enabled
by default. With the NVIDIA GPU Operator (v26.3.2) and the NVIDIA DRA Driver for GPUs
(v25.12.0) installed, GPU devices are advertised through ResourceSlices with rich
per-device attributes, and workloads allocate them through ResourceClaims. The sections
below show the API availability, the advertised devices and their attributes, and a
successful end-to-end allocation.

## 1. DRA APIs available

The DRA API group is served at GA (`resource.k8s.io/v1`):

```
$ oc api-resources --api-group=resource.k8s.io
NAME                     SHORTNAMES   APIVERSION           NAMESPACED   KIND
deviceclasses                         resource.k8s.io/v1   false        DeviceClass
resourceclaims                        resource.k8s.io/v1   true         ResourceClaim
resourceclaimtemplates                resource.k8s.io/v1   true         ResourceClaimTemplate
resourceslices                        resource.k8s.io/v1   false        ResourceSlice
```

The GPU DRA driver registers its device classes:

```
$ oc get deviceclass
NAME                  AGE
gpu.nvidia.com        32m
mig.nvidia.com        32m
vfio.gpu.nvidia.com   32m
```

```
$ oc get pods -n nvidia-dra-driver-gpu
NAME                                         READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-kubelet-plugin-twk54   1/1     Running   0          32m
```

## 2. Devices advertised with fine-grained attributes

Each physical GPU is advertised as an individual device in a ResourceSlice. Beyond a
simple count, every device exposes structured attributes (architecture, driver and CUDA
versions, compute capability, memory capacity, and a unique UUID) that can be used for
attribute-based selection in a ResourceClaim:

```
$ oc get resourceslices
NAME                                                   NODE                              DRIVER           POOL                              AGE
harpatil4220c-g9kfm-gpu-c-rd8r8-gpu.nvidia.com-wvrlc   harpatil4220c-g9kfm-gpu-c-rd8r8   gpu.nvidia.com   harpatil4220c-g9kfm-gpu-c-rd8r8   32m
```

```
$ oc get resourceslices -o yaml
apiVersion: resource.k8s.io/v1
kind: ResourceSlice
metadata:
  name: harpatil4220c-g9kfm-gpu-c-rd8r8-gpu.nvidia.com-wvrlc
  ...
spec:
  driver: gpu.nvidia.com
  nodeName: harpatil4220c-g9kfm-gpu-c-rd8r8
  pool:
    generation: 1
    name: harpatil4220c-g9kfm-gpu-c-rd8r8
    resourceSliceCount: 1
  devices:
  - name: gpu-0
    attributes:
      addressingMode: {string: HMM}
      architecture: {string: Turing}
      brand: {string: Nvidia}
      cudaComputeCapability: {version: 7.5.0}
      cudaDriverVersion: {version: 13.0.0}
      driverVersion: {version: 580.126.20}
      productName: {string: Tesla T4}
      resource.kubernetes.io/pciBusID: {string: "0000:00:04.0"}
      resource.kubernetes.io/pcieRoot: {string: pci0000:00}
      type: {string: gpu}
      uuid: {string: GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174}
    capacity:
      memory: {value: 15Gi}
  - name: gpu-1
    attributes:
      addressingMode: {string: HMM}
      architecture: {string: Turing}
      brand: {string: Nvidia}
      cudaComputeCapability: {version: 7.5.0}
      cudaDriverVersion: {version: 13.0.0}
      driverVersion: {version: 580.126.20}
      productName: {string: Tesla T4}
      resource.kubernetes.io/pciBusID: {string: "0000:00:05.0"}
      resource.kubernetes.io/pcieRoot: {string: pci0000:00}
      type: {string: gpu}
      uuid: {string: GPU-204349eb-0405-e927-278b-8804172f1be3}
    capacity:
      memory: {value: 15Gi}
```

> Output condensed: `metadata` is elided and attribute maps are shown in flow style
> for readability; all device attributes are listed.

## 3. Allocation through a ResourceClaim

The following manifest requests one GPU of device class `gpu.nvidia.com` through a
ResourceClaim and runs a container that verifies device access:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: dra-test
---
apiVersion: resource.k8s.io/v1
kind: ResourceClaim
metadata:
  name: gpu-claim
  namespace: dra-test
spec:
  devices:
    requests:
      - name: gpu
        exactly:
          deviceClassName: gpu.nvidia.com
          allocationMode: ExactCount
          count: 1
---
apiVersion: v1
kind: Pod
metadata:
  name: dra-gpu-test
  namespace: dra-test
spec:
  restartPolicy: Never
  securityContext:
    runAsNonRoot: true
    seccompProfile:
      type: RuntimeDefault
  resourceClaims:
    - name: gpu
      resourceClaimName: gpu-claim
  containers:
    - name: gpu-test
      image: nvidia/cuda:12.9.0-base-ubuntu24.04
      command: ["bash", "-c", "ls -l /dev/nvidia* && echo '---' && nvidia-smi -L && echo 'DRA GPU allocation successful'"]
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
      resources:
        claims:
          - name: gpu
```

```
$ oc apply -f dra-gpu-test.yaml
namespace/dra-test created
resourceclaim.resource.k8s.io/gpu-claim created
pod/dra-gpu-test created
```

```
$ oc get pod dra-gpu-test -n dra-test
NAME           READY   STATUS      RESTARTS   AGE
dra-gpu-test   0/1     Completed   0          19m
```

```
$ oc logs dra-gpu-test -n dra-test
crw-rw-rw-. 1 1000760000 root 195, 254 Jun  9 18:35 /dev/nvidia-modeset
crw-rw-rw-. 1 1000760000 root 511,   0 Jun  9 18:35 /dev/nvidia-uvm
crw-rw-rw-. 1 1000760000 root 511,   1 Jun  9 18:35 /dev/nvidia-uvm-tools
crw-rw-rw-. 1 1000760000 root 195,   0 Jun  9 18:35 /dev/nvidia0
crw-rw-rw-. 1 1000760000 root 195, 255 Jun  9 18:35 /dev/nvidiactl
---
GPU 0: Tesla T4 (UUID: GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174)
DRA GPU allocation successful
```

The allocated device (`GPU-da96f367-...`) corresponds to `gpu-0` advertised in the
ResourceSlice above.

> The ResourceClaim returns to an unallocated state once the consuming pod completes, as
> the DRA controller releases the claimed device automatically.

## Result

**PASS.** OpenShift 4.22 supports the DRA APIs (`resource.k8s.io/v1`), advertises GPUs as
attribute-rich devices through ResourceSlices, and successfully allocates them to
workloads through ResourceClaims — satisfying the `dra_support` requirement.

## Cleanup

```
$ oc delete ns dra-test
```
