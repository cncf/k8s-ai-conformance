## Description

Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

## Evidence

Dynamic Resource Allocation provides flexible resource management for specialised hardware like GPUs, FPGAs and network-attached devices. This guide covers how to configure and validate DRA on a VKS v1.37 cluster.

### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator.

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

```shell
kubectl get nodes -o wide
```

```shell
NAME                                     STATUS   ROLES           AGE    VERSION            INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION             CONTAINER-RUNTIME
conformance-test-np1-sms5g-pt7t9-p42p6   Ready    <none>          6d1h   v1.37.1+vmware.4   172.30.0.3    <none>        Ubuntu 24.04.5 LTS   7.0.0-38-generic (amd64)   containerd://2.4.0+vmware.2-fips
conformance-test-np1-sms5g-pt7t9-vv2dd   Ready    <none>          6d1h   v1.37.1+vmware.4   172.30.0.4    <none>        Ubuntu 24.04.5 LTS   7.0.0-38-generic (amd64)   containerd://2.4.0+vmware.2-fips
conformance-test-pcdv2-gv6sq             Ready    control-plane   6d1h   v1.37.1+vmware.4   172.30.0.2    <none>        Ubuntu 24.04.5 LTS   7.0.0-38-generic (amd64)   containerd://2.4.0+vmware.2-fips
```

### Verify GPU Operator

Pods in the gpu-operator namespace should be in running or completed state (NVIDIA GPU Operator v25.10.1).

```shell
kubectl get pods -n gpu-operator
```

```shell
NAME                                                          READY   STATUS      RESTARTS   AGE
gpu-feature-discovery-22ml4                                   1/1     Running     0          5d
gpu-feature-discovery-rd7d5                                   1/1     Running     0          5d
gpu-operator-d57f6db6b-9srzg                                  1/1     Running     0          5d
gpu-operator-node-feature-discovery-gc-5f9c65779c-85d2x       1/1     Running     0          5d
gpu-operator-node-feature-discovery-master-596d8bc6c8-x6hlc   1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-5kv94              1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-kprf9              1/1     Running     0          5d
gpu-operator-node-feature-discovery-worker-mr5xf              1/1     Running     0          5d
nvidia-container-toolkit-daemonset-5nxjc                      1/1     Running     0          5d
nvidia-container-toolkit-daemonset-mswb8                      1/1     Running     0          5d
nvidia-cuda-validator-j95hw                                   0/1     Completed   0          5d
nvidia-cuda-validator-xbzc4                                   0/1     Completed   0          5d
nvidia-dcgm-exporter-gbtzd                                    1/1     Running     0          5d
nvidia-dcgm-exporter-wlfht                                    1/1     Running     0          5d
nvidia-device-plugin-daemonset-77fgr                          1/1     Running     0          5d
nvidia-device-plugin-daemonset-glxm5                          1/1     Running     0          5d
nvidia-driver-daemonset-rq87s                                 1/1     Running     0          5d
nvidia-driver-daemonset-twmhg                                 1/1     Running     0          5d
nvidia-operator-validator-28zsq                               1/1     Running     0          5d
nvidia-operator-validator-mjq6n                               1/1     Running     0          5d
```

The `ClusterPolicy` should report `ready`.

```shell
kubectl get clusterpolicy
```

```shell
NAME             STATUS   AGE
cluster-policy   ready    2026-10-29T22:27:36Z
```

Cluster nodes should carry labels with GPU information.

```shell
kubectl get node conformance-test-np1-sms5g-pt7t9-p42p6 -o json | jq '.metadata.labels'
```

The node labels demonstrate that VKS supports GPU driver and runtime lifecycle management through the NVIDIA GPU Operator. The operator automates the deployment and management of:

- **NVIDIA Driver**: Deployed via `nvidia-driver-daemonset`, with version information exposed in node labels (`nvidia.com/cuda.driver-version.full`, `nvidia.com/vgpu.host-driver-version`)
- **Container Runtime**: Deployed via `nvidia-container-toolkit-daemonset`, enabling GPU-accelerated containers
- **Verification Mechanism**: Node labels serve as the verification that drivers and runtime are properly installed and operational

The DRA ResourceSlice attributes (shown later in this document) also expose driver version information through the `driverVersion` attribute, providing additional verification of the driver lifecycle management.

(`feature.node.kubernetes.io/*` CPU feature labels omitted for brevity.)

