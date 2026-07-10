## Description

The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

## Evidence

### Prerequisites

* Provision a VKS v3.7.0 Cluster with v1.36.1 node pool (2 replicas), VM Class with vGPU profile and NVIDIA GPU Operator.
* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Kueue

```shell
## Create namespace for kueue install
kubectl create ns kueue-system

## Install via helm
helm install kueue oci://registry.k8s.io/kueue/charts/kueue --version=0.14.1 --namespace kueue-system --wait --timeout 300s
```

```shell
Pulled: registry.k8s.io/kueue/charts/kueue:0.14.1
Digest: sha256:b146879997b68f355b730da28413adb3fba1343d352f3dda4f9956b3a3bcd3ce
NAME: kueue
LAST DEPLOYED: Mon Jul  6 13:49:07 2026
NAMESPACE: kueue-system
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Wait until Kueue pods are ready

```shell
kubectl get deployments -n kueue-system
```

```shell
NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
kueue-controller-manager   1/1     1            1           2m9s
```

```shell
kubectl get pods -n kueue-system
```

```shell
NAME                                       READY   STATUS    RESTARTS   AGE
kueue-controller-manager-f4b4d6dc6-r74lz   1/1     Running   0          2m9s
```

Create two new namespaces

```shell
kubectl create namespace team-a

kubectl create namespace team-b
```

### Create the ResourceFlavor

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: default-flavor
EOF
```

Verify ResourceFlavor

```shell
kubectl get resourceflavor
```

```shell
NAME             AGE
default-flavor   0s
```

### Create the ClusterQueue

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: cluster-queue
spec:
  namespaceSelector: {}
  queueingStrategy: BestEffortFIFO
  resourceGroups:
  - coveredResources: ["cpu", "memory", "nvidia.com/gpu", "ephemeral-storage"]
    flavors:
    - name: "default-flavor"
      resources:
      - name: "cpu"
        nominalQuota: 4
      - name: "memory"
        nominalQuota: 2Gi
      - name: "nvidia.com/gpu"
        nominalQuota: 2
      - name: "ephemeral-storage"
        nominalQuota: 10Gi
EOF
```

Inspect ClusterQueue resource

```shell
kubectl get ClusterQueue
```

```shell
NAME            COHORT   PENDING WORKLOADS
cluster-queue            0
```

```shell
kubectl describe ClusterQueue
```

```shell
Name:         cluster-queue
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  kueue.x-k8s.io/v1beta1
Kind:         ClusterQueue
Spec:
  Flavor Fungibility:
    When Can Borrow:   Borrow
    When Can Preempt:  TryNextFlavor
  Preemption:
    Borrow Within Cohort:
      Policy:               Never
    Reclaim Within Cohort:  Never
    Within Cluster Queue:   Never
  Queueing Strategy:        BestEffortFIFO
  Resource Groups:
    Covered Resources:
      cpu
      memory
      nvidia.com/gpu
      ephemeral-storage
    Flavors:
      Name:  default-flavor
      Resources:
        Name:           cpu
        Nominal Quota:  4
        Name:           memory
        Nominal Quota:  2Gi
        Name:           nvidia.com/gpu
        Nominal Quota:  2
        Name:           ephemeral-storage
        Nominal Quota:  10Gi
  Stop Policy:          None
Status:
  Admitted Workloads:  0
  Conditions:
    Message:               Can admit new workloads
    Reason:                Ready
    Status:                True
    Type:                  Active
  Pending Workloads:    0
  Reserving Workloads:  0
Events:                 <none>
```

### Create LocalQueue

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  namespace: team-a
  name: lq-team-a
spec:
  clusterQueue: cluster-queue
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  namespace: team-b
  name: lq-team-b
spec:
  clusterQueue: cluster-queue
EOF
```

Inspect LocalQueue resource

```shell
kubectl get LocalQueue -A
```

```shell
NAMESPACE   NAME        CLUSTERQUEUE    PENDING WORKLOADS   ADMITTED WORKLOADS
team-a      lq-team-a   cluster-queue   0                   0
team-b      lq-team-b   cluster-queue   0                   0
```

```shell
kubectl describe LocalQueue lq-team-a -n team-a
```

