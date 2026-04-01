# Cluster Autoscaling

## Overview

Viettel AI Platform is deployed on bare-metal infrastructure (Ubuntu 22.04.5 LTS). The cluster autoscaler is not applicable in this context as the platform does not use cloud VMs that can be dynamically provisioned and deprovisioned. This is consistent with similar on-premises bare-metal AI platforms.

New GPU nodes are onboarded by bootstrapping servers with `kubeadm join`, after which the NVIDIA GPU Operator DaemonSet automatically installs drivers, configures the container runtime, and registers GPU resources.

## Cluster Nodes

```bash
$ kubectl get nodes -o wide
NAME          STATUS   ROLES           AGE   VERSION   INTERNAL-IP   OS-IMAGE              KERNEL-VERSION       CONTAINER-RUNTIME
ubuntu-sv16   Ready    control-plane   36h   v1.35.3   10.24.10.16   Ubuntu 22.04.5 LTS    5.15.0-173-generic   containerd://2.2.2
ubuntu-sv18   Ready    <none>          21h   v1.35.3   10.24.10.18   Ubuntu 22.04.5 LTS    5.15.0-173-generic   containerd://1.7.28
```

## Node GPU Resources

```bash
$ kubectl get nodes -o custom-columns='NODE:.metadata.name,GPU:.status.allocatable.nvidia\.com/gpu,CPU:.status.allocatable.cpu,MEM:.status.allocatable.memory'
NODE          GPU   CPU   MEM
ubuntu-sv16   2     96    394585252Ki
ubuntu-sv18   2     96    263434244Ki
```

Both nodes expose `nvidia.com/gpu: 2` (NVIDIA L40S) as allocatable resources. Total cluster GPU capacity: 4 NVIDIA L40S GPUs.

## Node Metrics

```bash
$ kubectl top nodes
NAME          CPU(cores)   CPU(%)   MEMORY(bytes)   MEMORY(%)
ubuntu-sv16   6905m        7%       46026Mi         11%
ubuntu-sv18   5206m        5%       5435Mi          2%
```

## HPA API Availability

```bash
$ kubectl api-resources --api-group=autoscaling
NAME                       SHORTNAMES   APIVERSION       NAMESPACED   KIND
horizontalpodautoscalers   hpa          autoscaling/v2   true         HorizontalPodAutoscaler
```

## Notes

For on-premises bare-metal deployments, the Cluster Autoscaler is not applicable since physical servers cannot be freely added and removed from the cluster like cloud VMs. The platform supports pod-level autoscaling via HPA with GPU metrics (see `pod_autoscaling.md`). Node expansion is performed through the standard Kubernetes node onboarding process with GPU Operator automation.
