# Google Kubernetes Engine (GKE) - Kubernetes AI Conformance (v1.37)

This directory contains the Kubernetes v1.37 AI Conformance submission and automated E2E test artifacts (`e2e.log`, `junit.xml`, `results.json`) for [Google Kubernetes Engine (GKE)](https://cloud.google.com/kubernetes-engine).

## How to Reproduce

### 1. Create GKE Cluster and DRA GPU Node Pool

Create a GKE cluster running Kubernetes `1.37` and a single-zone GPU node pool with Dynamic Resource Allocation (DRA) enabled and Cluster Autoscaler configured (`min-nodes=1`, `max-nodes=2`):

```bash
gcloud container clusters create ai-conformance-cluster \
  --region=us-central1 \
  --release-channel=rapid \
  --cluster-version=1.37.0-gke.3503000 \
  --machine-type=e2-medium \
  --num-nodes=1 \
  --autoscaling-profile=optimize-utilization

gcloud container node-pools create dra-autoscaling-pool \
  --cluster=ai-conformance-cluster \
  --region=us-central1 \
  --node-locations=us-central1-a \
  --machine-type=g2-standard-4 \
  --accelerator=type=nvidia-l4,count=1,gpu-driver-version=disabled \
  --num-nodes=1 \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=2 \
  --node-labels=test-role=dra-autoscaling,gke-no-default-nvidia-gpu-device-plugin=true,nvidia.com/gpu.present=true,cloud.google.com/gke-nvidia-gpu-dra-driver=true \
  --node-taints=nvidia.com/gpu=present:NoSchedule
```

Fetch the cluster credentials:

```bash
gcloud container clusters get-credentials ai-conformance-cluster --region=us-central1
```

### 2. Install NVIDIA GPU Drivers and DRA Driver (`v25.8.1`)

Install the COS preloaded NVIDIA GPU driver DaemonSet and the `nvidia-dra-driver-gpu` Helm chart:

```bash
kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/container-engine-accelerators/master/nvidia-driver-installer/cos/daemonset-preloaded.yaml

helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu \
  --version="25.8.1" \
  --create-namespace \
  --namespace nvidia-dra-driver-gpu \
  --set nvidiaDriverRoot=/home/kubernetes/bin/nvidia/ \
  --set gpuResourcesEnabledOverride=true \
  --set resources.computeDomains.enabled=false \
  --set controller.affinity=null \
  --set kubeletPlugin.priorityClassName="" \
  --set kubeletPlugin.tolerations[0].key=nvidia.com/gpu \
  --set kubeletPlugin.tolerations[0].operator=Exists \
  --set kubeletPlugin.tolerations[0].effect=NoSchedule
```

### 3. Install Kueue (`v0.18.2`) for Gang Scheduling

Install Kueue and configure a `ResourceFlavor`, `ClusterQueue`, test namespace, and `LocalQueue`:

```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/download/v0.18.2/manifests.yaml
kubectl wait --for=condition=Available deployment/kueue-controller-manager -n kueue-system --timeout=300s

cat <<EOF | kubectl apply -f -
apiVersion: kueue.x-k8s.io/v1beta1
kind: ResourceFlavor
metadata:
  name: e2e-flavor
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: ClusterQueue
metadata:
  name: e2e-cq
spec:
  namespaceSelector: {}
  resourceGroups:
  - coveredResources: ["cpu", "memory"]
    flavors:
    - name: e2e-flavor
      resources:
      - name: "cpu"
        nominalQuota: 10
      - name: "memory"
        nominalQuota: 10Gi
---
apiVersion: v1
kind: Namespace
metadata:
  name: ai-conformance-gang-scheduling
---
apiVersion: kueue.x-k8s.io/v1beta1
kind: LocalQueue
metadata:
  namespace: ai-conformance-gang-scheduling
  name: e2e-lq
spec:
  clusterQueue: e2e-cq
EOF
```

### 4. Run the AI Conformance Test Suite

Clone the `kubernetes-sigs/ai-conformance` repository and run the automated conformance tests in DRA mode:

```bash
git clone https://github.com/kubernetes-sigs/ai-conformance.git
cd ai-conformance

go test -v -timeout 60m ./test \
  -allocation-mode=dra \
  -autoscaler-node-pool-label=test-role=dra-autoscaling \
  -gang-scheduler-namespace=ai-conformance-gang-scheduling \
  -gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq \
  -json > results.json

jq -r 'select(.Output != null) | .Output' results.json > e2e.log
gotestsum --junitfile junit.xml --raw-command -- cat results.json
```

### 5. Clean Up

```bash
gcloud container clusters delete ai-conformance-cluster --region=us-central1 --quiet
```
