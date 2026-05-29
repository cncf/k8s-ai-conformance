

# Observability: Accelerator Metrics


## Verify NVIDIA Metrics
```bash
> kubectl get service -n gpu-operator

```

```

NAME                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
gpu-operator           ClusterIP   100.69.22.93    <none>        8080/TCP   7m11s
nvidia-dcgm-exporter   ClusterIP   100.66.246.82   <none>        9400/TCP   7m8s

```



Creating test namespace "nvidia-metrics-1774009423"
```bash
> kubectl run scrape-accelerator-metrics-1 -n nvidia-metrics-1774009423 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS http://nvidia-dcgm-exporter.gpu-operator.svc.cluster.local:9400/metrics

```

```

pod/scrape-accelerator-metrics-1 created

```

```bash
> kubectl wait -n nvidia-metrics-1774009423 pod/scrape-accelerator-metrics-1 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/scrape-accelerator-metrics-1 condition met

```

```bash
> kubectl logs -n nvidia-metrics-1774009423 scrape-accelerator-metrics-1

```

```

# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 210
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 405
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_GPU_TEMP GPU temperature (in C).
# TYPE DCGM_FI_DEV_GPU_TEMP gauge
DCGM_FI_DEV_GPU_TEMP{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 31
# HELP DCGM_FI_DEV_POWER_USAGE Power draw (in W).
# TYPE DCGM_FI_DEV_POWER_USAGE gauge
DCGM_FI_DEV_POWER_USAGE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 16.365000
# HELP DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Total energy consumption since boot (in mJ).
# TYPE DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION counter
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 4742776
# HELP DCGM_FI_DEV_PCIE_REPLAY_COUNTER Total number of PCIe retries.
# TYPE DCGM_FI_DEV_PCIE_REPLAY_COUNTER counter
DCGM_FI_DEV_PCIE_REPLAY_COUNTER{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 22563
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 470
# HELP DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS Number of remapped rows for uncorrectable errors
# TYPE DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS Number of remapped rows for correctable errors
# TYPE DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_ROW_REMAP_FAILURE Whether remapping of rows has failed
# TYPE DCGM_FI_DEV_ROW_REMAP_FAILURE gauge
DCGM_FI_DEV_ROW_REMAP_FAILURE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_PROF_GR_ENGINE_ACTIVE Ratio of time the graphics engine is active.
# TYPE DCGM_FI_PROF_GR_ENGINE_ACTIVE gauge
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Ratio of cycles the tensor (HMMA) pipe is active.
# TYPE DCGM_FI_PROF_PIPE_TENSOR_ACTIVE gauge
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_DRAM_ACTIVE Ratio of cycles the device memory interface is active sending or receiving data.
# TYPE DCGM_FI_PROF_DRAM_ACTIVE gauge
DCGM_FI_PROF_DRAM_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_PCIE_TX_BYTES The rate of data transmitted over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_TX_BYTES gauge
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 112603
# HELP DCGM_FI_PROF_PCIE_RX_BYTES The rate of data received over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_RX_BYTES gauge
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 17544

```



All required metrics found on attempt 1


Received metrics:
# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 210
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 405
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_GPU_TEMP GPU temperature (in C).
# TYPE DCGM_FI_DEV_GPU_TEMP gauge
DCGM_FI_DEV_GPU_TEMP{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 31
# HELP DCGM_FI_DEV_POWER_USAGE Power draw (in W).
# TYPE DCGM_FI_DEV_POWER_USAGE gauge
DCGM_FI_DEV_POWER_USAGE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 16.365000
# HELP DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Total energy consumption since boot (in mJ).
# TYPE DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION counter
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 4742776
# HELP DCGM_FI_DEV_PCIE_REPLAY_COUNTER Total number of PCIe retries.
# TYPE DCGM_FI_DEV_PCIE_REPLAY_COUNTER counter
DCGM_FI_DEV_PCIE_REPLAY_COUNTER{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 22563
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 470
# HELP DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS Number of remapped rows for uncorrectable errors
# TYPE DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS Number of remapped rows for correctable errors
# TYPE DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_ROW_REMAP_FAILURE Whether remapping of rows has failed
# TYPE DCGM_FI_DEV_ROW_REMAP_FAILURE gauge
DCGM_FI_DEV_ROW_REMAP_FAILURE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0
# HELP DCGM_FI_PROF_GR_ENGINE_ACTIVE Ratio of time the graphics engine is active.
# TYPE DCGM_FI_PROF_GR_ENGINE_ACTIVE gauge
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Ratio of cycles the tensor (HMMA) pipe is active.
# TYPE DCGM_FI_PROF_PIPE_TENSOR_ACTIVE gauge
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_DRAM_ACTIVE Ratio of cycles the device memory interface is active sending or receiving data.
# TYPE DCGM_FI_PROF_DRAM_ACTIVE gauge
DCGM_FI_PROF_DRAM_ACTIVE{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 0.000000
# HELP DCGM_FI_PROF_PCIE_TX_BYTES The rate of data transmitted over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_TX_BYTES gauge
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 112603
# HELP DCGM_FI_PROF_PCIE_RX_BYTES The rate of data received over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_RX_BYTES gauge
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="0",UUID="GPU-c2f76b33-7953-256e-8f33-1343419b1a1c",pci_bus_id="00000000:31:00.0",device="nvidia0",modelName="NVIDIA L4",Hostname="i-030791dc5a37cd257",DCGM_FI_DRIVER_VERSION="580.105.08"} 17544



Found expected metric: DCGM_FI_DEV_GPU_TEMP


Found expected metric: DCGM_FI_DEV_POWER_USAGE


Found expected metric: DCGM_FI_DEV_GPU_UTIL


Found expected metric: DCGM_FI_DEV_FB_USED


Deleting test namespace "nvidia-metrics-1774009423"


Namespace deletion took 6s