```shell
{
  "kubernetes.io/hostname": "conformance-test-np1-sms5g-pt7t9-p42p6",
  "nvidia.com/cuda.driver-version.full": "580.178.04",
  "nvidia.com/cuda.driver-version.major": "580",
  "nvidia.com/cuda.driver-version.minor": "178",
  "nvidia.com/cuda.driver-version.revision": "04",
  "nvidia.com/cuda.runtime-version.full": "13.0",
  "nvidia.com/gfd.timestamp": "1793313050",
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
  "nvidia.com/gpu.memory": "4096",
  "nvidia.com/gpu.mode": "compute",
  "nvidia.com/gpu.present": "true",
  "nvidia.com/gpu.product": "GRID-A100-4C",
  "nvidia.com/gpu.replicas": "1",
  "nvidia.com/gpu.sharing-strategy": "none",
  "nvidia.com/mig.capable": "false",
  "nvidia.com/mig.strategy": "single",
  "nvidia.com/mps.capable": "false",
  "nvidia.com/vgpu.host-driver-branch": "r596_84",
  "nvidia.com/vgpu.host-driver-version": "595.91.04",
  "nvidia.com/vgpu.present": "true",
  "run.tanzu.vmware.com/kubernetesDistributionVersion": "v1.37.1---vmware.4-vkr.1",
  "run.tanzu.vmware.com/tkr": "v1.37.1---vmware.4-vkr.1",
  "vks.vmware.com/nodepool": "np1"
}
```

(Identical labels, differing only in hostname and `gfd.timestamp`, were observed on the second worker node `conformance-test-np1-sms5g-pt7t9-vv2dd`.)

### Install DRA Driver

