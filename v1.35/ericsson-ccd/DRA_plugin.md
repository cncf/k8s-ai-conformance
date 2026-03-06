# NVIDIA DRA Plugin Installation and Validation Guide


## Overview

Dynamic Resource Allocation (DRA) is a Kubernetes feature that provides a flexible mechanism for requesting and sharing resources such as GPUs and other hardware accelerators. Unlike traditional device plugins, DRA allows workloads to describe resource requirements through ResourceClaim objects, enabling more expressive allocation, sharing across pods, and vendor-specific configuration via custom drivers.

DRA is a stable feature in Kubernetes 1.35 and is enabled by default. This guide shows the installation of the Nvidia GPU operator and the Nvidia DRA plugin on the Ericsson Cloud Container Distribution (ECCD) 2.34.0. The guide also demonstrates the allocation of the GPU resources in a pod using the DRA plugin.


## Step 1: Install the required drivers

```bash

# Install kernel development packages
sudo zypper in kernel-default-devel

# Install driver
arch=x86_64
distro=sles15
sudo zypper addrepo https://developer.download.nvidia.com/compute/cuda/repos/$distro/$arch/cuda-$distro.repo
sudo SUSEConnect --product sle-module-desktop-applications/15.7/x86_64
sudo SUSEConnect --product sle-module-development-tools/15.7/x86_64
sudo SUSEConnect --product PackageHub/15.7/x86_64
sudo zypper refresh

sudo zypper install nvidia-open

# Install Nvidia Container Toolkit
sudo zypper ar https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo
sudo zypper modifyrepo --enable nvidia-container-toolkit-experimental
export NVIDIA_CONTAINER_TOOLKIT_VERSION=1.18.2-1

sudo zypper --gpg-auto-import-keys install -y \
    nvidia-container-toolkit-${NVIDIA_CONTAINER_TOOLKIT_VERSION} \
    nvidia-container-toolkit-base-${NVIDIA_CONTAINER_TOOLKIT_VERSION} \
    libnvidia-container-tools-${NVIDIA_CONTAINER_TOOLKIT_VERSION} \
    libnvidia-container1-${NVIDIA_CONTAINER_TOOLKIT_VERSION}

#Configure Runtime
sudo nvidia-ctk runtime configure --runtime=containerd

# reboot the node
sudo reboot
```

## Step 2: Install the GPU Operator

Install the GPU Operator helm chart

  ```bash
  #Label the node
  kubectl label node <node> nvidia.com/dra-kubelet-plugin=true
  #Add the helm repo
  helm repo add nvidia https://helm.ngc.nvidia.com/nvidia \
      && helm repo update

  #Install the GPU Operator with the Nvidia Device Plugin disabled
  helm upgrade --install gpu-operator nvidia/gpu-operator \
    --version=v25.10.1 \
    --create-namespace \
    --namespace gpu-operator \
    --set devicePlugin.enabled=false \
    --set driver.manager.env[0].name=NODE_LABEL_FOR_GPU_POD_EVICTION \
    --set driver.manager.env[0].value="nvidia.com/dra-kubelet-plugin" \
    --set driver.enabled=false \
    --set toolkit.enabled=false

  ```

Verify that the operator pods are running

  ```bash
  kubectl get po -n gpu-operator
  NAME                                                          READY   STATUS      RESTARTS        AGE
  gpu-feature-discovery-sqsjv                                   1/1     Running     0               3d23h
  gpu-operator-d68f69784-sxgmg                                  1/1     Running     1 (3d21h ago)   3d23h
  gpu-operator-node-feature-discovery-gc-797c77fd98-4nrcd       1/1     Running     1 (3d21h ago)   3d23h
  gpu-operator-node-feature-discovery-master-5c64b6c596-hf9bw   1/1     Running     1 (3d21h ago)   3d23h
  gpu-operator-node-feature-discovery-worker-xzhn7              1/1     Running     1 (3d21h ago)   3d23h
  nvidia-cuda-validator-kfsbq                                   0/1     Completed   0               3d21h
  nvidia-dcgm-exporter-cwwhr                                    1/1     Running     0               3d23h
  nvidia-operator-validator-vdc9s                               1/1     Running     0               3d23h
  ```


## Step 3: Install the Nvidia DRA Plugin

```bash

# install nvidia-dra-driver-gpu helm chart
helm upgrade -i nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu \
  --version="25.12.0" \
  --namespace nvidia-dra-driver-gpu \
  --create-namespace \
  --set gpuResourcesEnabledOverride=true
```

## Step 4: Create the Resource Claim Template and Test App

```bash
#Get the Image tag for httpd image
export IMAGE_TAG=$(curl -s -q https://registry.eccd.local:5000/v2/httpd/tags/list | jq '.tags[]' -r | grep -v 'sha256')

cat <<EOF | envsubst | kubectl apply -f -
---
apiVersion: v1
kind: Namespace
metadata:
  name: testing
---
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: single-gpu
  namespace: testing
spec:
  spec:
    devices:
      requests:
      - name: req-0
        exactly:
          deviceClassName: gpu.nvidia.com
          allocationMode: ExactCount
          count: 1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: testing
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: test-app
          image: registry.eccd.local:5000/httpd:${IMAGE_TAG}
          command: ["/bin/sh", "-c"]
          args:
            - |
              while true; do
                sleep 30
              done
          resources:
            claims:
              - name: single-gpu
      resourceClaims:
        - name: single-gpu
          resourceClaimTemplateName: single-gpu
EOF
```

Verify that the pod is running and the ResourceClaim is allocated

```bash
kubectl get po,resourceclaim -n testing
NAME                            READY   STATUS    RESTARTS   AGE
pod/test-app-647c6f456b-sdsh9   1/1     Running   0          2m50s

NAME                                                                       STATE                AGE
resourceclaim.resource.k8s.io/test-app-647c6f456b-sdsh9-single-gpu-tc77s   allocated,reserved   2m50s
```

## Step 5: Verify GPU Allocation

Run the nvidia-smi tool in the pod to confirm allocation

```bash
kubectl exec -n testing test-app-647c6f456b-sdsh9 -- nvidia-smi -L
GPU 0: Tesla T4 (UUID: GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a)
```
