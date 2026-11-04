## Description

If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

## Evidence

Complete guide to configure Horizontal Pod Autoscaling using NVIDIA DCGM metrics with Prometheus Operator on VKS with Kubernetes v1.37.1.

### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator.
* Log in to the cluster as admin

* Install helm tool.

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Prometheus Operator

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts && \
    helm repo update prometheus-community
```

**Note**: The node-exporter DaemonSet of kube-prometheus-stack requires host access, so the pod security policy of the `prometheus` namespace is set to privileged.

```shell
kubectl create namespace prometheus
kubectl label --overwrite ns prometheus pod-security.kubernetes.io/enforce=privileged
```

```shell
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack --version 92.17.1 \
   --namespace prometheus \
   --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
   --wait --timeout 600s
```

```shell
NAME: kube-prometheus-stack
LAST DEPLOYED: Tue Nov  3 23:27:45 2026
NAMESPACE: prometheus
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
```

```shell
kubectl get pods -n prometheus
```

```shell
NAME                                                        READY   STATUS    RESTARTS   AGE
alertmanager-kube-prometheus-stack-alertmanager-0           2/2     Running   0          68s
kube-prometheus-stack-grafana-9586b9fcc-6xd8d               3/3     Running   0          74s
kube-prometheus-stack-kube-state-metrics-78dc75b868-spbw8   1/1     Running   0          74s
kube-prometheus-stack-operator-7545cddfdd-wv67p             1/1     Running   0          74s
kube-prometheus-stack-prometheus-node-exporter-244p8        1/1     Running   0          74s
kube-prometheus-stack-prometheus-node-exporter-q79jl        1/1     Running   0          74s
kube-prometheus-stack-prometheus-node-exporter-xcdk2        1/1     Running   0          74s
prometheus-kube-prometheus-stack-prometheus-0               2/2     Running   0          69s
```

```shell
kubectl get svc -n prometheus
```

```shell
NAME                                             TYPE        CLUSTER-IP        EXTERNAL-IP   PORT(S)                      AGE
alertmanager-operated                            ClusterIP   None              <none>        9093/TCP,9094/TCP,9094/UDP   68s
kube-prometheus-stack-alertmanager               ClusterIP   198.201.255.219   <none>        9093/TCP,8080/TCP            74s
kube-prometheus-stack-grafana                    ClusterIP   198.201.75.83     <none>        80/TCP                       74s
kube-prometheus-stack-kube-state-metrics         ClusterIP   198.201.17.184    <none>        8080/TCP                     74s
kube-prometheus-stack-operator                   ClusterIP   198.201.1.11      <none>        443/TCP                      74s
kube-prometheus-stack-prometheus                 ClusterIP   198.201.75.243    <none>        9090/TCP,8080/TCP            74s
kube-prometheus-stack-prometheus-node-exporter   ClusterIP   198.201.246.27    <none>        9100/TCP                     74s
prometheus-operated                              ClusterIP   None              <none>        9090/TCP                     69s
```

The Prometheus HTTP API is queried below through the Kubernetes API server service proxy, so no external load balancer is needed:

```shell
P=/api/v1/namespaces/prometheus/services/kube-prometheus-stack-prometheus:9090/proxy
```

### Scrape DCGM Exporter metrics

The `nvidia-dcgm-exporter` Service (installed by the GPU Operator) is not annotated with a `ServiceMonitor` by default, so one was created to let the Prometheus Operator discover it. The `honorLabels: true` setting is required — without it, Prometheus' service-discovery `namespace`/`pod` labels (which describe the DCGM exporter's own pod) collide with, and shadow (as `exported_namespace`/`exported_pod`), the `namespace`/`pod` labels that DCGM itself attaches identifying the GPU-consuming workload pod. The Prometheus Adapter (and hence the custom metrics API) needs the un-shadowed `pod`/`namespace` labels to map a metric back to the Pod that owns it.

```shell
cat <<EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: nvidia-dcgm-exporter
  namespace: gpu-operator
  labels:
    app: nvidia-dcgm-exporter
