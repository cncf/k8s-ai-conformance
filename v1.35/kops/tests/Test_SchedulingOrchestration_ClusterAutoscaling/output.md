

# Cluster Autoscaling for GPU Nodes


Creating test namespace "cluster-autoscaling-gpu-1774009686"


## Determine current GPU node count
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



Found 2 GPU nodes initially ([node/i-030791dc5a37cd257 node/i-034383b8fabf72012])


## Deploy GPU probe workload


Applying manifest "testdata/cluster-autoscaling-workload.yaml" to namespace "cluster-autoscaling-gpu-1774009686"
```bash
> kubectl apply -n cluster-autoscaling-gpu-1774009686 -f testdata/cluster-autoscaling-workload.yaml

```

```

deployment.apps/cluster-autoscaling-workload created

```



## Scale deployment to 3 replicas (initial GPU nodes: 2)
```bash
> kubectl scale deployment/cluster-autoscaling-workload -n cluster-autoscaling-gpu-1774009686 --replicas=3

```

```

deployment.apps/cluster-autoscaling-workload scaled

```



### Verify at least one pod is pending
```bash
> kubectl get pods -n cluster-autoscaling-gpu-1774009686 -l app=cluster-autoscaling-workload -o wide

```

```

NAME                                            READY   STATUS              RESTARTS   AGE   IP       NODE                  NOMINATED NODE   READINESS GATES
cluster-autoscaling-workload-7665f769f7-dj5sd   0/1     Pending             0          0s    <none>   i-034383b8fabf72012   <none>           <none>
cluster-autoscaling-workload-7665f769f7-rh97c   0/1     Pending             0          0s    <none>   <none>                <none>           <none>
cluster-autoscaling-workload-7665f769f7-zq2vq   0/1     ContainerCreating   0          0s    <none>   i-030791dc5a37cd257   <none>           <none>

```



## Wait for cluster to scale up
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



### Diagnostics at attempt 1
```bash
> kubectl get pods -n cluster-autoscaling-gpu-1774009686 -l app=cluster-autoscaling-workload -o wide

```

```

NAME                                            READY   STATUS              RESTARTS   AGE   IP       NODE                  NOMINATED NODE   READINESS GATES
cluster-autoscaling-workload-7665f769f7-dj5sd   0/1     ContainerCreating   0          0s    <none>   i-034383b8fabf72012   <none>           <none>
cluster-autoscaling-workload-7665f769f7-rh97c   0/1     Pending             0          0s    <none>   <none>                <none>           <none>
cluster-autoscaling-workload-7665f769f7-zq2vq   0/1     ContainerCreating   0          0s    <none>   i-030791dc5a37cd257   <none>           <none>

```

```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP     OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            16m   v1.35.3   172.20.61.76     3.22.41.71      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            12m   v1.35.3   172.20.121.7     3.141.6.123     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   17m   v1.35.3   172.20.30.71     13.59.192.81    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            12m   v1.35.3   172.20.166.146   3.141.190.20    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            16m   v1.35.3   172.20.127.65    18.117.80.144   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 1: GPU node count is still 2 (need > 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



Attempt 2: GPU node count is still 2 (need > 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



Attempt 3: GPU node count is still 2 (need > 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



Attempt 4: GPU node count is still 2 (need > 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Cluster scaled up: GPU nodes increased from 2 to 3 on attempt 5


## Wait for all replicas to be ready
```bash
> kubectl rollout status deployment/cluster-autoscaling-workload -n cluster-autoscaling-gpu-1774009686 --timeout=600s

```

```

Waiting for deployment "cluster-autoscaling-workload" rollout to finish: 2 of 3 updated replicas are available...
deployment "cluster-autoscaling-workload" successfully rolled out

```



### Verify GPU pods are running
```bash
> kubectl get pods -n cluster-autoscaling-gpu-1774009686 -l app=cluster-autoscaling-workload -o wide

```

```

NAME                                            READY   STATUS    RESTARTS   AGE    IP             NODE                  NOMINATED NODE   READINESS GATES
cluster-autoscaling-workload-7665f769f7-dj5sd   1/1     Running   0          5m8s   100.96.5.245   i-034383b8fabf72012   <none>           <none>
cluster-autoscaling-workload-7665f769f7-rh97c   1/1     Running   0          5m8s   100.96.7.40    i-07d9fe57777779c63   <none>           <none>
cluster-autoscaling-workload-7665f769f7-zq2vq   1/1     Running   0          5m8s   100.96.4.134   i-030791dc5a37cd257   <none>           <none>

```

```bash
> kubectl get pods -n cluster-autoscaling-gpu-1774009686 -l app=cluster-autoscaling-workload -o name

```

```

pod/cluster-autoscaling-workload-7665f769f7-dj5sd
pod/cluster-autoscaling-workload-7665f769f7-rh97c
pod/cluster-autoscaling-workload-7665f769f7-zq2vq

```

```bash
> kubectl logs -n cluster-autoscaling-gpu-1774009686 pod/cluster-autoscaling-workload-7665f769f7-dj5sd --tail=5

