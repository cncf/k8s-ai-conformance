# Omni AI Kubernetes cluster

[Omni](https://siderolabs.com/omni) allows you to create a Kubernetes cluster
with any machine in any environment. For this guide you will need to provide
your own machines connected to Omni — at minimum one control-plane node and
one worker node with an NVIDIA accelerator attached.

You will also need an Omni instance which you can sign up for at
[Omni SaaS](https://siderolabs.com/omni-signup).

The submitted evidence in this directory was collected against a two-node
cluster:

- **Control plane**: Raspberry Pi 5 (arm64)
- **Worker**: NVIDIA DGX Spark GB10 (arm64, single GPU exposed via the DRA
  driver `gpu.nvidia.com`)

## Create a cluster

Follow [the documentation](https://docs.siderolabs.com/omni/omni-cluster-setup/registering-machines/register-machines-with-omni)
to connect at least two machines to Omni.

Identify one machine to use as the control plane and one (or more) with an AI
accelerator to use as workers.

```
omnictl get machines
```

Export the machine UUIDs as variables for the cluster template.

```bash
export CP_UUID=38373431-6563-3733-3564-386565373200
export WORKER_UUID=d68614b6-bfde-11d3-8000-4cbb472dcf87
```

Create a cluster template.

```bash
cat << EOF > cluster-template.yaml
kind: Cluster
name: ai-conformance
kubernetes:
  version: v1.37.0
talos:
  version: v1.14.1
---
kind: ControlPlane
machines:
  - ${CP_UUID}
---
kind: Workers
machines:
  - ${WORKER_UUID}
---
kind: Machine
name: ${CP_UUID}
systemExtensions:
  - siderolabs/iscsi-tools
---
kind: Machine
name: ${WORKER_UUID}
systemExtensions:
  - siderolabs/nonfree-kmod-nvidia-production
  - siderolabs/nvidia-container-toolkit-production
kernelArgs:
  - arm64.nobti
patches:
  - idOverride: 500-gb10-nvidia-modules
    inline:
      machine:
        kernel:
          modules:
            - name: nvidia
              parameters:
                - NVreg_CoherentGPUMemoryMode=driver
            - name: nvidia_uvm
            - name: nvidia_drm
            - name: nvidia_modeset
EOF
```

Apply the cluster template to create a Kubernetes cluster.

```bash
omnictl cluster template sync -f cluster-template.yaml
```

## Kubeconfig

When the cluster is ready download the kubeconfig file:

```bash
omnictl kubeconfig -c ai-conformance
```

## NVIDIA DRA driver

Install the NVIDIA DRA driver to expose the GB10 GPU as a DRA-managed
resource:

```bash
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia
helm repo update
helm install dra-driver-nvidia-gpu \
  nvidia/dra-driver-nvidia-gpu \
  --namespace nvidia-system \
  --create-namespace
```

## Deploy Traefik Gateway API

Install Gateway API CRDs:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v3.5/docs/content/reference/dynamic-configuration/kubernetes-gateway-rbac.yml
```

Install Traefik via Helm:

```bash
cat << EOF > values.yaml
providers:
  kubernetesGateway:
    enabled: true
EOF

helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  -n traefik --create-namespace \
  -f values.yaml
```

## Deploy Kueue for gang scheduling

Install Kueue (v0.19.4 was used to collect the submitted test artifacts):

```bash
kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/download/v0.19.4/manifests.yaml
```

Configure a `ResourceFlavor`, `ClusterQueue`, and `LocalQueue` in a test
namespace so the AI-conformance gang scheduling test can admit workloads.

## Deploy KubeRay operator (optional)

Install KubeRay via the Helm chart to satisfy the `robust_controller`
requirement:

```bash
helm repo add kuberay https://ray-project.github.io/kuberay-helm/
helm repo update
helm install kuberay-operator kuberay/kuberay-operator --version 1.4.2 --namespace kube-system
```

## Run the AI Conformance test suite

Clone the upstream test suite and run against the cluster:

```bash
git clone https://github.com/kubernetes-sigs/ai-conformance
cd ai-conformance
export KUBECONFIG=/path/to/kubeconfig
go test -v -timeout 30m ./test \
  -gang-scheduler-namespace=e2e-gang \
  -gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq \
  | tee e2e.log
go test -v ./test -json > results.json \
  -gang-scheduler-namespace=e2e-gang \
  -gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq
go run gotest.tools/gotestsum@latest --junitfile junit.xml -- ./test \
  -gang-scheduler-namespace=e2e-gang \
  -gang-job-labels=kueue.x-k8s.io/queue-name=e2e-lq
```

`e2e.log`, `results.json`, and `junit.xml` in this directory were generated
by the commands above.
