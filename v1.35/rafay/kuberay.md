# KubeRay GPU Validation

This document demonstrates deploying the KubeRay operator via Rafay Platform and validating GPU workload scheduling using a RayJob on a Rafay MKS cluster.

## Prerequisites

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DRA (Dynamic Resource Allocation) driver support enabled is deployed through Blueprint as part of cluster provisioning.

## Configure KubeRay Operator via Rafay


### Step 1: Create Repository

Create a Helm repository so that the controller can retrieve the KubeRay Operator Helm chart automatically.

1. Select **Integrations -> Repositories**
2. Click **New Repository**
3. Enter the name `kuberay`
4. Select **Helm** for the type
5. Click **Create**
6. Enter `https://ray-project.github.io/kuberay-helm/` for the endpoint
7. Click **Save**
8. Optionally, click the **validate** button on the repo to confirm connectivity

### Step 2: Create kuberay-operator Add-On

Create a custom add-on for the kuberay-operator that will pull the Helm chart from the previously created repository.

1. Select **Infrastructure -> Add-Ons**
2. Click **New Add-On -> Create New Add-On**
3. Enter the name `kuberay-operator`
4. Select **Helm 3** for the type
5. Select **Pull files from repository**
6. Select **Helm** for the repository type
7. Select `kuberay` for the namespace
8. Click **Create**
9. Click **New Version**
10. Enter a version name
11. Select the previously created repository
12. Enter `kuberay-operator` for the chart name
13. Click **Save Changes**

### Step 3: Create Blueprint

Create a custom cluster blueprint that contains the kuberay-operator add-on. This blueprint can be applied to one or multiple clusters.

1. Select **Infrastructure -> Blueprints**
2. Click **New Blueprint**
3. Enter the name `kuberay`
4. Click **Save**
5. Enter a version name
6. Select **Minimal** for the base blueprint
7. In the add-ons section, click **Configure Add-Ons**
8. Click the **+** symbol next to the previously created add-on to add it to the blueprint
9. Click **Save Changes**

### Step 4: Apply Blueprint

Apply the previously created cluster blueprint to an existing cluster. The blueprint will deploy the kuberay-operator add-on to the cluster.

1. Select **Infrastructure -> Clusters**
2. Click the gear icon on the cluster card -> **Update Blueprint**
3. Select the previously created `kuberay` blueprint and version
4. Click **Save and Publish**

The controller will publish and reconcile the blueprint on the target cluster. This can take a few seconds to complete.

## Validation Steps

### Step 5: Verify KubeRay operator and GPU components are running

Confirm that the KubeRay operator, GPU operator, and all NVIDIA components are running on the cluster.

```
kubectl get po -A
NAMESPACE        NAME                                                              READY   STATUS      RESTARTS       AGE
ai-conformance   ai-conformance-gpu-operator-node-feature-discovery-gc-7bff8sb7l   1/1     Running     1 (131m ago)   25h
ai-conformance   ai-conformance-gpu-operator-node-feature-discovery-master-2rxkg   1/1     Running     1 (131m ago)   25h
ai-conformance   ai-conformance-gpu-operator-node-feature-discovery-worker-7fmvm   1/1     Running     1 (131m ago)   25h
ai-conformance   gpu-feature-discovery-jtn5m                                       1/1     Running     0              130m
ai-conformance   gpu-operator-7f7dfb9975-9j6ll                                     1/1     Running     1 (131m ago)   25h
ai-conformance   kuberay-operator-556ff8cf56-9jq86                                 1/1     Running     0              46s
ai-conformance   nvidia-container-toolkit-daemonset-p7dbp                          1/1     Running     0              130m
ai-conformance   nvidia-cuda-validator-2cdmv                                       0/1     Completed   0              126m
ai-conformance   nvidia-dcgm-4cdk6                                                 1/1     Running     0              130m
ai-conformance   nvidia-dcgm-exporter-jdglx                                        1/1     Running     0              130m
ai-conformance   nvidia-dra-driver-gpu-controller-6fd47d97cf-wnqqv                 1/1     Running     1 (131m ago)   25h
ai-conformance   nvidia-dra-driver-gpu-kubelet-plugin-khpgq                        2/2     Running     2 (131m ago)   24h
ai-conformance   nvidia-driver-daemonset-mwzs7                                     1/1     Running     1 (131m ago)   25h
ai-conformance   nvidia-operator-validator-hpwnb                                   1/1     Running     0              130m
```

