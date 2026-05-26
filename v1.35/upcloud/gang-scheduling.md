## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator

## Enable time-slicing
Follow guide at [util-time-slicing.md](util-time-slicing.md)

## Install Kueue

Run:
```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/download/v0.16.1/manifests.yaml
```

## Apply manifests to test

```bash
kubectl apply -f - <<EOF
apiVersion: kueue.x-k8s.io/v1beta2
kind: ResourceFlavor
metadata:
  name: "l40s-slice"
---
apiVersion: kueue.x-k8s.io/v1beta2
kind: ClusterQueue
metadata:
  name: "gpu-cluster-queue"
spec:
  namespaceSelector: {}
  resourceGroups:
  - coveredResources: ["nvidia.com/gpu.shared"]
    flavors:
    - name: "l40s-slice"
      resources:
      - name: "nvidia.com/gpu.shared"
        nominalQuota: 10
---
apiVersion: kueue.x-k8s.io/v1beta2
kind: LocalQueue
metadata:
  namespace: default
  name: "gpu-local-queue"
spec:
  clusterQueue: "gpu-cluster-queue"
EOF
```

## Run tests
Deploy following manifests to test how Kueue manages the workloads

### Job 1:
```bash
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: gpu-job-a
  labels:
    kueue.x-k8s.io/queue-name: gpu-local-queue
spec:
  template:
    spec:
      containers:
      - name: cuda-test
        image: nvidia/cuda:12.0.0-base-ubuntu22.04
        command: ["sleep", "60"]
        resources:
          requests:
            "nvidia.com/gpu.shared": 7
          limits:
            "nvidia.com/gpu.shared": 7
      restartPolicy: Never
EOF
```
### Job 2:
```bash
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: gpu-job-b
  labels:
    kueue.x-k8s.io/queue-name: gpu-local-queue
spec:
  template:
    spec:
      containers:
      - name: cuda-test
        image: nvidia/cuda:12.0.0-base-ubuntu22.04
        command: ["sleep", "60"]
        resources:
          requests:
            nvidia.com/gpu.shared: 5
          limits:
            nvidia.com/gpu.shared: 5
      restartPolicy: Never
EOF
```

## Check status

```bash
kubectl get workloads
NAME                  QUEUE             RESERVED IN         ADMITTED   FINISHED   AGE
job-gpu-job-a-da7a6   gpu-local-queue   gpu-cluster-queue   True                  14s
job-gpu-job-b-e2d2a   gpu-local-queue                                             8s
```

Describing node will also show that some slices are reserved
```bash
...
  nvidia.com/gpu         0           0
  nvidia.com/gpu.shared  7           7
...
```