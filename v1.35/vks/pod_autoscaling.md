## Description

If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

## Evidence

Complete guide to configure Horizontal Pod Autoscaling using NVIDIA DCGM metrics with Prometheus Operator on VKS v3.6.0 and Kubernetes v1.35.0.

### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

* Install Helm

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Install Prometheus Operator

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts && \
helm repo update
```

```shell
helm install prometheus-community/kube-prometheus-stack \
   --create-namespace --namespace prometheus \
   --generate-name \
   --set prometheus.service.type=LoadBalancer \
   --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
```

You may have to configure the prometheus service to be NodePort or LoadBalancer to view the dashboard and modify the `additionalScrapeConfigs` configMap as per [NVIDIA's documentation](https://docs.nvidia.com/datacenter/cloud-native/gpu-telemetry/latest/kube-prometheus.html).

In this guide, prometheus service was configured to by of type LoadBalancer to obtain an external IP for the server.

### Install Prometheus Adapter

```shell
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  --set prometheus.url="http://<kube-prometheus-stack-prometheus.default.svc.cluster.local or IP Address of Prometheus service>" --set prometheus.port="9090"
```


### Test Kubernetes Custom Metrics Endpoint

Check for the metrics exposed by DCGM Exporter such as `DCGM_FI_DEV_GPU_UTIL`.

```shell
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | jq -r . | grep DCGM
```

The output should be similar to the snippet below:

```shell
"namespaces/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"namespaces/DCGM_FI_DEV_SM_CLOCK"
"pods/DCGM_FI_DEV_MEM_CLOCK"
"pods/DCGM_FI_DEV_FB_USED"
"jobs.batch/DCGM_FI_DEV_SM_CLOCK"
"jobs.batch/DCGM_FI_DEV_MEMORY_TEMP"
"pods/DCGM_FI_DEV_DEC_UTIL"
"namespaces/DCGM_FI_DEV_DEC_UTIL"
"pods/DCGM_FI_DEV_XID_ERRORS"
"jobs.batch/DCGM_FI_DEV_FB_RESERVED"
"jobs.batch/DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"
"pods/DCGM_FI_DEV_FB_RESERVED"
"jobs.batch/DCGM_FI_DEV_FB_FREE"
"pods/DCGM_FI_DEV_MEM_COPY_UTIL"
"namespaces/DCGM_FI_DEV_FB_RESERVED"
"jobs.batch/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"namespaces/DCGM_FI_DEV_FB_FREE"
"namespaces/DCGM_FI_DEV_MEM_CLOCK"
"pods/DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"
"namespaces/DCGM_FI_DEV_XID_ERRORS"
"jobs.batch/DCGM_FI_DEV_GPU_UTIL"
"pods/DCGM_FI_DEV_MEMORY_TEMP"
"namespaces/DCGM_FI_DEV_MEMORY_TEMP"
"pods/DCGM_FI_DEV_GPU_UTIL"
"jobs.batch/DCGM_FI_DEV_FB_USED"
"namespaces/DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"
"namespaces/DCGM_FI_DEV_ENC_UTIL"
"pods/DCGM_FI_DEV_VGPU_LICENSE_STATUS"
"jobs.batch/DCGM_FI_DEV_MEM_CLOCK"
"jobs.batch/DCGM_FI_DEV_ENC_UTIL"
"jobs.batch/DCGM_FI_DEV_DEC_UTIL"
"namespaces/DCGM_FI_DEV_MEM_COPY_UTIL"
"namespaces/DCGM_FI_DEV_GPU_UTIL"
"pods/DCGM_FI_DEV_ENC_UTIL"
"namespaces/DCGM_FI_DEV_FB_USED"
"jobs.batch/DCGM_FI_DEV_XID_ERRORS"
"pods/DCGM_FI_DEV_FB_FREE"
"jobs.batch/DCGM_FI_DEV_MEM_COPY_UTIL"
"pods/DCGM_FI_DEV_SM_CLOCK"
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
          command: ["bash", "-c", "for (( c=1; c<=5000; c++ )); do ./vectorAdd; done"]
          resources:
            limits:
              nvidia.com/gpu: 1
EOF
```

Create HPA resource that utilises a custom metric such as `DCGM_FI_DEV_FB_FREE` to scale the deployment based on free GPU memory.

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
          averageValue: "10000"
EOF
```

Simulate load on GPU

```shell
kubectl exec -it deployment/cuda -- bash

for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done & \
for (( c=1; c<=5000; c++ )); do ./vectorAdd; done &
```

On inspecting hpa resource, you should be able to see the replicas being created for the deployment.

```shell
kubectl get hpa
```

The above command should display output as follows for successful Horizontal Pod Autoscaling.

```shell
NAME   REFERENCE         TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
cuda   Deployment/cuda   19762/10k   1         3         2          19m
```