### Step 6: Create the GPU validation script ConfigMap

Create a ConfigMap containing a Python script that will be executed as a RayJob. The script connects to the Ray cluster, dispatches a GPU task to the worker, runs a matrix multiplication on the T4 GPU, and reports the results.

```
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: ray-gpu-validate-script
  namespace: test
data:
  validate.py: |
    print("=== RAY GPU VALIDATION STARTING ===")
    import ray
    import time

    ray.init(address="auto")

    print("Waiting for GPU worker...")
    for i in range(60):
        resources = ray.cluster_resources()
        gpus = resources.get("GPU", 0)
        if gpus >= 1:
            break
        time.sleep(2)
    else:
        raise RuntimeError("No GPU worker available")

    print(f"Cluster resources: {ray.cluster_resources()}")

    @ray.remote(num_gpus=1)
    def gpu_task():
        import torch
        import socket

        gpu_available = torch.cuda.is_available()
        gpu_name = torch.cuda.get_device_name(0) if gpu_available else "N/A"
        gpu_memory = torch.cuda.get_device_properties(0).total_memory / 1e9 if gpu_available else 0

        x = torch.randn(1000, 1000, device="cuda")
        y = torch.randn(1000, 1000, device="cuda")
        z = torch.matmul(x, y)
        result_sum = z.sum().item()

        return {
            "hostname": socket.gethostname(),
            "gpu_available": gpu_available,
            "gpu_name": gpu_name,
            "gpu_memory_gb": round(gpu_memory, 2),
            "matrix_multiply_sum": round(result_sum, 4),
            "cuda_version": torch.version.cuda,
        }

    print("Running GPU task on worker...")
    result = ray.get(gpu_task.remote())

    print(f"Worker hostname: {result['hostname']}")
    print(f"GPU available: {result['gpu_available']}")
    print(f"GPU name: {result['gpu_name']}")
    print(f"GPU memory: {result['gpu_memory_gb']} GB")
    print(f"CUDA version: {result['cuda_version']}")
    print(f"Matrix multiply result sum: {result['matrix_multiply_sum']}")

    if result["gpu_available"] and "T4" in result["gpu_name"]:
        print("\nVALIDATION SUCCESS: Ray cluster can schedule and run GPU workloads")
    else:
        print("\nVALIDATION FAILED")

    print("=== RAY GPU VALIDATION COMPLETED ===")
EOF
```

### Step 7: Create and submit the RayJob

Create a RayJob resource that provisions an ephemeral Ray cluster with a CPU head node and a GPU worker node, then runs the validation script.

```
kubectl apply -f - <<EOF
apiVersion: ray.io/v1
kind: RayJob
metadata:
  name: ray-gpu-validate
  namespace: test
spec:
  entrypoint: python /tmp/scripts/validate.py
  rayClusterSpec:
    rayVersion: "2.49.2"
    headGroupSpec:
      rayStartParams:
        num-cpus: "1"
      template:
        spec:
          containers:
            - name: ray-head
              image: rayproject/ray-ml:2.49.2.7b0af3-py311-gpu
              resources:
                requests:
                  cpu: 500m
                  memory: 1Gi
                limits:
                  cpu: 1
                  memory: 2Gi
              volumeMounts:
                - name: script
                  mountPath: /tmp/scripts
          volumes:
            - name: script
              configMap:
                name: ray-gpu-validate-script
    workerGroupSpecs:
      - groupName: gpu-worker
        replicas: 1
        minReplicas: 1
        maxReplicas: 1
        rayStartParams:
          num-cpus: "1"
          num-gpus: "1"
        template:
          spec:
            containers:
              - name: ray-worker
                image: rayproject/ray-ml:2.49.2.7b0af3-py311-gpu
                resources:
                  requests:
                    cpu: 1
                    memory: 4Gi
                    nvidia.com/gpu: 1
                  limits:
                    cpu: 2
                    memory: 8Gi
                    nvidia.com/gpu: 1
  shutdownAfterJobFinishes: true
  ttlSecondsAfterFinished: 120
EOF
```

### Step 8: Verify RayJob pods are running

Confirm that the Ray head pod, GPU worker pod, and the job submitter pod are created and running.

