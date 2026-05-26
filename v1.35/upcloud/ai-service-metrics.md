## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator
* Have nvidia device plugin installed. guide: https://upcloud.com/docs/guides/gpu-workloads-managed-kubernetes/ 

## Install prometheus

Create a values.yaml
```yaml
prometheus:
  prometheusSpec:
    serviceMonitorSelectorNilUsesHelmValues: false
    serviceMonitorSelector: {}
    serviceMonitorNamespaceSelector: {}
```
Install kube-prometheus-stack using said values

```bash
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  -f values.yaml
```

## Verify it's working

```bash
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: ai-service-metrics-check
spec:
  template:
    spec:
      containers:
      - name: curl
        image: curlimages/curl:latest
        args: ["-s", "http://prometheus-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090/api/v1/query?query=DCGM_FI_DEV_GPU_TEMP"]
      restartPolicy: Never
  backoffLimit: 0
EOF
```
Check logs:

```bash
kubectl logs job.batch/ai-service-metrics-check
{"status":"success","data":{"resultType":"vector","result":[{"metric":{"DCGM_FI_DRIVER_VERSION":"595.58.03","Hostname":"gpu1-wl8sg","UUID":"GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43","__name__":"DCGM_FI_DEV_GPU_TEMP","container":"exporter","device":"nvidia0","endpoint":"metrics","gpu":"0","instance":"192.168.1.253:9400","job":"dcgm-exporter","modelName":"NVIDIA L40S","namespace":"default","pci_bus_id":"00000000:00:07.0","pod":"dcgm-exporter-zxgzh","service":"dcgm-exporter"},"value":[1775721747.082,"32"]}]}}
```