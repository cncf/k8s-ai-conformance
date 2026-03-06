# Accelerator Metrics Validation

## Overview

Accelerator metrics from NVIDIA GPUs can be exposed using the NVIDIA DCGM Exporter, which is deployed as part of the GPU Operator. The exporter collects GPU telemetry via the Data Center GPU Manager (DCGM) and exposes it as Prometheus-compatible metrics. This guide describes the installation of the Nvidia GPU operator and validation of metrics from the Nvidia DCGM Exporter on the Ericsson Cloud Container Distribution (ECCD) 2.34.0.

## Step 1: Install the required drivers and the GPU Operator

See Step 1 in DRA_plugin.md

## Step 2: Validate metrics from the Nvidia DCGM exporter

```bash
curl -s http://nvidia-dcgm-exporter.gpu-operator.svc.cluster.local:9400/metrics

DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 660
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 5000
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0
DCGM_FI_DEV_GPU_TEMP{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 74
DCGM_FI_DEV_POWER_USAGE{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 68.244000
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 15152422938
DCGM_FI_DEV_PCIE_REPLAY_COUNTER{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 91
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 73
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 9791
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 5120
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 447
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0.156205
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0.005822
DCGM_FI_PROF_DRAM_ACTIVE{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 0.092781
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 4791810
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="0",UUID="GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a",pci_bus_id="00000000:37:00.0",device="nvidia0",modelName="Tesla T4",Hostname="cp-bm14r7",DCGM_FI_DRIVER_VERSION="590.48.01"} 6935930
```