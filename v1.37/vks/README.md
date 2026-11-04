# VKS — Kubernetes AI Conformance v1.37 evidence

Platform: vSphere Kubernetes Service (VKS) v3.8.0, Kubernetes `v1.37.1+vmware.4` (VKR `v1.37.1---vmware.4-vkr.1`).

## Test environment

| Item | Value |
| --- | --- |
| Nodes | 1 control plane + 2 worker nodes (node pool `np1`), Ubuntu 24.04.5, containerd 2.4.0+vmware.2-fips |
| Accelerators | 1 × NVIDIA `GRID A100-4C` vGPU (4 GiB) per worker, guest driver 580.178.04, CUDA 13.0, host vGPU driver 595.91.04 |
| GPU stack | NVIDIA GPU Operator v25.10.1 (`ClusterPolicy` ready) |
| Client tools | kubectl v1.37.1+vmware.4, Helm v4.3.1 |
| Cluster DNS domain | `managedcluster1.local` |

## Results

| Evidence | Requirement(s) | Result |
| --- | --- | --- |
| [dra_support.md](dra_support.md) | driver_runtime_management | Pass |
| [secure_accelerator_access.md](secure_accelerator_access.md) | secure_accelerator_access | Pass (DRA allocation, a pod with no claim sees no GPU, allocation is exclusive) |
| [virtualized_accelerator.md](virtualized_accelerator.md) | virtualized_accelerator | Pass |
| [accelerator_metrics.md](accelerator_metrics.md) | accelerator_metrics | Pass |
| [gang_scheduling.md](gang_scheduling.md) | gang_scheduling | Pass (quota queueing plus all-or-nothing admission) |
| [robust_controller.md](robust_controller.md) | robust_controller | Pass |
| [pod_autoscaling.md](pod_autoscaling.md) | pod_autoscaling | Pass (scaled 1→2→1 on `DCGM_FI_DEV_FB_FREE`) |
