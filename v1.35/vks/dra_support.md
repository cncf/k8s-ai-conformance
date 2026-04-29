## Description 

Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

## Evidence

Dynamic Resource Allocation provides flexible resource management for specialised hardware like GPUs, FPGAs and network-attached devices. This guide covers how to configure DRA with VKS 3.6.0.

### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Verify GPU Operator

Pods in the gpu-operator namespace should be in running or completed state.

```shell
kubectl get pods -n gpu-operator
```

```shell
NAME                                                          READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-dm9n4                                   1/1     Running     0          2m59s
gpu-feature-discovery-dqtqf                                   1/1     Running     0          3m6s
gpu-feature-discovery-gwskj                                   1/1     Running     0          3m6s
gpu-feature-discovery-wg8pk                                   1/1     Running     0          3m6s
gpu-operator-d68f69784-6fhzr                                  1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-gc-797c77fd98-lb5qz       1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-master-5c64b6c596-q2pk2   1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-worker-dl4rg              1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-worker-j7hkq              1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-worker-jm9w6              1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-worker-mbj6f              1/1     Running     0          3m32s
gpu-operator-node-feature-discovery-worker-n7sn7              1/1     Running     0          3m32s
nvidia-container-toolkit-daemonset-4m6p8                      1/1     Running     0          2m59s
nvidia-container-toolkit-daemonset-bxq56                      1/1     Running     0          3m6s
nvidia-container-toolkit-daemonset-q5w7p                      1/1     Running     0          3m6s
nvidia-container-toolkit-daemonset-trm75                      1/1     Running     0          3m6s
nvidia-cuda-validator-5htcn                                   0/1     Completed   0          9s
nvidia-cuda-validator-6cxgj                                   0/1     Completed   0          44s
nvidia-cuda-validator-7tw9w                                   0/1     Completed   0          26s
nvidia-cuda-validator-ffkh5                                   0/1     Completed   0          33s
nvidia-dcgm-exporter-gmxdq                                    1/1     Running     0          2m59s
nvidia-dcgm-exporter-jssg8                                    1/1     Running     0          3m5s
nvidia-dcgm-exporter-xd9rm                                    1/1     Running     0          3m6s
nvidia-dcgm-exporter-zhxww                                    1/1     Running     0          3m6s
nvidia-device-plugin-daemonset-5mv8c                          1/1     Running     0          3m6s
nvidia-device-plugin-daemonset-ljfkm                          1/1     Running     0          2m59s
nvidia-device-plugin-daemonset-m4trl                          1/1     Running     0          3m6s
nvidia-device-plugin-daemonset-ppmq5                          1/1     Running     0          3m5s
nvidia-driver-daemonset-82r8f                                 1/1     Running     0          3m19s
nvidia-driver-daemonset-dkb82                                 1/1     Running     0          3m19s
nvidia-driver-daemonset-t2rps                                 1/1     Running     0          3m19s
nvidia-driver-daemonset-txh24                                 1/1     Running     0          3m19s
nvidia-operator-validator-gjrzv                               1/1     Running     0          3m6s
nvidia-operator-validator-rwlhl                               1/1     Running     0          2m59s
nvidia-operator-validator-tqnm2                               1/1     Running     0          3m6s
nvidia-operator-validator-zf5dx                               1/1     Running     0          3m6s
```

Cluster nodes should carry labels with GPU information.

```shell
kubectl get node -o json | jq '.items[].metadata.labels'
```

The node labels demonstrate that VKS supports GPU driver and runtime lifecycle management through the NVIDIA GPU Operator. The operator automates the deployment and management of:

- **NVIDIA Driver**: Deployed via `nvidia-driver-daemonset`, with version information exposed in node labels (`nvidia.com/cuda.driver-version.full`, `nvidia.com/vgpu.host-driver-version`)
- **Container Runtime**: Deployed via `nvidia-container-toolkit-daemonset`, enabling GPU-accelerated containers
- **Verification Mechanism**: Node labels serve as the verification that drivers and runtime are properly installed and operational

The DRA ResourceSlice attributes (shown later in this document) also expose driver version information through the `driverVersion` attribute, providing additional verification of the driver lifecycle management.

