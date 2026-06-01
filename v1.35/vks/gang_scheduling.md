## Description 

The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

## Evidence

### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Kueue

```shell
## Create namespace for kueue install
kubectl create ns  kueue-system

## Install via helm
helm install kueue oci://registry.k8s.io/kueue/charts/kueue --version=0.14.1 --namespace kueue-system --wait --timeout 300s
```


Wait until Kueue pods are ready

```shell
kubectl get deployments -n kueue-system
NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
kueue-controller-manager   1/1     1            1           91s

kubectl get pods -n kueue-system
NAME                                        READY   STATUS    RESTARTS   AGE
kueue-controller-manager-855894b6c5-7ljfk   1/1     Running   0          102s
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
NAME             AGE
default-flavor   10s

kubectl describe resourceflavor
Name:         default-flavor
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  kueue.x-k8s.io/v1beta1
Kind:         ResourceFlavor
Metadata:
  Creation Timestamp:  2026-03-17T06:49:45Z
  Finalizers:
    kueue.x-k8s.io/resource-in-use
  Generation:        1
  Resource Version:  284547
  UID:               f95336ba-31b4-4a0e-a0fb-484860b9c4e5
Spec:
Events:  <none>
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
        nominalQuota: 10
      - name: "memory"
        nominalQuota: 2Gi
      - name: "nvidia.com/gpu"
        nominalQuota: 10
      - name: "ephemeral-storage"
        nominalQuota: 10Gi
EOF
```

Inspect ClusterQueue resource

```shell
kubectl get ClusterQueue
NAME            COHORT   PENDING WORKLOADS
cluster-queue            0

kubectl describe ClusterQueue
Name:         cluster-queue
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  kueue.x-k8s.io/v1beta1
Kind:         ClusterQueue
Metadata:
  Creation Timestamp:  2026-03-17T06:50:58Z
  Finalizers:
    kueue.x-k8s.io/resource-in-use
  Generation:        1
  Resource Version:  284821
  UID:               0224e254-5529-4477-ac54-bf8b2531eb24
Spec:
  Flavor Fungibility:
    When Can Borrow:   Borrow
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
        Nominal Quota:  10
        Name:           memory
        Nominal Quota:  2Gi
        Name:           nvidia.com/gpu
        Nominal Quota:  10
        Name:           ephemeral-storage
        Nominal Quota:  10Gi
  Stop Policy:          None
Status:
  Admitted Workloads:  0
  Conditions:
    Last Transition Time:  2026-03-17T06:50:58Z
    Message:               Can admit new workloads
    Observed Generation:   1
    Reason:                Ready
    Status:                True
    Type:                  Active
  Flavors Reservation:
    Name:  default-flavor
    Resources:
      Borrowed:  0
      Name:      cpu
      Total:     0
      Borrowed:  0
      Name:      ephemeral-storage
      Total:     0
      Borrowed:  0
      Name:      memory
      Total:     0
      Borrowed:  0
      Name:      nvidia.com/gpu
      Total:     0
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
NAMESPACE   NAME        CLUSTERQUEUE    PENDING WORKLOADS   ADMITTED WORKLOADS
team-a      lq-team-a   cluster-queue   0                   0
team-b      lq-team-b   cluster-queue   0                   0

kubectl describe LocalQueue lq-team-a -n team-a
Name:         lq-team-a
Namespace:    team-a
Labels:       <none>
Annotations:  <none>
API Version:  kueue.x-k8s.io/v1beta1
Kind:         LocalQueue
Metadata:
  Creation Timestamp:  2026-03-17T06:54:13Z
  Generation:          1
  Resource Version:    285568
  UID:                 1a5225cb-569c-4af1-8339-85c460ddcae8
Spec:
  Cluster Queue:  cluster-queue
  Stop Policy:    None
Status:
  Admitted Workloads:  0
  Conditions:
    Last Transition Time:  2026-03-17T06:54:13Z
    Message:               Can submit new workloads to localQueue
    Observed Generation:   1
    Reason:                Ready
    Status:                True
    Type:                  Active
  Flavor Usage:
    Name:  default-flavor
    Resources:
      Name:   cpu
      Total:  0
      Name:   ephemeral-storage
      Total:  0
      Name:   memory
      Total:  0
      Name:   nvidia.com/gpu
      Total:  0
  Flavors:
    Name:  default-flavor
    Resources:
      cpu
      ephemeral-storage
      memory
      nvidia.com/gpu
  Flavors Reservation:
    Name:  default-flavor
    Resources:
      Name:             cpu
      Total:            0
      Name:             ephemeral-storage
      Total:            0
      Name:             memory
      Total:            0
      Name:             nvidia.com/gpu
      Total:            0
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
  ttlSecondsAfterFinished: 60
  parallelism: 3
  completions: 3
  suspend: true
  template:
    spec:
      nodeSelector:
        nvidia.com/gpu.product: "GRID-A100-8C"
      containers:
      - name: dummy-job
        image: gcr.io/k8s-staging-perf-tests/sleep:latest
        args: ["10s"]
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
            ephemeral-storage: "512Mi"
            nvidia.com/gpu: "1"
          limits:
            cpu: "500m"
            memory: "512Mi"
            ephemeral-storage: "512Mi"
            nvidia.com/gpu: "1"
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
  ttlSecondsAfterFinished: 60
  parallelism: 3
  completions: 3
  suspend: true
  template:
    spec:
      nodeSelector:
        nvidia.com/gpu.product: "GRID-A100-8C"
      containers:
      - name: dummy-job
        image: gcr.io/k8s-staging-perf-tests/sleep:latest
        args: ["10s"]
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
            ephemeral-storage: "512Mi"
            nvidia.com/gpu: "1"
          limits:
            cpu: "500m"
            memory: "512Mi"
            ephemeral-storage: "512Mi"
            nvidia.com/gpu: "1"
      restartPolicy: Never
EOF
```

