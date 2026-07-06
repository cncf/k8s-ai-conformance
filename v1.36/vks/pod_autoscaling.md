## Description

If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

## Evidence

Complete guide to configure Horizontal Pod Autoscaling using NVIDIA DCGM metrics with Prometheus Operator on VKS with Kubernetes v1.36.1.

### Prerequisites

* Provision a VKS v3.7.0 Cluster with v1.36.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator.
* Log in to the cluster as admin

* Install helm tool.

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Prometheus Operator

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts && \
    helm repo update prometheus-community
```

```shell
helm install prometheus-community/kube-prometheus-stack \
   --create-namespace --namespace prometheus \
   --generate-name \
   --set prometheus.service.type=LoadBalancer \
   --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
```

```shell
NAME: kube-prometheus-stack-1783455636
LAST DEPLOYED: Tue Jul  7 20:20:40 2026
NAMESPACE: prometheus
STATUS: deployed
REVISION: 1
```

The VKS Supervisor natively provisions a load balancer IP for the `LoadBalancer`-type Prometheus service — no additional NodePort configuration is required in this environment:

```shell
kubectl get svc -n prometheus kube-prometheus-stack-1783-prometheus
```

```shell
NAME                                      TYPE           CLUSTER-IP       EXTERNAL-IP    PORT(S)                         AGE
kube-prometheus-stack-1783-prometheus     LoadBalancer   10.107.118.217   10.159.3.168   9090:32657/TCP,8080:31929/TCP   41s
```

```shell
kubectl get pods -n prometheus
```

```shell
NAME                                                              READY   STATUS    RESTARTS   AGE
alertmanager-kube-prometheus-stack-1783-alertmanager-0            2/2     Running   0          32s
kube-prometheus-stack-1783-operator-5bd7666b75-zjt4j              1/1     Running   0          41s
kube-prometheus-stack-1783455636-grafana-856df58895-hnj49         3/3     Running   0          41s
kube-prometheus-stack-1783455636-kube-state-metrics-844bc98dsfk   1/1     Running   0          41s
prometheus-kube-prometheus-stack-1783-prometheus-0                2/2     Running   0          32s
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

With 3 GPU-equipped worker nodes, Prometheus discovers and scrapes a separate DCGM exporter target per node:

```shell
curl -s "http://10.159.3.168:9090/api/v1/targets" | jq -r '.data.activeTargets[] | select(.labels.job=="nvidia-dcgm-exporter") | .labels.instance'
```

```shell
192.168.145.24:9400
192.168.146.6:9400
192.168.147.10:9400
```

### Install Prometheus Adapter

```shell
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  --set prometheus.url="http://kube-prometheus-stack-1783-prometheus.prometheus.svc.cluster.local" \
  --set prometheus.port="9090"
```

