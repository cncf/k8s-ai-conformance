## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator

Further labeling node is required for nvdp/nvidia-device-plugin chart
```bash
kubectl label node <GPU-NODE> gpu=NVIDIA-L40S feature.node.kubernetes.io/pci-10de.present=true
```

time-slicing-config.yaml
```yaml
version: v1
sharing:
  timeSlicing:
    renameByDefault: true
    resources:
      - name: nvidia.com/gpu
        replicas: 10
```

Enable time slicing with provided config. If you are running device plugin from chart `nvidia/nvidia-device-plugin`, delete it first.

```bash
helm upgrade -i nvdp nvdp/nvidia-device-plugin  \
  --namespace nvidia-device-plugin --create-namespace \
  --set-file config.map.config=time-slicing-config.yaml \
  --set nodeSelector.gpu="NVIDIA-L40S" \
  --version 0.17.1
```
