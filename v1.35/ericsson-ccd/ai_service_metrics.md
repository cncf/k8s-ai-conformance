# AI Service Metrics Validation

## Overview

This guide validates that the metrics exposed by the NVIDIA DCGM Exporter can be scraped and consumed by Victoria Metrics.
The monitoring solution in Ericsson Cloud Container Distribution 2.34.0 is Victoria Metrics.

## Step 1: Install the GPU Operator

Refer to Step 2 in DRA_plugin.md

## Step 2: Configure the Scrape Target for the DCGM Exporter in Victoria Metrics

Patch the vmagent ConfigMap to add the scrape target for the `gpu-operator` namespace, allowing scraping of the HTTP endpoint for the `nvidia-dcgm-exporter` service.

```bash
kubectl get cm eric-victoria-metrics-vmagent -n monitoring -o json | python3 -c "
import sys, json, yaml
cm = json.load(sys.stdin)
scrape = yaml.safe_load(cm['data']['scrape.yml'])
for job in scrape['scrape_configs']:
    if job['job_name'] == 'kubernetes-service-endpoints':
        if 'gpu-operator' not in job['kubernetes_sd_configs'][0]['namespaces']['names']:
            job['kubernetes_sd_configs'][0]['namespaces']['names'].append('gpu-operator')
        for rule in job['relabel_configs']:
            if rule.get('action') == 'keep' and rule.get('regex') == 'https;.*' and '__meta_kubernetes_pod_label_k8s_app' in rule.get('source_labels', []):
                rule['regex'] = 'https;.*|http;nvidia-dcgm-exporter'
                rule['source_labels'] = ['__scheme__', '__meta_kubernetes_pod_label_app']
        break
cm['data']['scrape.yml'] = yaml.dump(scrape, default_flow_style=False, width=1000)
for k in ['resourceVersion','uid','creationTimestamp','managedFields']: cm.get('metadata',{}).pop(k,None)
json.dump(cm, sys.stdout)
" | kubectl apply -f -
```


## Step 3: Restart vmagent to Reload the ConfigMap

```bash
kubectl -n monitoring rollout restart deploy eric-victoria-metrics-agent
```


## Step 4: Verify the DCGM Metrics Are Being Scraped

```bash

#Query VMSelect for the DCGM Metrics
curl -s https://eric-victoria-metrics-cluster-vmselect.monitoring.svc.cluster.local:8481/select/0/prometheus/api/v1/query?query=DCGM_FI_DEV_PCIE_REPLAY_COUNTER | jq .

{
  "status": "success",
  "isPartial": false,
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "DCGM_FI_DEV_PCIE_REPLAY_COUNTER",
          "DCGM_FI_DRIVER_VERSION": "590.48.01",
          "Hostname": "cp-bm14r7",
          "UUID": "GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",
          "app": "nvidia-dcgm-exporter",
          "device": "nvidia0",
          "gpu": "0",
          "instance": "192.168.139.179:9400",
          "job": "kubernetes-service-endpoints",
          "kubernetes_name": "nvidia-dcgm-exporter",
          "kubernetes_namespace": "gpu-operator",
          "modelName": "Tesla T4",
          "nodename": "cp-bm14r7",
          "pci_bus_id": "00000000:37:00.0"
        },
        "value": [
          1772536457,
          "0"
        ]
      }
    ]
  },
  "stats": {
    "seriesFetched": "1",
    "executionTimeMsec": 1
  }
}
```
