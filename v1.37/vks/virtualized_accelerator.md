## Description

If the platform supports virtualized accelerators (e.g., vGPU), it should demonstrate that these resources can be scheduled and utilized by workloads.

## Evidence

VKS supports virtualized GPU resources through VMware vSphere with NVIDIA vGPU technology. This enables multiple workloads to share physical GPU resources while maintaining isolation and performance.

### vGPU Configuration

VKS exposes vGPU devices through VM Classes that are configured with NVIDIA vGPU profiles. These profiles define the amount of GPU memory and compute resources allocated to each virtual machine.

**VM Class with vGPU Profile:**

The cluster's worker node pool is configured with the VM Class, which carries a vGPU profile:
- `GRID A100-4C`: 4GB vGPU profile

This profile is based on NVIDIA GRID technology running on VMware vSphere, providing virtualized access to an NVIDIA A100 GPU.

### vGPU Evidence in Node Labels

The presence of vGPU support is verified through node labels that are automatically populated by the NVIDIA GPU Operator:

```shell
kubectl get nodes -o json | jq '.items[].metadata.labels' | grep vgpu
```

```shell
  "nvidia.com/vgpu.host-driver-branch": "r596_84",
  "nvidia.com/vgpu.host-driver-version": "595.91.04",
  "nvidia.com/vgpu.present": "true",
  "nvidia.com/vgpu.host-driver-branch": "r596_84",
  "nvidia.com/vgpu.host-driver-version": "595.91.04",
  "nvidia.com/vgpu.present": "true",
```

```shell
kubectl get nodes -o json | jq '.items[].metadata.labels' | grep 'gpu.product'
```

```shell
  "nvidia.com/gpu.product": "GRID-A100-4C",
  "nvidia.com/gpu.product": "GRID-A100-4C",
```

Key indicators of vGPU support:
- `nvidia.com/vgpu.present: "true"` - Confirms vGPU is available on the node
- `nvidia.com/vgpu.host-driver-version` - Shows the host driver version managing the vGPU
- `nvidia.com/gpu.product` - Shows GRID-prefixed product names, indicating vGPU profiles

### vGPU in DRA ResourceSlice Attributes

The NVIDIA DRA driver (v25.12.0, installed as described in [dra_support.md](dra_support.md)) exposes the vGPU devices as schedulable DRA resources:

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-6d4c57447b-725bh   1/1     Running   0          16s
nvidia-dra-driver-gpu-kubelet-plugin-mtcfs          2/2     Running   0          16s
nvidia-dra-driver-gpu-kubelet-plugin-xjl98          2/2     Running   0          16s
```

The ResourceSlice attributes show the brand as `NvidiaVCS` (NVIDIA vCompute Server), which is the vGPU-specific identifier:

```shell
kubectl get resourceslices -o yaml | grep -A 24 "brand:"
```

```shell
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.178.4
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-ba9adcc0-9809-47bb-85a6-b46a00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-sms5g-pt7t9-p42p6
    pool:
--
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.178.4
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-60e223c0-015b-436c-ba00-26da00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-sms5g-pt7t9-vv2dd
    pool:
