### Giant Swarm Platform

Giant Swarm Platform is a managed Kubernetes platform developed by [Giant Swarm](https://www.giantswarm.io).

### How to Reproduce

#### Create Cluster

First access the [Giant Swarm Platform](https://docs.giantswarm.io/getting-started/), and login to platform API.
After successful login, select [Create a  cluster](https://docs.giantswarm.io/getting-started/provision-your-first-workload-cluster/)  with the specific DRA values.

```yaml
global:
  connectivity:
    availabilityZoneUsageLimit: 3
    network: {}
    topology: {}
  controlPlane: {}
  metadata:
    name: $CLUSTER
    $organization: fer
    preventDeletion: false
  nodePools:
    nodepool0:
      instanceType: m5.xlarge
      maxSize: 2
      minSize: 1
      rootVolumeSizeGB: 8
    nodepool1:
      instanceType: p4d.24xlarge
      maxSize: 2
      minSize: 1
      rootVolumeSizeGB: 15
      instanceWarmup: 600
      minHealthyPercentage: 90
      customNodeTaints:
      - key: "nvidia.com/gpu"
        value: "Exists"
        effect: "NoSchedule"
  providerSpecific: {}
  release:
    version: 34.0.0
```

# AI platform components

The following components should be installed to complete the AI setup:

## 1. NVIDIA GPU Operator

**Purpose**: Manages NVIDIA GPU resources in Kubernetes clusters.

**Installation via Giant Swarm App Platform**:

```sh
kubectl gs template app \
  --catalog giantswarm \
  --name gpu-operator \
  --cluster-name $CLUSTER \
  --target-namespace kube-system \
  --version 1.3.0 \
  --organization $ORGANIZATION | kubectl apply -f -
```

## 2. NVIDIA DRA Driver GPU

**Purpose**: Provides Dynamic Resource Allocation (DRA) support for NVIDIA GPUs.

**Chart**: Giant Swarm publishes a downstream build of the NVIDIA DRA driver, [`dra-driver-nvidia-gpu`](https://github.com/giantswarm/dra-driver-nvidia-gpu), patched to work on Flatcar Container Linux. The upstream chart from NGC fails on Flatcar because the kubelet-plugin container cannot resolve absolute symlinks inside its driver-root bind mount; this build adds optional extra search paths and host mounts to resolve them.

**Installation via Giant Swarm App Platform**:

```sh
# Provide the kubelet plugin toleration so it can land on GPU nodes
cat > dra-values.yaml <<EOF
kubeletPlugin:
  tolerations:
  - key: "nvidia.com/gpu"
    operator: "Exists"
    effect: "NoSchedule"
resources:
  gpus:
    enabled: false
EOF

kubectl create configmap dra-driver-nvidia-gpu-user-values \
  --from-file=values=dra-values.yaml \
  --namespace=org-$ORGANIZATION

kubectl gs template app \
  --catalog=giantswarm \
  --cluster-name=$CLUSTER \
  --organization=$ORGANIZATION \
  --name=dra-driver-nvidia-gpu \
  --app-name=dra-driver-nvidia-gpu \
  --user-configmap=dra-driver-nvidia-gpu-user-values \
  --target-namespace=kube-system \
  --version=25.3.2-flatcar.1 | kubectl apply -f -
```

## 3. Kuberay Operator

**Purpose**: Manages Ray clusters for distributed AI/ML workloads.

**Installation via Giant Swarm App Platform**:

```sh
kubectl gs template app \
  --catalog giantswarm \
  --name kuberay-operator \
  --cluster-name $CLUSTER \
  --target-namespace kube-system \
  --version 1.1.0 \
  --organization $ORGANIZATION | kubectl apply -f -
```

## 4. Kueue

**Purpose**: Provides job queueing and resource management for batch workloads.

**Installation via Flux HelmRelease**:

```sh
# First create the Kueue Helm Repository
kubectl gs template app \
  --catalog=giantswarm \
  --cluster-name$CLUSTER\
  --organization=ORGANIZATION \
  --name=kueue \
  --target-namespace=kueue-system \
  --version=0.2.0 | kubectl apply -f -
```

## 5. Gateway API

**Purpose**: Provides advanced traffic management capabilities for inference services.

**Installation via Giant Swarm App Platform**:

```sh
kubectl gs template app \
  --catalog giantswarm \
  --name gateway-api-bundle \
  --cluster-name $CLUSTER \
  --target-namespace kube-system \
  --version 1.15.1 \
  --organization $ORGANIZATION | kubectl apply -f -
```

## 6. AWS EFS CSI Driver

**Purpose**: Enables persistent storage using AWS Elastic File System for shared AI model storage.

**Installation via Giant Swarm App Platform**:

```sh
kubectl gs template app \
  --catalog giantswarm \
  --name aws-efs-csi-driver \
  --cluster-name $CLUSTER \
  --target-namespace kube-system \
  --version 3.3.0 \
  --organization $ORGANIZATION | kubectl apply -f -
```

## 7. JobSet

**Purpose**: Manages sets of Jobs for distributed training workloads.

**Installation via Flux HelmRelease**:

```sh
kubectl apply -f - <<EOF
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: $CLUSTER-jobset
  namespace: org-$ORGANIZATION
spec:
  interval: 5m
  chart:
    spec:
      chart: oci://registry.k8s.io/jobset/charts/jobset
      version: "0.12.0"
  targetNamespace: kube-system
  kubeConfig:
    secretRef:
      name: $CLUSTER-kubeconfig
      key: value
EOF
```

## 9. Prometheus Adapter

**Purpose**: Enables custom metrics for Horizontal Pod Autoscaler, including AI/ML specific metrics.

**Installation via Flux HelmRelease**:

```sh
kubectl gs template app \
  --catalog=giantswarm \
  --cluster-name=$CLUSTER \
  --org $ORGANIZATION \
  --name=keda \
  --target-namespace=keda-system \
  --version=5.0.2 | kubectl apply -f -
```

## 10. Sonobuoy Configuration

**Purpose**: Applies PolicyExceptions and configurations needed for AI conformance testing.

**Installation**: Applied directly to the workload cluster using the kubeconfig:

```sh
# Download and apply the configuration
kubectl --kubeconfig=/path/to/workload-cluster-kubeconfig apply -f https://gist.githubusercontent.com/pipo02mix/80415c1182a5920af46a85c7adf90a8a/raw/d75d7593194fb2a3beba0549f946cb6f8a5a5f46/sonobuoy-rews.yaml
```

All these components work together to provide a complete AI/ML platform on Kubernetes with GPU support, workload management, monitoring, and conformance testing capabilities.

#### Run conformance Test by Sonobuoy

Login to the control-plane of the cluster created by Giant Swarm Platform.

Start the conformance tests:

```sh
sonobuoy run --plugin https://raw.githubusercontent.com/pipo02mix/ai-conformance/c0f5f45e131445e1cf833276ca66e251b1b200e9/sonobuoy-plugin.yaml
````

Monitor the conformance tests by tracking the sonobuoy logs, and wait for the line: "no-exit was specified, sonobuoy is now blocking"

```sh
stern -n sonobuoy sonobuoy
```

Retrieve result:

```sh
outfile=$(sonobuoy retrieve)
sonobuoy results $outfile
```
