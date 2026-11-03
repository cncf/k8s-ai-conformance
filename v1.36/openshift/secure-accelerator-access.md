# Secure Accelerator Access — Conformance Evidence

**Requirement (`secure_accelerator_access`, MUST):** Ensure that access to accelerators
from within containers is properly isolated and mediated by the Kubernetes resource
management framework (device plugin or DRA) and container runtime, preventing
unauthorized access or interference between workloads.

**Cluster:** OpenShift Container Platform 5.0 RC.7 (AWS, `g5.xlarge`, 1x NVIDIA A10G per node)
**Kubernetes version:** v1.36
**Date:** 2026-10-26

## Summary

On OpenShift 5.0, GPU access is granted exclusively through the Kubernetes device plugin
framework and the NVIDIA container runtime. A workload receives a device only by requesting
it via `nvidia.com/gpu` resource limits; the container runtime then injects exactly the
claimed device. The automated conformance test `TestSecureAcceleratorAccess` verifies that
(1) accelerator access is mediated by the GPU Operator stack and the resource framework,
(2) a workload not requesting GPU resources sees zero accelerators, and (3) in a
multi-container pod, only the container granted the GPU can access it.

## 1. Accelerator access is mediated by the platform

The NVIDIA GPU Operator v26.7.2 manages the driver, container runtime configuration, and
device plugin as DaemonSets bound to GPU nodes. Driver version 615.71.09 is compiled via
the OpenShift Driver Toolkit for RHCOS 10.2 (kernel 6.12.0-211.66.1.el10_2.x86_64).

The GPU Operator components on each GPU node:
- `nvidia-driver-daemonset` — compiles and loads kernel modules
- `nvidia-container-toolkit-daemonset` — configures the container runtime
- `nvidia-device-plugin-daemonset` — registers `nvidia.com/gpu` extended resource
- `nvidia-operator-validator` — gates node readiness

Access is brokered through the device plugin: a pod must include
`resources.limits: nvidia.com/gpu: 1` to obtain a device, and the device is injected by
the NVIDIA container runtime rather than mounted directly from the host.

## 2. Automated test results (e2e.log)

The upstream conformance test `TestSecureAcceleratorAccess` passed with three subtests:

```
--- PASS: TestSecureAcceleratorAccess (119.56s)
    --- PASS: TestSecureAcceleratorAccess/PositiveAccessTest (36.90s)
    --- PASS: TestSecureAcceleratorAccess/NegativeIsolationTest (36.89s)
    --- PASS: TestSecureAcceleratorAccess/MultiContainerIsolationTest (37.45s)
```

### PositiveAccessTest
A pod requesting 1 GPU (`nvidia.com/gpu: 1`) saw exactly 1 accelerator device:
```
RESULT: ACCELERATOR_COUNT=1
```

### NegativeIsolationTest
A pod on the same node without GPU resource requests saw zero accelerators:
```
RESULT: ACCELERATOR_COUNT=0
```

### MultiContainerIsolationTest
In a multi-container pod, the `authorized` container (with GPU request) saw 1 device,
while the `unauthorized` container (without GPU request) saw 0:
```
authorized:   RESULT: ACCELERATOR_COUNT=1
unauthorized: RESULT: ACCELERATOR_COUNT=0
```

## 3. Security context

Test pods ran under the OpenShift `restricted-v2` Security Context Constraints (SCC)
with no host device mounts and no privileged access, confirming that GPU access is
fully mediated by the device plugin and container runtime framework.
