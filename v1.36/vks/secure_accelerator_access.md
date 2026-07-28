## Description 

Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Evidence

### Test 1: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads

#### Prerequisites

* Provision a VKS v3.7.0 Cluster with v1.36.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

#### Install DRA Driver

Add the following helm repo to initiate driver installation

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

Verify that the DRA controller and kubelet-plugin pods are running.

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-76864c66cb-dz9hc   1/1     Running   0          3d17h
nvidia-dra-driver-gpu-kubelet-plugin-qbfsj          2/2     Running   0          3d17h
nvidia-dra-driver-gpu-kubelet-plugin-sdp5b          2/2     Running   0          3d17h
```

#### Deploy a workload

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

Verify that the ResourceClaim resources are created

```shell
kubectl get resourceclaims -n gpu-test1
```

```shell
NAME                                               STATE                AGE
dra-gpu-example-798fccd4d-sfl2p-single-gpu-d2jzk   allocated,reserved   22s
```

A detailed look at the `ResourceClaim` shows the device allocated and the node it is bound to, mediated entirely by the DRA driver and the Kubernetes resource management framework.

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

Verify that the deployment is up and running

```shell
kubectl get pods -n gpu-test1
```

```shell
NAME                              READY   STATUS    RESTARTS   AGE
dra-gpu-example-798fccd4d-sfl2p   1/1     Running   0          22s
```

Verify that the pod has access to GPU resources

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example
```

```shell
Mon Jul  6 13:16:25 UTC 2026
GPU 0: GRID A100-10C (UUID: GPU-44d2411e-1604-478f-aa96-5a3c00000000)
```

The GPU UUID reported by `nvidia-smi -L` inside the container matches the `gpu-0` device UUID (`GPU-44d2411e-1604-478f-aa96-5a3c00000000`) allocated to node `conformance-test-np1-qcqc8-bpdtl-wv2gs` in the `ResourceClaim` shown above, confirming that access to the accelerator is properly isolated and mediated end-to-end by the DRA resource management framework and container runtime on VKS v1.36.

### Test 2: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads

Run the DRA E2E test to verify that configs and devices are mapped to the right containers preventing unauthorized access or interference between workloads

```shell
mkdir k8s.io && cd k8s.io && git clone https://github.com/kubernetes/kubernetes
cd kubernetes
git checkout v1.36.1
```

```shell
$ make WHAT="github.com/onsi/ginkgo/v2/ginkgo k8s.io/kubernetes/test/e2e/e2e.test" && KUBERNETES_PROVIDER=local hack/ginkgo-e2e.sh -ginkgo.focus='must map configs and devices to the right containers'
Setting up for KUBERNETES_PROVIDER="local".
Skeleton Provider: prepare-e2e not implemented
KUBE_MASTER_IP: 
KUBE_MASTER: 
  I0707 06:23:04.590328   49880 e2e.go:109] Starting e2e run "e344aa72-d38b-4ec3-a07c-0ae4eeca9bcf" on Ginkgo node 1
Running Suite: Kubernetes e2e suite - /home/hardik/ai-conformance-136/k8s.io/kubernetes/_output/bin
===================================================================================================
Random Seed: 1783405383 - will randomize all specs

Will run 1 of 7579 specs
•

Ran 1 of 7579 Specs in 16.775 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 7578 Skipped
PASS

Ginkgo ran 1 suite in 18.13566611s
Test Suite Passed
```
