# GPU Sharing — Conformance Evidence

**Requirement (`gpu_sharing`, SHOULD):** For accelerators that support static GPU
sharing, provide well-defined mechanisms for at least one GPU sharing strategy to improve
utilization for workloads that do not require a full dedicated GPU. If hardware-level
partitioning is supported, then these fractional GPU resources should be exposed as
schedulable resources. If software-based sharing (e.g. time-slicing) is supported, then
oversubscription of GPUs should be allowed.

**Cluster:** OpenShift Container Platform 4.22.0 (GCP, `n1-standard-8`, 2x NVIDIA Tesla T4)
**Kubernetes version:** v1.35.5
**Date:** 2026-06-09

## Summary

OpenShift 4.22, via the NVIDIA GPU Operator, supports software-based GPU sharing through
time-slicing, which allows oversubscription of physical GPUs. This evidence configures
time-slicing with 4 replicas per GPU, shows the node advertising 8 schedulable GPUs from
2 physical devices, and demonstrates four workloads concurrently sharing the two physical
GPUs.

> The two T4 GPUs in this cluster are not MIG-capable, so hardware-level partitioning
> (MIG) is not exercised here; the requirement is satisfied by demonstrating one sharing
> strategy. MIG on MIG-capable accelerators (e.g. A100/H100) is configured through the
> same GPU Operator and is referenced in the NVIDIA OpenShift documentation.

## 1. Enable time-slicing

A device-plugin configuration requesting 4 time-sliced replicas per GPU:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: time-slicing-config
  namespace: nvidia-gpu-operator
data:
  any: |-
    version: v1
    flags:
      migStrategy: none
    sharing:
      timeSlicing:
        resources:
        - name: nvidia.com/gpu
          replicas: 4
```

The GPU Operator's ClusterPolicy is pointed at the configuration, after which the device
plugin re-rolls automatically:

```
$ oc patch clusterpolicy gpu-cluster-policy --type merge \
    -p '{"spec":{"devicePlugin":{"config":{"name":"time-slicing-config","default":"any"}}}}'
clusterpolicy.nvidia.com/gpu-cluster-policy patched
```

## 2. Oversubscription: 2 physical GPUs advertised as 8 schedulable

Before enabling time-slicing the node advertised one schedulable unit per physical GPU:

```
$ oc get node harpatil4220c-g9kfm-gpu-c-rd8r8 -o jsonpath='{.status.capacity.nvidia\.com/gpu}'
2
```

After enabling time-slicing (4 replicas x 2 physical GPUs) it advertises 8:

```
$ oc get node harpatil4220c-g9kfm-gpu-c-rd8r8 -o jsonpath='{.status.capacity.nvidia\.com/gpu}'
8
```

## 3. Workloads share the physical GPUs

A Deployment of 4 replicas, each requesting `nvidia.com/gpu: 1`, is scheduled onto the
node — more GPU-requesting pods than there are physical GPUs:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ts-share
  namespace: gpu-sharing-test
spec:
  replicas: 4
  selector:
    matchLabels: {app: ts-share}
  template:
    metadata:
      labels: {app: ts-share}
    spec:
      securityContext:
        runAsNonRoot: true
        seccompProfile: {type: RuntimeDefault}
      containers:
        - name: cuda
          image: nvidia/cuda:12.9.0-base-ubuntu24.04
          command: ["bash","-c","nvidia-smi -L; sleep 3600"]
          securityContext:
            allowPrivilegeEscalation: false
            capabilities: {drop: ["ALL"]}
          resources:
            limits:
              nvidia.com/gpu: 1
```

All four pods are admitted and running on the single two-GPU node:

```
$ oc get pods -n gpu-sharing-test -o wide
NAME                        READY   STATUS    RESTARTS   AGE   IP            NODE                              ...
ts-share-6774cbb654-btf44   1/1     Running   0          11s   10.130.2.38   harpatil4220c-g9kfm-gpu-c-rd8r8   ...
ts-share-6774cbb654-cplqx   1/1     Running   0          11s   10.130.2.36   harpatil4220c-g9kfm-gpu-c-rd8r8   ...
ts-share-6774cbb654-jwfs6   1/1     Running   0          11s   10.130.2.37   harpatil4220c-g9kfm-gpu-c-rd8r8   ...
ts-share-6774cbb654-zkxn6   1/1     Running   0          11s   10.130.2.39   harpatil4220c-g9kfm-gpu-c-rd8r8   ...
```

Each pod reports the physical GPU it was assigned. Across the four pods only the two
physical GPU UUIDs appear — each physical GPU is shared by two pods:

```
$ for p in $(oc get pods -n gpu-sharing-test -o name); do oc logs -n gpu-sharing-test $p; done
GPU 0: Tesla T4 (UUID: GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174)
GPU 0: Tesla T4 (UUID: GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174)
GPU 0: Tesla T4 (UUID: GPU-204349eb-0405-e927-278b-8804172f1be3)
GPU 0: Tesla T4 (UUID: GPU-204349eb-0405-e927-278b-8804172f1be3)

# UUID tally across the 4 pods:
   2 GPU-204349eb-0405-e927-278b-8804172f1be3
   2 GPU-da96f367-6e47-8dc4-fd29-4f0b97b9f174
```

## Result

**PASS.** OpenShift 4.22 supports software-based GPU sharing via time-slicing through the
NVIDIA GPU Operator. Two physical GPUs are oversubscribed to 8 schedulable units, and
four workloads run concurrently while sharing the two physical GPUs (two workloads per
GPU) — satisfying the `gpu_sharing` requirement.

## Cleanup

```
$ oc delete ns gpu-sharing-test
$ oc patch clusterpolicy gpu-cluster-policy --type json \
    -p '[{"op":"remove","path":"/spec/devicePlugin/config"}]'
$ oc delete configmap time-slicing-config -n nvidia-gpu-operator
```