```

```

NVIDIA L4, 33, 0 %, 0 MiB, 23034 MiB
2026-03-20T12:32:11+00:00 GPU probe alive on cluster-autoscaling-workload-7665f769f7-dj5sd
NVIDIA L4, 33, 0 %, 0 MiB, 23034 MiB
2026-03-20T12:33:11+00:00 GPU probe alive on cluster-autoscaling-workload-7665f769f7-dj5sd
NVIDIA L4, 33, 0 %, 0 MiB, 23034 MiB

```

```bash
> kubectl logs -n cluster-autoscaling-gpu-1774009686 pod/cluster-autoscaling-workload-7665f769f7-rh97c --tail=5

```

```

2026-03-20T12:33:14+00:00 GPU probe alive on cluster-autoscaling-workload-7665f769f7-rh97c
NVIDIA L4, 39, 0 %, 0 MiB, 23034 MiB

```

```bash
> kubectl logs -n cluster-autoscaling-gpu-1774009686 pod/cluster-autoscaling-workload-7665f769f7-zq2vq --tail=5

```

```

NVIDIA L4, 30, 0 %, 0 MiB, 23034 MiB
2026-03-20T12:32:07+00:00 GPU probe alive on cluster-autoscaling-workload-7665f769f7-zq2vq
NVIDIA L4, 29, 0 %, 0 MiB, 23034 MiB
2026-03-20T12:33:07+00:00 GPU probe alive on cluster-autoscaling-workload-7665f769f7-zq2vq
NVIDIA L4, 29, 0 %, 0 MiB, 23034 MiB

```

&check; Cluster autoscaler scaled up GPU nodes from 2 to accommodate 3 GPU pods



## Scale down and verify cluster scale-down
```bash
> kubectl scale deployment/cluster-autoscaling-workload -n cluster-autoscaling-gpu-1774009686 --replicas=0

```

```

deployment.apps/cluster-autoscaling-workload scaled

```



Waiting for cluster to scale down (this may take several minutes)...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 1
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            21m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            18m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   22m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            18m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   Ready    node            4m    v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            21m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 1: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 2: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 3: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 4: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 5: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 6
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE     VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            24m     v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            20m     v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   24m     v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            20m     v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   Ready    node            6m31s   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            24m     v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 6: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 7: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 8: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 9: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 10: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 11
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE    VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            26m    v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            23m    v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   27m    v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            23m    v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   Ready    node            9m2s   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            26m    v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 11: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 12: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 13: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 14: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 15: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 16
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            29m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            25m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   29m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            25m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   Ready    node            11m   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            29m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 16: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 17: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 18: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 19: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 20: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 21
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS   ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready    node            31m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready    node            28m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready    control-plane   32m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready    node            28m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   Ready    node            14m   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0ab9100a694772f64   Ready    node            31m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 21: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 22: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 23: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 24: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 25: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 26
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS                        ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready                         node            34m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready                         node            30m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready                         control-plane   34m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready                         node            30m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   NotReady,SchedulingDisabled   node            16m   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://Unknown
i-0ab9100a694772f64   Ready                         node            34m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 26: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 27: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 28: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 29: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 30: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 31
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS                        ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready                         node            36m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready                         node            33m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready                         control-plane   37m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready                         node            33m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   NotReady,SchedulingDisabled   node            19m   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://Unknown
i-0ab9100a694772f64   Ready                         node            36m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 31: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 32: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 33: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 34: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



Attempt 35: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012
node/i-07d9fe57777779c63

```



### Scale-down diagnostics at attempt 36
```bash
> kubectl get nodes -o wide

```

```

NAME                  STATUS                        ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP      OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
i-02fb36e88acb587c7   Ready                         node            39m   v1.35.3   172.20.61.76     3.22.41.71       Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-030791dc5a37cd257   Ready                         node            35m   v1.35.3   172.20.121.7     3.141.6.123      Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-0308c4a041ac88ead   Ready                         control-plane   40m   v1.35.3   172.20.30.71     13.59.192.81     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-034383b8fabf72012   Ready                         node            35m   v1.35.3   172.20.166.146   3.141.190.20     Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6
i-07d9fe57777779c63   NotReady,SchedulingDisabled   node            21m   v1.35.3   172.20.26.144    18.119.104.226   Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://Unknown
i-0ab9100a694772f64   Ready                         node            39m   v1.35.3   172.20.127.65    18.117.80.144    Ubuntu 24.04.3 LTS   6.14.0-1018-aws   containerd://2.1.6

```



Attempt 36: GPU node count is still 3 (need <= 2), waiting 30s...
```bash
> kubectl get nodes -l nvidia.com/gpu.present=true -o name

```

```

node/i-030791dc5a37cd257
node/i-034383b8fabf72012

```



Cluster scaled down: GPU nodes decreased to 2 on attempt 37
&check; Cluster autoscaler scaled down GPU nodes back to 2



Deleting test namespace "cluster-autoscaling-gpu-1774009686"


Namespace deletion took 6s