spec:
  selector:
    matchLabels:
      app: nvidia-dcgm-exporter
  endpoints:
  - port: gpu-metrics
    interval: 15s
    honorLabels: true
EOF
```

```shell
servicemonitor.monitoring.coreos.com/nvidia-dcgm-exporter created
```

With 2 GPU-equipped worker nodes, Prometheus discovers and scrapes a separate DCGM exporter target per node:

```shell
kubectl get --raw "$P/api/v1/targets?state=active" | jq -r '.data.activeTargets[] | select(.labels.job=="nvidia-dcgm-exporter") | .labels.instance + "  health=" + .health'
```

```shell
193.0.1.12:9400  health=up
193.0.2.15:9400  health=up
```

### Install Prometheus Adapter

```shell
helm install prometheus-adapter prometheus-community/prometheus-adapter --version 5.4.0 \
  --namespace prometheus \
  --set prometheus.url="http://kube-prometheus-stack-prometheus.prometheus.svc.managedcluster1.local" \
  --set prometheus.port="9090" \
  --wait --timeout 300s
```

```shell
NAME: prometheus-adapter
LAST DEPLOYED: Tue Nov  3 23:29:44 2026
NAMESPACE: prometheus
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Verify that the Prometheus Adapter is running and registered as the custom metrics API.

```shell
kubectl get pods -n prometheus -l app.kubernetes.io/name=prometheus-adapter
```

```shell
NAME                                  READY   STATUS    RESTARTS   AGE
prometheus-adapter-7d88565499-vjkzd   1/1     Running   0          31s
```

```shell
kubectl get apiservice v1beta1.custom.metrics.k8s.io
```

```shell
NAME                            SERVICE                         AVAILABLE   AGE
v1beta1.custom.metrics.k8s.io   prometheus/prometheus-adapter   True        31s
```

### Create a deployment that utilises GPU

```shell
kubectl create namespace hpa-test
```

```shell
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cuda
  namespace: hpa-test
  labels:
    app: cuda
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cuda
  template:
    metadata:
      labels:
        app: cuda
    spec:
      containers:
        - name: cuda
          image: "registry.k8s.io/cuda-vector-add:v0.1"
          command: ["bash", "-c", "for (( c=1; c<=5000; c++ )); do ./vectorAdd; done; sleep 3600"]
          resources:
            limits:
              nvidia.com/gpu: 1
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["ALL"]
            runAsNonRoot: true
            runAsUser: 1000
            seccompProfile:
              type: RuntimeDefault
EOF
```

```shell
kubectl get pods -n hpa-test -l app=cuda -o wide
```

