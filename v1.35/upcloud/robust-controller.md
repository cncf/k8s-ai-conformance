## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator
* Have nvidia device plugin installed. guide: https://upcloud.com/docs/guides/gpu-workloads-managed-kubernetes/ 

## Installing Kuberay

Helm installation
```bash
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
helm install kuberay-operator kuberay/kuberay-operator --version 1.6.0
```
Verify that operator is running
```bash
kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
kuberay-operator-84777d9679-npk6m   1/1     Running   0          24s
```

### Deploy a workload
```bash
kubectl apply -f - <<EOF
apiVersion: ray.io/v1
kind: RayJob
metadata:
  name: kuberay-gpu-test
spec:
  entrypoint: |
    python -c "
    import ray
    import subprocess
    ray.init()
    @ray.remote(num_gpus=1)
    def check_gpu():
        return subprocess.check_output(['nvidia-smi']).decode('utf-8')

    print('--- Ray Resources ---')
    print(ray.cluster_resources())
    print('--- NVIDIA SMI Output ---')
    print(ray.get(check_gpu.remote()))
    "
  rayClusterSpec:
    rayVersion: '2.9.0'
    headGroupSpec:
      template:
        spec:
          containers:
            - name: ray-head
              image: rayproject/ray:2.9.0-gpu
              resources:
                limits:
                  cpu: "1"
                  memory: "2Gi"
                  nvidia.com/gpu: "1"
                requests:
                  cpu: "1"
                  memory: "2Gi"
                  nvidia.com/gpu: "1"
EOF
```

### Find pod and check logs

```bash
kubectl get pods
NAME                                READY   STATUS      RESTARTS   AGE
kuberay-gpu-test-4rtkn              0/1     Completed   0          4m57s
kuberay-gpu-test-5k29t-head-wkbkq   1/1     Running     0          5m35s
kuberay-operator-84777d9679-8s6rv   1/1     Running     0          45m
```

```bash
kubectl logs kuberay-gpu-test-4rtkn

2026-04-20 02:16:14,484	INFO cli.py:36 -- Job submission server address: http://kuberay-gpu-test-5k29t-head-svc.default.svc.cluster.local:8265
2026-04-20 02:16:14,864	SUCC cli.py:60 -- ---------------------------------------------------
2026-04-20 02:16:14,864	SUCC cli.py:61 -- Job 'kuberay-gpu-test-mgjhs' submitted successfully
2026-04-20 02:16:14,864	SUCC cli.py:62 -- ---------------------------------------------------
2026-04-20 02:16:14,864	INFO cli.py:285 -- Next steps
2026-04-20 02:16:14,864	INFO cli.py:286 -- Query the logs of the job:
2026-04-20 02:16:14,864	INFO cli.py:288 -- ray job logs kuberay-gpu-test-mgjhs
2026-04-20 02:16:14,864	INFO cli.py:290 -- Query the status of the job:
2026-04-20 02:16:14,864	INFO cli.py:292 -- ray job status kuberay-gpu-test-mgjhs
2026-04-20 02:16:14,864	INFO cli.py:294 -- Request the job to be stopped:
2026-04-20 02:16:14,864	INFO cli.py:296 -- ray job stop kuberay-gpu-test-mgjhs
2026-04-20 02:16:15,894	INFO cli.py:36 -- Job submission server address: http://kuberay-gpu-test-5k29t-head-svc.default.svc.cluster.local:8265
2026-04-20 02:16:15,545	INFO worker.py:1405 -- Using address 192.168.5.125:6379 set in the environment variable RAY_ADDRESS
2026-04-20 02:16:15,545	INFO worker.py:1540 -- Connecting to existing Ray cluster at address: 192.168.5.125:6379...
2026-04-20 02:16:15,553	INFO worker.py:1715 -- Connected to Ray cluster. View the dashboard at http://192.168.5.125:8265
--- Ray Resources ---
{'node:192.168.5.125': 1.0, 'CPU': 1.0, 'node:__internal_head__': 1.0, 'object_store_memory': 421087641.0, 'accelerator_type:L40S': 1.0, 'GPU': 1.0, 'memory': 2147483648.0}
--- NVIDIA SMI Output ---
Mon Apr 20 02:16:15 2026
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.58.03              Driver Version: 595.58.03      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA L40S                    On  |   00000000:00:07.0 Off |                    0 |
| N/A   31C    P8             25W /  350W |       0MiB /  46068MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|  No running processes found                                                             |
+-----------------------------------------------------------------------------------------+

2026-04-20 02:16:17,917	SUCC cli.py:60 -- --------------------------------------
2026-04-20 02:16:17,917	SUCC cli.py:61 -- Job 'kuberay-gpu-test-mgjhs' succeeded
2026-04-20 02:16:17,917	SUCC cli.py:62 -- --------------------------------------
```

