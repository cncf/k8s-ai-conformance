## Description

The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

## Evidence

### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

```shell
kubectl get nodes
```

```shell
NAME                                     STATUS   ROLES           AGE    VERSION
conformance-test-np1-sms5g-pt7t9-p42p6   Ready    <none>          6d2h   v1.37.1+vmware.4
conformance-test-np1-sms5g-pt7t9-vv2dd   Ready    <none>          6d2h   v1.37.1+vmware.4
conformance-test-pcdv2-gv6sq             Ready    control-plane   6d2h   v1.37.1+vmware.4
```

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install KubeRay Operator

Add the following helm repos and install the kuberay-operator

```shell
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
kubectl create namespace ray-system
helm install kuberay-operator kuberay/kuberay-operator --version 1.7.1 --namespace ray-system \
  --set 'env[0].name=CLUSTER_DOMAIN' \
  --set 'env[0].value=managedcluster1.local' \
  --wait --timeout 300s
# CLUSTER_DOMAIN must match `spec.clusterNetwork.serviceDomain` in the cluster object.
```

```shell
NAME: kuberay-operator
LAST DEPLOYED: Tue Nov  3 23:21:14 2026
NAMESPACE: ray-system
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Verify that kuberay-operator pod is running and in Ready state, and that its CRDs are registered

```shell
kubectl get pods -n ray-system -l app.kubernetes.io/name=kuberay-operator
```

```shell
NAME                                READY   STATUS    RESTARTS   AGE
kuberay-operator-58f6d5cd74-x8trg   1/1     Running   0          17s
```

```shell
kubectl get crd | grep ray.io
```

```shell
rayclusters.ray.io                                            Namespaced   v1(storage),v1alpha1        2026-11-03T23:21:08Z
raycronjobs.ray.io                                            Namespaced   v1(storage)                 2026-11-03T23:21:09Z
rayjobs.ray.io                                                Namespaced   v1(storage),v1alpha1        2026-11-03T23:21:10Z
rayservices.ray.io                                            Namespaced   v1(storage),v1alpha1        2026-11-03T23:21:11Z
```

### Install RayCluster

Install raycluster using the following helmchart

```shell
helm install -f ray-cluster.yaml raycluster kuberay/ray-cluster --version 1.7.1 --namespace ray-system
```

```shell
NAME: raycluster
LAST DEPLOYED: Tue Nov  3 23:21:31 2026
NAMESPACE: ray-system
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

**Note**: Here, `ray-cluster.yaml` overrides the chart defaults with the following values: the pod/container security contexts of head and worker are set as required because this cluster enforces the `restricted` Pod Security Admission profile.

```yaml
image:
  repository: rayproject/ray
  tag: 2.52.0
head:
  podSecurityContext:
    runAsNonRoot: true
    runAsUser: 1000
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
    runAsUser: 1000
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
kubectl get pods -n ray-system -o wide
```

```shell
NAME                                          READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
kuberay-operator-58f6d5cd74-x8trg             1/1     Running   0          78s   193.0.1.71   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
raycluster-kuberay-head-78vm9                 1/1     Running   0          61s   193.0.2.62   conformance-test-np1-sms5g-pt7t9-p42p6   <none>           <none>
raycluster-kuberay-workergroup-worker-9n5mw   1/1     Running   0          61s   193.0.1.72   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl get rayclusters -n ray-system
```

```shell
NAME                 DESIRED WORKERS   AVAILABLE WORKERS   CPUS   MEMORY   GPUS   STATUS   AGE
raycluster-kuberay   1                 1                   2      6Gi      0      ready    61s
```

### Run KubeRay Job

Deploy a ray job for sample workload

```shell
kubectl apply -n ray-system -f ray-job.sample.yaml
```

```shell
rayjob.ray.io/rayjob-sample created
configmap/ray-job-code-sample created
```

