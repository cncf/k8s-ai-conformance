# Accelerator Metrics Validation

## Overview

Accelerator metrics from NVIDIA GPUs are exposed using the NVIDIA DCGM Exporter, which is deployed as part of the GPU Operator. The exporter collects GPU telemetry via the Data Center GPU Manager (DCGM) and exposes it as Prometheus-compatible metrics. This guide describes the validation of GPU metrics on Viettel AI Platform.

## Step 1: Verify GPU Operator and DCGM Exporter Pods

```bash
$ kubectl get pods -n gpu-operator
NAME                                                          READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-kpxlb                                   1/1     Running     0          28h
gpu-feature-discovery-rw4hg                                   1/1     Running     0          5h8m
gpu-operator-55f7fcd79-qhwss                                  1/1     Running     1          28h
gpu-operator-node-feature-discovery-gc-585b876f9c-btngm       1/1     Running     1          28h
gpu-operator-node-feature-discovery-master-7f6684fb45-9jgvm   1/1     Running     1          28h
gpu-operator-node-feature-discovery-worker-2vmxj              1/1     Running     1          28h
gpu-operator-node-feature-discovery-worker-qq82n              1/1     Running     0          5h8m
nvidia-container-toolkit-daemonset-q7jhl                      1/1     Running     0          5h8m
nvidia-container-toolkit-daemonset-rgqwj                      1/1     Running     0          28h
nvidia-cuda-validator-7hghr                                   0/1     Completed   0          27h
nvidia-cuda-validator-x72v5                                   0/1     Completed   0          5h7m
nvidia-dcgm-exporter-cn9dw                                    1/1     Running     0          28h
nvidia-dcgm-exporter-n7jh7                                    1/1     Running     0          5h8m
nvidia-device-plugin-daemonset-mz74d                          1/1     Running     0          5h8m
nvidia-device-plugin-daemonset-wcwpp                          1/1     Running     0          28h
nvidia-operator-validator-9j4p2                               1/1     Running     0          5h7m
nvidia-operator-validator-fttz6                               1/1     Running     0          28h
```

## Step 2: Verify DCGM Exporter Service

```bash
$ kubectl get svc -n gpu-operator nvidia-dcgm-exporter
NAME                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
nvidia-dcgm-exporter   ClusterIP   10.107.142.195   <none>        9400/TCP   28h
```

## Step 3: Query DCGM Metrics

```bash
$ kubectl run dcgm-curl --image=curlimages/curl --restart=Never --rm -i --command -- \
    curl -s http://10.107.142.195:9400/metrics | \
    grep -E '^DCGM_FI_DEV_(GPU_UTIL|FB_USED|FB_FREE|GPU_TEMP|POWER_USAGE|MEM_COPY_UTIL|SM_CLOCK|MEM_CLOCK)|^DCGM_FI_PROF'
```

**Output:**
```
DCGM_FI_DEV_SM_CLOCK{gpu="0",UUID="GPU-a8affe7e-7930-ae20-3fc6-4b0edf3ebbf3",pci_bus_id="00000000:0D:00.0",device="nvidia0",modelName="NVIDIA L40S",Hostname="ubuntu-sv16",DCGM_FI_DRIVER_VERSION="580.126.09"} 210
DCGM_FI_DEV_SM_CLOCK{gpu="1",UUID="GPU-72ad4d20-cfed-f08d-c894-f620962ca3e4",pci_bus_id="00000000:B4:00.0",device="nvidia1",modelName="NVIDIA L40S",Hostname="ubuntu-sv16",DCGM_FI_DRIVER_VERSION="580.126.09"} 210
DCGM_FI_DEV_MEM_CLOCK{gpu="0",...,modelName="NVIDIA L40S",...} 405
DCGM_FI_DEV_MEM_CLOCK{gpu="1",...,modelName="NVIDIA L40S",...} 405
DCGM_FI_DEV_GPU_TEMP{gpu="0",...,modelName="NVIDIA L40S",...} 32
DCGM_FI_DEV_GPU_TEMP{gpu="1",...,modelName="NVIDIA L40S",...} 32
DCGM_FI_DEV_POWER_USAGE{gpu="0",...,modelName="NVIDIA L40S",...} 34.156000
DCGM_FI_DEV_POWER_USAGE{gpu="1",...,modelName="NVIDIA L40S",...} 34.164000
DCGM_FI_DEV_GPU_UTIL{gpu="0",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_DEV_GPU_UTIL{gpu="1",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="0",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="1",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_DEV_FB_FREE{gpu="0",...,modelName="NVIDIA L40S",...} 45457
DCGM_FI_DEV_FB_FREE{gpu="1",...,modelName="NVIDIA L40S",...} 45457
DCGM_FI_DEV_FB_USED{gpu="0",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_DEV_FB_USED{gpu="1",...,modelName="NVIDIA L40S",...} 0
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="0",...,modelName="NVIDIA L40S",...} 0.000000
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="1",...,modelName="NVIDIA L40S",...} 0.000000
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="0",...,modelName="NVIDIA L40S",...} 0.000000
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="1",...,modelName="NVIDIA L40S",...} 0.000000
```

The DCGM Exporter reports metrics for 2x NVIDIA L40S GPUs on `ubuntu-sv16` (driver version `580.126.09`).

## Metrics Summary

| Metric                            | Description                          | Value (idle) |
| --------------------------------- | ------------------------------------ | ------------ |
| `DCGM_FI_DEV_SM_CLOCK`            | Streaming Multiprocessor clock (MHz) | 210          |
| `DCGM_FI_DEV_MEM_CLOCK`           | Memory clock (MHz)                   | 405          |
| `DCGM_FI_DEV_GPU_TEMP`            | GPU temperature (°C)                 | 32           |
| `DCGM_FI_DEV_POWER_USAGE`         | Power draw (Watts)                   | 34.15        |
| `DCGM_FI_DEV_GPU_UTIL`            | GPU utilization (%)                  | 0            |
| `DCGM_FI_DEV_MEM_COPY_UTIL`       | Memory copy utilization (%)          | 0            |
| `DCGM_FI_DEV_FB_FREE`             | Framebuffer free (MiB)               | 45457        |
| `DCGM_FI_DEV_FB_USED`             | Framebuffer used (MiB)               | 0            |
| `DCGM_FI_PROF_GR_ENGINE_ACTIVE`   | Graphics engine active ratio         | 0.0          |
| `DCGM_FI_PROF_PIPE_TENSOR_ACTIVE` | Tensor core active ratio             | 0.0          |