```

The `brand: NvidiaVCS` attribute specifically indicates that these are virtualized GPU resources managed through [NVIDIA's vCompute Server technology](https://www.nvidia.com/content/dam/en-zz/Solutions/design-visualization/solutions/resources/documents1/nvidia-virtual-compute-server-solution-overview.pdf). The vGPU is also exposed through the device plugin as the extended resource `nvidia.com/gpu`:

```shell
kubectl get nodes -o custom-columns=NAME:.metadata.name,GPU_ALLOCATABLE:.status.allocatable.nvidia\.com/gpu
```

```shell
NAME                                     GPU_ALLOCATABLE
conformance-test-np1-sms5g-pt7t9-p42p6   1
conformance-test-np1-sms5g-pt7t9-vv2dd   1
conformance-test-pcdv2-gv6sq             <none>
```

### Workload Scheduling with vGPU

Workloads can request vGPU resources through standard Kubernetes resource requests. The same `ResourceClaimTemplate` + `Deployment` used in [dra_support.md](dra_support.md) (image `ubuntu:22.04`, `deviceClassName: gpu.nvidia.com`) was deployed in namespace `gpu-test1`:

```shell
kubectl get pods -n gpu-test1 -o wide
```

```shell
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
dra-gpu-example-6cf7c799b6-kxtb6   1/1     Running   0          8s    193.0.1.58   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl get resourceclaims -n gpu-test1 -o custom-columns=NAME:.metadata.name,POOL:.status.allocation.devices.results[0].pool,DEVICE:.status.allocation.devices.results[0].device
```

```shell
NAME                                                POOL                                     DEVICE
dra-gpu-example-6cf7c799b6-kxtb6-single-gpu-m4ffl   conformance-test-np1-sms5g-pt7t9-vv2dd   gpu-0
```

The pod logs confirm access to the vGPU device:

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example --tail=2
```

```shell
Tue Nov  3 23:13:31 UTC 2026
GPU 0: GRID A100-4C (UUID: GPU-60e223c0-015b-436c-ba00-26da00000000)
```

From inside the container, `nvidia-smi` reports the vGPU profile and the virtualization mode:

```shell
kubectl exec -n gpu-test1 deploy/dra-gpu-example -- nvidia-smi
```

```shell
Tue Nov  3 23:13:40 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 580.178.04             Driver Version: 580.178.04     CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  GRID A100-4C                   On  |   00000000:03:00.0 Off |                    0 |
| N/A   N/A    P0            N/A  /  N/A  |       0MiB /   4096MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|  No running processes found                                                             |
+-----------------------------------------------------------------------------------------+
```

```shell
kubectl exec -n gpu-test1 deploy/dra-gpu-example -- bash -c 'nvidia-smi -q | grep -E "Product Name|Product Brand|Virtualization Mode|Host VGPU Mode|vGPU Software Licensed Product|License Status"'
```

```shell
    Product Name                                       : GRID A100-4C
    Product Brand                                      : NVIDIA Virtual Compute Server
    GPU Virtualization Mode
        Virtualization Mode                            : VGPU
        Host VGPU Mode                                 : N/A
    vGPU Software Licensed Product
        Product Name                                   : NVIDIA Virtual Compute Server
        License Status                                 : Unlicensed
```

The UUID reported inside the container (`GPU-60e223c0-015b-436c-ba00-26da00000000`) matches the `gpu-0` device UUID published by node `conformance-test-np1-sms5g-pt7t9-vv2dd` in the ResourceSlice shown above, and `Virtualization Mode: VGPU` confirms the device is a vGPU, correctly scheduled and exposed to the workload via DRA.

**Note**: The vGPU guest driver in this lab reports `License Status: Unlicensed` because a firewall in the lab blocks the nodes from reaching the NVIDIA license server. This restricts CUDA compute, but not device discovery, scheduling, allocation or telemetry.

### Configuration Documentation

For detailed instructions on configuring vGPU on VKS, refer to the official VMware documentation:

https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Cross-References

The following evidence files demonstrate vGPU functionality across different AI conformance requirements:

- **dra_support.md**: Shows the vGPU device (GRID A100-4C) exposed through DRA ResourceSlices and successfully allocated to a workload
- **pod_autoscaling.md**: Demonstrates HPA functionality with vGPU-backed pods using DCGM metrics
- **gang_scheduling.md**: Shows Kueue scheduling GPU workloads across the vGPU-enabled node pool
- **secure_accelerator_access.md**: Validates proper isolation and access control for vGPU resources

### Cleanup

```shell
kubectl delete namespace gpu-test1
helm uninstall nvidia-dra-driver-gpu -n nvidia-dra-driver-gpu
kubectl delete namespace nvidia-dra-driver-gpu
kubectl delete crd computedomaincliques.resource.nvidia.com computedomains.resource.nvidia.com
kubectl get resourceslice -o name | xargs -r kubectl delete
```
