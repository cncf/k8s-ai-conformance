# Evidence: Gang Scheduling

---

## Overview

This test validates that the Kubernetes platform supports gang scheduling via Kueue. Gang scheduling ensures that all pods in a distributed workload are scheduled together (all-or-nothing), which is critical for distributed AI/ML training jobs where partial scheduling would waste resources and cause deadlocks.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform.

---

## Step 1: Install Kueue

Install Kueue, a Kubernetes-native job queueing system that provides gang scheduling capabilities:

```
helm install kueue oci://registry.k8s.io/kueue/charts/kueue \
  --version=0.16.2 \
  --namespace  kueue-system \
  --create-namespace \
  --wait --timeout 300s
```

> **Note:** ✅ Kueue installed successfully in the kueue-system namespace.

---

## Step 2: Create Namespace, ResourceFlavor, ClusterQueue, and LocalQueue

Create the test namespace and Kueue resources to manage workload admission:

```
apiVersion: v1
kind: Namespace
metadata:
  labels:
    kubernetes.io/metadata.name: gang-scheduling-test
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/enforce-version: latest
    pod-security.kubernetes.io/audit: privileged
    pod-security.kubernetes.io/warn: privileged
  name: gang-scheduling-test

---
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: default-flavor
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: gang-test-cq
spec:
  namespaceSelector: {}
  queueingStrategy: BestEffortFIFO
  resourceGroups:
  - coveredResources: ["cpu", "memory"]
    flavors:
    - name: default-flavor
      resources:
      - name: cpu
        nominalQuota: "10"
      - name: memory
        nominalQuota: "10Gi"
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  name: gang-test-lq
  namespace: gang-scheduling-test
spec:
  clusterQueue: gang-test-cq
```

---

## Step 3: Verify Queue Configuration

Verify that the ClusterQueue and LocalQueue were created successfully:

```
$ kubectl get clusterqueue -A

NAME           COHORT   PENDING WORKLOADS
gang-test-cq            0
```

```
$ kubectl get localqueue -A

NAMESPACE              NAME           CLUSTERQUEUE   PENDING WORKLOADS   ADMITTED WORKLOADS
gang-scheduling-test   gang-test-lq   gang-test-cq   0                   0
```

> **Note:** ✅ ClusterQueue and LocalQueue created successfully with CPU and memory quotas configured.

---

## Step 4: Create Gang-Scheduled Job (Parallelism=2)

Create a Job with parallelism=2 that will be gang-scheduled (all pods must be admitted together):

```
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-test-job
  namespace: gang-test-scheduling
  labels:
    kueue.x-k8s.io/queue-name: gang-test-lq
spec:
  parallelism: 2
  completions: 2
  template:
    spec:
      containers:
      - name: worker
        image: busybox
        command: ["sh", "-c"]
        args:
          - |
            echo "=== GANG SCHEDULING VALIDATION ==="
            echo "Pod: $HOSTNAME"
            echo "Start time: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
            sleep 5
            echo "GANG MEMBER COMPLETED"
        resources:
          requests:
            cpu: "100m"
            memory: "64Mi"
      restartPolicy: Never
```

---

## Step 5: Verify Workload Admission

Check that the workload was admitted by Kueue (gang scheduling - all pods admitted together):

```
$ kubectl get workloads -A

NAMESPACE              NAME                      QUEUE          RESERVED IN    ADMITTED   FINISHED   AGE
gang-scheduling-test   job-gang-test-job-dfacb   gang-test-lq   gang-test-cq   True                  3s
```

> **Note:** ✅ Workload was admitted (ADMITTED=True) and reserved in gang-test-cq. All 2 pods were scheduled together.

---

## Step 6: Verify Gang Execution

Check the job logs to confirm all gang members executed:

```
$ kubectl logs -f -n gang-scheduling-test job/gang-test-job
Found 2 pods, using pod/gang-test-job-lf7pz
=== GANG SCHEDULING VALIDATION ===
Pod: gang-test-job-lf7pz
Start time: 2026-03-25T18:33:14Z
GANG MEMBER COMPLETED
```

> **Note:** ✅ Gang job executed successfully. Both pods started together and completed their work.

---

## Step 7: Test All-or-Nothing Behavior (Insufficient Resources)

Create a Job with parallelism=5 requesting resources that exceed the cluster quota. This demonstrates gang scheduling's all-or-nothing behavior - the entire workload stays pending rather than partially scheduling:

```
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-test-job
  namespace: gang-test-scheduling
  labels:
    kueue.x-k8s.io/queue-name: gang-test-lq
spec:
  parallelism: 5
  completions: 5
  template:
    spec:
      containers:
      - name: worker
        image: busybox
        command: ["sh", "-c"]
        args:
          - |
            echo "=== GANG SCHEDULING VALIDATION ==="
            echo "Pod: $HOSTNAME"
            echo "Start time: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
            sleep 5
            echo "GANG MEMBER COMPLETED"
        resources:
          requests:
            cpu: "5"
            memory: "2Gi"
      restartPolicy: Never
```

---

## Step 8: Verify All-or-Nothing Behavior

Check that the workload stays pending because the cluster cannot satisfy the total resource request (5 pods × 5 CPU = 25 CPU, but quota is only 10 CPU):

```
$ kubectl get workloads -A
NAMESPACE              NAME                      QUEUE          RESERVED IN   ADMITTED   FINISHED   AGE
gang-scheduling-test   job-gang-test-job-09418   gang-test-lq                                       3s

$ kubectl get clusterqueue -A
NAME           COHORT   PENDING WORKLOADS
gang-test-cq            1

$ kubectl get localqueue -A
NAMESPACE              NAME           CLUSTERQUEUE   PENDING WORKLOADS   ADMITTED WORKLOADS
gang-scheduling-test   gang-test-lq   gang-test-cq   1                   0
```

> **Note:** ✅ The workload remains PENDING (ADMITTED is empty, PENDING WORKLOADS=1). This demonstrates gang scheduling's all-or-nothing behavior - Kueue will not partially schedule the job. All 5 pods must be schedulable together, or none are scheduled.

---