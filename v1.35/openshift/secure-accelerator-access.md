# Secure Accelerator Access — Conformance Evidence

**Requirement (`secure_accelerator_access`, MUST):** Ensure that access to accelerators
from within containers is properly isolated and mediated by the Kubernetes resource
management framework (device plugin or DRA) and container runtime, preventing
unauthorized access or interference between workloads.

**Cluster:** OpenShift Container Platform 4.22.0 (GCP, `n1-standard-8`, 2x NVIDIA Tesla T4)
**Kubernetes version:** v1.35.5
**Date:** 2026-06-09

## Summary

On OpenShift 4.22, GPU access is granted exclusively through the Kubernetes resource
management framework and the NVIDIA container runtime. A workload receives a device only
by claiming it (here via a DRA ResourceClaim); the container runtime then injects exactly
the claimed device. This evidence shows that (1) accelerator access is mediated by the
GPU Operator stack and the resource framework, and (2) a workload granted one GPU is
isolated from the other GPU on the same node, with no path to unauthorized access.

## 1. Accelerator access is mediated by the platform

The NVIDIA GPU Operator manages the driver, container runtime configuration, and device
lifecycle, and reports a healthy state:

```
$ oc get clusterpolicy gpu-cluster-policy
NAME                 STATUS   AGE
gpu-cluster-policy   ready    44m
```

```
$ oc get pods -n nvidia-gpu-operator
NAME                                           READY   STATUS      RESTARTS      AGE
gpu-feature-discovery-rdg5s                    1/1     Running     0             41m
gpu-operator-85656dd79d-2zt64                  1/1     Running     0             44m
nvidia-container-toolkit-daemonset-7bqvx       1/1     Running     0             41m
nvidia-cuda-validator-brppm                    0/1     Completed   0             37m
nvidia-dcgm-exporter-2pkvm                     1/1     Running     3 (36m ago)   41m
nvidia-dcgm-m9bhs                              1/1     Running     0             41m
nvidia-device-plugin-daemonset-fmhdb           1/1     Running     0             41m
nvidia-driver-daemonset-9.8.20260520-0-kmbqm   2/2     Running     0             41m
nvidia-node-status-exporter-hxgtq              1/1     Running     0             41m
nvidia-operator-validator-tzwql                1/1     Running     0             41m
```

Access is brokered through the resource management framework: a pod must reference a
ResourceClaim to obtain a device, and the device is injected by the NVIDIA container
runtime rather than mounted directly from the host. The isolation test pod below
declares its GPU through `resourceClaims` and carries no host device mounts:

```
$ oc get pod isolation-test -n secure-access-test -o jsonpath='{.spec.resourceClaims}'
[{"name":"gpu","resourceClaimName":"isolated-gpu"}]
```

```
$ oc get pod isolation-test -n secure-access-test -o json | jq '[.spec.volumes[] | select(.hostPath)] | length'
0
```

The pod is admitted under OpenShift's default, most restrictive Security Context
Constraint, confirming no elevated host privileges are required for accelerator access:

```
$ oc get pod isolation-test -n secure-access-test -o jsonpath='{.metadata.annotations.openshift\.io/scc}'
restricted-v2
```

## 2. Workload isolation

The GPU node hosts two Tesla T4 GPUs. The following workload claims a single GPU and
inspects which devices are visible inside the container:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: secure-access-test
---
apiVersion: resource.k8s.io/v1
kind: ResourceClaim
metadata:
  name: isolated-gpu
  namespace: secure-access-test
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
  name: isolation-test
  namespace: secure-access-test
spec:
  restartPolicy: Never
  securityContext:
    runAsNonRoot: true
    seccompProfile:
      type: RuntimeDefault
  resourceClaims:
    - name: gpu
      resourceClaimName: isolated-gpu
  containers:
    - name: gpu-test
      image: nvidia/cuda:12.9.0-base-ubuntu24.04
      command: ["bash","-c","echo '=== Visible NVIDIA devices ==='; ls -l /dev/nvidia*; echo; echo '=== nvidia-smi output ==='; nvidia-smi -L; echo; echo '=== GPU count visible to container ==='; nvidia-smi -L | wc -l; echo 'Secure accelerator access test completed'"]
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
      resources:
        claims:
          - name: gpu
```

```
$ oc apply -f secure-accelerator-access.yaml
namespace/secure-access-test created
resourceclaim.resource.k8s.io/isolated-gpu created
pod/isolation-test created
```

```
$ oc get pod isolation-test -n secure-access-test
NAME             READY   STATUS      RESTARTS   AGE
isolation-test   0/1     Completed   0          19m
```

```
$ oc logs isolation-test -n secure-access-test
=== Visible NVIDIA devices ===
crw-rw-rw-. 1 1000770000 root 195, 254 Jun  9 18:35 /dev/nvidia-modeset
crw-rw-rw-. 1 1000770000 root 511,   0 Jun  9 18:35 /dev/nvidia-uvm
crw-rw-rw-. 1 1000770000 root 511,   1 Jun  9 18:35 /dev/nvidia-uvm-tools
crw-rw-rw-. 1 1000770000 root 195,   0 Jun  9 18:35 /dev/nvidia0
crw-rw-rw-. 1 1000770000 root 195, 255 Jun  9 18:35 /dev/nvidiactl

=== nvidia-smi output ===
GPU 0: Tesla T4 (UUID: GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174)

=== GPU count visible to container ===
1
Secure accelerator access test completed
```

Although the node has two GPUs, the container can see and access only the single GPU it
was granted (`/dev/nvidia0`, `GPU-da96f367-...`). The second GPU (`GPU-204349eb-...`) is
not present in the container's device namespace, so the workload has no way to reach or
interfere with accelerators allocated to other workloads.

## Result

**PASS.** Accelerator access on OpenShift 4.22 is mediated by the Kubernetes resource
management framework and the NVIDIA container runtime, requires no privileged host
access, and confines each workload to exactly the device it was granted — satisfying the
`secure_accelerator_access` requirement.

## Cleanup

```
$ oc delete ns secure-access-test
```
