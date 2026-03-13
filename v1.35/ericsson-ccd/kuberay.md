# KubeRay Installation and Validation Guide

This guide covers uploading Ray images to Kubernetes registry, installing KubeRay operator and cluster, and validating with two methods.
(Method 1: Using Existing RayCluster, Method 2: Using RayJob with Ephemeral Cluster)

Note: All the commands in the following steps are to be executed on a control plane node of the Ericsson CCD deployment.

## Overview
KubeRay is a Kubernetes operator that simplifies deploying, managing, and scaling Ray applications (AI/ML workloads) on Kubernetes. It automates the lifecycle of Ray clusters, including auto-scaling, fault tolerance, and support for heterogeneous hardware (GPUs). It supports features like RayService for high-availability model serving.


## Step 1: Install KubeRay Operator

```bash
# Add Helm repository
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update

# Create namespace
kubectl create namespace kuberay-system

# Install operator with control plane toleration
helm install kuberay-operator kuberay/kuberay-operator \
  --version 1.5.1 \
  -n kuberay-system \
  --set tolerations[0].key=node-role.kubernetes.io/control-plane \
  --set tolerations[0].operator=Exists \
  --set tolerations[0].effect=NoSchedule

# Verify operator installation
kubectl get pods -n kuberay-system
```

## Step 2: Install RayCluster

```bash
# Install RayCluster with control plane toleration
helm install raycluster-test kuberay/ray-cluster \
  --version 1.5.1 \
  -n default \
  --set head.tolerations[0].key=node-role.kubernetes.io/control-plane \
  --set head.tolerations[0].operator=Exists \
  --set head.tolerations[0].effect=NoSchedule \
  --set worker.tolerations[0].key=node-role.kubernetes.io/control-plane \
  --set worker.tolerations[0].operator=Exists \
  --set worker.tolerations[0].effect=NoSchedule \
  --set worker.replicas=2

# Wait for cluster to be ready
kubectl get raycluster -n default
kubectl get pods -n default -l ray.io/cluster=raycluster-test-kuberay
```

## Validation Method 1: Using Existing RayCluster

This method validates the existing persistent RayCluster by running distributed workload directly in the cluster.

```bash
# Get head pod name
HEAD_POD=$(kubectl get pods -n default \
  -l ray.io/cluster=raycluster-test-kuberay,ray.io/node-type=head \
  -o jsonpath='{.items[0].metadata.name}')

# Run distributed validation script
kubectl exec ${HEAD_POD} -n default -- python -c "
import ray
import time

ray.init(address='auto')

@ray.remote
def distributed_task(x):
    time.sleep(1)
    return x * x

# Run distributed tasks
futures = [distributed_task.remote(i) for i in range(10)]
results = ray.get(futures)

print(f'Distributed task results: {results}')
print(f'Cluster resources: {ray.cluster_resources()}')
print('VALIDATION SUCCESS: Existing cluster validated')
"
```

**Expected Output:**
```
Distributed task results: [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
Cluster resources: {'CPU': 2.0, 'memory': 536870912, 'node:10.244.0.5': 1.0, ...}
VALIDATION SUCCESS: Existing cluster validated
```

## Validation Method 2: Using RayJob with Ephemeral Cluster

This method creates a RayJob that provisions its own ephemeral cluster, runs validation, and cleans up automatically.

### Step 1: Create Validation Script ConfigMap