```shell
Name:         lq-team-a
Namespace:    team-a
API Version:  kueue.x-k8s.io/v1beta1
Kind:         LocalQueue
Spec:
  Cluster Queue:  cluster-queue
  Stop Policy:    None
Status:
  Admitted Workloads:  0
  Conditions:
    Message:               Can submit new workloads to localQueue
    Reason:                Ready
    Status:                True
    Type:                  Active
  Flavors:
    Name:  default-flavor
    Resources:
      cpu
      ephemeral-storage
      memory
      nvidia.com/gpu
  Pending Workloads:    0
  Reserving Workloads:  0
Events:                 <none>
```

### Create jobs

```shell
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  namespace: team-a
  generateName: sample-job-team-a-
  annotations:
    kueue.x-k8s.io/queue-name: lq-team-a
spec:
  ttlSecondsAfterFinished: 120
  parallelism: 1
  completions: 3
  suspend: true
  template:
    spec:
      nodeSelector:
        nvidia.com/gpu.product: "GRID-A100-10C"
      containers:
      - name: dummy-job
        image: gcr.io/k8s-staging-perf-tests/sleep:latest
        args: ["10s"]
        resources:
          requests:
            cpu: "500m"
            memory: "256Mi"
            ephemeral-storage: "256Mi"
            nvidia.com/gpu: "1"
          limits:
            cpu: "500m"
            memory: "256Mi"
            ephemeral-storage: "256Mi"
            nvidia.com/gpu: "1"
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
          runAsNonRoot: true
          runAsUser: 1000
          seccompProfile:
            type: RuntimeDefault
      restartPolicy: Never
EOF

cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  namespace: team-b
  generateName: sample-job-team-b-
  annotations:
    kueue.x-k8s.io/queue-name: lq-team-b
spec:
  ttlSecondsAfterFinished: 120
  parallelism: 1
  completions: 3
  suspend: true
  template:
    spec:
      nodeSelector:
        nvidia.com/gpu.product: "GRID-A100-10C"
      containers:
      - name: dummy-job
        image: gcr.io/k8s-staging-perf-tests/sleep:latest
        args: ["10s"]
        resources:
          requests:
            cpu: "500m"
            memory: "256Mi"
            ephemeral-storage: "256Mi"
            nvidia.com/gpu: "1"
          limits:
            cpu: "500m"
            memory: "256Mi"
            ephemeral-storage: "256Mi"
            nvidia.com/gpu: "1"
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
          runAsNonRoot: true
          runAsUser: 1000
          seccompProfile:
            type: RuntimeDefault
      restartPolicy: Never
EOF
```

Three jobs per team were submitted 10 seconds apart:

```shell
job.batch/sample-job-team-a-mbd4h created
job.batch/sample-job-team-b-fp5jw created
job.batch/sample-job-team-a-sqtjs created
job.batch/sample-job-team-b-zr4zj created
job.batch/sample-job-team-a-nptv9 created
job.batch/sample-job-team-b-pc9bl created
```

### Monitor the jobs being scheduled and executed

Immediately after creation, only one job per team was admitted (2 GPUs = full ClusterQueue quota); the rest were correctly held `Suspended`:

```shell
kubectl get jobs -n team-a
```

```shell
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-a-mbd4h   Running     2/3           39s        40s
sample-job-team-a-nptv9   Suspended   0/3                      19s
sample-job-team-a-sqtjs   Suspended   0/3                      29s
```

```shell
kubectl get jobs -n team-b
```

```shell
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-b-fp5jw   Running     2/3           39s        39s
sample-job-team-b-pc9bl   Suspended   0/3                      19s
sample-job-team-b-zr4zj   Suspended   0/3                      19s
```

The ClusterQueue events confirm quota-based admission and back-pressure on the suspended workloads:

```shell
kubectl get events -n team-a --sort-by='.lastTimestamp'
```

