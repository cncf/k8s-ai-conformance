# AI conformance for Nscale Managed Kubernetes (NKS)

This directory holds the Kubernetes AI Conformance v1.37 evidence for Nscale
Managed Kubernetes (NKS). NKS v1.37 is a [certified Kubernetes
platform](https://github.com/cncf/k8s-conformance/tree/master/v1.37/nks).

The evidence was collected on a freshly provisioned NKS cluster running
Kubernetes v1.37.2 (platform release `v1.37.2-nks.1`):

- a highly available (3 replica) control plane, managed by Nscale;
- a compute node pool of 3 workers;
- a reservation-backed GPU node pool of 2 bare-metal hosts, each with
  8 NVIDIA B200 GPUs (16 in total), exposed through Dynamic Resource Allocation
  (DRA) by the NVIDIA DRA driver (`gpu.nvidia.com`).

The tests come from the upstream
[AI Conformance test suite](https://github.com/kubernetes-sigs/ai-conformance/tree/main/test)
at commit `fd9ba2d8e445068a037a769ee90d9a8165be7677`.

| Requirement | Test | Result |
| --- | --- | --- |
| `ai_inference` | `TestGatewayAPIInference` | pass |
| `secure_accelerator_access` | `TestSecureAcceleratorAccess` | pass |
| `gang_scheduling` | `TestGangScheduling` | pass |
| `robust_controller` | `TestRobustCRDControllerOperation` | pass |
| `accelerator_metrics` | `TestAcceleratorPerformanceMetrics` | pass |
| `ai_service_metrics` | `TestAIServiceMetrics` | pass |
| `pod_autoscaling` | `TestHPAAutoscaling` | pass |
| `driver_runtime_management` | `TestAcceleratorDriverRuntimeManagement` | pass |
| `workload_sandboxing` | `TestWorkloadSandboxing` | skip (no sandboxed RuntimeClass) |
| `cluster_autoscaling` | `TestAcceleratorClusterAutoscaling` | not run (N/A, see `PRODUCT.yaml`) |

`results.json` is the raw `go test -json` stream, `e2e.log` is the test output
it carries, and `junit.xml` is generated from `results.json`.

`high_performance_networking` has no upstream test. `high-performance-networking.log`
is a read-only capture from the same cluster. It shows the `rdma.nscale.com`
DeviceClass and the InfiniBand interfaces that DRANET, delivered by the NKS
hardware add-on, publishes on each GPU host as `dra.net` devices in
ResourceSlices:

```sh
kubectl get deviceclass rdma.nscale.com -o yaml
kubectl get resourceslices -o json \
  | jq '.items[] | select(.spec.driver == "dra.net") | {node: .spec.nodeName, devices: [.spec.devices[].name]}'
```

## Create an NKS cluster

The following steps use the
[Nscale CLI](https://docs.nscale.com/docs/compute/managed-kubernetes).

### 1. Set a CLI context

Create a context that pins the organization, project and region. Region IDs
are listed with `nscale regions list`. The region must be the one that holds
your GPU reservation.

```sh
nscale contexts set --name nks-ai-conformance --user default \
    --org <organization-id> \
    --project <project-id> \
    --region <region-id> \
    --ci
```

### 2. Choose a Kubernetes release

List the available NKS releases and pick the v1.37 release, which matches
`platformVersion` in `PRODUCT.yaml`:

```sh
nscale k8s release list --context nks-ai-conformance
```

### 3. Create the cluster

Describe the cluster in a JSON file. It references an existing network in the
project (`networkId`, listed with `nscale networks list`) and the release
chosen above. The API server is given a public IP so the test tooling can
reach it:

```json
{
  "metadata": { "name": "ai-conformance" },
  "spec": {
    "networkId": "<network-id>",
    "platformReleaseId": "<v1.37 release id>",
    "apiServer": { "publicIP": true, "allowedCidrs": ["0.0.0.0/0"] }
  }
}
```

```sh
nscale k8s cluster create --context nks-ai-conformance --file ./cluster.json --yes --wait
```

This provisions a highly available control plane with Cilium as the CNI. The
NKS hardware add-on is enabled by default; it delivers the NVIDIA GPU stack
(the NVIDIA DRA driver and DCGM exporter, in the `nvidia-platform` namespace)
to the cluster's GPU nodes. The output includes the cluster ID, which later
steps use.

### 4. Add a compute node pool

The AI suite's add-ons and controllers run on ordinary compute workers. Pick a
flavor with `nscale flavors list`:

```json
{
  "metadata": { "name": "default" },
  "spec": {
    "clusterId": "<cluster-id>",
    "provisioningMode": "compute",
    "replicas": 3,
    "compute": { "flavorId": "<flavor-id>" }
  }
}
```

```sh
nscale k8s nodepool create --context nks-ai-conformance --file ./compute-pool.json --yes --wait
```

### 5. Add a GPU node pool backed by a reservation

GPU capacity on Nscale comes from a
[reservation](https://docs.nscale.com/docs/compute/reservations) of bare-metal
GPU hosts. The suite needs at least 2 GPUs that can be allocated at the same
time, because `TestHPAAutoscaling` scales to two Pods holding one GPU each.
Find the reservation in the cluster's organization, project and region:

```sh
nscale reservations list --context nks-ai-conformance
```

Create a node pool from it. The image, host kind and GPU model come from the
platform release and the reservation. A reservation-backed pool cannot be
resized, and its labels and taints cannot be changed later, so set them now if
you need any. The tested pool had 2 replicas and no taints:

```json
{
  "metadata": { "name": "gpu" },
  "spec": {
    "clusterId": "<cluster-id>",
    "provisioningMode": "reservation",
    "replicas": 2,
    "reservation": { "reservationId": "<reservation-id>" }
  }
}
```

```sh
nscale k8s nodepool create --context nks-ai-conformance --file ./gpu-pool.json --yes --wait
```

### 6. Get a kubeconfig

Fetch a kubeconfig for the cluster (using the **cluster** ID) and point
`kubectl` at it. The kubeconfig uses `nscale k8s token` as a client-go exec
plugin, so `kubectl` obtains fresh tokens from the CLI as needed:

```sh
nscale k8s kubeconfig get --context nks-ai-conformance --cluster <cluster-id> > kubeconfig.yaml
export KUBECONFIG=$(pwd)/kubeconfig.yaml
```

### 7. Verify the cluster

Check that every worker is `Ready` on v1.37.2, that each GPU host publishes a
`gpu.nvidia.com` ResourceSlice, and that the DCGM exporter Service exists:

```sh
kubectl get nodes -o wide
kubectl get resourceslices
kubectl get pods,services -n nvidia-platform
```

## Install the workload add-ons

The AI suite exercises components a customer installs into their own cluster.
None of them is part of the NKS platform; each is installed with the upstream
suite's own setup scripts, at the versions pinned in the suite commit:

| Component | Version | Used by |
| --- | --- | --- |
| cert-manager | v1.19.2 | Kubeflow Trainer webhooks |
| Kubeflow Trainer | 2.3.0 | `TestRobustCRDControllerOperation` |
| Kueue | v0.18.2 | `TestGangScheduling` |
| Gateway API CRDs and agentgateway | v1.6.3 and v1.7.0 | `TestGatewayAPIInference` |
| kube-prometheus-stack | 91.9.0 | `TestAIServiceMetrics`, `TestHPAAutoscaling` |
| prometheus-adapter | 5.3.0 | `TestHPAAutoscaling` |

Each script appends the test flags its component needs to
`E2E_TEST_ARGS_FILE`. The scripts need `kubectl` and `helm`, and run from the
suite's repository root:

```sh
git clone https://github.com/kubernetes-sigs/ai-conformance
cd ai-conformance
git checkout fd9ba2d8e445068a037a769ee90d9a8165be7677

export E2E_TEST_ARGS_FILE=$(pwd)/test-args
: > "${E2E_TEST_ARGS_FILE}"
for script in 10-cert-manager 40-kubeflow-trainer 50-gang-scheduler \
    50-agentgateway 50-kube-prometheus-stack 50-prometheus-adapter; do
  bash "test/e2e/setup.d/${script}.sh"
done
```

The agentgateway script creates the `agentgateway-clusterip` GatewayClass,
whose Gateways are exposed as ClusterIP Services, so the test reaches them
from inside the cluster without a cloud load balancer.

The DCGM exporter is not installed: the one NKS delivers is used instead. Point
the accelerator metrics test at it:

```sh
cat >> "${E2E_TEST_ARGS_FILE}" <<EOF
-accelerator-metrics-namespace=nvidia-platform
-accelerator-metrics-service-name=nvidia-dcgm-exporter-dra
EOF
```

`test-args` then holds these flags:

```text
-gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq
-gang-scheduler-name=kueue
-gang-scheduler-namespace=ai-conformance-gang-scheduling
-gateway-api-class=agentgateway-clusterip
-service-metrics-prometheus-service=monitoring/prometheus-operated:9090
-service-metrics-scrape-labels=release=kube-prometheus-stack
-hpa-custom-metric-name=ai_conformance_hpa_queue_depth
-accelerator-metrics-namespace=nvidia-platform
-accelerator-metrics-service-name=nvidia-dcgm-exporter-dra
```

## Run the AI conformance tests

Run every test except `TestAcceleratorClusterAutoscaling`. `results.json`
records the raw event stream, and `e2e.log` is its test output:

```sh
go test -json -count=1 -timeout=90m \
  -run '^(TestGatewayAPIInference|TestSecureAcceleratorAccess|TestGangScheduling|TestRobustCRDControllerOperation|TestAcceleratorPerformanceMetrics|TestAIServiceMetrics|TestHPAAutoscaling|TestAcceleratorDriverRuntimeManagement|TestWorkloadSandboxing)$' \
  ./test -args -kubeconfig="${KUBECONFIG}" $(cat "${E2E_TEST_ARGS_FILE}") \
  > results.json
jq -j 'select(.Output != null) | .Output' results.json > e2e.log
```

The dependency install takes about 4 minutes and the tests about 15. Use a freshly provisioned cluster for each
submission run.