```bash
kubectl create namespace ray-job

kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: ray-validate-script
  namespace: ray-job
data:
  validate.py: |
    print("=== VALIDATE.PY STARTING ===")
    import ray
    import os
    import socket
    import time

    # Initialize Ray and wait for workers
    ray.init(address="auto")

    # Wait for workers to be available
    print("Waiting for workers to be available...")
    for i in range(30):  # Wait up to 30 seconds
        resources = ray.cluster_resources()
        total_cpus = resources.get("CPU", 0)
        print(f"Available CPUs: {total_cpus}")
        if total_cpus >= 3:  # Head + 2 workers = 3 CPUs minimum
            break
        time.sleep(1)
    else:
        print("Warning: Not all workers may be available")

    # Force distribution by using more tasks than head node CPUs
    @ray.remote(num_cpus=1)
    def task(x):
        import socket
        import os
        import time
        time.sleep(1)  # Small delay to ensure distribution
        return {
            "result": x * 2,
            "node": socket.gethostname(),
            "pid": os.getpid()
        }

    # Submit more tasks to force distribution
    print("Submitting distributed tasks...")
    # Use fewer tasks to ensure they can run concurrently
    futures = [task.remote(i) for i in range(20)]  # 20 tasks for 3 CPUs
    results = ray.get(futures)

    nodes = set(r["node"] for r in results)
    workers = set((r["node"], r["pid"]) for r in results)
    print(f"Results: {[r['result'] for r in results]}")
    print(f"Nodes used: {len(nodes)} -> {list(nodes)}")
    print(f"Workers used: {len(workers)}")
    print(f"Cluster resources: {ray.cluster_resources()}")
    print(f"Available resources: {ray.available_resources()}")

    if len(nodes) >= 2 and len(workers) >= 2:
        print("VALIDATION SUCCESS: Distributed across multiple nodes and workers")
    else:
        print(f"VALIDATION FAILED: Only {len(nodes)} nodes and {len(workers)} workers used")
    print("=== VALIDATE.PY COMPLETED ===")
EOF
```

### Step 2: Create and Run RayJob

```bash
kubectl apply -f - <<EOF
apiVersion: ray.io/v1
kind: RayJob
metadata:
  name: ray-distributed-validate
  namespace: ray-job
spec:
  entrypoint: python /tmp/scripts/validate.py
  rayClusterSpec:
    rayVersion: "2.46.0"
    headGroupSpec:
      rayStartParams:
        num-cpus: "1"
      template:
        spec:
          tolerations:
            - key: node-role.kubernetes.io/control-plane
              operator: Exists
              effect: NoSchedule
          containers:
            - name: ray-head
              image: rayproject/ray:2.46.0
              volumeMounts:
                - name: script
                  mountPath: /tmp/scripts
          volumes:
            - name: script
              configMap:
                name: ray-validate-script
    workerGroupSpecs:
      - groupName: workers
        replicas: 2
        minReplicas: 2
        maxReplicas: 2
        rayStartParams:
          num-cpus: "1"
        template:
          spec:
            affinity:
              podAntiAffinity:
                preferredDuringSchedulingIgnoredDuringExecution:
                - weight: 100
                  podAffinityTerm:
                    labelSelector:
                      matchExpressions:
                      - key: ray.io/node-type
                        operator: In
                        values: ["worker"]
                    topologyKey: kubernetes.io/hostname
            tolerations:
              - key: node-role.kubernetes.io/control-plane
                operator: Exists
                effect: NoSchedule
            containers:
              - name: ray-worker
                image: rayproject/ray:2.46.0
  shutdownAfterJobFinishes: true
  ttlSecondsAfterFinished: 30
EOF
```

### Step 3: Monitor RayJob Completion

```bash
# Watch job status
kubectl get rayjob ray-distributed-validate -n ray-job -w

# Check job logs
kubectl logs -n ray-job job/ray-distributed-validate --tail=100
```

**Expected Output:**
```
=== VALIDATE.PY STARTING ===
Waiting for workers to be available...
Available CPUs: 3
Submitting distributed tasks...
Results: [0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38]
Nodes used: 2 -> ['ray-head-xxx', 'ray-worker-xxx']
Workers used: 3
Cluster resources: {'CPU': 3.0, 'memory': 1073741824, ...}
VALIDATION SUCCESS: Distributed across multiple nodes and workers
=== VALIDATE.PY COMPLETED ===
```

## Cleanup

```bash
# Method 1: Delete existing cluster
helm uninstall raycluster-test -n default

# Method 2: RayJob auto-cleans, but delete namespace if needed
kubectl delete namespace ray-job

# Uninstall operator
helm uninstall kuberay-operator -n kuberay-system
kubectl delete namespace kuberay-system
```