```
kubectl get pods -n test
NAME                                             READY   STATUS      RESTARTS   AGE
ray-llm-validate-f2w4g                           0/1     Completed   0          2m20s
ray-llm-validate-hjnmh-gpu-worker-worker-9wvsb   1/1     Running     0          2m48s
ray-llm-validate-hjnmh-head-vbqgj                1/1     Running     0          2m48s
```

### Step 9: Verify the RayJob completed successfully

Check the job logs to confirm that the GPU was detected, the matrix multiplication ran on the T4 GPU, and the validation succeeded.

```
kubectl logs -f -n test ray-llm-validate-f2w4g
2026-03-25 15:05:36,409 - INFO - NumExpr defaulting to 4 threads.
2026-03-25 15:05:37,109	INFO cli.py:41 -- Job submission server address: http://ray-llm-validate-hjnmh-head-svc.test.svc.cluster.local:8265
2026-03-25 15:05:37,914	SUCC cli.py:65 -- ---------------------------------------------------
2026-03-25 15:05:37,914	SUCC cli.py:66 -- Job 'ray-llm-validate-xf4pd' submitted successfully
2026-03-25 15:05:37,915	SUCC cli.py:67 -- ---------------------------------------------------
2026-03-25 15:05:37,915	INFO cli.py:291 -- Next steps
2026-03-25 15:05:37,915	INFO cli.py:292 -- Query the logs of the job:
2026-03-25 15:05:37,915	INFO cli.py:294 -- ray job logs ray-llm-validate-xf4pd
2026-03-25 15:05:37,915	INFO cli.py:296 -- Query the status of the job:
2026-03-25 15:05:37,915	INFO cli.py:298 -- ray job status ray-llm-validate-xf4pd
2026-03-25 15:05:37,915	INFO cli.py:300 -- Request the job to be stopped:
2026-03-25 15:05:37,915	INFO cli.py:302 -- ray job stop ray-llm-validate-xf4pd
2026-03-25 15:05:39,387 - INFO - NumExpr defaulting to 4 threads.
2026-03-25 15:05:40,196	INFO cli.py:41 -- Job submission server address: http://ray-llm-validate-hjnmh-head-svc.test.svc.cluster.local:8265
2026-03-25 15:05:37,611	INFO job_manager.py:531 -- Runtime env is setting up.
=== RAY GPU VALIDATION STARTING ===
2026-03-25 15:05:39,916	INFO worker.py:1630 -- Using address 10.244.84.170:6379 set in the environment variable RAY_ADDRESS
2026-03-25 15:05:39,924	INFO worker.py:1771 -- Connecting to existing Ray cluster at address: 10.244.84.170:6379...
2026-03-25 15:05:39,943	INFO worker.py:1942 -- Connected to Ray cluster. View the dashboard at http://10.244.84.170:8265
Waiting for GPU worker...
Cluster resources: {'CPU': 2.0, 'object_store_memory': 1611150950.0, 'memory': 6442450944.0, 'node:__internal_head__': 1.0, 'node:10.244.84.170': 1.0, 'node:10.244.84.137': 1.0, 'GPU': 1.0, 'accelerator_type:T4': 1.0}
Running GPU task on worker...
Worker hostname: ray-llm-validate-hjnmh-gpu-worker-worker-9wvsb
GPU available: True
GPU name: Tesla T4
GPU memory: 15.64 GB
CUDA version: 12.1
Matrix multiply result sum: 22322.4336

VALIDATION SUCCESS: Ray cluster can schedule and run GPU workloads
=== RAY GPU VALIDATION COMPLETED ===
2026-03-25 15:05:48,241	SUCC cli.py:65 -- --------------------------------------
2026-03-25 15:05:48,241	SUCC cli.py:66 -- Job 'ray-llm-validate-xf4pd' succeeded
2026-03-25 15:05:48,241	SUCC cli.py:67 -- --------------------------------------
```

## Summary

The KubeRay validation successfully:

- Deployed the KubeRay operator via Rafay Blueprint to the MKS cluster
- Created a RayJob that provisioned an ephemeral Ray cluster with a GPU worker
- Scheduled and executed a GPU workload (matrix multiplication) on the Tesla T4 GPU
- Confirmed GPU availability (15.64 GB VRAM, CUDA 12.1) and correct task distribution to the worker node

This confirms that KubeRay can schedule and run GPU workloads on the Rafay MKS cluster.