To install [NVIDIA DRA Driver](https://github.com/NVIDIA/k8s-dra-driver-gpu/wiki/Installation), add the following helm repo

```shell
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update
```

**Note**: The pod security policy in the `nvidia-dra-driver-gpu` namespace must be set to privileged for the driver's kubelet-plugin DaemonSet. The namespace is created and labelled before installation.

```shell
kubectl create namespace nvidia-dra-driver-gpu
kubectl label --overwrite ns nvidia-dra-driver-gpu pod-security.kubernetes.io/enforce=privileged
```

```shell
namespace/nvidia-dra-driver-gpu created
namespace/nvidia-dra-driver-gpu labeled
```

Install the driver

```shell
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu --version=25.12.0 --namespace nvidia-dra-driver-gpu --set nvidiaDriverRoot=/run/nvidia/driver --set resources.gpus.enabled=true --set gpuResourcesEnabledOverride=true --wait --timeout 300s
```

```shell
NAME: nvidia-dra-driver-gpu
LAST DEPLOYED: Tue Nov  3 23:07:48 2026
NAMESPACE: nvidia-dra-driver-gpu
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

### Verify Driver Installation

The DRA controller and kubelet-plugin pods should be running and in `Ready` state.

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-6d4c57447b-lp4tk   1/1     Running   0          17s
nvidia-dra-driver-gpu-kubelet-plugin-jbp8c          2/2     Running   0          17s
nvidia-dra-driver-gpu-kubelet-plugin-zz64m          2/2     Running   0          17s
```

(Two `kubelet-plugin` pods, one per GPU-equipped worker node — the control plane has no GPU and does not run a kubelet-plugin pod.)

DeviceClass resources should be created.

```shell
kubectl get deviceclasses
```

```shell
NAME                                        AGE
compute-domain-daemon.nvidia.com            17s
compute-domain-default-channel.nvidia.com   17s
gpu.nvidia.com                              17s
mig.nvidia.com                              17s
vfio.gpu.nvidia.com                         17s
```

ResourceSlice resources for the above DeviceClasses should be created.

```shell
kubectl get resourceslice
```

```shell
NAME                                                              NODE                                     DRIVER                      POOL                                     AGE
conformance-test-np1-sms5g-pt7t9-p42p6-compute-domain.nvid457rh   conformance-test-np1-sms5g-pt7t9-p42p6   compute-domain.nvidia.com   conformance-test-np1-sms5g-pt7t9-p42p6   14s
conformance-test-np1-sms5g-pt7t9-p42p6-gpu.nvidia.com-2tf4l       conformance-test-np1-sms5g-pt7t9-p42p6   gpu.nvidia.com              conformance-test-np1-sms5g-pt7t9-p42p6   15s
conformance-test-np1-sms5g-pt7t9-vv2dd-compute-domain.nvid2drmh   conformance-test-np1-sms5g-pt7t9-vv2dd   compute-domain.nvidia.com   conformance-test-np1-sms5g-pt7t9-vv2dd   14s
conformance-test-np1-sms5g-pt7t9-vv2dd-gpu.nvidia.com-28k4d       conformance-test-np1-sms5g-pt7t9-vv2dd   gpu.nvidia.com              conformance-test-np1-sms5g-pt7t9-vv2dd   14s
```

A detailed look at the `gpu.nvidia.com` ResourceSlices:

```shell
kubectl get resourceslices -o json | jq '.items[] | select(.spec.driver=="gpu.nvidia.com")'
```

```yaml
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-11-03T23:07:50Z"
    name: conformance-test-np1-sms5g-pt7t9-p42p6-gpu.nvidia.com-2tf4l
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: conformance-test-np1-sms5g-pt7t9-p42p6
      uid: 56a0a297-b610-4d2a-adf6-b34bf65c0990
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
          version: 580.178.4
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-ba9adcc0-9809-47bb-85a6-b46a00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-sms5g-pt7t9-p42p6
    pool:
      generation: 1
      name: conformance-test-np1-sms5g-pt7t9-p42p6
      resourceSliceCount: 1
- apiVersion: resource.k8s.io/v1
  kind: ResourceSlice
  metadata:
    creationTimestamp: "2026-11-03T23:07:51Z"
    name: conformance-test-np1-sms5g-pt7t9-vv2dd-gpu.nvidia.com-28k4d
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Node
      name: conformance-test-np1-sms5g-pt7t9-vv2dd
      uid: 1ab95080-26f1-408f-b758-7cfa94b1de4d
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
          version: 580.178.4
        productName:
          string: GRID A100-4C
        resource.kubernetes.io/pciBusID:
          string: "0000:03:00.0"
        resource.kubernetes.io/pcieRoot:
          string: pci0000:03
        type:
          string: gpu
        uuid:
          string: GPU-60e223c0-015b-436c-ba00-26da00000000
      capacity:
        memory:
          value: 4Gi
      name: gpu-0
    driver: gpu.nvidia.com
    nodeName: conformance-test-np1-sms5g-pt7t9-vv2dd
    pool:
      generation: 1
      name: conformance-test-np1-sms5g-pt7t9-vv2dd
      resourceSliceCount: 1
```

(The JSON output has been rendered as YAML for readability. The `compute-domain.nvidia.com` ResourceSlices, one per node, are omitted above for brevity — they were also present, matching the `deviceclasses` list.)

### Deploy a workload

To deploy a workload that utilises DRA,

- create a `ResourceClaimTemplate` containing the `DeviceClass` requests

- create a deployment that references the `ResourceClaimTemplate`

```shell
kubectl create namespace gpu-test1
```

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

```shell
resourceclaimtemplate.resource.k8s.io/gpu-claim-template created
deployment.apps/dra-gpu-example created
```

Verify that the ResourceClaim resources are created.

```shell
kubectl get resourceclaims -n gpu-test1
```

```shell
NAME                                                STATE                AGE
dra-gpu-example-6cf7c799b6-wtzvq-single-gpu-498f6   allocated,reserved   9s
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
    creationTimestamp: "2026-11-03T23:08:15Z"
    finalizers:
    - resource.kubernetes.io/delete-protection
    generateName: dra-gpu-example-6cf7c799b6-wtzvq-single-gpu-
    name: dra-gpu-example-6cf7c799b6-wtzvq-single-gpu-498f6
    namespace: gpu-test1
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Pod
      name: dra-gpu-example-6cf7c799b6-wtzvq
      uid: 3fc549b8-6f24-45df-b888-79ebf46c37cd
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
      allocationTimestamp: "2026-11-03T23:08:15Z"
      devices:
        results:
        - device: gpu-0
          driver: gpu.nvidia.com
          pool: conformance-test-np1-sms5g-pt7t9-vv2dd
          request: single-gpu
      nodeSelector:
        nodeSelectorTerms:
        - matchFields:
          - key: metadata.name
            operator: In
            values:
            - conformance-test-np1-sms5g-pt7t9-vv2dd
    reservedFor:
    - name: dra-gpu-example-6cf7c799b6-wtzvq
      resource: pods
      uid: 3fc549b8-6f24-45df-b888-79ebf46c37cd
kind: List
```

Verify that the deployment is up and running.

```shell
kubectl get pods -n gpu-test1 -o wide
```

```shell
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
dra-gpu-example-6cf7c799b6-wtzvq   1/1     Running   0          9s    193.0.1.53   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

Finally, verify the GPU allocated through DRA is visible and usable inside the container.

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example
```

```shell
Tue Nov  3 23:08:16 UTC 2026
GPU 0: GRID A100-4C (UUID: GPU-60e223c0-015b-436c-ba00-26da00000000)
```

The GPU UUID reported by `nvidia-smi -L` inside the container matches the `gpu-0` device UUID (`GPU-60e223c0-015b-436c-ba00-26da00000000`) published in the ResourceSlice of node `conformance-test-np1-sms5g-pt7t9-vv2dd` and allocated to the pod in the `ResourceClaim` shown above, confirming end-to-end DRA-based GPU allocation is working correctly on VKS v1.37.

### Cleanup

```shell
kubectl delete namespace gpu-test1
helm uninstall nvidia-dra-driver-gpu -n nvidia-dra-driver-gpu
kubectl delete namespace nvidia-dra-driver-gpu
# Helm does not remove CRDs; the kubelet plugin's ResourceSlices are removed explicitly as well
kubectl delete crd computedomaincliques.resource.nvidia.com computedomains.resource.nvidia.com
kubectl get resourceslice -o name | xargs -r kubectl delete
```
