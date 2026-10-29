# Driver and Runtime Management — Conformance Evidence

**Requirement (`driver_runtime_management`, SHOULD):** Provide a verifiable mechanism for
ensuring that compatible accelerator drivers and corresponding container runtime
configurations are correctly installed and maintained on nodes with accelerators.

**Cluster:** OpenShift Container Platform 5.0 RC.7 (AWS, `g5.xlarge`, 1x NVIDIA A10G per node)
**Kubernetes version:** v1.36
**Date:** 2026-10-26

## Summary

On OpenShift 5.0 the NVIDIA GPU Operator v26.7.2 installs and maintains the full
accelerator driver and container-runtime stack on GPU nodes, and gates node readiness
on a set of validations. The operator is installed from the certified-operators catalog
(v4.22 channel) and manages the complete lifecycle including driver compilation,
runtime configuration, and device plugin registration.

## 1. Driver and runtime components are installed and maintained

The GPU Operator deploys and reconciles the driver, the NVIDIA Container Toolkit (runtime
configuration), and the device plugin as DaemonSets bound to GPU nodes:

- `nvidia-driver-daemonset-10.2.20261015-0` — driver 615.71.09, compiled via OpenShift
  Driver Toolkit for kernel 6.12.0-211.66.1.el10_2.x86_64 (RHCOS 10.2)
- `nvidia-container-toolkit-daemonset` — configures containerd/cri-o for GPU passthrough
- `nvidia-device-plugin-daemonset` — registers `nvidia.com/gpu` extended resource
- `gpu-feature-discovery` — labels nodes with GPU attributes (compute capability, memory, etc.)
- `nvidia-dcgm` + `nvidia-dcgm-exporter` — GPU metrics via DCGM
- `nvidia-operator-validator` — validates the full stack before marking node ready

## 2. Node Feature Discovery labels GPU attributes

NFD (Node Feature Discovery) operator detects the NVIDIA PCI device and labels nodes:
```
feature.node.kubernetes.io/pci-10de.present=true
```

The GPU Operator then applies detailed labels:
```
nvidia.com/gpu.present=true
nvidia.com/gpu.product=NVIDIA-A10G
nvidia.com/gpu.memory=23028
nvidia.com/gpu.compute.major=8
nvidia.com/gpu.compute.minor=6
nvidia.com/gpu.family=ampere
nvidia.com/cuda.driver-version.full=615.71.09
nvidia.com/cuda.runtime-version.full=13.4
```

## 3. Validation gates node readiness

The `nvidia-operator-validator` DaemonSet runs init containers that sequentially verify:
1. Driver is loaded and functional
2. Container toolkit is configured
3. Device plugin is registered
4. CUDA workloads can execute

Only after all validations pass does the node report GPU capacity via `nvidia.com/gpu: 1`
in its allocatable resources.

## 4. ClusterPolicy configuration

The GPU Operator is configured via a ClusterPolicy custom resource (`gpu-cluster-policy`)
that specifies:
- `driver.version: 615.71.09` — required for RHEL 10 / RHCOS 10.2 kernel compatibility
- Automatic DaemonSet deployment to nodes with `nvidia.com/gpu.deploy.*=true` labels
- Tolerations for GPU-tainted nodes