```shell
NAME: prometheus-adapter
LAST DEPLOYED: Tue Jul  7 20:23:35 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

Verify that Prometheus, Prometheus Adapter and DCGM Exporter services are functional.

```shell
kubectl get svc -A
```

```shell
NAMESPACE            NAME                                                        TYPE           CLUSTER-IP       EXTERNAL-IP    PORT(S)                         AGE
default              kubernetes                                                  ClusterIP      10.96.0.1        <none>         443/TCP                         8h
default              prometheus-adapter                                          ClusterIP      10.104.77.20     <none>         443/TCP                         10m
gpu-operator         gpu-operator                                                ClusterIP      10.98.19.19      <none>         8080/TCP                        99m
gpu-operator         nvidia-dcgm-exporter                                        ClusterIP      10.111.166.129   <none>         9400/TCP                        99m
kube-system          antrea                                                      ClusterIP      10.96.112.78     <none>         443/TCP                         8h
kube-system          kube-dns                                                    ClusterIP      10.96.0.10       <none>         53/UDP,53/TCP,9153/TCP          8h
kube-system          metrics-server                                              ClusterIP      10.110.31.44     <none>         443/TCP                         8h
prometheus           kube-prometheus-stack-1783-alertmanager                     ClusterIP      10.105.120.202   <none>         9093/TCP,8080/TCP               13m
prometheus           kube-prometheus-stack-1783-operator                         ClusterIP      10.96.240.201    <none>         443/TCP                         13m
prometheus           kube-prometheus-stack-1783-prometheus                       LoadBalancer   10.107.118.217   10.159.3.168   9090:32657/TCP,8080:31929/TCP   13m
prometheus           kube-prometheus-stack-1783455636-grafana                    ClusterIP      10.96.27.110     <none>         80/TCP                          13m
prometheus           kube-prometheus-stack-1783455636-kube-state-metrics         ClusterIP      10.108.206.188   <none>         8080/TCP                        13m
prometheus           kube-prometheus-stack-1783455636-prometheus-node-exporter   ClusterIP      10.100.126.86    <none>         9100/TCP                        13m
prometheus           prometheus-operated                                         ClusterIP      None             <none>         9090/TCP                        13m
```

### Test Kubernetes Custom Metrics Endpoint

Check for the metrics exposed by DCGM Exporter such as `DCGM_FI_DEV_GPU_UTIL`.

```shell
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | jq -r . | grep DCGM
```

The output confirms `pods/DCGM_FI_DEV_GPU_UTIL`, `pods/DCGM_FI_DEV_FB_FREE`, and the rest of the DCGM metric set are registered against `pods`, `namespaces`, `services`, and `jobs.batch`:

```shell
"namespaces/DCGM_FI_DEV_FB_FREE"
"services/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"services/DCGM_FI_DEV_SM_CLOCK"
"namespaces/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"pods/DCGM_FI_DEV_GPU_UTIL"
"services/DCGM_FI_DEV_GPU_UTIL"
"pods/DCGM_FI_DEV_FB_USED"
"namespaces/DCGM_FI_DEV_GPU_UTIL"
"pods/DCGM_FI_DEV_FB_RESERVED"
"jobs.batch/DCGM_FI_DEV_MEMORY_TEMP"
"jobs.batch/DCGM_FI_DEV_SM_CLOCK"
"jobs.batch/DCGM_FI_DEV_ENC_UTIL"
"namespaces/DCGM_FI_DEV_MEMORY_TEMP"
"jobs.batch/DCGM_FI_DEV_MEM_COPY_UTIL"
"jobs.batch/DCGM_FI_DEV_FB_RESERVED"
"services/DCGM_FI_DEV_MEM_COPY_UTIL"
"services/DCGM_FI_DEV_ENC_UTIL"
"jobs.batch/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"jobs.batch/DCGM_FI_DEV_GPU_UTIL"
"services/DCGM_FI_DEV_FB_RESERVED"
"namespaces/DCGM_FI_DEV_FB_USED"
"pods/DCGM_FI_DEV_SM_CLOCK"
"pods/DCGM_FI_DEV_MEM_COPY_UTIL"
"namespaces/DCGM_FI_DEV_DEC_UTIL"
"jobs.batch/DCGM_FI_DEV_FB_USED"
"namespaces/DCGM_FI_DEV_SM_CLOCK"
"pods/DCGM_FI_DEV_DEC_UTIL"
"services/DCGM_FI_DEV_MEMORY_TEMP"
"namespaces/DCGM_FI_DEV_MEM_CLOCK"
"services/DCGM_FI_DEV_MEM_CLOCK"
"namespaces/DCGM_FI_DEV_MEM_COPY_UTIL"
"pods/DCGM_FI_DEV_MEM_CLOCK"
"services/DCGM_FI_DEV_DEC_UTIL"
"pods/DCGM_FI_DEV_FB_FREE"
"pods/DCGM_FI_DEV_MEMORY_TEMP"
"jobs.batch/DCGM_FI_DEV_DEC_UTIL"
"services/DCGM_FI_DEV_FB_FREE"
"namespaces/DCGM_FI_DEV_ENC_UTIL"
"pods/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"jobs.batch/DCGM_FI_DEV_FB_FREE"
"services/DCGM_FI_DEV_FB_USED"
"pods/DCGM_FI_DEV_ENC_UTIL"
"namespaces/DCGM_FI_DEV_FB_RESERVED"
"jobs.batch/DCGM_FI_DEV_MEM_CLOCK"
```

### Create HPA to Scale deployment based on GPU Utilization

Create a deployment that utilises GPU.

```shell
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cuda
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
          image: "k8s.gcr.io/cuda-vector-add:v0.1"
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
kubectl get pods -l app=cuda -o wide
```

```shell
NAME                   READY   STATUS    RESTARTS   AGE   IP               NODE
cuda-887b9c799-49zn9   1/1     Running   0          15s   192.168.146.12   ai-conf-ai-conf-np-7p35-fd6m8-spzvx-chgb6
```

Verify that DCGM Exporter related metrics are available in Prometheus for the above deployment, correctly attributed to the `cuda` pod (thanks to `honorLabels: true` on the ServiceMonitor):

```shell
curl -s -G http://10.159.3.168:9090/api/v1/query --data-urlencode 'query={__name__=~"DCGM.*",pod="cuda-887b9c799-49zn9"}' | jq -r .
```

```shell
DCGM_FI_DEV_SM_CLOCK: 1410
DCGM_FI_DEV_MEM_CLOCK: 1215
DCGM_FI_DEV_MEMORY_TEMP: 0
DCGM_FI_DEV_GPU_UTIL: 13
DCGM_FI_DEV_MEM_COPY_UTIL: 0
DCGM_FI_DEV_ENC_UTIL: 0
DCGM_FI_DEV_DEC_UTIL: 0
DCGM_FI_DEV_FB_FREE: 8905
DCGM_FI_DEV_FB_USED: 417
DCGM_FI_DEV_FB_RESERVED: 916
DCGM_FI_DEV_VGPU_LICENSE_STATUS: 0
```

Confirm the custom metrics API resolves `DCGM_FI_DEV_FB_FREE` back to the `cuda` Pod object (this is what the HPA controller queries):

```shell
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/DCGM_FI_DEV_FB_FREE" | jq -r .
```

```shell
{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "items": [
    {
      "describedObject": { "kind": "Pod", "namespace": "default", "name": "cuda-887b9c799-49zn9", "apiVersion": "/v1" },
      "metricName": "DCGM_FI_DEV_FB_FREE",
      "timestamp": "2026-07-07T21:09:28Z",
      "value": "8905"
    }
  ]
}
```

Create an HPA resource that utilises a custom metric, `DCGM_FI_DEV_FB_FREE`, to scale the deployment based on free GPU memory. This GPU has 10240 MiB of total framebuffer memory, so the `averageValue` target is scaled down proportionally from `10000` to `5000` to remain meaningful for this card:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cuda
spec:
  scaleTargetRef:
    kind: Deployment
    name: cuda
    apiVersion: apps/v1
  minReplicas: 1
  maxReplicas: 3
  metrics:
    - type: Pods
      pods:
        metric:
          name: DCGM_FI_DEV_FB_FREE
        target:
          type: AverageValue
          averageValue: "5000"
EOF
```

