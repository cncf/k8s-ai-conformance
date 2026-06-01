## Description 

Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Evidence

### Test 1: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads

#### Prerequisites

* Provision a VKS v3.6.0 Cluster with v1.35.0 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

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
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu --version="25.3.2" --create-namespace --namespace nvidia-dra-driver-gpu --set nvidiaDriverRoot=/run/nvidia/driver --set resources.gpus.enabled=true --set gpuResourcesEnabledOverride=true
```

The Driver installation creates the namespace `nvidia-dra-driver-gpu`. 

**Note**: The pod security policy in nvidia-dra-driver-gpu namespace should be set to privileged to proceed with driver installation and deployment of workloads.

```shell
kubectl label --overwrite ns nvidia-dra-driver-gpu pod-security.kubernetes.io/enforce=privileged
namespace/nvidia-dra-driver-gpu labeled
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
dra-gpu-example-7c7d5856fb-mbn76-single-gpu-2tmqp   allocated,reserved   22s
```

Verify that the deployment is up and running

```shell
kubectl get pods -n gpu-test1
```

```shell
NAME                               READY   STATUS    RESTARTS   AGE
dra-gpu-example-7c7d5856fb-mbn76   1/1     Running   0          81s
```

Verify that the pod has access to GPU resources

```shell
kubectl logs dra-gpu-example-7c7d5856fb-mbn76 -n gpu-test1
```

```shell
Tue Mar 17 13:18:00 UTC 2026
GPU 0: GRID A100-40C (UUID: GPU-0946f70c-659b-4c39-8b78-7e7600000000)
```

### Test 2: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads

Run the DRA E2E test to verify that configs and devices are mapped to the right containers preventing unauthorized access or interference between workloads

```shell
mkdir k8s.io && cd k8s.io && git clone https://github.com/kubernetes/kubernetes
cd kubernetes
git checkout v1.35.0
```

```shell
$ make WHAT="github.com/onsi/ginkgo/v2/ginkgo k8s.io/kubernetes/test/e2e/e2e.test" && KUBERNETES_PROVIDER=local hack/ginkgo-e2e.sh -ginkgo.focus='must map configs and devices to the right containers'
go: downloading go1.25.5 (linux/amd64)
+++ [0317 13:23:16] Building go targets for linux/amd64
    github.com/onsi/ginkgo/v2/ginkgo (non-static)
    k8s.io/kubernetes/test/e2e/e2e.test (test)
Setting up for KUBERNETES_PROVIDER="local".
Skeleton Provider: prepare-e2e not implemented
KUBE_MASTER_IP:
KUBE_MASTER:
  I0317 13:31:04.997382  299708 e2e.go:109] Starting e2e run "78e7893b-b600-49a8-93a3-f5c3093c878e" on Ginkgo node 1
Running Suite: Kubernetes e2e suite - /home/hardik/ai-conformance/secure-accelerator-access/k8s.io/kubernetes/_output/bin
=========================================================================================================================
Random Seed: 1773754264 - will randomize all specs

Will run 1 of 7348 specs
•

Ran 1 of 7348 Specs in 14.804 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 7347 Skipped
PASS

Ginkgo ran 1 suite in 16.129434516s
Test Suite Passed
```
