## Description

For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

## Evidence

### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

NVIDIA GPU Operator includes [DCGM Exporter](https://github.com/NVIDIA/dcgm-exporter?tab=readme-ov-file#quickstart-on-kubernetes) which exposes metrics from the accelerator hardware.

Verify that GPU Operator is running.

```shell
kubectl get pods -n gpu-operator
NAME                                                          READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-stsqc                                   1/1     Running     0          75m
gpu-operator-d68f69784-wltls                                  1/1     Running     0          75m
gpu-operator-node-feature-discovery-gc-797c77fd98-5d6rx       1/1     Running     0          75m
gpu-operator-node-feature-discovery-master-5c64b6c596-9d5db   1/1     Running     0          75m
gpu-operator-node-feature-discovery-worker-75s2f              1/1     Running     0          75m
gpu-operator-node-feature-discovery-worker-9jvh7              1/1     Running     0          75m
nvidia-container-toolkit-daemonset-nqzmh                      1/1     Running     0          75m
nvidia-cuda-validator-br89d                                   0/1     Completed   0          72m
nvidia-dcgm-exporter-bbn6x                                    1/1     Running     0          75m
nvidia-device-plugin-daemonset-n62lm                          1/1     Running     0          31m
nvidia-driver-daemonset-dkz47                                 1/1     Running     0          75m
nvidia-operator-validator-p7p59                               1/1     Running     0          75m
```

Verify that DCGM Exporter is hosting metrics from the underlying GPU.

The metrics endpoint should display an output as follows:

```shell
curl -L http://192.168.1.16:9400/metrics

# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 1410
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 1215
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 37739
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 3220
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-0946f70c-659b-4c39-8b78-7e7600000000",pci_bus_id="00000000:03:00.0",device="nvidia0",modelName="GRID A100-40C",Hostname="test-ai-1-node-pool-40c-d5ds8-br9r6-9tb2c",DCGM_FI_DRIVER_VERSION="580.126.09"} 1
* Connection #0 to host 192.168.1.16 left intact
```