```shell
{
  "nvidia.com/cuda.driver-version.full": "580.126.09",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "126",
  "nvidia.com/cuda.driver-version.revision": "09",
  "nvidia.com/cuda.driver.major": "580",
  "nvidia.com/cuda.driver.minor": "126",
  "nvidia.com/cuda.driver.rev": "09",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/cuda.runtime-version.major": "13",
  "nvidia.com/cuda.runtime-version.minor": "0",
  "nvidia.com/cuda.runtime.major": "13",
  "nvidia.com/cuda.runtime.minor": "0",
  "nvidia.com/gfd.timestamp": "1773657194",
  "nvidia.com/gpu-driver-upgrade-state": "upgrade-done",
  "nvidia.com/gpu.compute.major": "8",
  "nvidia.com/gpu.compute.minor": "0",
  "nvidia.com/gpu.count": "1",
  "nvidia.com/gpu.deploy.container-toolkit": "true",
  "nvidia.com/gpu.deploy.dcgm": "true",
  "nvidia.com/gpu.deploy.dcgm-exporter": "true",
  "nvidia.com/gpu.deploy.device-plugin": "true",
  "nvidia.com/gpu.deploy.driver": "true",
  "nvidia.com/gpu.deploy.gpu-feature-discovery": "true",
  "nvidia.com/gpu.deploy.node-status-exporter": "true",
  "nvidia.com/gpu.deploy.nvsm": "",
  "nvidia.com/gpu.deploy.operator-validator": "true",
  "nvidia.com/gpu.family": "ampere",
  "nvidia.com/gpu.machine": "VMware201",
  "nvidia.com/gpu.memory": "40960",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-40C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r582_12",
  "nvidia.com/vgpu.host-driver-version": "580.126.08",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.35.0---vmware.2-vkr.4",
  "run.tanzu.vmware.com/tkr": "v1.35.0---vmware.2-vkr.4",
  "vks.vmware.com/nodepool": "node-pool-40c"
}
{
  "nvidia.com/cuda.driver-version.full": "580.126.09",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "126",
  "nvidia.com/cuda.driver-version.revision": "09",
  "nvidia.com/cuda.driver.major": "580",
  "nvidia.com/cuda.driver.minor": "126",
  "nvidia.com/cuda.driver.rev": "09",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/cuda.runtime-version.major": "13",
  "nvidia.com/cuda.runtime-version.minor": "0",
  "nvidia.com/cuda.runtime.major": "13",
  "nvidia.com/cuda.runtime.minor": "0",
  "nvidia.com/gfd.timestamp": "1773657205",
  "nvidia.com/gpu-driver-upgrade-state": "upgrade-done",
  "nvidia.com/gpu.compute.major": "8",
  "nvidia.com/gpu.compute.minor": "0",
  "nvidia.com/gpu.count": "1",
  "nvidia.com/gpu.deploy.container-toolkit": "true",
  "nvidia.com/gpu.deploy.dcgm": "true",
  "nvidia.com/gpu.deploy.dcgm-exporter": "true",
  "nvidia.com/gpu.deploy.device-plugin": "true",
  "nvidia.com/gpu.deploy.driver": "true",
  "nvidia.com/gpu.deploy.gpu-feature-discovery": "true",
  "nvidia.com/gpu.deploy.node-status-exporter": "true",
  "nvidia.com/gpu.deploy.nvsm": "",
  "nvidia.com/gpu.deploy.operator-validator": "true",
  "nvidia.com/gpu.family": "ampere",
  "nvidia.com/gpu.machine": "VMware201",
  "nvidia.com/gpu.memory": "4096",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-4C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r582_12",
  "nvidia.com/vgpu.host-driver-version": "580.126.08",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.35.0---vmware.2-vkr.4",
  "run.tanzu.vmware.com/tkr": "v1.35.0---vmware.2-vkr.4",
  "vks.vmware.com/nodepool": "node-pool-4c"

}
{
 "nvidia.com/cuda.driver-version.full": "580.126.09",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "126",
  "nvidia.com/cuda.driver-version.revision": "09",
  "nvidia.com/cuda.driver.major": "580",
  "nvidia.com/cuda.driver.minor": "126",
  "nvidia.com/cuda.driver.rev": "09",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/cuda.runtime-version.major": "13",
  "nvidia.com/cuda.runtime-version.minor": "0",
  "nvidia.com/cuda.runtime.major": "13",
  "nvidia.com/cuda.runtime.minor": "0",
  "nvidia.com/gfd.timestamp": "1773657211",
  "nvidia.com/gpu-driver-upgrade-state": "upgrade-done",
  "nvidia.com/gpu.compute.major": "8",
  "nvidia.com/gpu.compute.minor": "0",
  "nvidia.com/gpu.count": "1",
  "nvidia.com/gpu.deploy.container-toolkit": "true",
  "nvidia.com/gpu.deploy.dcgm": "true",
  "nvidia.com/gpu.deploy.dcgm-exporter": "true",
  "nvidia.com/gpu.deploy.device-plugin": "true",
  "nvidia.com/gpu.deploy.driver": "true",
  "nvidia.com/gpu.deploy.gpu-feature-discovery": "true",
  "nvidia.com/gpu.deploy.node-status-exporter": "true",
  "nvidia.com/gpu.deploy.nvsm": "",
  "nvidia.com/gpu.deploy.operator-validator": "true",
  "nvidia.com/gpu.family": "ampere",
  "nvidia.com/gpu.machine": "VMware201",
  "nvidia.com/gpu.memory": "4096",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-4C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r582_12",
  "nvidia.com/vgpu.host-driver-version": "580.126.08",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.35.0---vmware.2-vkr.4",
  "run.tanzu.vmware.com/tkr": "v1.35.0---vmware.2-vkr.4",
  "vks.vmware.com/nodepool": "node-pool-4c"
}
{
 "nvidia.com/cuda.driver-version.full": "580.126.09",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "126",
  "nvidia.com/cuda.driver-version.revision": "09",
  "nvidia.com/cuda.driver.major": "580",
  "nvidia.com/cuda.driver.minor": "126",
  "nvidia.com/cuda.driver.rev": "09",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/cuda.runtime-version.major": "13",
  "nvidia.com/cuda.runtime-version.minor": "0",
  "nvidia.com/cuda.runtime.major": "13",
  "nvidia.com/cuda.runtime.minor": "0",
  "nvidia.com/gfd.timestamp": "1773657228",
  "nvidia.com/gpu-driver-upgrade-state": "upgrade-done",
  "nvidia.com/gpu.compute.major": "8",
  "nvidia.com/gpu.compute.minor": "0",
  "nvidia.com/gpu.count": "1",
  "nvidia.com/gpu.deploy.container-toolkit": "true",
  "nvidia.com/gpu.deploy.dcgm": "true",
  "nvidia.com/gpu.deploy.dcgm-exporter": "true",
  "nvidia.com/gpu.deploy.device-plugin": "true",
  "nvidia.com/gpu.deploy.driver": "true",
  "nvidia.com/gpu.deploy.gpu-feature-discovery": "true",
  "nvidia.com/gpu.deploy.node-status-exporter": "true",
  "nvidia.com/gpu.deploy.nvsm": "",
  "nvidia.com/gpu.deploy.operator-validator": "true",
  "nvidia.com/gpu.family": "ampere",
  "nvidia.com/gpu.machine": "VMware201",
  "nvidia.com/gpu.memory": "4096",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-4C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r582_12",
  "nvidia.com/vgpu.host-driver-version": "580.126.08",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.35.0---vmware.2-vkr.4",
  "run.tanzu.vmware.com/tkr": "v1.35.0---vmware.2-vkr.4",
  "vks.vmware.com/nodepool": "node-pool-4c"
}
```