```shell
Normal    QuotaReserved      job-sample-job-team-a-mbd4h-1b468   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 1s
Normal    Admitted           job-sample-job-team-a-mbd4h-1b468   Admitted by ClusterQueue cluster-queue, wait time since reservation was 0s
Normal    Started            sample-job-team-a-mbd4h             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-mbd4h             Job resumed
Warning   Pending            job-sample-job-team-a-sqtjs-4fbeb   couldn't assign flavors to pod set main: insufficient unused quota for nvidia.com/gpu in flavor default-flavor, 1 more needed
Warning   Pending            job-sample-job-team-a-nptv9-60f04   couldn't assign flavors to pod set main: insufficient unused quota for nvidia.com/gpu in flavor default-flavor, 1 more needed
Normal    Completed          sample-job-team-a-mbd4h             Job completed
Normal    FinishedWorkload   sample-job-team-a-mbd4h             Workload 'team-a/job-sample-job-team-a-mbd4h-1b468' is declared finished
Normal    QuotaReserved      job-sample-job-team-a-sqtjs-4fbeb   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 35s
Normal    Admitted           job-sample-job-team-a-sqtjs-4fbeb   Admitted by ClusterQueue cluster-queue, wait time since reservation was 0s
Normal    Started            sample-job-team-a-sqtjs             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-sqtjs             Job resumed
Normal    Completed          sample-job-team-a-sqtjs             Job completed
Normal    FinishedWorkload   sample-job-team-a-sqtjs             Workload 'team-a/job-sample-job-team-a-sqtjs-4fbeb' is declared finished
Normal    QuotaReserved      job-sample-job-team-a-nptv9-60f04   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 67s
Normal    Admitted           job-sample-job-team-a-nptv9-60f04   Admitted by ClusterQueue cluster-queue, wait time since reservation was 0s
Normal    Started            sample-job-team-a-nptv9             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-nptv9             Job resumed
Normal    Completed          sample-job-team-a-nptv9             Job completed
Normal    FinishedWorkload   sample-job-team-a-nptv9             Workload 'team-a/job-sample-job-team-a-nptv9-60f04' is declared finished
```

As each running job freed its `nvidia.com/gpu` quota unit on completion, Kueue admitted the next queued workload in FIFO order — `sqtjs` and `zr4zj` were admitted once `mbd4h`/`fp5jw` completed, then `nptv9` and `pc9bl` were admitted once `sqtjs`/`zr4zj` completed:

```shell
kubectl get jobs -n team-a
```

```shell
NAME                      STATUS     COMPLETIONS   DURATION   AGE
sample-job-team-a-mbd4h   Complete   3/3           44s        2m31s
sample-job-team-a-nptv9   Complete   3/3           40s        2m10s
sample-job-team-a-sqtjs   Complete   3/3           42s        2m20s
```

```shell
kubectl get jobs -n team-b
```

```shell
NAME                      STATUS     COMPLETIONS   DURATION   AGE
sample-job-team-b-fp5jw   Complete   3/3           45s        2m30s
sample-job-team-b-pc9bl   Complete   3/3           42s        2m10s
sample-job-team-b-zr4zj   Complete   3/3           42s        2m20s
```

All 6 jobs (18 pods total, 3 completions each) ran to completion, and the ClusterQueue quota returned to fully unused once the queue drained:

```shell
kubectl get clusterqueue
```

```shell
NAME            COHORT   PENDING WORKLOADS
cluster-queue            0
```

```shell
kubectl describe clusterqueue cluster-queue
```

```shell
...
  Flavors Usage:
    Name:  default-flavor
    Resources:
      Borrowed:         0
      Name:             cpu
      Total:            0
      Borrowed:         0
      Name:             ephemeral-storage
      Total:            0
      Borrowed:         0
      Name:             memory
      Total:            0
      Borrowed:         0
      Name:             nvidia.com/gpu
      Total:            0
  Pending Workloads:    0
  Reserving Workloads:  0
```

This demonstrates that VKS v1.36 supports installing and successfully operating Kueue as a gang-scheduling / quota-based-queueing solution: workloads from two independent tenant namespaces (`team-a`, `team-b`) shared a single `ClusterQueue` with a fixed `nvidia.com/gpu` quota, were correctly suspended when quota was exhausted, and were admitted in FIFO order as GPU capacity was released by completing jobs — all backed by real GPU scheduling on the underlying vGPU-enabled node pool.
