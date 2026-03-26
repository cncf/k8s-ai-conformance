# AI Service Metrics

This document demonstrates how to configure Prometheus to scrape GPU metrics from the DCGM exporter and verify that metrics are being collected.

## Prerequisites

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DRA (Dynamic Resource Allocation) driver support enabled is deployed through Blueprint as part of cluster provisioning. The GPU Operator will install all necessary components including the NVIDIA device plugin, container toolkit, and DRA driver required for this conformance test.

> **Note:** Creating a Rafay v1.35 MKS cluster with the default blueprint comes with Prometheus components preinstalled in the `rafay-infra` namespace. No additional Prometheus setup is required.

## Steps

### Step 1: Verify Prometheus components are running

Confirm that Prometheus server, adapter, alertmanager, and related components are running in the `rafay-infra` namespace.

```
$ kubectl get po -n rafay-infra

NAME                                                   READY   STATUS    RESTARTS      AGE
ebs-csi-controller-5595b45c7c-msc4q                    5/5     Running   5 (72m ago)   24h
ebs-csi-controller-5595b45c7c-slvnj                    5/5     Running   6 (71m ago)   24h
ebs-csi-node-r7h7k                                     3/3     Running   4 (71m ago)   24h
openebs-localpv-localpv-provisioner-5c45bcd7ff-d4lnk   1/1     Running   1 (72m ago)   24h
rafay-prometheus-adapter-c7c76bc8-68ktz                1/1     Running   1 (72m ago)   24h
rafay-prometheus-alertmanager-0                        2/2     Running   2 (72m ago)   24h
rafay-prometheus-helm-exporter-5677684675-fbxjk        1/1     Running   1 (72m ago)   24h
rafay-prometheus-kube-state-metrics-78cb655786-fkf5n   1/1     Running   1 (72m ago)   24h
rafay-prometheus-metrics-server-85f5657dcb-2gthz       1/1     Running   1 (72m ago)   24h
rafay-prometheus-node-exporter-bkpbp                   1/1     Running   1 (72m ago)   24h
rafay-prometheus-server-0                              2/2     Running   0             2m37s
snapshot-controller-7947694867-29kk8                   1/1     Running   1 (72m ago)   24h
snapshot-controller-7947694867-cxmmq                   1/1     Running   1 (72m ago)   24h

```

### Step 2: Add a scrape config for the DCGM exporter

Add the following scrape configuration to the Prometheus server so it scrapes GPU metrics from the DCGM exporter service.

```
scrape_configs:
  - job_name: dcgm-exporter
    static_configs:
    - targets: ['nvidia-dcgm-exporter.ai-conformance.svc:9400']
```

### Step 3: Verify the DCGM exporter target is active

Port-forward the Prometheus server and query the targets API to confirm the DCGM exporter is being scraped successfully. The `health` field should show `up`.

```
curl -k http://localhost:9090/api/v1/targets | grep dcgm | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0Handling connection for 9090
100   697  100   697    0     0   197k      0 --:--:-- --:--:-- --:--:--  226k
{
  "status": "success",
  "data": {
    "activeTargets": [
      {
        "discoveredLabels": {
          "__address__": "nvidia-dcgm-exporter.ai-conformance.svc:9400",
          "__metrics_path__": "/metrics",
          "__scheme__": "http",
          "__scrape_interval__": "1m",
          "__scrape_timeout__": "10s",
          "job": "dcgm-exporter"
        },
        "labels": {
          "instance": "nvidia-dcgm-exporter.ai-conformance.svc:9400",
          "job": "dcgm-exporter"
        },
        "scrapePool": "dcgm-exporter",
        "scrapeUrl": "http://nvidia-dcgm-exporter.ai-conformance.svc:9400/metrics",
        "globalUrl": "http://nvidia-dcgm-exporter.ai-conformance.svc:9400/metrics",
        "lastError": "",
        "lastScrape": "2026-03-25T19:37:54.898031524Z",
        "lastScrapeDuration": 0.005745036,
        "health": "up",
        "scrapeInterval": "1m",
        "scrapeTimeout": "10s"
      }
    ],
    "droppedTargets": []
  }
}
```

### Step 4: Query a GPU metric from Prometheus

Query a specific DCGM metric (e.g., `DCGM_FI_DEV_PCIE_REPLAY_COUNTER`) to confirm GPU metrics are being collected and returned by Prometheus.

```
curl -k http://localhost:9090/api/v1/query?query=DCGM_FI_DEV_PCIE_REPLAY_COUNTER | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0Handling connection for 9090
100   429  100   429    0     0   112k      0 --:--:-- --:--:-- --:--:--  139k
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "DCGM_FI_DRIVER_VERSION": "580.105.08",
          "Hostname": "ip-10-200-3-41",
          "UUID": "GPU-f4b6eae2-9160-3e9b-d87b-ffb86f9c8eef",
          "__name__": "DCGM_FI_DEV_PCIE_REPLAY_COUNTER",
          "device": "nvidia0",
          "gpu": "0",
          "instance": "nvidia-dcgm-exporter.ai-conformance.svc:9400",
          "job": "dcgm-exporter",
          "modelName": "Tesla T4",
          "pci_bus_id": "00000000:00:1E.0"
        },
        "value": [
          1774467530.676,
          "0"
        ]
      }
    ]
  }
}
```


