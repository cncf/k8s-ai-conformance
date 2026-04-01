# AI Service Metrics Validation

## Overview

This guide validates that Prometheus-format metrics exposed by AI services (KServe, Envoy AI Gateway, NVIDIA DCGM Exporter) can be discovered and collected on Viettel AI Platform. The platform uses Grafana Alloy as the metrics collector and Loki as the storage backend.

## Step 1: Verify Monitoring Stack

```bash
$ kubectl get pods -n loki
NAME                        READY   STATUS    RESTARTS      AGE
alloy-7754987576-zk4nk      2/2     Running   2 (27h ago)   35h
loki-0                      2/2     Running   2 (27h ago)   35h
loki-canary-bqngz           1/1     Running   0             5h9m
loki-canary-ttf67           1/1     Running   1 (27h ago)   35h
```

Grafana Alloy (OpenTelemetry Collector) is running with 2 replicas including the main collector container.

## Step 2: Verify AI Inference Service Metrics Endpoints

```bash
$ kubectl get pods -n kserve
NAME                                         READY   STATUS    RESTARTS      AGE
kserve-controller-manager-786fcddc9-hrcv5    2/2     Running   2 (25h ago)   33h
llmisvc-controller-manager-d7bb56889-27g9m   1/1     Running   1 (25h ago)   33h

$ kubectl get pods -n envoy-ai-gateway-system
NAME                                     READY   STATUS    RESTARTS      AGE
ai-gateway-controller-6495f95969-nsdth   1/1     Running   1 (27h ago)   36h
```

## Step 3: Verify DCGM Metrics Are Accessible

```bash
$ kubectl get svc -n gpu-operator nvidia-dcgm-exporter
NAME                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
nvidia-dcgm-exporter   ClusterIP   10.107.142.195   <none>        9400/TCP   28h

$ kubectl run metrics-test --image=curlimages/curl --restart=Never --rm -i --command -- \
    curl -s http://10.107.142.195:9400/metrics | grep "^DCGM_FI_DEV_GPU_UTIL" | head -2
DCGM_FI_DEV_GPU_UTIL{gpu="0",...,modelName="NVIDIA L40S",Hostname="ubuntu-sv16",...} 0
DCGM_FI_DEV_GPU_UTIL{gpu="1",...,modelName="NVIDIA L40S",Hostname="ubuntu-sv16",...} 0
```

## Step 4: Verify Prometheus Stack (Custom Metrics for HPA)

Prometheus is deployed in the `monitoring` namespace, scraping DCGM Exporter and exposing GPU metrics via the Kubernetes Custom Metrics API through Prometheus Adapter:

```bash
$ kubectl get pods -n monitoring | grep -E 'NAME|server|adapter'
NAME                                   READY   STATUS    RESTARTS   AGE
prometheus-adapter-bdf748cc-vtfdx      1/1     Running   0          5m45s
prometheus-server-d766bd656-gv7wc      2/2     Running   0          7m21s
```

This enables GPU utilization (`DCGM_FI_DEV_GPU_UTIL`) to drive HorizontalPodAutoscaler decisions for AI inference workloads.

## Metrics Sources

| Component            | Endpoint                                         | Format         | Collector                  |
| -------------------- | ------------------------------------------------ | -------------- | -------------------------- |
| NVIDIA DCGM Exporter | `nvidia-dcgm-exporter.gpu-operator:9400/metrics` | Prometheus     | Grafana Alloy + Prometheus |
| KServe controller    | `/metrics` on port 8080                          | Prometheus     | Grafana Alloy              |
| Envoy AI Gateway     | Envoy admin stats `/stats/prometheus`            | Prometheus     | Grafana Alloy              |
| Prometheus Adapter   | `custom.metrics.k8s.io/v1beta1`                  | Kubernetes API | HPA                        |
