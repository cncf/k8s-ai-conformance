## Description

Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

## Evidence

Dynamic Resource Allocation provides flexible resource management for specialised hardware like GPUs, FPGAs and network-attached devices. This guide covers how to configure and validate DRA on a VKS v1.36 cluster.

### Prerequisites

* Provision a VKS v3.7.0 Cluster with v1.36.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator.

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

### Verify GPU Operator

Pods in the gpu-operator namespace should be in running or completed state.

```shell
kubectl get pods -n gpu-operator
```

```shell
NAME                                                         READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-kc5zk                                  1/1     Running     0          3d18h
gpu-feature-discovery-rk8kw                                  1/1     Running     0          3d18h
gpu-operator-7d7969bcd5-zkfzd                                1/1     Running     0          3d18h
gpu-operator-node-feature-discovery-gc-5cb78546ff-wdtgw      1/1     Running     0          3d18h
gpu-operator-node-feature-discovery-master-5ffff79b5-2qjmt   1/1     Running     0          3d18h
gpu-operator-node-feature-discovery-worker-2l65g             1/1     Running     0          3d18h
gpu-operator-node-feature-discovery-worker-g4zmm             1/1     Running     0          3d18h
gpu-operator-node-feature-discovery-worker-vnwgq             1/1     Running     0          3d18h
nvidia-container-toolkit-daemonset-r89mn                     1/1     Running     0          3d18h
nvidia-container-toolkit-daemonset-w89bd                     1/1     Running     0          3d18h
nvidia-cuda-validator-5slqv                                  0/1     Completed   0          3d17h
nvidia-cuda-validator-xzpf8                                  0/1     Completed   0          3d17h
nvidia-dcgm-exporter-fxvql                                   1/1     Running     0          3d18h
nvidia-dcgm-exporter-gk2v9                                   1/1     Running     0          3d18h
nvidia-device-plugin-daemonset-bbvvq                         1/1     Running     0          3d18h
nvidia-device-plugin-daemonset-mth9d                         1/1     Running     0          3d18h
nvidia-driver-daemonset-86lcq                                1/1     Running     0          3d18h
nvidia-driver-daemonset-qbhpr                                1/1     Running     0          3d18h
nvidia-operator-validator-cfbqm                              1/1     Running     0          3d18h
nvidia-operator-validator-kwcf7                              1/1     Running     0          3d18h
```

The `ClusterPolicy` should report `ready`.

```shell
kubectl get clusterpolicy
```

```shell
NAME             STATUS   AGE
cluster-policy   ready    2026-07-02T19:15:34Z
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
  "kubernetes.io/hostname": "conformance-test-np1-qcqc8-bpdtl-rcj8j",
  "nvidia.com/cuda.driver-version.full": "580.126.09",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "126",
  "nvidia.com/cuda.driver-version.revision": "09",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/gfd.timestamp": "1783019888",
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
  "nvidia.com/gpu.deploy.operator-validator": "true",
  "nvidia.com/gpu.family": "ampere",
  "nvidia.com/gpu.machine": "VMware201",
  "nvidia.com/gpu.memory": "10240",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-10C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r596_25",
  "nvidia.com/vgpu.host-driver-version": "595.71.03",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.36.1---vmware.4-vkr.5",
  "run.tanzu.vmware.com/tkr": "v1.36.1---vmware.4-vkr.5",
  "vks.vmware.com/nodepool": "np1"
}
```

(Identical labels, differing only in GPU UUID and `gfd.timestamp`, were observed on the second worker node.)

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
NAME: nvidia-dra-driver-gpu
LAST DEPLOYED: Thu Jul  2 19:31:06 2026
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

