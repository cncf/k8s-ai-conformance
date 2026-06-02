

# Driver and Runtime Management Verification


## Identifying GPU Nodes
```bash
> kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.nvidia\.com/gpu}{"\n"}{end}'

```

```

i-02fb36e88acb587c7	
i-030791dc5a37cd257	1
i-0308c4a041ac88ead	
i-034383b8fabf72012	1
i-0ab9100a694772f64

```



* Node i-030791dc5a37cd257 has 1 NVIDIA GPUs


* Node i-034383b8fabf72012 has 1 NVIDIA GPUs
&check; Found 2 GPU node(s): i-030791dc5a37cd257, i-034383b8fabf72012



## Verifying GPU Operator DaemonSet


* Found DaemonSet: gpu-feature-discovery (Ready: 2/2)


* Found DaemonSet: nvidia-container-toolkit-daemonset (Ready: 2/2)


* Found DaemonSet: nvidia-dcgm-exporter (Ready: 2/2)


* Found DaemonSet: nvidia-device-plugin-daemonset (Ready: 2/2)


* Found DaemonSet: nvidia-device-plugin-mps-control-daemon (Ready: 0/0)


* Found DaemonSet: nvidia-driver-daemonset (Ready: 2/2)


* Found DaemonSet: nvidia-gpu-operator-node-feature-discovery-worker (Ready: 5/5)


* Found DaemonSet: nvidia-mig-manager (Ready: 0/0)


* Found DaemonSet: nvidia-operator-validator (Ready: 2/2)
&check; Driver DaemonSet nvidia-driver-daemonset is healthy: 2/2 pods ready



## Verifying Driver Installation with Diagnostic Job


Creating test namespace "testacceleratorsdriverruntimemanagement-1774009319"


Applying manifest "testdata/driver-check.yaml" to namespace "testacceleratorsdriverruntimemanagement-1774009319"
```bash
> kubectl apply -n testacceleratorsdriverruntimemanagement-1774009319 -f testdata/driver-check.yaml

```

```

job.batch/driver-check created

```

```bash
> kubectl wait -n testacceleratorsdriverruntimemanagement-1774009319 --for=condition=Complete job.batch/driver-check --timeout=5m

```

```

job.batch/driver-check condition met

```

```bash
> kubectl logs --namespace testacceleratorsdriverruntimemanagement-1774009319 job/driver-check

```

```

=== NVIDIA Driver Check ===
NVIDIA-SMI version  : 580.105.08
NVML version        : 580.105
DRIVER version      : 580.105.08
CUDA Version        : 13.0

=== GPU Detection ===
GPU 0: NVIDIA L4 (UUID: GPU-c2f76b33-7953-256e-8f33-1343419b1a1c)

=== Driver Information ===
driver_version, name, uuid
580.105.08, NVIDIA L4, GPU-c2f76b33-7953-256e-8f33-1343419b1a1c

SUCCESS: NVIDIA driver and runtime verified

```

&check; NVIDIA driver and runtime successfully verified on GPU nodes



## Checking for DRA Driver/Runtime Version Exposure


* DRA DeviceClass 'gpu.nvidia.com' found


* Note: DRA-based driver version verification not yet implemented in this test


* Future enhancement: Query driver/runtime versions via DRA APIs
&check; Driver Runtime Management conformance test PASSED



Deleting test namespace "testacceleratorsdriverruntimemanagement-1774009319"


Namespace deletion took 6s