Simulate jobs every 10 seconds using the following helper script:

```shell
#!/bin/bash
# Helper script to simulate job scheduling
while :
do
  kubectl create -f ${1}
  kubectl create -f ${2}
  sleep ${3:-10}
done
```

Run the script:

```shell
bash create_jobs.sh job-team-a.yaml job-team-b.yaml 10
```

Monitor the jobs being scheduled and executed

```shell
job.batch/sample-job-team-a-5b4mh created
job.batch/sample-job-team-a-6g5zn created
job.batch/sample-job-team-a-k42x2 created
job.batch/sample-job-team-a-lhg6p created
job.batch/sample-job-team-a-srv4f created
job.batch/sample-job-team-a-zhqb4 created
job.batch/sample-job-team-b-c8q4z created
job.batch/sample-job-team-b-rqlfm created
job.batch/sample-job-team-b-tc9hs created
job.batch/sample-job-team-b-wcfms created
job.batch/sample-job-team-b-z7qsg created
job.batch/sample-job-team-b-zrfpn created


kubectl get jobs -n team-a
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-a-5b4mh   Suspended   0/3                      55m
sample-job-team-a-6g5zn   Suspended   0/3                      54m
sample-job-team-a-k42x2   Suspended   0/3                      54m
sample-job-team-a-lhg6p   Complete    3/3           55m        56m
sample-job-team-a-srv4f   Suspended   0/3                      53m
sample-job-team-a-zhqb4   Suspended   0/3                      52m


kubectl get jobs -n team-b
NAME                      STATUS      COMPLETIONS   DURATION   AGE
sample-job-team-b-c8q4z   Suspended   0/3                      54m
sample-job-team-b-rqlfm   Running     1/3           29s        55m
sample-job-team-b-tc9hs   Suspended   0/3                      52m
sample-job-team-b-wcfms   Complete    3/3           52s        56m
sample-job-team-b-z7qsg   Suspended   0/3                      53m
sample-job-team-b-zrfpn   Suspended   0/3                      54m
```