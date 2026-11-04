## Description

The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

## Evidence

### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool (2 replicas), VM Class with vGPU profile and NVIDIA GPU Operator.
* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Kueue

```shell
## Create namespace for kueue install
kubectl create ns kueue-system

## Install via helm
helm install kueue oci://registry.k8s.io/kueue/charts/kueue --version=0.20.4 --namespace kueue-system --wait --timeout 300s
```

```shell
Pulled: registry.k8s.io/kueue/charts/kueue:0.20.4
Digest: sha256:f583edffec554671dba5c0d4df0676b35984b7c678b0f722298a53db562bfd8b
NAME: kueue
LAST DEPLOYED: Tue Nov  3 23:16:08 2026
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
kueue-controller-manager   1/1     1            1           16s
```

```shell
kubectl get pods -n kueue-system
```

```shell
NAME                                        READY   STATUS    RESTARTS   AGE
kueue-controller-manager-546c9f459f-vhdqk   1/1     Running   0          16s
```

Kueue v0.20 serves the `kueue.x-k8s.io/v1beta2` API, which is used for all objects below:

```shell
kubectl api-resources --api-group=kueue.x-k8s.io
```

```shell
NAME                         SHORTNAMES                         APIVERSION               NAMESPACED   KIND
admissionchecks              ac                                 kueue.x-k8s.io/v1beta2   false        AdmissionCheck
clusterqueues                cq                                 kueue.x-k8s.io/v1beta2   false        ClusterQueue
cohorts                      co                                 kueue.x-k8s.io/v1beta2   false        Cohort
localqueues                  queue,queues,lq                    kueue.x-k8s.io/v1beta2   true         LocalQueue
multikueueclusters           mkc                                kueue.x-k8s.io/v1beta2   false        MultiKueueCluster
multikueueconfigs            mkconf                             kueue.x-k8s.io/v1beta2   false        MultiKueueConfig
provisioningrequestconfigs   prc                                kueue.x-k8s.io/v1beta2   false        ProvisioningRequestConfig
resourceflavors              flavor,flavors,rf                  kueue.x-k8s.io/v1beta2   false        ResourceFlavor
topologies                   topo                               kueue.x-k8s.io/v1beta2   false        Topology
workloadpriorityclasses      wpc                                kueue.x-k8s.io/v1beta2   false        WorkloadPriorityClass
workloads                    kwl,kueueworkload,kueueworkloads   kueue.x-k8s.io/v1beta2   true         Workload
```

Create two new namespaces

```shell
kubectl create namespace team-a

kubectl create namespace team-b
```

### Create the ResourceFlavor, ClusterQueue and LocalQueues

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kueue.x-k8s.io/v1beta2
kind: ResourceFlavor
metadata:
  name: default-flavor
---
apiVersion: kueue.x-k8s.io/v1beta2
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
---
apiVersion: kueue.x-k8s.io/v1beta2
kind: LocalQueue
metadata:
  namespace: team-a
  name: lq-team-a
spec:
  clusterQueue: cluster-queue
---
apiVersion: kueue.x-k8s.io/v1beta2
kind: LocalQueue
metadata:
  namespace: team-b
  name: lq-team-b
spec:
  clusterQueue: cluster-queue
EOF
```

```shell
resourceflavor.kueue.x-k8s.io/default-flavor created
clusterqueue.kueue.x-k8s.io/cluster-queue created
localqueue.kueue.x-k8s.io/lq-team-a created
localqueue.kueue.x-k8s.io/lq-team-b created
```

Verify ResourceFlavor

```shell
kubectl get resourceflavor
```

```shell
NAME             AGE
default-flavor   5s
```

Inspect ClusterQueue resource

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
Name:         cluster-queue
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  kueue.x-k8s.io/v1beta2
Kind:         ClusterQueue
Spec:
  Flavor Fungibility:
    When Can Borrow:   MayStopSearch
    When Can Preempt:  TryNextFlavor
  Namespace Selector:
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
    Last Transition Time:  2026-11-03T23:16:35Z
    Message:               Can admit new workloads
    Observed Generation:   1
    Reason:                Ready
    Status:                True
    Type:                  Active
  Pending Workloads:    0
  Reserving Workloads:  0
Events:                 <none>
```

Inspect LocalQueue resource

