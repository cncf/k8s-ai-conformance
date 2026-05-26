## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator


## Install NVIDIA gpu operator
```bash
helm install --wait gpu-operator \
    -n gpu-operator --create-namespace \
    nvidia/gpu-operator \
    --set driver.enabled=false \
    --set devicePlugin.enabled=false \
    --version=v26.3.1
```
And then install DRA driver
```bash
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu \
    --version="25.12.0" \
    --create-namespace \
    --namespace nvidia-dra-driver-gpu \
    --set gpuResourcesEnabledOverride=true
```

```bash
kubectl apply -f - <<EOF
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: gpu-claim-template
spec:
  spec:
    devices:
      requests:
      - name: single-gpu
        exactly:
          deviceClassName: gpu.nvidia.com
          allocationMode: ExactCount
          count: 1
EOF
```

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dra-example
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dra-example
  template:
    metadata:
      labels:
        app: dra-example
    spec:
      containers:
        - name: cuda
          image: nvidia/cuda:12.4.1-base-ubuntu22.04
          command: ["bash", "-lc", "nvidia-smi"]
          resources:
            claims:
            - name: single-gpu
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: gpu-claim-template
EOF
```
Check pod logs to verify that gpu is available
```bash
kubectl logs deployment/dra-example
Wed Apr 29 15:19:22 2026
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.58.03              Driver Version: 595.58.03      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA L40S                    On  |   00000000:00:07.0 Off |                    0 |
| N/A   33C    P0             59W /  350W |       0MiB /  46068MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|  No running processes found                                                             |
+-----------------------------------------------------------------------------------------+
```

Then try without gpu resource defined:

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dra-example
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dra-example
  template:
    metadata:
      labels:
        app: dra-example
    spec:
      containers:
        - name: cuda
          image: nvidia/cuda:12.4.1-base-ubuntu22.04
          command: ["bash", "-lc", "nvidia-smi"]
          #resources:
          #  claims:
          #  - name: single-gpu
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: gpu-claim-template
EOF
```

Verify that gpu is not available
```bash
kubectl logs deployment/dra-example
bash: line 1: nvidia-smi: command not found
```
