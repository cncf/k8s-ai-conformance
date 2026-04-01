# Gang Scheduling Validation (Kueue)

## Overview

Gang scheduling ensures that a group of related pods are scheduled together or not at all. This prevents resource fragmentation and deadlocks in distributed AI workloads. Kueue is the Kubernetes SIGs job queueing system that provides this capability through its workload admission and quota management pipeline.

This guide covers installing Kueue and validating gang scheduling (all-or-nothing) behavior on Viettel AI Platform.

## Step 1: Install Kueue

```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/download/v0.10.2/manifests.yaml
```

**Output (last 5 lines):**
```
service/kueue-webhook-service serverside-applied
deployment.apps/kueue-controller-manager serverside-applied
apiservice.apiregistration.k8s.io/v1beta1.visibility.kueue.x-k8s.io serverside-applied
mutatingwebhookconfiguration.admissionregistration.k8s.io/kueue-mutating-webhook-configuration serverside-applied
validatingwebhookconfiguration.admissionregistration.k8s.io/kueue-validating-webhook-configuration serverside-applied
```

## Step 2: Verify Kueue Pods and APIs

```bash
$ kubectl get pods -n kueue-system
NAME                                        READY   STATUS    RESTARTS   AGE
kueue-controller-manager-5b9d6dc7df-pm4mv   2/2     Running   0          37s

$ kubectl api-resources --api-group=kueue.x-k8s.io
NAME                         SHORTNAMES          APIVERSION                NAMESPACED   KIND
admissionchecks                                  kueue.x-k8s.io/v1beta1    false        AdmissionCheck
clusterqueues                cq                  kueue.x-k8s.io/v1beta1    false        ClusterQueue
cohorts                                          kueue.x-k8s.io/v1alpha1   false        Cohort
localqueues                  queue,queues,lq     kueue.x-k8s.io/v1beta1    true         LocalQueue
multikueueclusters                               kueue.x-k8s.io/v1beta1    false        MultiKueueCluster
multikueueconfigs                                kueue.x-k8s.io/v1beta1    false        MultiKueueConfig
provisioningrequestconfigs                       kueue.x-k8s.io/v1beta1    false        ProvisioningRequestConfig
resourceflavors              flavor,flavors,rf   kueue.x-k8s.io/v1beta1    false        ResourceFlavor
topologies                                       kueue.x-k8s.io/v1alpha1   false        Topology
workloadpriorityclasses                          kueue.x-k8s.io/v1beta1    false        WorkloadPriorityClass
workloads                    wl                  kueue.x-k8s.io/v1beta1    true         Workload
```

## Step 3: Create Namespace and Queue Infrastructure

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  labels:
    kubernetes.io/metadata.name: gang-scheduling-test
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/enforce-version: latest
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
EOF
```

Verify the ClusterQueue is ready:

```bash
$ kubectl get clusterqueue gang-test-cq -o yaml | grep -A5 'conditions:'
conditions:
  - lastTransitionTime: "2026-04-01T09:28:57Z"
    message: Can admit new workloads
    observedGeneration: 1
    reason: Ready
    status: "True"
```

## Validation Method 1: Gang Scheduling Job (All-or-Nothing)

```bash
kubectl apply -f - <<'EOF'
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-test-job
  namespace: gang-scheduling-test
  labels:
    kueue.x-k8s.io/queue-name: gang-test-lq
spec:
  completions: 3
  parallelism: 3
  template:
    spec:
      containers:
      - name: worker
        image: ubuntu:22.04
        command: ["sh", "-c"]
        args:
        - |
          echo '=== GANG SCHEDULING VALIDATION ==='
          echo "Pod: $HOSTNAME"
          echo "Start time: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
          sleep 5
          echo 'GANG MEMBER COMPLETED'
        resources:
          requests:
            cpu: "100m"
            memory: "64Mi"
          limits:
            cpu: "200m"
            memory: "128Mi"
      restartPolicy: Never
  backoffLimit: 0
EOF
```

Verify all 3 pods were admitted and completed together:

```bash
$ kubectl get workloads -n gang-scheduling-test
NAME                      QUEUE          RESERVED IN    ADMITTED   FINISHED   AGE
job-gang-test-job-b942a   gang-test-lq   gang-test-cq   True       True       83s

$ kubectl get pods -n gang-scheduling-test
NAME                  READY   STATUS      RESTARTS   AGE
gang-test-job-9hnn6   0/1     Completed   0          78s
gang-test-job-kswl4   0/1     Completed   0          78s
gang-test-job-nw877   0/1     Completed   0          78s

$ kubectl logs -n gang-scheduling-test -l job-name=gang-test-job --prefix
[pod/gang-test-job-kswl4/worker] === GANG SCHEDULING VALIDATION ===
[pod/gang-test-job-kswl4/worker] Pod: gang-test-job-kswl4
[pod/gang-test-job-kswl4/worker] Start time: 2026-04-01T09:29:24Z
[pod/gang-test-job-kswl4/worker] GANG MEMBER COMPLETED
[pod/gang-test-job-nw877/worker] === GANG SCHEDULING VALIDATION ===
[pod/gang-test-job-nw877/worker] Pod: gang-test-job-nw877
[pod/gang-test-job-nw877/worker] Start time: 2026-04-01T09:29:23Z
[pod/gang-test-job-nw877/worker] GANG MEMBER COMPLETED
[pod/gang-test-job-9hnn6/worker] === GANG SCHEDULING VALIDATION ===
[pod/gang-test-job-9hnn6/worker] Pod: gang-test-job-9hnn6
[pod/gang-test-job-9hnn6/worker] Start time: 2026-04-01T09:29:25Z
[pod/gang-test-job-9hnn6/worker] GANG MEMBER COMPLETED
```

All 3 gang members started within 2 seconds of each other and completed successfully. The Workload shows `ADMITTED: True` and `FINISHED: True`.

## Cleanup

```bash
kubectl delete job --all -n gang-scheduling-test
kubectl delete localqueue gang-test-lq -n gang-scheduling-test
kubectl delete clusterqueue gang-test-cq
kubectl delete resourceflavor default-flavor
kubectl delete namespace gang-scheduling-test
```
