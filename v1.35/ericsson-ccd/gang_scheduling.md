# Kueue Gang Scheduling Installation and Validation Guide

This guide covers installing Kueue via Helm with CCD internal registry,
configuring queue infrastructure, and validating gang scheduling (all-or-nothing)
behavior on CCD clusters.

## Overview

Gang scheduling ensures that a group of related pods are scheduled together
or not at all. This prevents resource fragmentation and deadlocks in distributed
workloads. Kueue is the Kubernetes SIGs job queueing system that provides this
capability through its workload admission and quota management pipeline.

## Prerequisites

- CCD cluster with Kubernetes 1.29+ (CCD 2.34.x is compatible)
- kubectl and helm access to the cluster


## Step 1: Install Kueue via Helm

```bash
helm install kueue oci://registry.k8s.io/kueue/charts/kueue \
  --version=0.16.2 \
  --namespace  kueue-system \
  --create-namespace \
  --wait --timeout 300s

# Verify
kubectl get pods -n kueue-system
kubectl get crd | grep kueue
```

## Step 2: Create Test Namespace and Queue Infrastructure

```bash
# Create namespace with pod security labels
cat <<EOF | kubectl apply -f -
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
EOF

# Create queue infrastructure
cat <<EOF | kubectl apply -f -
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

kubectl get clusterqueue gang-test-cq
```

## Validation Method 1: Gang Scheduling Job (All-or-Nothing)

Uses an image from the CCD internal registry (httpd is always available).

```bash
export IMAGE_TAG=$(curl -s -q https://registry.eccd.local:5000/v2/httpd/tags/list | jq '.tags[]' -r | grep -v 'sha256')
```

```bash
cat <<EOF | kubectl apply -f -
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
      tolerations:
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule
      containers:
      - name: worker
        image: registry.eccd.local:5000/httpd:${IMAGE_TAG}
        command: ["sh", "-c"]
        args:
        - |
          echo "=== GANG SCHEDULING VALIDATION ==="
          echo "Pod: \$HOSTNAME"
          echo "Start time: \$(date -u +%Y-%m-%dT%H:%M:%SZ)"
          sleep 5
          echo "GANG MEMBER COMPLETED"
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

# Verify
kubectl get workloads -n gang-scheduling-test
kubectl get pods -n gang-scheduling-test -w
kubectl logs -n gang-scheduling-test -l job-name=gang-test-job
```

**Expected Output:**
```
=== GANG SCHEDULING VALIDATION ===
Pod: gang-test-job-xxxxx
Start time: 2026-02-23T10:25:35Z
GANG MEMBER COMPLETED
```

## Validation Method 2: Gang Rejection (Insufficient Resources)

```bash
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-oversized-job
  namespace: gang-scheduling-test
  labels:
    kueue.x-k8s.io/queue-name: gang-test-lq
spec:
  completions: 5
  parallelism: 5
  suspend: true
  template:
    spec:
      tolerations:
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule
      containers:
      - name: worker
        image: registry.eccd.local:5000/httpd:${IMAGE_TAG}
        command: ["sleep", "10"]
        resources:
          requests:
            cpu: "5"
            memory: "1Gi"
      restartPolicy: Never
  backoffLimit: 0
EOF

# Verify workload stays pending (25 CPU requested > 10 CPU quota)
kubectl get workloads -n gang-scheduling-test
```

## Cleanup

```bash
kubectl delete job --all -n gang-scheduling-test
kubectl delete workloads --all -n gang-scheduling-test
kubectl delete localqueue gang-test-lq -n gang-scheduling-test
kubectl delete clusterqueue gang-test-cq
kubectl delete resourceflavor default-flavor
kubectl delete namespace gang-scheduling-test
helm uninstall kueue -n kueue-system
```