```shell
NAME                    READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
cuda-86c7cd677d-gk7xj   1/1     Running   0          46s   193.0.1.84   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

**Note**: The vGPU in this lab is unlicensed (a firewall blocks the nodes from reaching the NVIDIA license server; see [virtualized_accelerator.md](virtualized_accelerator.md)), so CUDA compute is restricted and `vectorAdd` fails (`Failed to allocate device vector A (error code all CUDA-capable devices are busy or unavailable)!`). No GPU load is generated and utilization stays at 0%, so the HPA scales on `DCGM_FI_DEV_FB_FREE` instead.

Verify that DCGM Exporter related metrics are available in Prometheus for the above deployment, correctly attributed to the `cuda` pod (thanks to `honorLabels: true` on the ServiceMonitor):

```shell
kubectl get --raw "$P/api/v1/query?query=%7B__name__%3D~%22DCGM.*%22%2Cpod%3D%22cuda-86c7cd677d-gk7xj%22%7D" | jq -r '.data.result[] | .metric.__name__ + ": " + .value[1] + "   (namespace=" + .metric.namespace + ", pod=" + .metric.pod + ", Hostname=" + .metric.Hostname + ")"'
```

(The URL-encoded query is `{__name__=~"DCGM.*",pod="cuda-86c7cd677d-gk7xj"}`.)

```shell
DCGM_FI_DEV_SM_CLOCK: 1410   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_MEM_CLOCK: 1215   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_MEMORY_TEMP: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_GPU_UTIL: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_MEM_COPY_UTIL: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_ENC_UTIL: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_DEC_UTIL: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_FB_FREE: 3691   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_FB_USED: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_FB_RESERVED: 404   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
DCGM_FI_DEV_VGPU_LICENSE_STATUS: 0   (namespace=hpa-test, pod=cuda-86c7cd677d-gk7xj, Hostname=conformance-test-np1-sms5g-pt7t9-vv2dd)
```

### Test Kubernetes Custom Metrics Endpoint

Check for the metrics exposed by DCGM Exporter such as `DCGM_FI_DEV_GPU_UTIL`.

```shell
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | jq -r '.resources[].name' | grep DCGM | sort
```

The output confirms `pods/DCGM_FI_DEV_GPU_UTIL`, `pods/DCGM_FI_DEV_FB_FREE`, and the rest of the DCGM metric set are registered against `pods`, `namespaces`, `services`, and `jobs.batch`:

```shell
jobs.batch/DCGM_FI_DEV_DEC_UTIL
jobs.batch/DCGM_FI_DEV_ENC_UTIL
jobs.batch/DCGM_FI_DEV_FB_FREE
jobs.batch/DCGM_FI_DEV_FB_RESERVED
jobs.batch/DCGM_FI_DEV_FB_USED
jobs.batch/DCGM_FI_DEV_GPU_UTIL
jobs.batch/DCGM_FI_DEV_MEMORY_TEMP
jobs.batch/DCGM_FI_DEV_MEM_CLOCK
jobs.batch/DCGM_FI_DEV_MEM_COPY_UTIL
jobs.batch/DCGM_FI_DEV_SM_CLOCK
jobs.batch/DCGM_FI_DEV_VGPU_LICENSE_STATUS
namespaces/DCGM_FI_DEV_DEC_UTIL
namespaces/DCGM_FI_DEV_ENC_UTIL
namespaces/DCGM_FI_DEV_FB_FREE
namespaces/DCGM_FI_DEV_FB_RESERVED
namespaces/DCGM_FI_DEV_FB_USED
namespaces/DCGM_FI_DEV_GPU_UTIL
namespaces/DCGM_FI_DEV_MEMORY_TEMP
namespaces/DCGM_FI_DEV_MEM_CLOCK
namespaces/DCGM_FI_DEV_MEM_COPY_UTIL
namespaces/DCGM_FI_DEV_SM_CLOCK
namespaces/DCGM_FI_DEV_VGPU_LICENSE_STATUS
pods/DCGM_FI_DEV_DEC_UTIL
pods/DCGM_FI_DEV_ENC_UTIL
pods/DCGM_FI_DEV_FB_FREE
pods/DCGM_FI_DEV_FB_RESERVED
pods/DCGM_FI_DEV_FB_USED
pods/DCGM_FI_DEV_GPU_UTIL
pods/DCGM_FI_DEV_MEMORY_TEMP
pods/DCGM_FI_DEV_MEM_CLOCK
pods/DCGM_FI_DEV_MEM_COPY_UTIL
pods/DCGM_FI_DEV_SM_CLOCK
pods/DCGM_FI_DEV_VGPU_LICENSE_STATUS
services/DCGM_FI_DEV_DEC_UTIL
services/DCGM_FI_DEV_ENC_UTIL
services/DCGM_FI_DEV_FB_FREE
services/DCGM_FI_DEV_FB_RESERVED
services/DCGM_FI_DEV_FB_USED
services/DCGM_FI_DEV_GPU_UTIL
services/DCGM_FI_DEV_MEMORY_TEMP
services/DCGM_FI_DEV_MEM_CLOCK
services/DCGM_FI_DEV_MEM_COPY_UTIL
services/DCGM_FI_DEV_SM_CLOCK
services/DCGM_FI_DEV_VGPU_LICENSE_STATUS
```

Confirm the custom metrics API resolves `DCGM_FI_DEV_FB_FREE` back to the `cuda` Pod object (this is what the HPA controller queries):

```shell
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/hpa-test/pods/*/DCGM_FI_DEV_FB_FREE" | jq .
```

```shell
{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "metadata": {},
  "items": [
    {
      "describedObject": {
        "kind": "Pod",
        "namespace": "hpa-test",
        "name": "cuda-86c7cd677d-gk7xj",
        "apiVersion": "/v1"
      },
      "metricName": "DCGM_FI_DEV_FB_FREE",
      "timestamp": "2026-11-03T23:32:25Z",
      "value": "3691",
      "selector": null
    }
  ]
}
```

### Create HPA to Scale deployment based on a GPU metric

Create an HPA resource that utilises a custom metric, `DCGM_FI_DEV_FB_FREE`, to scale the deployment based on free GPU memory. This vGPU profile has 4096 MiB of total framebuffer memory, so the `averageValue` target is set to `2000` (roughly half of the total framebuffer). The cluster has 2 GPUs (one per worker node), so `maxReplicas` is 2:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cuda
  namespace: hpa-test
spec:
  scaleTargetRef:
    kind: Deployment
    name: cuda
    apiVersion: apps/v1
  minReplicas: 1
  maxReplicas: 2
  metrics:
    - type: Pods
      pods:
        metric:
          name: DCGM_FI_DEV_FB_FREE
        target:
          type: AverageValue
          averageValue: "2000"
EOF
```

```shell
horizontalpodautoscaler.autoscaling/cuda created
```

The GPU's free framebuffer memory (3691 MiB) was above the 2000 MiB target, so the HPA scaled the `cuda` Deployment from 1 up to the `maxReplicas` ceiling of 2 — and, with 2 GPUs available across the node pool, both replicas were actually scheduled and became `Running`, with no `Pending` pods:

```shell
kubectl get hpa cuda -n hpa-test -w
```

```shell
NAME   REFERENCE         TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
cuda   Deployment/cuda   3691/2k   1     2     1     0s
cuda   Deployment/cuda   3691/2k   1     2     1     10s
cuda   Deployment/cuda   3691/2k   1     2     2     20s
cuda   Deployment/cuda   3691/2k   1     2     2     30s
```

```shell
kubectl get pods -n hpa-test -l app=cuda -o wide
```

```shell
NAME                    READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
cuda-86c7cd677d-gk7xj   1/1     Running   0          96s   193.0.1.84   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
cuda-86c7cd677d-x85w6   1/1     Running   0          40s   193.0.2.70   conformance-test-np1-sms5g-pt7t9-p42p6   <none>           <none>
```

Each replica landed on a distinct GPU-equipped worker node — the scheduler placed the pods 1-per-node since `nvidia.com/gpu.sharing-strategy` is `none` (exclusive GPU assignment) and each node only has 1 allocatable GPU.

```shell
kubectl describe hpa cuda -n hpa-test
```

```shell
Name:                             cuda
Namespace:                        hpa-test
Labels:                           <none>
Annotations:                      <none>
CreationTimestamp:                Tue, 03 Nov 2026 23:32:35 +0000
Reference:                        Deployment/cuda
Metrics:                          ( current / target )
  "DCGM_FI_DEV_FB_FREE" on pods:  3691 / 2k
Min replicas:                     1
Max replicas:                     2
Deployment pods:                  2 current / 2 desired
Conditions:
  Type            Status  Reason            Message
  ----            ------  ------            -------
  AbleToScale     True    ReadyForNewScale  recommended size matches current size
  ScalingActive   True    ValidMetricFound  the HPA was able to successfully calculate a replica count from pods metric DCGM_FI_DEV_FB_FREE
  ScalingLimited  True    TooManyReplicas   the desired replica count is more than the maximum replica count
  ScaledToZero    False   NotScaledToZero   the HPA controller did not scale the workload to zero
Events:
  Type    Reason             Age   From                       Message
  ----    ------             ----  ----                       -------
  Normal  SuccessfulRescale  41s   horizontal-pod-autoscaler  New size: 2; reason: pods metric DCGM_FI_DEV_FB_FREE above target
```

`ScalingLimited: True / TooManyReplicas` confirms the replica count is capped by the HPA's own `maxReplicas: 2` spec — not by insufficient cluster GPU capacity, since the node pool has exactly 2 allocatable GPUs to match.

```shell
kubectl get pods -A --field-selector=status.phase=Pending
```

```shell
No resources found
```

### Scale down based on the same GPU metric

To show that the HPA also scales GPU pods down based on the custom metric, the target was raised above the observed per-pod value (to 8000 MiB) and the scale-down stabilization window was set to 0s so the result is visible immediately:

```shell
kubectl patch hpa cuda -n hpa-test --type=merge -p '{"spec":{"metrics":[{"type":"Pods","pods":{"metric":{"name":"DCGM_FI_DEV_FB_FREE"},"target":{"type":"AverageValue","averageValue":"8000"}}}],"behavior":{"scaleDown":{"stabilizationWindowSeconds":0}}}}'
```

```shell
horizontalpodautoscaler.autoscaling/cuda patched
```

```shell
kubectl get hpa cuda -n hpa-test -w
```

```shell
NAME   REFERENCE         TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
cuda   Deployment/cuda   3691/8k   1     2     2     51s
cuda   Deployment/cuda   3691/8k   1     2     1     61s
cuda   Deployment/cuda   3691/8k   1     2     1     71s
cuda   Deployment/cuda   3691/8k   1     2     1     81s
```

```shell
kubectl get pods -n hpa-test -l app=cuda -o wide
```

```shell
NAME                    READY   STATUS    RESTARTS   AGE     IP           NODE                                     NOMINATED NODE   READINESS GATES
cuda-86c7cd677d-gk7xj   1/1     Running   0          2m18s   193.0.1.84   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl describe hpa cuda -n hpa-test
```

```shell
...
Metrics:                          ( current / target )
  "DCGM_FI_DEV_FB_FREE" on pods:  3691 / 8k
Min replicas:                     1
Max replicas:                     2
...
Deployment pods:       1 current / 1 desired
Conditions:
  Type            Status  Reason              Message
  ----            ------  ------              -------
  AbleToScale     True    ReadyForNewScale    recommended size matches current size
  ScalingActive   True    ValidMetricFound    the HPA was able to successfully calculate a replica count from pods metric DCGM_FI_DEV_FB_FREE
  ScalingLimited  False   DesiredWithinRange  the desired count is within the acceptable range
  ScaledToZero    False   NotScaledToZero     the HPA controller did not scale the workload to zero
Events:
  Type    Reason             Age   From                       Message
  ----    ------             ----  ----                       -------
  Normal  SuccessfulRescale  82s   horizontal-pod-autoscaler  New size: 2; reason: pods metric DCGM_FI_DEV_FB_FREE above target
  Normal  SuccessfulRescale  31s   horizontal-pod-autoscaler  New size: 1; reason: All metrics below target
```

This demonstrates that VKS v1.37 correctly integrates the Kubernetes HorizontalPodAutoscaler with custom, accelerator-specific metrics for GPU-backed workloads: the `custom.metrics.k8s.io` API served live per-Pod DCGM metrics (`DCGM_FI_DEV_FB_FREE`, `DCGM_FI_DEV_GPU_UTIL`, and the full DCGM metric set) sourced from Prometheus via the Prometheus Adapter, and the HPA controller correctly calculated replica counts from those metrics — scaling the GPU Deployment from 1 up to its `maxReplicas` ceiling of 2 (both replicas genuinely scheduled `Running`, one per GPU-equipped worker node, zero `Pending` pods) and back down to 1 when the metric fell below target.

### Cleanup

```shell
kubectl delete hpa cuda -n hpa-test; kubectl delete deployment cuda -n hpa-test
kubectl delete namespace hpa-test
kubectl delete servicemonitor nvidia-dcgm-exporter -n gpu-operator
helm uninstall prometheus-adapter -n prometheus
helm uninstall kube-prometheus-stack -n prometheus
# kubelet Service created by the Prometheus Operator (not owned by the Helm release)
kubectl delete svc -n kube-system -l app.kubernetes.io/managed-by=prometheus-operator
# Helm does not remove CRDs
kubectl get crd -o name | grep monitoring.coreos.com | xargs -r kubectl delete
kubectl delete namespace prometheus
```