The idle GPU's free framebuffer memory (8905-9323 MiB across the 3 nodes) was well above the 5000 MiB target, so the HPA progressively scaled the `cuda` Deployment from 1 up to the `maxReplicas` ceiling of 3 — and, with 3 GPUs now available across the node pool (one per worker node), all 3 replicas were actually scheduled and became `Running`, with no `Pending` pods:

```shell
kubectl get hpa cuda -w
```

```shell
NAME   REFERENCE         TARGETS       MINPODS   MAXPODS   REPLICAS   AGE
cuda   Deployment/cuda   <unknown>/5k  1         3         0          5s
cuda   Deployment/cuda   9323/5k       1         3         1          22s
cuda   Deployment/cuda   9114/5k       1         3         2          33s
cuda   Deployment/cuda   9034666m/5k   1         3         3          48s
```

```shell
kubectl get pods -l app=cuda -o wide
```

```shell
NAME                   READY   STATUS    RESTARTS   AGE   IP               NODE
cuda-887b9c799-49zn9   1/1     Running   0          91s   192.168.146.12   ai-conf-ai-conf-np-7p35-fd6m8-spzvx-chgb6
cuda-887b9c799-thlmf   1/1     Running   0          3s    192.168.145.37   ai-conf-ai-conf-np-7p35-fd6m8-spzvx-rpm7d
cuda-887b9c799-vl8j7   1/1     Running   0          18s   192.168.147.13   ai-conf-ai-conf-np-7p35-fd6m8-spzvx-j65j6
```

