# Evidence: Accelerator Metrics

---

## Overview

This test validates that the Kubernetes platform exposes GPU metrics including utilization, memory usage, temperature, and power draw via NVIDIA DCGM (Data Center GPU Manager) Exporter. These metrics are exposed in Prometheus format for monitoring and alerting purposes.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The NVIDIA GPU Operator with DCGM Exporter is deployed through Blueprint as part of cluster provisioning. The GPU Operator will install all necessary components including DCGM for GPU metrics collection.

---

## Step 1: Verify DCGM Exporter Installation

First, verify that the NVIDIA DCGM Exporter is installed and running as part of the GPU Operator:

```
$ kubectl get pods -n ai-conformance | grep dcgm

nvidia-dcgm-4zbs7                                                 1/1     Running     1 (45m ago)   48m
nvidia-dcgm-exporter-bhqpl                                        1/1     Running     1 (45m ago)   48m
```

> **Note:** ✅ DCGM and DCGM Exporter pods are running in the ai-conformance namespace.

---

## Step 2: Access DCGM Metrics Endpoint

Access the DCGM Exporter pod to retrieve GPU metrics:

```
$ kubectl exec -it -n ai-conformance nvidia-dcgm-4zbs7 -- bash
```

---

## Step 3: Verify Core GPU Metrics (REQUIRED)

The following core metrics are required for conformance:

### GPU Utilization (REQUIRED)

```
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
```

> **Note:** ✅ GPU Utilization metric (DCGM_FI_DEV_GPU_UTIL) is available.

### GPU Memory Usage (REQUIRED)

```
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 14912
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 447
```

> **Note:** ✅ GPU Memory metrics (DCGM_FI_DEV_FB_FREE, DCGM_FI_DEV_FB_USED, DCGM_FI_DEV_FB_RESERVED) are available.

---

## Step 4: Verify Additional GPU Metrics (OPTIONAL)

The following additional metrics are exposed when available from the hardware:

### GPU Temperature

```
# HELP DCGM_FI_DEV_GPU_TEMP GPU temperature (in C).
# TYPE DCGM_FI_DEV_GPU_TEMP gauge
DCGM_FI_DEV_GPU_TEMP{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 24
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
```

> **Note:** ✅ GPU Temperature metrics (DCGM_FI_DEV_GPU_TEMP, DCGM_FI_DEV_MEMORY_TEMP) are available.

### GPU Power Draw

```
# HELP DCGM_FI_DEV_POWER_USAGE Power draw (in W).
# TYPE DCGM_FI_DEV_POWER_USAGE gauge
DCGM_FI_DEV_POWER_USAGE{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 14.512000
# HELP DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Total energy consumption since boot (in mJ).
# TYPE DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION counter
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 48287135
```

> **Note:** ✅ GPU Power metrics (DCGM_FI_DEV_POWER_USAGE, DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION) are available.

### Clock Frequencies

```
# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 300
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 405
```

> **Note:** ✅ Clock frequency metrics (DCGM_FI_DEV_SM_CLOCK, DCGM_FI_DEV_MEM_CLOCK) are available.

### PCIe Bandwidth

```
# HELP DCGM_FI_PROF_PCIE_TX_BYTES The rate of data transmitted over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_TX_BYTES gauge
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 29285
# HELP DCGM_FI_PROF_PCIE_RX_BYTES The rate of data received over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_RX_BYTES gauge
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 30378
```

> **Note:** ✅ PCIe bandwidth metrics (DCGM_FI_PROF_PCIE_TX_BYTES, DCGM_FI_PROF_PCIE_RX_BYTES) are available.

### Encoder/Decoder Utilization

```
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
```

> **Note:** ✅ Encoder/Decoder utilization metrics are available.

### Tensor Core and DRAM Activity (Profiling Metrics)

```
# HELP DCGM_FI_PROF_GR_ENGINE_ACTIVE Ratio of time the graphics engine is active.
# TYPE DCGM_FI_PROF_GR_ENGINE_ACTIVE gauge
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Ratio of cycles the tensor (HMMA) pipe is active.
# TYPE DCGM_FI_PROF_PIPE_TENSOR_ACTIVE gauge
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_DRAM_ACTIVE Ratio of cycles the device memory interface is active sending or receiving data.
# TYPE DCGM_FI_PROF_DRAM_ACTIVE gauge
DCGM_FI_PROF_DRAM_ACTIVE{gpu="0",UUID="GPU-d7cbb0cd-d4bd-c2f7-dfaf-221df5e958ea",pci_bus_id="00000000:00:1E.0",device="nvidia0",modelName="Tesla T4",Hostname="ip-10-200-3-41",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000001
```

> **Note:** ✅ Profiling metrics for tensor core and DRAM activity are available.

---

**Summary:** The platform exposes comprehensive GPU metrics via NVIDIA DCGM Exporter in Prometheus format. All required metrics (GPU utilization and memory usage) are available, along with additional metrics for temperature, power draw, clock frequencies, and PCIe bandwidth. These metrics enable effective monitoring and observability of GPU workloads.