### Install DRA Driver

To install [NVIDIA DRA Driver](https://github.com/NVIDIA/k8s-dra-driver-gpu/wiki/Installation), add the following helm repo

```shell
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update
```

Install the driver

```shell
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu --version="25.12.0" --create-namespace --namespace nvidia-dra-driver-gpu --set nvidiaDriverRoot=/run/nvidia/driver --set resources.gpus.enabled=true --set gpuResourcesEnabledOverride=true
```

```shell
The output of the above command should be along the following lines
NAME: nvidia-dra-driver-gpu
LAST DEPLOYED: Mon Mar 16 10:40:45 2026
NAMESPACE: nvidia-dra-driver-gpu
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

The Driver installation creates the namespace `nvidia-dra-driver-gpu`. 

**Note**: The pod security policy in nvidia-dra-driver-gpu namespace should be set to privileged to proceed with driver installation and deployment of workloads.

```shell
kubectl label --overwrite ns nvidia-dra-driver-gpu pod-security.kubernetes.io/enforce=privileged
namespace/nvidia-dra-driver-gpu labeled
```

### Verify Driver Installation

The DRA controller and kubelet-plugin pods should be running and in `Ready` state

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
The output of the above command should be along the lines of the following snippet
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-6fd47d97cf-l67fz   1/1     Running   0          24s
nvidia-dra-driver-gpu-kubelet-plugin-2kh62          2/2     Running   0          24s
nvidia-dra-driver-gpu-kubelet-plugin-4kwbt          2/2     Running   0          24s
nvidia-dra-driver-gpu-kubelet-plugin-7sgkh          2/2     Running   0          24s
nvidia-dra-driver-gpu-kubelet-plugin-w9rq4          2/2     Running   0          24s
```

DeviceClass resources should be created

```shell
kubectl get deviceclasses
```

```shell
NAME                                        AGE
compute-domain-daemon.nvidia.com            85s
compute-domain-default-channel.nvidia.com   85s
gpu.nvidia.com                              85s
mig.nvidia.com                              85s
vfio.gpu.nvidia.com                         85s
```

ResourceSlice resources for the above DeviceClasses should be created

```shell
kubectl get resourceslice 
```

```shell
NAME                                                              NODE                                        DRIVER                      POOL                                        AGE
test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-compute-domain.ngrqr6   test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn   compute-domain.nvidia.com   test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn   111s
test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-gpu.nvidia.com-cjcwz    test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn   gpu.nvidia.com              test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn   112s
test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-compute-domain.nvnrjx5   test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf    compute-domain.nvidia.com   test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf    109s
test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-gpu.nvidia.com-z4694     test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf    gpu.nvidia.com              test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf    109s
test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-compute-domain.nvjktdp   test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k    compute-domain.nvidia.com   test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k    110s
test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-gpu.nvidia.com-rszc8     test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k    gpu.nvidia.com              test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k    111s
test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-compute-domain.nvcw22q   test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4    compute-domain.nvidia.com   test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4    111s
test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-gpu.nvidia.com-hx25p     test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4    gpu.nvidia.com              test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4    111s
```

A detailed look at resourceslice resources can be seen for more information.

```shell
kubectl get resourceslices -o yaml
```

```shell
apiVersion: v1
items:
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:55Z"
    generateName: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-compute-domain.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-compute-domain.ngrqr6
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
      uid: 1942fcb0-8341-4a70-ba0c-e651b46eb4d2
    resourceVersion: "23267"
    uid: 75f9c32e-5564-4886-99b3-aaaf73e0ded7
  spec:
    devices:
    - attributes:
        id:
          int: 0
        type:
          string: channel
      name: channel-0
    - attributes:
        id:
          int: 0
        type:
          string: daemon
      name: daemon-0
    driver: compute-domain.nvidia.com
    nodeName: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
    pool:
      generation: 1
      name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:54Z"
    generateName: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-gpu.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn-gpu.nvidia.com-cjcwz
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
      uid: 1942fcb0-8341-4a70-ba0c-e651b46eb4d2
    resourceVersion: "23258"
    uid: eed45ad5-b247-495c-bfe9-f26d456d23c9
  spec:
    devices:
    - attributes:
        addressingMode:
          string: None
        architecture:
          string: Ampere
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.126.9
        productName:
          string: GRID A100-40C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-b65409f1-9131-4932-bbf9-065400000000
      capacity:
        memory:
          value: 40Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
    pool:
      generation: 1
      name: test-ai-0-node-pool-40c-kn48j-7n6h9-gvhpn
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:57Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-compute-domain.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-compute-domain.nvnrjx5
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
      uid: 8e0579d4-4f32-48aa-b65f-fdaf5849f5b1
    resourceVersion: "23304"
    uid: 7e527125-f78c-4cd4-8b56-9f4031aa8157
  spec:
    devices:
    - attributes:
        id:
          int: 0
        type:
          string: channel
      name: channel-0
    - attributes:
        id:
          int: 0
        type:
          string: daemon
      name: daemon-0
    driver: compute-domain.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:57Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-gpu.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf-gpu.nvidia.com-z4694
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
      uid: 8e0579d4-4f32-48aa-b65f-fdaf5849f5b1
    resourceVersion: "23303"
    uid: 70a521d2-30c3-43a5-b9cb-acf53a70d185
  spec:
    devices:
    - attributes:
        addressingMode:
          string: None
        architecture:
          string: Ampere
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.126.9
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-4f7c0c2e-47b4-4eb2-ae41-4f3e00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-gcxqf
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:56Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-compute-domain.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-compute-domain.nvjktdp
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
      uid: 2db241c9-d4d8-4a81-9432-221e90da8870
    resourceVersion: "23281"
    uid: e9c47767-6148-414a-ad37-931230ccfc48
  spec:
    devices:
    - attributes:
        id:
          int: 0
        type:
          string: daemon
      name: daemon-0
    - attributes:
        id:
          int: 0
        type:
          string: channel
      name: channel-0
    driver: compute-domain.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:55Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-gpu.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k-gpu.nvidia.com-rszc8
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
      uid: 2db241c9-d4d8-4a81-9432-221e90da8870
    resourceVersion: "23275"
    uid: fbdd8c4e-d59d-4302-a97f-3f1f01e65d12
  spec:
    devices:
    - attributes:
        addressingMode:
          string: None
        architecture:
          string: Ampere
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.126.9
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-66e8b699-e84c-465b-9464-9dad00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-lrj4k
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:55Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-compute-domain.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-compute-domain.nvcw22q
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
      uid: 2313b8bf-ea88-4a46-8dc9-3dcf5c31fb05
    resourceVersion: "23270"
    uid: 3b1a7fd5-5352-4819-b341-0ae96d351c1b
  spec:
    devices:
    - attributes:
        id:
          int: 0
        type:
          string: daemon
      name: daemon-0
    - attributes:
        id:
          int: 0
        type:
          string: channel
      name: channel-0
    driver: compute-domain.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-03-16T10:40:55Z"
    generateName: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-gpu.nvidia.com-
    generation: 1
    name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4-gpu.nvidia.com-hx25p
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
      uid: 2313b8bf-ea88-4a46-8dc9-3dcf5c31fb05
    resourceVersion: "23269"
    uid: 34e35e70-6c9a-4650-8c42-1c2403d21696
  spec:
    devices:
    - attributes:
        addressingMode:
          string: None
        architecture:
          string: Ampere
        brand:
          string: NvidiaVCS
        cudaComputeCapability:
          version: 8.0.0
        cudaDriverVersion:
          version: 13.0.0
        driverVersion:
          version: 580.126.9
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-be9cd431-3f4d-4b86-9db5-74d000000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
    pool:
      generation: 1
      name: test-ai-0-node-pool-4c-zdfrn-rnkzj-w78x4
      resourceSliceCount: 1
kind: List
metadata:
  resourceVersion: ""
```

### Deploy a workload

To deploy a workload that utilises DRA, 

- create a `ResourceClaimTemplate` containing the `DeviceClass` requests

- create a deployment that references the `ResourceClaimTemplate`

```shell
cat <<EOF | kubectl apply -f -
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: gpu-claim-template
  namespace: gpu-test1
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

cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dra-gpu-example
  namespace: gpu-test1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dra-gpu-example
  template:
    metadata:
      labels:
        app: dra-gpu-example
    spec:
      containers:
      - name: ctr
        image: ubuntu:22.04
        command: ["bash", "-c"]
        args: ["while [ 1 ]; do date; echo $(nvidia-smi -L || echo Waiting...); sleep 60; done"]
        resources:
          claims:
          - name: single-gpu
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: gpu-claim-template
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
EOF
```

Verify that the ResourceClaim resources are created

```shell
kubectl get resourceclaims -n gpu-test1
```

```shell
NAME                                                STATE                AGE
dra-gpu-example-8df4b9554-bb8qn-single-gpu-gghvh    allocated,reserved   12s
```

Verify that the deployment is up and running

```shell
kubectl get pods -n gpu-test1
```

```shell
NAME                              READY   STATUS    RESTARTS   AGE
dra-gpu-example-8df4b9554-h8rv8   1/1     Running   0          30m
```