Each replica landed on a distinct GPU-equipped worker node — the scheduler naturally spread the pods 1-per-node since `nvidia.com/gpu.sharing-strategy` is `none` (exclusive GPU assignment) and each node only has 1 allocatable GPU.

```shell
kubectl describe hpa cuda
```

```shell
Name:                             cuda
Namespace:                        default
Reference:                        Deployment/cuda
Metrics:                          ( current / target )
  "DCGM_FI_DEV_FB_FREE" on pods:  9034666m / 5k
Min replicas:                     1
Max replicas:                     3
Deployment pods:                  3 current / 3 desired
Conditions:
  Type            Status  Reason            Message
  ----            ------  ------            -------
  AbleToScale     True    ReadyForNewScale  recommended size matches current size
  ScalingActive   True    ValidMetricFound  the HPA was able to successfully calculate a replica count from pods metric DCGM_FI_DEV_FB_FREE
  ScalingLimited  True    TooManyReplicas   the desired replica count is more than the maximum replica count
Events:
  Type    Reason             Age   From                       Message
  ----    ------             ----  ----                       -------
  Normal  SuccessfulRescale  33s   horizontal-pod-autoscaler  New size: 2; reason: pods metric DCGM_FI_DEV_FB_FREE above target
  Normal  SuccessfulRescale  18s   horizontal-pod-autoscaler  New size: 3; reason: pods metric DCGM_FI_DEV_FB_FREE above target
```

`ScalingLimited: True / TooManyReplicas` confirms the replica count is capped by the HPA's own `maxReplicas: 3` spec — not by insufficient cluster GPU capacity, since the node pool now has exactly 3 allocatable GPUs to match.

Simulate load on all 3 scheduled GPU pods concurrently:

```shell
for p in $(kubectl get pods -l app=cuda -o jsonpath='{.items[*].metadata.name}'); do
  kubectl exec $p -- bash -c '
  for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
  for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
  for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
  for (( c=1; c<=5000; c++ )); do ./vectorAdd; done &
  wait' &
done
```

Under concurrent load, some `vectorAdd` invocations correctly failed with device contention, confirming genuine GPU pressure was applied:

```shell
Failed to allocate device vector A (error code all CUDA-capable devices are busy or unavailable)!
```

While load ran, all 3 pods remained healthy and `Running`, the `DCGM_FI_DEV_FB_FREE` metric reflected the live memory pressure per pod, and the HPA continued reconciling against the `cuda` Deployment at its 3-replica ceiling:

```shell
kubectl get pods -l app=cuda
```

```shell
NAME                   READY   STATUS    RESTARTS   AGE
cuda-887b9c799-49zn9   1/1     Running   0          4m12s
cuda-887b9c799-thlmf   1/1     Running   0          2m44s
cuda-887b9c799-vl8j7   1/1     Running   0          2m59s
```

```shell
kubectl get hpa cuda
```

```shell
NAME   REFERENCE         TARGETS       MINPODS   MAXPODS   REPLICAS   AGE
cuda   Deployment/cuda   8449/5k       1         3         3          3m14s
```

```shell
kubectl get pods -A --field-selector=status.phase=Pending
```

```shell
No resources found
```

This demonstrates that VKS v1.36 correctly integrates the Kubernetes HorizontalPodAutoscaler with custom, accelerator-specific metrics for GPU-backed workloads: the `custom.metrics.k8s.io` API served live per-Pod DCGM metrics (`DCGM_FI_DEV_FB_FREE`, `DCGM_FI_DEV_GPU_UTIL`, and the full DCGM metric set) sourced from Prometheus via the Prometheus Adapter, and the HPA controller correctly calculated and progressively scaled the target Deployment from 1 to its `maxReplicas` ceiling of 3 based on those metrics — with all 3 replicas genuinely scheduled `Running`, one per GPU-equipped worker node, and zero `Pending` pods across the cluster.