```shell
kubectl get localqueue -A
```

```shell
NAMESPACE   NAME        CLUSTERQUEUE    PENDING WORKLOADS   ADMITTED WORKLOADS
team-a      lq-team-a   cluster-queue   0                   0
team-b      lq-team-b   cluster-queue   0                   0
```

```shell
kubectl describe localqueue lq-team-a -n team-a
```

```shell
Name:         lq-team-a
Namespace:    team-a
API Version:  kueue.x-k8s.io/v1beta2
Kind:         LocalQueue
Spec:
  Cluster Queue:  cluster-queue
  Stop Policy:    None
Status:
  Admitted Workloads:  0
  Conditions:
    Last Transition Time:  2026-11-03T23:16:35Z
    Message:               Can submit new workloads to localQueue
    Observed Generation:   1
    Reason:                Ready
    Status:                True
    Type:                  Active
  Pending Workloads:    0
  Reserving Workloads:  0
Events:                 <none>
```

(Flavor reservation/usage counters, all `0` at this point, are omitted for brevity.)

### Create jobs

The team-a Job manifest is shown below; the team-b manifest is identical except for `namespace: team-b` and the queue label `lq-team-b`.

```shell
cat <<EOF > job-team-a.yaml
apiVersion: batch/v1
kind: Job
metadata:
  namespace: team-a
  generateName: sample-job-team-a-
  labels:
    kueue.x-k8s.io/queue-name: lq-team-a
spec:
  ttlSecondsAfterFinished: 600
  parallelism: 1
  completions: 3
  suspend: true
  template:
    spec:
      nodeSelector:
        nvidia.com/gpu.product: "GRID-A100-4C"
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
for i in 1 2 3; do kubectl create -f job-team-a.yaml; kubectl create -f job-team-b.yaml; sleep 10; done
```

```shell
job.batch/sample-job-team-a-9r7x8 created
job.batch/sample-job-team-b-b4qnk created
job.batch/sample-job-team-a-xzdj9 created
job.batch/sample-job-team-b-bxpw9 created
job.batch/sample-job-team-a-rvbsw created
job.batch/sample-job-team-b-6gvqt created
```

### Monitor the jobs being scheduled and executed

Immediately after creation, only one job per team was admitted (2 GPUs = full ClusterQueue quota); the rest were correctly held `Suspended`:

```shell
kubectl get jobs -n team-a
```

```shell
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-a-9r7x8   Running     1/3           29s        29s
sample-job-team-a-rvbsw   Suspended   0/3                      9s
sample-job-team-a-xzdj9   Suspended   0/3                      19s
```

```shell
kubectl get jobs -n team-b
```

```shell
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-b-6gvqt   Suspended   0/3                      8s
sample-job-team-b-b4qnk   Running     1/3           29s        29s
sample-job-team-b-bxpw9   Suspended   0/3                      19s
```

```shell
kubectl get workloads -A
```

```shell
NAMESPACE   NAME                                QUEUE       RESERVED IN     ADMITTED   FINISHED   AGE
team-a      job-sample-job-team-a-9r7x8-520bd   lq-team-a   cluster-queue   True                  29s
team-a      job-sample-job-team-a-rvbsw-f10cb   lq-team-a                                         9s
team-a      job-sample-job-team-a-xzdj9-15dfb   lq-team-a                   False                 19s
team-b      job-sample-job-team-b-6gvqt-6fabd   lq-team-b                                         8s
team-b      job-sample-job-team-b-b4qnk-262ea   lq-team-b   cluster-queue   True                  29s
team-b      job-sample-job-team-b-bxpw9-5b27c   lq-team-b                                         19s
```

```shell
kubectl get clusterqueue
```

```shell
NAME            COHORT   PENDING WORKLOADS
cluster-queue            4
```

As each running job freed its `nvidia.com/gpu` quota unit on completion, Kueue admitted the next queued workload in FIFO order. The team-a events confirm quota-based admission (the queued workloads waited 34s and 65s for quota):

```shell
kubectl get events -n team-a --sort-by='.lastTimestamp'
```

