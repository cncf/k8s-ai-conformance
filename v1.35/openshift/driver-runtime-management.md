# Driver and Runtime Management — Conformance Evidence

**Requirement (`driver_runtime_management`, SHOULD):** Provide a verifiable mechanism for
ensuring that compatible accelerator drivers and corresponding container runtime
configurations are correctly installed and maintained on nodes with accelerators. Once
the accelerator supports exposing driver and runtime version information as part of DRA,
then the platform should use the DRA mechanism for verification.

**Cluster:** OpenShift Container Platform 4.22.0 (GCP, `n1-standard-8`, 2x NVIDIA Tesla T4)
**Kubernetes version:** v1.35.5
**Date:** 2026-06-09

## Summary

On OpenShift 4.22 the NVIDIA GPU Operator installs and maintains the full accelerator
driver and container-runtime stack on GPU nodes, and gates node readiness on a set of
validations. In addition, the NVIDIA DRA Driver for GPUs publishes the installed driver
and CUDA versions as device attributes in ResourceSlices, providing a DRA-native,
machine-readable mechanism for verifying driver and runtime state.

## 1. Driver and runtime components are installed and maintained

The GPU Operator deploys and reconciles the driver, the NVIDIA Container Toolkit (runtime
configuration), and the device plugin as DaemonSets bound to GPU nodes:

```
$ oc get ds -n nvidia-gpu-operator
NAME                                      DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   AGE
nvidia-container-toolkit-daemonset        1         1         1       1            1           60m
nvidia-device-plugin-daemonset            1         1         1       1            1           60m
nvidia-driver-daemonset-9.8.20260520-0    1         1         1       1            1           60m
nvidia-operator-validator                 1         1         1       1            1           60m
```

```
$ oc get pods -n nvidia-gpu-operator
NAME                                           READY   STATUS      RESTARTS   AGE
nvidia-container-toolkit-daemonset-7bqvx       1/1     Running     0          60m
nvidia-cuda-validator-brppm                    0/1     Completed   0          56m
nvidia-driver-daemonset-9.8.20260520-0-kmbqm   2/2     Running     0          61m
nvidia-operator-validator-tzwql                1/1     Running     0          60m
```

> Output abridged to the driver/runtime components relevant to this requirement (and
> the `NODE SELECTOR` column is omitted for width). The operator also deploys
> `gpu-feature-discovery`, `nvidia-dcgm`, `nvidia-dcgm-exporter`, `nvidia-mig-manager`,
> `nvidia-device-plugin-mps-control-daemon`, and `nvidia-node-status-exporter`
> DaemonSets in the same namespace; the full pod listing appears in
> [secure-accelerator-access.md](secure-accelerator-access.md).

The driver DaemonSet is selected onto nodes by OS version
(`feature.node.kubernetes.io/system-os_release.OSTREE_VERSION=9.8.20260520-0`), ensuring
the driver build matches the node's RHCOS version — i.e. the operator maintains a
*compatible* driver for each node.

## 2. Installed driver version

The driver loaded on the node, queried from the driver container:

```
$ oc exec -n nvidia-gpu-operator nvidia-driver-daemonset-9.8.20260520-0-kmbqm \
    -c nvidia-driver-ctr -- nvidia-smi --query-gpu=name,driver_version,index --format=csv
name, driver_version, index
Tesla T4, 580.126.20, 0
Tesla T4, 580.126.20, 1
```

## 3. Readiness is gated on validation

The `nvidia-operator-validator` verifies the driver, toolkit, CUDA, and device-plugin are
correctly installed before the node is marked ready for GPU workloads:

```
$ oc get pods -n nvidia-gpu-operator -l app=nvidia-operator-validator
NAME                              READY   STATUS    RESTARTS   AGE
nvidia-operator-validator-tzwql   1/1     Running   0          60m

$ oc logs -n nvidia-gpu-operator -l app=nvidia-operator-validator -c nvidia-operator-validator
all validations are successful
```

The `nvidia-cuda-validator` (Completed, above) additionally confirms the runtime can
launch a CUDA workload end to end.

## 4. Driver and CUDA versions exposed through DRA

The NVIDIA DRA Driver for GPUs advertises the driver version, CUDA driver version, and
compute capability as device attributes in each ResourceSlice. This provides the
DRA-native verification mechanism described in the requirement — driver and runtime
information is queryable through the `resource.k8s.io` API:

```
$ oc get resourceslices -o yaml
...
  devices:
  - name: gpu-0
    attributes:
      architecture: {string: Turing}
      cudaComputeCapability: {version: 7.5.0}
      cudaDriverVersion: {version: 13.0.0}
      driverVersion: {version: 580.126.20}
      productName: {string: Tesla T4}
...
```

The `driverVersion` (`580.126.20`) reported through DRA matches the driver reported by
`nvidia-smi` on the node, confirming the DRA attributes are an accurate source of truth
for verification.

## Result

**PASS.** OpenShift 4.22, via the NVIDIA GPU Operator, installs and maintains a compatible
driver and container-runtime configuration on GPU nodes and gates node readiness on
successful validation. Driver and CUDA versions are additionally exposed through the
DRA `resource.k8s.io` API, satisfying the `driver_runtime_management` requirement
including its forward-looking DRA-based verification mechanism.