**Note**: Here, ray-job.sample.yaml is the [upstream sample manifest](https://raw.githubusercontent.com/ray-project/kuberay/v1.7.1/ray-operator/config/samples/ray-job.sample.yaml) with the following modifications:

- The same `podSecurityContext`/`securityContext` shown above was added to the head and worker pod templates in `rayClusterSpec`, and to a `submitterPodTemplate` (the head, worker and job-submitter pods are otherwise rejected by the `restricted` Pod Security Admission policy enforced on this cluster).
Verify the job status and logs. The job logs should be similar to the snippet below

```shell
kubectl get rayjob -n ray-system rayjob-sample
```

```shell
NAME            JOB STATUS   DEPLOYMENT STATUS   RAY CLUSTER NAME      START TIME             END TIME               AGE
rayjob-sample   SUCCEEDED    Complete            rayjob-sample-nzndp   2026-11-03T23:24:19Z   2026-11-03T23:25:12Z   61s
```

```shell
kubectl get pods -n ray-system
```

```shell
NAME                                           READY   STATUS      RESTARTS   AGE
kuberay-operator-58f6d5cd74-x8trg              1/1     Running     0          4m6s
raycluster-kuberay-head-78vm9                  1/1     Running     0          3m49s
raycluster-kuberay-workergroup-worker-9n5mw    1/1     Running     0          3m49s
rayjob-sample-n8brp                            0/1     Completed   0          38s
rayjob-sample-nzndp-head-trh9l                 1/1     Running     0          60s
rayjob-sample-nzndp-small-group-worker-q6d88   1/1     Running     0          60s
```

```shell
kubectl logs -n ray-system -l job-name=rayjob-sample
```

```shell
Running entrypoint for job rayjob-sample-9kmct: python /home/ray/samples/sample_code.py
2026-11-03 15:24:59,023	INFO worker.py:1696 -- Using address 193.0.1.74:6379 set in the environment variable RAY_ADDRESS
2026-11-03 15:24:59,027	INFO worker.py:1837 -- Connecting to existing Ray cluster at address: 193.0.1.74:6379...
2026-11-03 15:24:59,043	INFO worker.py:2014 -- Connected to Ray cluster. View the dashboard at 193.0.1.74:8265 
test_counter got 1
test_counter got 2
test_counter got 3
test_counter got 4
test_counter got 5
2026-11-03 15:25:09,560	SUCC cli.py:65 -- -----------------------------------
2026-11-03 15:25:09,560	SUCC cli.py:66 -- Job 'rayjob-sample-9kmct' succeeded
2026-11-03 15:25:09,560	SUCC cli.py:67 -- -----------------------------------
```

(A Ray `FutureWarning` and a transient raylet "infeasible resource requests" notice printed before the counter output are omitted above.)

The RayJob controller created a dedicated RayCluster (`rayjob-sample-nzndp`), a submitter Job, ran the Ray workload with its runtime environment (`requests==2.26.0`, `pendulum==2.1.2`, `counter_name` env var — all asserted by the sample code) and reported `SUCCEEDED / Complete`.

### Verify reconciliation of custom resources

**Self-healing:** deleting a worker pod of the `raycluster-kuberay` RayCluster causes the operator to recreate it.

```shell
kubectl delete -n ray-system pod/raycluster-kuberay-workergroup-worker-9n5mw
```

```shell
pod "raycluster-kuberay-workergroup-worker-9n5mw" deleted from ray-system namespace
```

```shell
kubectl get pods -n ray-system -l ray.io/cluster=raycluster-kuberay
```

```shell
NAME                                          READY   STATUS    RESTARTS   AGE
raycluster-kuberay-head-78vm9                 1/1     Running   0          4m20s
raycluster-kuberay-workergroup-worker-xst8g   0/1     Running   0          20s
```

**Spec changes:** scaling the worker group in the RayCluster CR is reconciled into a second worker pod.

```shell
kubectl patch raycluster raycluster-kuberay -n ray-system --type=json -p='[{"op":"replace","path":"/spec/workerGroupSpecs/0/replicas","value":2},{"op":"replace","path":"/spec/workerGroupSpecs/0/maxReplicas","value":3}]'
```

```shell
raycluster.ray.io/raycluster-kuberay patched
```

```shell
kubectl get pods -n ray-system -l ray.io/cluster=raycluster-kuberay
```

```shell
NAME                                          READY   STATUS    RESTARTS   AGE
raycluster-kuberay-head-78vm9                 1/1     Running   0          4m36s
raycluster-kuberay-workergroup-worker-5vwsk   0/1     Running   0          16s
raycluster-kuberay-workergroup-worker-xst8g   1/1     Running   0          36s
```

```shell
kubectl get rayclusters -n ray-system raycluster-kuberay
```

```shell
NAME                 DESIRED WORKERS   AVAILABLE WORKERS   CPUS   MEMORY   GPUS   STATUS   AGE
raycluster-kuberay   2                 2                   3      7Gi      0      ready    4m36s
```

**Admission validation:** an invalid RayCluster (missing the required `headGroupSpec`) is rejected at admission by the CRD's OpenAPI schema validation. (The KubeRay operator Helm chart v1.7.1 does not deploy an admission webhook; RayCluster validation is enforced through the CRD schema and by the operator during reconciliation.)

```shell
cat <<EOF | kubectl apply -n ray-system -f -
apiVersion: ray.io/v1
kind: RayCluster
metadata:
  name: invalid-raycluster
spec:
  rayVersion: "2.52.0"
EOF
```

```shell
The RayCluster "invalid-raycluster" is invalid: 
* spec.headGroupSpec: Required value
* <nil>: Invalid value: null: some validation rules were not checked because the object was invalid; correct the existing errors to complete validation
```

```shell
kubectl get rayclusters -n ray-system
```

```shell
NAME                  DESIRED WORKERS   AVAILABLE WORKERS   CPUS    MEMORY   GPUS   STATUS   AGE
raycluster-kuberay    2                 2                   3       7Gi      0      ready    4m36s
rayjob-sample-nzndp   1                 1                   1200m   2Gi      0      ready    108s
```

This demonstrates that on VKS v1.37 the KubeRay operator installs and runs correctly, its CRDs are served with schema validation enforced at admission, and its custom resources (`RayCluster`, `RayJob`) are reconciled reliably — including self-healing after pod deletion and convergence after spec changes.

### Cleanup

Ray custom resources are removed before the operator so the operator can process their finalizers:

```shell
kubectl delete -n ray-system -f ray-job.sample.yaml
helm uninstall raycluster -n ray-system
helm uninstall kuberay-operator -n ray-system
kubectl get crd -o name | grep '\.ray\.io' | xargs -r kubectl delete   # Helm does not remove CRDs
kubectl delete namespace ray-system
```
