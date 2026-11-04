## Description

For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

## Evidence

### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

NVIDIA GPU Operator includes [DCGM Exporter](https://github.com/NVIDIA/dcgm-exporter?tab=readme-ov-file#quickstart-on-kubernetes) which exposes metrics from the accelerator hardware.

Verify that GPU Operator is running.

```shell
kubectl get pods -n gpu-operator
```

```shell
NAME                                                          READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-22ml4                                   1/1     Running     0          5d
gpu-feature-discovery-rd7d5                                   1/1     Running     0          5d
gpu-operator-d57f6db6b-9srzg                                  1/1     Running     0          5d
gpu-operator-node-feature-discovery-gc-5f9c65779c-85d2x       1/1     Running     0          5d
gpu-operator-node-feature-discovery-master-596d8bc6c8-x6hlc   1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-5kv94              1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-kprf9              1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-mr5xf              1/1     Running     0          5d
nvidia-container-toolkit-daemonset-5nxjc                      1/1     Running     0          5d
nvidia-container-toolkit-daemonset-mswb8                      1/1     Running     0          5d
nvidia-cuda-validator-j95hw                                   0/1     Completed   0          5d
nvidia-cuda-validator-xbzc4                                   0/1     Completed   0          5d
nvidia-dcgm-exporter-gbtzd                                    1/1     Running     0          5d
nvidia-dcgm-exporter-wlfht                                    1/1     Running     0          5d
nvidia-device-plugin-daemonset-77fgr                          1/1     Running     0          5d
nvidia-device-plugin-daemonset-glxm5                          1/1     Running     0          5d
nvidia-driver-daemonset-rq87s                                 1/1     Running     0          5d
nvidia-driver-daemonset-twmhg                                 1/1     Running     0          5d
nvidia-operator-validator-28zsq                               1/1     Running     0          5d
nvidia-operator-validator-mjq6n                               1/1     Running     0          5d
```

The DCGM Exporter image and its ClusterIP service (one endpoint per GPU node):

```shell
kubectl get ds -n gpu-operator nvidia-dcgm-exporter -o jsonpath='{.spec.template.spec.containers[0].image}'
```

```shell
nvcr.io/nvidia/k8s/dcgm-exporter:4.4.2-4.7.0-distroless
```

```shell
kubectl get svc,endpoints -n gpu-operator nvidia-dcgm-exporter
```

```shell
NAME                           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/nvidia-dcgm-exporter   ClusterIP   198.201.34.71   <none>        9400/TCP   5d

NAME                             ENDPOINTS                         AGE
endpoints/nvidia-dcgm-exporter   193.0.1.12:9400,193.0.2.15:9400   5d
```

Verify that DCGM Exporter is hosting metrics from the underlying GPU. A temporary pod was run inside the cluster (in a dedicated `metrics-test` namespace, with a restricted-PSA compliant security context) to query the `nvidia-dcgm-exporter` ClusterIP service on port `9400`:

```shell
kubectl create namespace metrics-test

kubectl run -n metrics-test curl-test --restart=Never \
  --image=curlimages/curl:latest \
  --overrides='{"spec":{"securityContext":{"runAsNonRoot":true,"runAsUser":1000,"seccompProfile":{"type":"RuntimeDefault"}},"containers":[{"name":"curl-test","image":"curlimages/curl:latest","command":["sleep","3600"],"securityContext":{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]}}}]}}'

kubectl exec -n metrics-test curl-test -- curl -s http://nvidia-dcgm-exporter.gpu-operator.svc:9400/metrics
```

The metrics endpoint displayed the following output:

```shell
# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 1410
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 1215
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 3691
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 404
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
```

### Per-accelerator metrics from every GPU node

Each GPU node runs its own DCGM Exporter pod; scraping each endpoint directly shows that the core utilization and memory metrics are reported separately for every accelerator (distinct `UUID` / `Hostname` labels):

```shell
kubectl get pods -n gpu-operator -l app=nvidia-dcgm-exporter -o wide
```

```shell
NAME                         READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
nvidia-dcgm-exporter-gbtzd   1/1     Running   0          5d    193.0.2.15   conformance-test-np1-sms5g-pt7t9-p42p6   <none>           <none>
nvidia-dcgm-exporter-wlfht   1/1     Running   0          5d    193.0.1.12   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl exec -n metrics-test curl-test -- sh -c "curl -s http://193.0.2.15:9400/metrics | grep -E '^DCGM_FI_DEV_(GPU_UTIL|FB_USED|FB_FREE|MEM_COPY_UTIL)'"
```

```shell
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 3691
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-ba9adcc0-9809-47bb-85a6-b46a00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-p42p6",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
```

```shell
kubectl exec -n metrics-test curl-test -- sh -c "curl -s http://193.0.1.12:9400/metrics | grep -E '^DCGM_FI_DEV_(GPU_UTIL|FB_USED|FB_FREE|MEM_COPY_UTIL)'"
```

```shell
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-60e223c0-015b-436c-ba00-26da00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-vv2dd",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-60e223c0-015b-436c-ba00-26da00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-vv2dd",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-60e223c0-015b-436c-ba00-26da00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-vv2dd",DCGM_FI_DRIVER_VERSION="580.178.04"} 3691
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-60e223c0-015b-436c-ba00-26da00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-4C",Hostname="conformance-test-np1-sms5g-pt7t9-vv2dd",DCGM_FI_DRIVER_VERSION="580.178.04"} 0
```

The output confirms that DCGM Exporter exposes the core set of per-accelerator metrics required for conformance, including GPU utilization (`DCGM_FI_DEV_GPU_UTIL`), memory utilization (`DCGM_FI_DEV_MEM_COPY_UTIL`), framebuffer memory used/free/reserved (`DCGM_FI_DEV_FB_USED`, `DCGM_FI_DEV_FB_FREE`, `DCGM_FI_DEV_FB_RESERVED`), clock frequencies, temperature, and vGPU license status, via a standardized Prometheus-compatible `/metrics` endpoint, for each `GRID A100-4C` vGPU-backed device in the cluster. Power draw is reported as `N/A` by the vGPU guest (see `nvidia-smi` output in [virtualized_accelerator.md](virtualized_accelerator.md)) and is therefore not exposed. The same metrics are collected by Prometheus and consumed through the Kubernetes custom metrics API in [pod_autoscaling.md](pod_autoscaling.md).

### Cleanup

```shell
kubectl delete namespace metrics-test
```
