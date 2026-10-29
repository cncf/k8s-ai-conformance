# GPU Sharing — Conformance Evidence

**Requirement (`gpu_sharing`, SHOULD):** For accelerators that support static GPU sharing,
provide well-defined mechanisms for at least one GPU sharing strategy to improve
utilization for workloads that do not require a full dedicated GPU.

**Cluster:** OpenShift Container Platform 5.0 RC.7 (AWS, `g5.xlarge`, 1x NVIDIA A10G per node)
**Kubernetes version:** v1.36
**Date:** 2026-10-26

## Summary

OpenShift 5.0 supports GPU sharing via the NVIDIA GPU Operator, which provides two
mechanisms: software-based time-slicing and hardware-level Multi-Instance GPU (MIG)
partitioning (on supported hardware such as A100/H100).

## 1. Time-slicing

Time-slicing allows multiple workloads to share a single physical GPU by oversubscribing
the `nvidia.com/gpu` extended resource. Configuration is done through the ClusterPolicy:

```yaml
spec:
  devicePlugin:
    config:
      name: time-slicing-config
      default: any
```

With a ConfigMap defining the slicing policy:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: time-slicing-config
  namespace: nvidia-gpu-operator
data:
  any: |-
    version: v1
    sharing:
      timeSlicing:
        resources:
        - name: nvidia.com/gpu
          replicas: 4
```

This exposes each physical GPU as 4 schedulable units, allowing 4 workloads to
share a single GPU via CUDA time-slicing.

## 2. Multi-Instance GPU (MIG)

For MIG-capable accelerators (A100, H100), the GPU Operator supports hardware-level
partitioning that creates isolated GPU instances with dedicated compute, memory, and
cache resources:

```yaml
spec:
  mig:
    strategy: mixed  # or "single"
```

MIG profiles (e.g., `1g.5gb`, `2g.10gb`, `3g.20gb`) are configured per-node via labels
and exposed as distinct schedulable resources (e.g., `nvidia.com/mig-1g.5gb`).

## 3. Documentation references

- [Time-slicing GPUs in OpenShift](https://docs.nvidia.com/datacenter/cloud-native/openshift/latest/time-slicing-gpus-in-openshift.html)
- [MIG on OpenShift](https://docs.nvidia.com/datacenter/cloud-native/openshift/latest/mig-ocp.html)

Note: The NVIDIA A10G (Ampere) used in this conformance cluster does not support MIG.
Time-slicing is the applicable sharing mechanism for this GPU type.
