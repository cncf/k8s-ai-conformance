## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator
* Have nvidia device plugin installed. guide: https://upcloud.com/docs/guides/gpu-workloads-managed-kubernetes/ 

# install CDGM exporter
```bash
helm repo add gpu-helm-charts https://nvidia.github.io/dcgm-exporter/helm-charts
helm repo update
helm install dcgm-exporter gpu-helm-charts/dcgm-exporter --set nodeSelector.gpu='NVIDIA-L40S'
```

### Run job to check that metrics work

```bash
kubectl apply -f - <<EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: dcgm-quick-check
spec:
  template:
    spec:
      containers:
      - name: curl
        image: curlimages/curl:latest
        args: ["-s", "http://dcgm-exporter:9400/metrics"]
      restartPolicy: Never
  backoffLimit: 0
EOF
```

### Check logs

```bash
$ kubectl logs job/dcgm-quick-check

# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 210
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
DCGM_FI_DEV_MEM_CLOCK{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 405
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
DCGM_FI_DEV_MEMORY_TEMP{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_GPU_TEMP GPU temperature (in C).
# TYPE DCGM_FI_DEV_GPU_TEMP gauge
DCGM_FI_DEV_GPU_TEMP{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 32
# HELP DCGM_FI_DEV_POWER_USAGE Power draw (in W).
# TYPE DCGM_FI_DEV_POWER_USAGE gauge
DCGM_FI_DEV_POWER_USAGE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 35.686000
# HELP DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Total energy consumption since boot (in mJ).
# TYPE DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION counter
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 5551899580
# HELP DCGM_FI_DEV_PCIE_REPLAY_COUNTER Total number of PCIe retries.
# TYPE DCGM_FI_DEV_PCIE_REPLAY_COUNTER counter
DCGM_FI_DEV_PCIE_REPLAY_COUNTER{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
DCGM_FI_DEV_GPU_UTIL{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
DCGM_FI_DEV_ENC_UTIL{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
DCGM_FI_DEV_DEC_UTIL{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
DCGM_FI_DEV_FB_FREE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 45459
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
DCGM_FI_DEV_FB_USED{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_FB_RESERVED Framebuffer memory reserved (in MiB).
# TYPE DCGM_FI_DEV_FB_RESERVED gauge
DCGM_FI_DEV_FB_RESERVED{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 608
# HELP DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS Number of remapped rows for uncorrectable errors
# TYPE DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS Number of remapped rows for correctable errors
# TYPE DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS counter
DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_ROW_REMAP_FAILURE Whether remapping of rows has failed
# TYPE DCGM_FI_DEV_ROW_REMAP_FAILURE gauge
DCGM_FI_DEV_ROW_REMAP_FAILURE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0
# HELP DCGM_FI_PROF_GR_ENGINE_ACTIVE Ratio of time the graphics engine is active.
# TYPE DCGM_FI_PROF_GR_ENGINE_ACTIVE gauge
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0.000000
# HELP DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Ratio of cycles the tensor (HMMA) pipe is active.
# TYPE DCGM_FI_PROF_PIPE_TENSOR_ACTIVE gauge
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0.000000
# HELP DCGM_FI_PROF_DRAM_ACTIVE Ratio of cycles the device memory interface is active sending or receiving data.
# TYPE DCGM_FI_PROF_DRAM_ACTIVE gauge
DCGM_FI_PROF_DRAM_ACTIVE{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 0.000131
# HELP DCGM_FI_PROF_PCIE_TX_BYTES The rate of data transmitted over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_TX_BYTES gauge
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 175084
# HELP DCGM_FI_PROF_PCIE_RX_BYTES The rate of data received over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_RX_BYTES gauge
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="0",UUID="GPU-5b4e07f6-3a86-4340-039d-25f2a83cdc43",pci_bus_id="00000000:00:07.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="gpu1-wl8sg",DCGM_FI_DRIVER_VERSION="595.58.03"} 19831
```