The DRA controller and kubelet-plugin pods should be running and in `Ready` state.

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-76864c66cb-dz9hc   1/1     Running   0          3d17h
nvidia-dra-driver-gpu-kubelet-plugin-qbfsj          2/2     Running   0          3d17h
nvidia-dra-driver-gpu-kubelet-plugin-sdp5b          2/2     Running   0          3d17h
```

(Two `kubelet-plugin` pods, one per GPU-equipped worker node — the single-node control plane has no GPU and does not run a kubelet-plugin pod.)

DeviceClass resources should be created.

```shell
kubectl get deviceclasses
```

```shell
NAME                                        AGE
compute-domain-daemon.nvidia.com            3d17h
compute-domain-default-channel.nvidia.com   3d17h
gpu.nvidia.com                              3d17h
mig.nvidia.com                              3d17h
vfio.gpu.nvidia.com                         3d17h
```

ResourceSlice resources for the above DeviceClasses should be created.

```shell
kubectl get resourceslice
```

```shell
NAME                                                              NODE                                      DRIVER                      POOL                                      AGE
conformance-test-np1-qcqc8-bpdtl-rcj8j-compute-domain.nvid9wzc5   conformance-test-np1-qcqc8-bpdtl-rcj8j   compute-domain.nvidia.com   conformance-test-np1-qcqc8-bpdtl-rcj8j   3d17h
conformance-test-np1-qcqc8-bpdtl-rcj8j-gpu.nvidia.com-jfrc8       conformance-test-np1-qcqc8-bpdtl-rcj8j   gpu.nvidia.com              conformance-test-np1-qcqc8-bpdtl-rcj8j   3d17h
conformance-test-np1-qcqc8-bpdtl-wv2gs-compute-domain.nvidjkpdq   conformance-test-np1-qcqc8-bpdtl-wv2gs   compute-domain.nvidia.com   conformance-test-np1-qcqc8-bpdtl-wv2gs   3d17h
conformance-test-np1-qcqc8-bpdtl-wv2gs-gpu.nvidia.com-6kw4r       conformance-test-np1-qcqc8-bpdtl-wv2gs   gpu.nvidia.com              conformance-test-np1-qcqc8-bpdtl-wv2gs   3d17h
```

A detailed look at resourceslice resources can be seen for more information.

```shell
kubectl get resourceslices -o yaml
```

```yaml
apiVersion: v1
items:
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-07-02T19:33:59Z"
    name: conformance-test-np1-qcqc8-bpdtl-rcj8j-gpu.nvidia.com-jfrc8
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: conformance-test-np1-qcqc8-bpdtl-rcj8j
      uid: 14336b0f-e28a-41e5-8333-fd36ddcb3cda
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
          string: GRID A100-10C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-0b580e61-4daa-462a-95e3-84dd00000000
      capacity:
        memory:
          value: 10Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-qcqc8-bpdtl-rcj8j
    pool:
      generation: 1
      name: conformance-test-np1-qcqc8-bpdtl-rcj8j
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-07-02T19:33:59Z"
    name: conformance-test-np1-qcqc8-bpdtl-wv2gs-gpu.nvidia.com-6kw4r
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: conformance-test-np1-qcqc8-bpdtl-wv2gs
      uid: ef039bfd-d5dc-42af-92b7-b24402360be6
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
          string: GRID A100-10C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-44d2411e-1604-478f-aa96-5a3c00000000
      capacity:
        memory:
          value: 10Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-qcqc8-bpdtl-wv2gs
    pool:
      generation: 1
      name: conformance-test-np1-qcqc8-bpdtl-wv2gs
      resourceSliceCount: 1
kind: List
```

(The `compute-domain.nvidia.com` ResourceSlices, one channel/daemon device pair per node, are omitted above for brevity — they were also present, matching the `deviceclasses` list.)

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
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
          runAsNonRoot: true
          runAsUser: 1000
          seccompProfile:
            type: RuntimeDefault
      resourceClaims:
      - name: single-gpu
        resourceClaimTemplateName: gpu-claim-template
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
EOF
```

Verify that the ResourceClaim resources are created.

```shell
kubectl get resourceclaims -n gpu-test1
```

```shell
NAME                                               STATE                AGE
dra-gpu-example-798fccd4d-sfl2p-single-gpu-d2jzk   allocated,reserved   22s
```

```shell
kubectl get resourceclaims -n gpu-test1 -o yaml
```

```yaml
apiVersion: v1
items:
- apiVersion: resource.k8s.io/v1
  kind: ResourceClaim
  metadata:
    annotations:
      resource.kubernetes.io/pod-claim-name: single-gpu
    creationTimestamp: "2026-07-06T13:16:19Z"
    generateName: dra-gpu-example-798fccd4d-sfl2p-single-gpu-
    name: dra-gpu-example-798fccd4d-sfl2p-single-gpu-d2jzk
    namespace: gpu-test1
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Pod
      name: dra-gpu-example-798fccd4d-sfl2p
      uid: d36ca24f-9949-4c33-9d47-c06846f775b1
  spec:
    devices:
      requests:
      - exactly:
          allocationMode: ExactCount
          count: 1
          deviceClassName: gpu.nvidia.com
        name: single-gpu
  status:
    allocation:
      allocationTimestamp: "2026-07-06T13:16:19Z"
      devices:
        results:
        - device: gpu-0
          driver: gpu.nvidia.com
          pool: conformance-test-np1-qcqc8-bpdtl-wv2gs
          request: single-gpu
      nodeSelector:
        nodeSelectorTerms:
        - matchFields:
          - key: metadata.name
            operator: In
            values:
            - conformance-test-np1-qcqc8-bpdtl-wv2gs
    reservedFor:
    - name: dra-gpu-example-798fccd4d-sfl2p
      resource: pods
      uid: d36ca24f-9949-4c33-9d47-c06846f775b1
kind: List
```

Verify that the deployment is up and running.

```shell
kubectl get pods -n gpu-test1
```

```shell
NAME                              READY   STATUS    RESTARTS   AGE
dra-gpu-example-798fccd4d-sfl2p   1/1     Running   0          22s
```

Finally, verify the GPU allocated through DRA is visible and usable inside the container.

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example
```

```shell
Mon Jul  6 13:16:25 UTC 2026
GPU 0: GRID A100-10C (UUID: GPU-44d2411e-1604-478f-aa96-5a3c00000000)
```

The GPU UUID reported by `nvidia-smi -L` inside the container matches the `gpu-0` device UUID (`GPU-44d2411e-1604-478f-aa96-5a3c00000000`) allocated to node `conformance-test-np1-qcqc8-bpdtl-wv2gs` in the `ResourceSlice`/`ResourceClaim` shown above, confirming end-to-end DRA-based GPU allocation is working correctly on VKS v1.36.
