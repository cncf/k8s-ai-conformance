## Description

If the platform supports virtualized accelerators (e.g., vGPU), it should demonstrate that these resources can be scheduled and utilized by workloads.

## Evidence

VKS supports virtualized GPU resources through VMware vSphere with NVIDIA vGPU technology. This enables multiple workloads to share physical GPU resources while maintaining isolation and performance.

### vGPU Configuration

VKS exposes vGPU devices through VM Classes that are configured with NVIDIA vGPU profiles. These profiles define the amount of GPU memory and compute resources allocated to each virtual machine.

**Example VM Class with vGPU Profile:**

The cluster used in this demonstration is configured with VM classes that include vGPU profiles:
- `GRID A100-40C`: 40GB vGPU profile
- `GRID A100-4C`: 4GB vGPU profile

These profiles are based on NVIDIA GRID technology running on VMware vSphere, providing virtualized access to NVIDIA A100 GPUs.

### vGPU Evidence in Node Labels

The presence of vGPU support is verified through node labels that are automatically populated by the NVIDIA GPU Operator:

```shell
kubectl get nodes -o json | jq '.items[].metadata.labels' | grep vgpu
```

```shell
"nvidia.com/vgpu.present": "true"
"nvidia.com/vgpu.host-driver-branch": "r582_12"
"nvidia.com/vgpu.host-driver-version": "580.126.08"
"nvidia.com/gpu.product": "GRID-A100-40C"
"nvidia.com/gpu.product": "GRID-A100-4C"
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
  string: GRID A100-40C
```

The `brand: NvidiaVCS` attribute specifically indicates that these are virtualized GPU resources managed through [NVIDIA's vCompute Server technology](https://www.nvidia.com/content/dam/en-zz/Solutions/design-visualization/solutions/resources/documents1/nvidia-virtual-compute-server-solution-overview.pdf).

### Workload Scheduling with vGPU

Workloads can request vGPU resources through standard Kubernetes resource requests. The DRA framework and NVIDIA GPU Operator handle the scheduling and allocation of vGPU devices to containers.

**Example from dra_support.md:**

The DRA deployment examples in `dra_support.md` demonstrate workloads successfully scheduled on vGPU-enabled nodes:

```shell
kubectl get pods -n gpu-test1
NAME                              READY   STATUS    RESTARTS   AGE
dra-gpu-example-7c7d5856fb-mbn76   1/1     Running   0          81s
```

The pod logs confirm access to the vGPU device:

```shell
kubectl logs dra-gpu-example-7c7d5856fb-mbn76 -n gpu-test1
Tue Mar 17 13:18:00 UTC 2026
GPU 0: GRID A100-40C (UUID: GPU-0946f70c-659b-4c39-8b78-7e7600000000)
```

### Configuration Documentation

For detailed instructions on configuring vGPU on VKS, refer to the official VMware documentation:

https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Cross-References

The following evidence files demonstrate vGPU functionality across different AI conformance requirements:

- **dra_support.md**: Shows vGPU devices (GRID A100-40C, GRID A100-4C) exposed through DRA ResourceSlices and successfully allocated to workloads
- **pod_autoscaling.md**: Demonstrates HPA functionality with vGPU-backed pods using DCGM metrics
- **gang_scheduling.md**: Shows Kueue scheduling workloads across vGPU-enabled node pools
- **secure_accelerator_access.md**: Validates proper isolation and access control for vGPU resources