```shell
TYPE      REASON             OBJECT                              MESSAGE
Normal    QuotaReserved      job-sample-job-team-a-9r7x8-520bd   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 1s
Normal    Admitted           job-sample-job-team-a-9r7x8-520bd   Admitted by ClusterQueue cluster-queue, wait time since reservation was 1s
Normal    Started            sample-job-team-a-9r7x8             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-9r7x8             Job resumed
Normal    Completed          sample-job-team-a-9r7x8             Job completed
Normal    FinishedWorkload   sample-job-team-a-9r7x8             Workload 'team-a/job-sample-job-team-a-9r7x8-520bd' is declared finished
Normal    QuotaReserved      job-sample-job-team-a-xzdj9-15dfb   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 34s
Normal    Admitted           job-sample-job-team-a-xzdj9-15dfb   Admitted by ClusterQueue cluster-queue, wait time since reservation was 1s
Normal    Started            sample-job-team-a-xzdj9             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-xzdj9             Job resumed
Normal    Completed          sample-job-team-a-xzdj9             Job completed
Normal    FinishedWorkload   sample-job-team-a-xzdj9             Workload 'team-a/job-sample-job-team-a-xzdj9-15dfb' is declared finished
Normal    QuotaReserved      job-sample-job-team-a-rvbsw-f10cb   Quota reserved in ClusterQueue cluster-queue, wait time since queued was 65s
Normal    Admitted           job-sample-job-team-a-rvbsw-f10cb   Admitted by ClusterQueue cluster-queue, wait time since reservation was 0s
Normal    Started            sample-job-team-a-rvbsw             Admitted by clusterQueue cluster-queue
Normal    Resumed            sample-job-team-a-rvbsw             Job resumed
Normal    Completed          sample-job-team-a-rvbsw             Job completed
Normal    FinishedWorkload   sample-job-team-a-rvbsw             Workload 'team-a/job-sample-job-team-a-rvbsw-f10cb' is declared finished
```

(Events are grouped per workload here for readability; many events share the same one-second timestamp, so events with the same timestamp may interleave in the raw `--sort-by` ordering. Per-pod `Container started` events are omitted.)

All 6 jobs (18 pods total, 3 completions each) ran to completion:

```shell
kubectl get jobs -n team-a
```

```shell
NAME                      STATUS     COMPLETIONS   DURATION   AGE
sample-job-team-a-9r7x8   Complete   3/3           42s        2m7s
sample-job-team-a-rvbsw   Complete   3/3           40s        107s
sample-job-team-a-xzdj9   Complete   3/3           40s        117s
```

```shell
kubectl get jobs -n team-b
```

```shell
NAME                      STATUS     COMPLETIONS   DURATION   AGE
sample-job-team-b-6gvqt   Complete   3/3           40s        106s
sample-job-team-b-b4qnk   Complete   3/3           43s        2m7s
sample-job-team-b-bxpw9   Complete   3/3           40s        117s
```

```shell
kubectl get workloads -A
```

```shell
NAMESPACE   NAME                                QUEUE       RESERVED IN     ADMITTED   FINISHED   AGE
team-a      job-sample-job-team-a-9r7x8-520bd   lq-team-a   cluster-queue   True       True       2m7s
team-a      job-sample-job-team-a-rvbsw-f10cb   lq-team-a   cluster-queue   True       True       107s
team-a      job-sample-job-team-a-xzdj9-15dfb   lq-team-a   cluster-queue   True       True       117s
team-b      job-sample-job-team-b-6gvqt-6fabd   lq-team-b   cluster-queue   True       True       106s
team-b      job-sample-job-team-b-b4qnk-262ea   lq-team-b   cluster-queue   True       True       2m7s
team-b      job-sample-job-team-b-bxpw9-5b27c   lq-team-b   cluster-queue   True       True       117s
```

The ClusterQueue quota returned to fully unused once the queue drained:

```shell
kubectl get clusterqueue
```

```shell
NAME            COHORT   PENDING WORKLOADS
cluster-queue            0
```

```shell
kubectl describe clusterqueue cluster-queue | sed -n '/Flavors Usage/,/Reserving Workloads/p'
```

