## Description

If the platform supports virtualized accelerators (e.g., vGPU), it should demonstrate that these resources can be scheduled and utilized by workloads.

## Evidence

VKS supports virtualized GPU resources through VMware vSphere with NVIDIA vGPU technology. This enables multiple workloads to share physical GPU resources while maintaining isolation and performance.

### vGPU Configuration

VKS exposes vGPU devices through VM Classes that are configured with NVIDIA vGPU profiles. These profiles define the amount of GPU memory and compute resources allocated to each virtual machine.

**VM Class with vGPU Profile:**

The cluster's worker node pool is configured with the VM Class, which carries a vGPU profile:
- `GRID A100-10C`: 10GB vGPU profile

This profile is based on NVIDIA GRID technology running on VMware vSphere, providing virtualized access to an NVIDIA A100 GPU.

### vGPU Evidence in Node Labels

The presence of vGPU support is verified through node labels that are automatically populated by the NVIDIA GPU Operator:

```shell
kubectl get nodes -o json | jq '.items[].metadata.labels' | grep vgpu
```

```shell
"nvidia.com/vgpu.host-driver-branch": "r596_25",
"nvidia.com/vgpu.host-driver-version": "595.71.03",
"nvidia.com/vgpu.present": "true",
"nvidia.com/vgpu.host-driver-branch": "r596_25",
"nvidia.com/vgpu.host-driver-version": "595.71.03",
"nvidia.com/vgpu.present": "true",
```

```shell
kubectl get nodes -o json | jq '.items[].metadata.labels' | grep "gpu.product"
```

```shell
"nvidia.com/gpu.product": "GRID-A100-10C",
"nvidia.com/gpu.product": "GRID-A100-10C",
```

Key indicators of vGPU support:
- `nvidia.com/vgpu.present: "true"` - Confirms vGPU is available on the node
- `nvidia.com/vgpu.host-driver-version` - Shows the host driver version managing the vGPU
- `nvidia.com/gpu.product` - Shows GRID-prefixed product names, indicating vGPU profiles

### vGPU in DRA ResourceSlice Attributes

The Dynamic Resource Allocation (DRA) framework exposes vGPU devices as schedulable resources. The ResourceSlice attributes show the brand as `NvidiaVCS` (NVIDIA vCompute Server), which is the vGPU-specific identifier:

```shell
kubectl get resourceslices -o yaml | grep -A 20 "brand:"
```

```shell
brand:
  string: NvidiaVCS
cudaComputeCapability:
  version: 8.0.0
cudaDriverVersion:
  version: 13.0.0
driverVersion:
  version: 580.126.9
productName:
  string: GRID A100-10C
resource.kubernetes.io/pciBusID:
  string: "0000:03:00.0"
resource.kubernetes.io/pcieRoot:
  string: pci0000:03
type:
  string: gpu
uuid:
  string: GPU-0b580e61-4daa-462a-95e3-84dd00000000
capacity:
  memory:
    value: 10Gi
--
brand:
  string: NvidiaVCS
cudaComputeCapability:
  version: 8.0.0
cudaDriverVersion:
  version: 13.0.0
driverVersion:
  version: 580.126.9
productName:
  string: GRID A100-10C
resource.kubernetes.io/pciBusID:
  string: "0000:03:00.0"
resource.kubernetes.io/pcieRoot:
  string: pci0000:03
type:
  string: gpu
uuid:
  string: GPU-44d2411e-1604-478f-aa96-5a3c00000000
capacity:
  memory:
    value: 10Gi
```

The `brand: NvidiaVCS` attribute specifically indicates that these are virtualized GPU resources managed through [NVIDIA's vCompute Server technology](https://www.nvidia.com/content/dam/en-zz/Solutions/design-visualization/solutions/resources/documents1/nvidia-virtual-compute-server-solution-overview.pdf).

### Workload Scheduling with vGPU

Workloads can request vGPU resources through standard Kubernetes resource requests. The DRA framework and NVIDIA GPU Operator handle the scheduling and allocation of vGPU devices to containers.

**Example from dra_support.md:**

The DRA deployment in `dra_support.md` demonstrates a workload successfully scheduled on a vGPU-enabled node:

```shell
kubectl get pods -n gpu-test1
```

```shell
NAME                              READY   STATUS    RESTARTS   AGE
dra-gpu-example-798fccd4d-sfl2p   1/1     Running   0          16m
```

The pod logs confirm access to the vGPU device:

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example --tail=5
```

```shell
GPU 0: GRID A100-10C (UUID: GPU-44d2411e-1604-478f-aa96-5a3c00000000)
Mon Jul  6 13:31:26 UTC 2026
GPU 0: GRID A100-10C (UUID: GPU-44d2411e-1604-478f-aa96-5a3c00000000)
Mon Jul  6 13:32:26 UTC 2026
GPU 0: GRID A100-10C (UUID: GPU-44d2411e-1604-478f-aa96-5a3c00000000)
```

The UUID reported inside the container (`GPU-44d2411e-1604-478f-aa96-5a3c00000000`) matches the `gpu-0` device UUID allocated to node `conformance-test-np1-qcqc8-bpdtl-wv2gs` in the ResourceSlice shown above, confirming the vGPU device was correctly scheduled and is usable by the workload via DRA.

### Configuration Documentation

For detailed instructions on configuring vGPU on VKS, refer to the official VMware documentation:

https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Cross-References

The following evidence files demonstrate vGPU functionality across different AI conformance requirements:

- **dra_support.md**: Shows the vGPU device (GRID A100-10C) exposed through DRA ResourceSlices and successfully allocated to a workload
- **pod_autoscaling.md**: Demonstrates HPA functionality with vGPU-backed pods using DCGM metrics
- **gang_scheduling.md**: Shows Kueue scheduling workloads across vGPU-enabled node pools
- **secure_accelerator_access.md**: Validates proper isolation and access control for vGPU resources
