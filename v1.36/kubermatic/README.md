# Kubermatic Kubernetes Platform: AI Conformance v1.36

`e2e.log`, `junit.xml` and `results.json` come from one run of the upstream
[AI Conformance test suite](https://github.com/kubernetes-sigs/ai-conformance/tree/b66dc7982aa8c9ad2ff277cb0a31e0417b376114/test)
at commit `b66dc7982aa8c9ad2ff277cb0a31e0417b376114` against a KKP user cluster.
This page describes how to reproduce that run.

## Environment

- KKP Enterprise Edition v2.31.
- One user cluster on AWS, Kubernetes v1.36.5, Cilium CNI.
- One MachineDeployment with a `g4dn.xlarge` node (one NVIDIA T4), Ubuntu 24.04.

## Cluster setup

`KUBECONFIG` points to the user cluster.

1. Install the NVIDIA GPU Operator and Kueue from the KKP default application
   catalog:

   ```bash
   kubectl create namespace nvidia-gpu-operator
   kubectl create namespace kueue
   ```

   ```yaml
   apiVersion: apps.kubermatic.k8c.io/v1
   kind: ApplicationInstallation
   metadata:
     name: nvidia-gpu-operator
     namespace: nvidia-gpu-operator
   spec:
     applicationRef:
       name: nvidia-gpu-operator
       version: v26.3.3
     namespace:
       name: nvidia-gpu-operator
       create: true
     valuesBlock: |
       node-feature-discovery:
         fullnameOverride: gpu-operator-node-feature-discovery
   ---
   apiVersion: apps.kubermatic.k8c.io/v1
   kind: ApplicationInstallation
   metadata:
     name: kueue
     namespace: kueue
   spec:
     applicationRef:
       name: kueue
       version: v0.19.6
     namespace:
       name: kueue
       create: true
     valuesBlock: |
       nameOverride: kueue
       fullnameOverride: kueue
   ```

   Wait until the node reports `nvidia.com/gpu` as allocatable. The GPU
   Operator builds the driver on the node, which takes a few minutes.

2. Install Kubeflow Trainer with Helm:

   ```bash
   helm upgrade -i kubeflow-trainer oci://ghcr.io/kubeflow/charts/kubeflow-trainer \
     --version 2.3.0 --namespace kubeflow-system --create-namespace \
     --set runtimes.defaultEnabled=true --wait
   ```

3. Enable user cluster monitoring (MLA) for the cluster. The seed needs the
   User Cluster MLA stack installed. This deploys the monitoring agent in
   `mla-system`, which sends metrics to Cortex on the seed.

4. Create the Kueue queues for the gang scheduling test:

   ```yaml
   apiVersion: kueue.x-k8s.io/v1beta2
   kind: ResourceFlavor
   metadata:
     name: e2e-flavor
   ---
   apiVersion: kueue.x-k8s.io/v1beta2
   kind: ClusterQueue
   metadata:
     name: e2e-cq
   spec:
     namespaceSelector: {}
     resourceGroups:
     - coveredResources: ["cpu", "memory"]
       flavors:
       - name: e2e-flavor
         resources:
         - name: "cpu"
           nominalQuota: 10
         - name: "memory"
           nominalQuota: 10Gi
   ---
   apiVersion: v1
   kind: Namespace
   metadata:
     name: ai-conformance-gang-scheduling
   ---
   apiVersion: kueue.x-k8s.io/v1beta2
   kind: LocalQueue
   metadata:
     name: e2e-lq
     namespace: ai-conformance-gang-scheduling
   spec:
     clusterQueue: e2e-cq
   ```

## Metrics through MLA

`mla-scrape-manifest.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: monitoring-scraping-ai-conformance-{{.RunID}}
  namespace: mla-system
data:
  config.yaml: |
    - job_name: ai-conformance-{{.RunID}}
      scrape_interval: {{.ScrapeInterval}}
      kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names:
          - {{.Namespace}}
      relabel_configs:
      - source_labels: [__meta_kubernetes_service_name]
        regex: {{.ServiceName}}
        action: keep
      - source_labels: [__meta_kubernetes_endpoint_port_name]
        regex: {{.PortName}}
        action: keep
```

On the seed, during the test:

```bash
kubectl -n cluster-<cluster-id> port-forward svc/mla-gateway 8081:80
```

## Running the tests

From a checkout of `kubernetes-sigs/ai-conformance` at the commit above, with
`KUBECONFIG` pointing to the user cluster:

```bash
go test -v ./test -timeout 60m -json -args \
  -gang-scheduler-name=kueue \
  -gang-scheduler-namespace=ai-conformance-gang-scheduling \
  -gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq \
  -operator=kubeflow-trainer \
  -service-metrics-query-url=http://localhost:8081/api/prom \
  -service-metrics-scrape-manifest=mla-scrape-manifest.yaml \
  > results.json

jq -rj 'select(.Action=="output") | .Output' results.json > e2e.log
go run gotest.tools/gotestsum@latest --junitfile junit.xml --raw-command -- cat results.json
```