```shell
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

### All-or-nothing (gang) admission

To demonstrate all-or-nothing semantics for a multi-pod workload, two additional Jobs were created from the same template with `parallelism` = `completions` = N (so all N pods must run concurrently, each requesting 1 GPU) and a 60s run time:

- `gang-job-3pods`: N=3 → needs 3 GPUs, which exceeds the ClusterQueue's `nvidia.com/gpu` quota of 2
- `gang-job-2pods`: N=2 → needs exactly the 2 available GPUs

```shell
kubectl create -f gang-job-3pods.yaml
```

```shell
job.batch/gang-job-3pods created
```

After 30 seconds the 3-pod gang has not been partially started — it stays `Suspended` and **no pods** were created, rather than running 2 of its 3 pods:

```shell
kubectl get job gang-job-3pods -n team-a
```

```shell
NAME             STATUS      COMPLETIONS   DURATION   AGE
gang-job-3pods   Suspended   0/3                      30s
```

```shell
kubectl get pods -n team-a -l job-name=gang-job-3pods
```

```shell
No resources found in team-a namespace.
```

```shell
kubectl get workloads -n team-a -l kueue.x-k8s.io/job-uid=<gang-job-3pods uid> -o jsonpath='{range .items[*]}{.metadata.name}{": "}{range .status.conditions[*]}{.type}={.status} ({.message}){"; "}{end}{"\n"}{end}'
```

```shell
job-gang-job-3pods-086ab: QuotaReserved=False (couldn't assign flavors to pod set main: insufficient quota for nvidia.com/gpu in flavor default-flavor, previously considered podsets requests (0) + current podset request (3) > maximum capacity (2)); Admitted=False (The workload has no reservation); PodsReady=False (Not all pods are ready or succeeded); 
```

The 2-pod gang, which fits the quota, is admitted as a unit and both pods start simultaneously on the two GPU nodes, while the 3-pod gang continues to hold zero pods:

```shell
kubectl create -f gang-job-2pods.yaml
```

```shell
job.batch/gang-job-2pods created
```

```shell
kubectl get job gang-job-2pods -n team-a
```

```shell
NAME             STATUS    COMPLETIONS   DURATION   AGE
gang-job-2pods   Running   0/2           25s        25s
```

```shell
kubectl get pods -n team-a -l job-name=gang-job-2pods -o wide
```

```shell
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
gang-job-2pods-2s89m   1/1     Running   0          25s   193.0.2.61   conformance-test-np1-sms5g-pt7t9-p42p6   <none>           <none>
gang-job-2pods-4kshn   1/1     Running   0          25s   193.0.1.70   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl get job gang-job-3pods -n team-a; kubectl get pods -n team-a -l job-name=gang-job-3pods
```

```shell
NAME             STATUS      COMPLETIONS   DURATION   AGE
gang-job-3pods   Suspended   0/3                      55s
No resources found in team-a namespace.
```

```shell
kubectl get workloads -n team-a
```

```shell
NAME                                QUEUE       RESERVED IN     ADMITTED   FINISHED   AGE
job-gang-job-2pods-c0e82            lq-team-a   cluster-queue   True                  26s
job-gang-job-3pods-086ab            lq-team-a                   False                 56s
job-sample-job-team-a-9r7x8-520bd   lq-team-a   cluster-queue   True       True       3m16s
job-sample-job-team-a-rvbsw-f10cb   lq-team-a   cluster-queue   True       True       2m56s
job-sample-job-team-a-xzdj9-15dfb   lq-team-a   cluster-queue   True       True       3m6s
```

This demonstrates that VKS v1.37 supports installing and successfully operating Kueue as a gang-scheduling / quota-based-queueing solution: workloads from two independent tenant namespaces (`team-a`, `team-b`) shared a single `ClusterQueue` with a fixed `nvidia.com/gpu` quota, were correctly suspended when quota was exhausted, and were admitted in FIFO order as GPU capacity was released; and multi-pod workloads are admitted all-or-nothing — a gang that cannot be fully satisfied creates no pods at all, while a gang that fits starts all of its pods together — all backed by real GPU scheduling on the underlying vGPU-enabled node pool.

### Cleanup

```shell
kubectl delete jobs --all -n team-a; kubectl delete jobs --all -n team-b
kubectl delete localqueue lq-team-a -n team-a; kubectl delete localqueue lq-team-b -n team-b
kubectl delete clusterqueue cluster-queue; kubectl delete resourceflavor default-flavor
kubectl delete namespace team-a team-b
helm uninstall kueue -n kueue-system
kubectl delete namespace kueue-system
```
