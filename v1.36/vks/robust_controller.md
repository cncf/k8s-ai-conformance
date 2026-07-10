## Description

The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

## Evidence

### Prerequisites

* Provision a VKS v3.7.0 Cluster with v1.36.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

```shell
kubectl get nodes
```

```shell
NAME                                     STATUS   ROLES           AGE     VERSION
conformance-test-kjs2q-nwp9v             Ready    control-plane   4d13h   v1.36.1+vmware.4
conformance-test-np1-qcqc8-bpdtl-rcj8j   Ready    <none>          4d13h   v1.36.1+vmware.4
conformance-test-np1-qcqc8-bpdtl-wv2gs   Ready    <none>          4d13h   v1.36.1+vmware.4
```

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install KubeRay Operator

Add the following helm repos and install the kuberay-operator

```shell
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
helm install kuberay-operator kuberay/kuberay-operator --version 1.4.2 \
  --set 'env[0].name=CLUSTER_DOMAIN' \
  --set 'env[0].value=managedcluster1.local'  # Match this value to `spec.clusterNetwork.serviceDomain` in cluster object.
```

```shell
NAME: kuberay-operator
LAST DEPLOYED: Tue Jul  7 08:33:23 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Verify that kuberay-operator pod is running and in Ready state

```shell
kubectl get po -l app.kubernetes.io/name=kuberay-operator
```
```shell
NAME                                READY   STATUS    RESTARTS   AGE
kuberay-operator-757897cbf4-zrqhl   1/1     Running   0          15m
```

### Install RayCluster

Install raycluster using the following helmchart

```shell
helm install -f ray-cluster.yaml raycluster kuberay/ray-cluster --version 1.4.2
```

```shell
NAME: raycluster
LAST DEPLOYED: Tue Jul  7 07:40:07 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

**Note**: Here, the ray-cluster.yaml is the default values.yaml of the chart with modifications to the podsecurity and container security context of containers head and worker as below (required because this cluster enforces the `restricted` Pod Security Admission profile).

```yaml
head:
  podSecurityContext:
    runAsNonRoot: true
    seccompProfile:
      type: RuntimeDefault
  securityContext:
    allowPrivilegeEscalation: false
    capabilities:
      drop:
        - ALL

worker:
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

Verify that raycluster pods and resources are in Ready state

```shell
kubectl get po
```

```shell
NAME                                          READY   STATUS    RESTARTS   AGE
kuberay-operator-757897cbf4-zrqhl             1/1     Running   0          15m
raycluster-kuberay-head-ks8mq                 1/1     Running   0          68m
raycluster-kuberay-workergroup-worker-xrjvl   1/1     Running   0          14m
```

```shell
kubectl get rayclusters
```

```shell
NAME                 DESIRED WORKERS   AVAILABLE WORKERS   CPUS   MEMORY   GPUS   STATUS   AGE
raycluster-kuberay   1                 1                   2      3G       0      ready    68m
```

### Run KubeRay Job

Deploy a ray job for sample workload

```shell
kubectl apply -f ray-job.sample.yaml
```

```shell
rayjob.ray.io/rayjob-sample created
configmap/ray-job-code-sample created
```

**Note**: Here, ray-job.sample.yaml is the [upstream sample manifest](https://raw.githubusercontent.com/ray-project/kuberay/34ea80e0f51f80fb092cdc17ca75d4139449edef/ray-operator/config/samples/ray-job.sample.yaml) with the following modifications:

- The same `podSecurityContext`/`securityContext` shown above was added to the head and worker pod templates in `rayClusterSpec`, and to a `submitterPodTemplate` (the head, worker and job-submitter pods are otherwise rejected by the `restricted` Pod Security Admission policy enforced on this cluster).

Verify the job status and logs. The job logs should be similar to the snippet below

```shell
kubectl get rayjob rayjob-sample
```

```shell
NAME            JOB STATUS   DEPLOYMENT STATUS   RAY CLUSTER NAME      START TIME             END TIME               AGE
rayjob-sample   SUCCEEDED    Complete            rayjob-sample-xtfq4   2026-07-07T08:45:00Z   2026-07-07T08:45:54Z   3m43s
```

```shell
kubectl logs -l=job-name=rayjob-sample
```

```shell
2026-07-07 01:45:40,858	INFO worker.py:1694 -- Connecting to existing Ray cluster at address: 193.0.1.37:6379...
2026-07-07 01:45:40,867	INFO worker.py:1879 -- Connected to Ray cluster. View the dashboard at 193.0.1.37:8265
test_counter got 1
test_counter got 2
test_counter got 3
test_counter got 4
test_counter got 5
2026-07-07 01:45:51,566	SUCC cli.py:65 -- -----------------------------------
2026-07-07 01:45:51,566	SUCC cli.py:66 -- Job 'rayjob-sample-2jsn7' succeeded
2026-07-07 01:45:51,566	SUCC cli.py:67 -- -----------------------------------
```
