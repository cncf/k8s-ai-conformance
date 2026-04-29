## Description

The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

## Evidence

### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install KubeRay Operator

Add the following helm repos and install the kuberay-operator

```shell
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
helm install kuberay-operator kuberay/kuberay-operator --version 1.4.2
```

```shell
NAME: kuberay-operator
LAST DEPLOYED: Tue Mar 17 11:49:16 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Verify that kuberay-operator pod is running and in Ready state

```shell
kubectl get po  -w
```
```shell
NAME                               READY   STATUS    RESTARTS   AGE
kuberay-operator-6d6bc7f499-wrfkw  1/1     Running   0          20s
```

### Install RayCluster

Install raycluster using the following helmchart

```shell
helm install -f ray-cluster.yaml raycluster kuberay/ray-cluster --version 1.4.2
```

```shell
NAME: raycluster
LAST DEPLOYED: Tue Mar 17 12:30:05 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

**Note**: Here, the ray-cluster.yaml is the default values.yaml of the chart with modifications to the podsecurity and container security context of containers head and worker as below.

```shell
podSecurityContext: 
  runAsNonRoot: true
  seccompProfile:
    type: RuntimeDefault
securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
```

Verify that raycluster resources are in Ready state

```shell
kubectl get rayclusters
```

```shell
NAME                 DESIRED WORKERS   AVAILABLE WORKERS   CPUS   MEMORY   GPUS   STATUS   AGE
raycluster-kuberay   1                 1                   2      3G       0      ready    26s
```

### Run KubeRay Job

Deploy a ray job for sample workload
```shell
kubectl apply -f https://raw.githubusercontent.com/ray-project/kuberay/34ea80e0f51f80fb092cdc17ca75d4139449edef/ray-operator/config/samples/ray-job.sample.yaml
```

The job logs should be similar to the snippet below
```shell
kubectl logs -l=job-name=rayjob-sample
```

```shell
2026-03-17 05:59:44,657    INFO cli.py:41 -- Job submission server address: http://rayjob-sample-x4w9n-head-svc.default.svc.cluster.local:8265
2026-03-17 05:59:45,258    SUCC cli.py:65 -- ------------------------------------------------
2026-03-17 05:59:45,258    SUCC cli.py:66 -- Job 'rayjob-sample-snlrz' submitted successfully
2026-03-17 05:59:45,258    SUCC cli.py:67 -- ------------------------------------------------
2026-03-17 05:59:45,258    INFO cli.py:291 -- Next steps
2026-03-17 05:59:45,258    INFO cli.py:292 -- Query the logs of the job:
2026-03-17 05:59:45,258    INFO cli.py:294 -- ray job logs rayjob-sample-snlrz
2026-03-17 05:59:45,258    INFO cli.py:296 -- Query the status of the job:
2026-03-17 05:59:45,258    INFO cli.py:298 -- ray job status rayjob-sample-snlrz
2026-03-17 05:59:45,258    INFO cli.py:300 -- Request the job to be stopped:
2026-03-17 05:59:45,258    INFO cli.py:302 -- ray job stop rayjob-sample-snlrz
2026-03-17 05:59:46,805    INFO cli.py:41 -- Job submission server address: http://rayjob-sample-x4w9n-head-svc.default.svc.cluster.local:8265
2026-03-17 05:59:44,991    INFO job_manager.py:531 -- Runtime env is setting up.
2026-03-17 05:59:49,828    INFO worker.py:1554 -- Using address 192.168.1.88:6379 set in the environment variable RAY_ADDRESS
2026-03-17 05:59:49,828    INFO worker.py:1694 -- Connecting to existing Ray cluster at address: 192.168.1.88:6379...
2026-03-17 05:59:49,839    INFO worker.py:1879 -- Connected to Ray cluster. View the dashboard at 192.168.1.88:8265
test_counter got 1
test_counter got 2
test_counter got 3
test_counter got 4
test_counter got 5
2026-03-17 05:59:57,885    SUCC cli.py:65 -- -----------------------------------
2026-03-17 05:59:57,885    SUCC cli.py:66 -- Job 'rayjob-sample-snlrz' succeeded
2026-03-17 05:59:57,885    SUCC cli.py:67 -- -----------------------------------
```
