## Description

For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

## Evidence

### Prerequisites

* Provision a VKS Cluster with v1.36.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

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
gpu-feature-discovery-kc5zk                                   1/1     Running     0          4d
gpu-feature-discovery-rk8kw                                   1/1     Running     0          4d
gpu-operator-7d7969bcd5-zkfzd                                  1/1     Running     0          4d
gpu-operator-node-feature-discovery-gc-5cb78546ff-wdtgw       1/1     Running     0          4d
gpu-operator-node-feature-discovery-master-5ffff79b5-2qjmt    1/1     Running     0          4d
gpu-operator-node-feature-discovery-worker-2l65g              1/1     Running     0          4d
gpu-operator-node-feature-discovery-worker-g4zmm              1/1     Running     0          4d
gpu-operator-node-feature-discovery-worker-vnwgq              1/1     Running     0          4d
nvidia-container-toolkit-daemonset-r89mn                      1/1     Running     0          4d
nvidia-container-toolkit-daemonset-w89bd                      1/1     Running     0          4d
nvidia-cuda-validator-5slqv                                   0/1     Completed   0          4d
nvidia-cuda-validator-xzpf8                                   0/1     Completed   0          4d
nvidia-dcgm-exporter-fxvql                                    1/1     Running     0          4d
nvidia-dcgm-exporter-gk2v9                                    1/1     Running     0          4d
nvidia-device-plugin-daemonset-bbvvq                          1/1     Running     0          4d
nvidia-device-plugin-daemonset-mth9d                          1/1     Running     0          4d
nvidia-driver-daemonset-86lcq                                 1/1     Running     0          4d
nvidia-driver-daemonset-qbhpr                                 1/1     Running     0          4d
nvidia-operator-validator-cfbqm                               1/1     Running     0          4d
nvidia-operator-validator-kwcf7                               1/1     Running     0          4d
```

Verify that DCGM Exporter is hosting metrics from the underlying GPU. A temporary pod was run inside the cluster to query the `nvidia-dcgm-exporter` ClusterIP service on port `9400`:

```shell
kubectl run -n gpu-operator curl-test --restart=Never \
  --image=curlimages/curl:latest --command -- sleep 3600

kubectl exec -n gpu-operator curl-test -- curl -s http://nvidia-dcgm-exporter:9400/metrics
```

The metrics endpoint displayed the following output:

```shell
# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 1410
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 1215
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 9323
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 916
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-44d2411e-1604-478f-aa96-5a3c00000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-10C",Hostname="conformance-test-np1-qcqc8-bpdtl-wv2gs",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
```

The output confirms that DCGM Exporter exposes the core set of per-accelerator metrics required for conformance, including GPU utilization (`DCGM_FI_DEV_GPU_UTIL`), memory utilization (`DCGM_FI_DEV_MEM_COPY_UTIL`), framebuffer memory used/free/reserved (`DCGM_FI_DEV_FB_USED`, `DCGM_FI_DEV_FB_FREE`, `DCGM_FI_DEV_FB_RESERVED`), clock frequencies, temperature, and vGPU license status, via a standardized Prometheus-compatible `/metrics` endpoint on the `GRID A100-10C` vGPU-backed device.
