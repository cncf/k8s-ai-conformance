## Description 

Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Evidence

### Test 1: Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads

#### Prerequisites

* Provision a VKS v3.8.0 Cluster with v1.37.1 node pool, VM Class with vGPU profile and NVIDIA GPU Operator

* Log in to the cluster as admin

References:

- https://techdocs.broadcom.com/us/en/vmware-cis/private-ai/foundation-with-nvidia/9-0/private-ai-foundation-9-x/deploying-ai-workloads-on-tkg-clusters/deploy-a-gpu-accelerated-tkg-cluster-with-kubectl-connected.html

#### Install DRA Driver

Add the following helm repo to initiate driver installation

```shell
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update
```

**Note**: The pod security policy in nvidia-dra-driver-gpu namespace should be set to privileged to proceed with driver installation and deployment of workloads.

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
LAST DEPLOYED: Tue Nov  3 23:10:27 2026
NAMESPACE: nvidia-dra-driver-gpu
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
TEST SUITE: None
```

Verify that the DRA controller and kubelet-plugin pods are running.

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-6d4c57447b-qqwqd   1/1     Running   0          17s
nvidia-dra-driver-gpu-kubelet-plugin-47qqd          2/2     Running   0          17s
nvidia-dra-driver-gpu-kubelet-plugin-9kdqb          2/2     Running   0          17s
```

Each GPU worker node publishes exactly one `gpu.nvidia.com` device with a distinct UUID:

```shell
kubectl get resourceslice -o custom-columns=NAME:.metadata.name,NODE:.spec.nodeName,DRIVER:.spec.driver,UUID:.spec.devices[0].attributes.uuid.string
```

```shell
NAME                                                              NODE                                     DRIVER                      UUID
conformance-test-np1-sms5g-pt7t9-p42p6-compute-domain.nvidmg8rp   conformance-test-np1-sms5g-pt7t9-p42p6   compute-domain.nvidia.com   <none>
conformance-test-np1-sms5g-pt7t9-p42p6-gpu.nvidia.com-f28pk       conformance-test-np1-sms5g-pt7t9-p42p6   gpu.nvidia.com              GPU-ba9adcc0-9809-47bb-85a6-b46a00000000
conformance-test-np1-sms5g-pt7t9-vv2dd-compute-domain.nvidv9f96   conformance-test-np1-sms5g-pt7t9-vv2dd   compute-domain.nvidia.com   <none>
conformance-test-np1-sms5g-pt7t9-vv2dd-gpu.nvidia.com-m4jr7       conformance-test-np1-sms5g-pt7t9-vv2dd   gpu.nvidia.com              GPU-60e223c0-015b-436c-ba00-26da00000000
```

#### Deploy a workload

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

Verify that the ResourceClaim resources are created

```shell
kubectl get resourceclaims -n gpu-test1
```

```shell
NAME                                                STATE                AGE
dra-gpu-example-6cf7c799b6-6j66p-single-gpu-jcgx6   allocated,reserved   9s
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
    creationTimestamp: "2026-11-03T23:10:44Z"
    finalizers:
    - resource.kubernetes.io/delete-protection
    generateName: dra-gpu-example-6cf7c799b6-6j66p-single-gpu-
    name: dra-gpu-example-6cf7c799b6-6j66p-single-gpu-jcgx6
    namespace: gpu-test1
    ownerReferences:
    - apiVersion: v1
      controller: true
      kind: Pod
      name: dra-gpu-example-6cf7c799b6-6j66p
      uid: e64640c4-c20f-4c59-a1bc-3faf5c9aac82
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
      allocationTimestamp: "2026-11-03T23:10:44Z"
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
    - name: dra-gpu-example-6cf7c799b6-6j66p
      resource: pods
      uid: e64640c4-c20f-4c59-a1bc-3faf5c9aac82
kind: List
```

Verify that the deployment is up and running

```shell
kubectl get pods -n gpu-test1 -o wide
```

```shell
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
dra-gpu-example-6cf7c799b6-6j66p   1/1     Running   0          9s    193.0.1.55   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

Verify that the pod has access to GPU resources

```shell
kubectl logs -n gpu-test1 deployment/dra-gpu-example
```

```shell
Tue Nov  3 23:10:45 UTC 2026
GPU 0: GRID A100-4C (UUID: GPU-60e223c0-015b-436c-ba00-26da00000000)
```

The GPU UUID reported by `nvidia-smi -L` inside the container matches the `gpu-0` device UUID (`GPU-60e223c0-015b-436c-ba00-26da00000000`) of node `conformance-test-np1-sms5g-pt7t9-vv2dd` allocated in the `ResourceClaim` shown above, confirming that access to the accelerator is properly isolated and mediated end-to-end by the DRA resource management framework and container runtime on VKS v1.37.

### Test 2: Verify that workloads cannot access accelerators that were not allocated to them

#### Isolation check A: a pod without a ResourceClaim on the same GPU node sees no GPU

A pod that requests no accelerator is pinned to the same node (`conformance-test-np1-sms5g-pt7t9-vv2dd`) on which the DRA workload above holds the GPU.

```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: no-gpu-claim
  namespace: gpu-test1
spec:
  nodeName: conformance-test-np1-sms5g-pt7t9-vv2dd
  restartPolicy: Never
  containers:
  - name: ctr
    image: ubuntu:22.04
    command: ["bash", "-c"]
    args: ["ls -l /dev/nvidia* 2>&1; nvidia-smi -L 2>&1 || echo 'nvidia-smi: no GPU visible to this container'; sleep 3600"]
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop: ["ALL"]
      runAsNonRoot: true
      runAsUser: 1000
      seccompProfile:
        type: RuntimeDefault
  tolerations:
  - key: "nvidia.com/gpu"
    operator: "Exists"
    effect: "NoSchedule"
EOF
```

```shell
kubectl get pod no-gpu-claim -n gpu-test1 -o wide
```

```shell
NAME           READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
no-gpu-claim   1/1     Running   0          4s    193.0.1.56   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
```

```shell
kubectl logs -n gpu-test1 no-gpu-claim
```

```shell
ls: cannot access '/dev/nvidia*': No such file or directory
bash: line 1: nvidia-smi: command not found
nvidia-smi: no GPU visible to this container
```

By contrast, the container that holds the ResourceClaim on the same node has the device nodes injected by the container runtime:

```shell
kubectl exec -n gpu-test1 deploy/dra-gpu-example -- bash -c 'ls -l /dev/nvidia*'
```

```shell
crw-rw-rw- 1 1000 root 195, 254 Nov  3 23:10 /dev/nvidia-modeset
crw-rw-rw- 1 1000 root 235,   0 Nov  3 23:10 /dev/nvidia-uvm
crw-rw-rw- 1 1000 root 235,   1 Nov  3 23:10 /dev/nvidia-uvm-tools
crw-rw-rw- 1 1000 root 195,   0 Nov  3 23:10 /dev/nvidia0
crw-rw-rw- 1 1000 root 195, 255 Nov  3 23:10 /dev/nvidiactl
```

The un-allocated pod has neither the GPU device nodes nor the driver user-space tools injected, even though it shares the node with a GPU workload — accelerator access is granted only through the Kubernetes resource management framework.

#### Isolation check B: allocated devices are exclusive to their claim

The cluster has two GPUs (one per worker). The DRA deployment is scaled to 3 replicas, each with its own ResourceClaim.

```shell
kubectl scale deployment dra-gpu-example -n gpu-test1 --replicas=3
```

```shell
kubectl get pods -n gpu-test1 -l app=dra-gpu-example -o wide
```

```shell
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE                                     NOMINATED NODE   READINESS GATES
dra-gpu-example-6cf7c799b6-6j66p   1/1     Running   0          74s   193.0.1.55   conformance-test-np1-sms5g-pt7t9-vv2dd   <none>           <none>
dra-gpu-example-6cf7c799b6-mmwp8   1/1     Running   0          45s   193.0.2.50   conformance-test-np1-sms5g-pt7t9-p42p6   <none>           <none>
dra-gpu-example-6cf7c799b6-zhmhp   0/1     Pending   0          45s   <none>       <none>                                   <none>           <none>
```

```shell
kubectl get resourceclaims -n gpu-test1 -o custom-columns=NAME:.metadata.name,STATE:.status.allocation.devices.results[0].pool,DEVICE:.status.allocation.devices.results[0].device,RESERVEDFOR:.status.reservedFor[0].name
```

```shell
NAME                                                STATE                                    DEVICE   RESERVEDFOR
dra-gpu-example-6cf7c799b6-6j66p-single-gpu-jcgx6   conformance-test-np1-sms5g-pt7t9-vv2dd   gpu-0    dra-gpu-example-6cf7c799b6-6j66p
dra-gpu-example-6cf7c799b6-mmwp8-single-gpu-5j9z4   conformance-test-np1-sms5g-pt7t9-p42p6   gpu-0    dra-gpu-example-6cf7c799b6-mmwp8
dra-gpu-example-6cf7c799b6-zhmhp-single-gpu-2j9c7   <none>                                   <none>   <none>
```

Each running pod sees only its own allocated device:

```shell
for p in $(kubectl get pods -n gpu-test1 -l app=dra-gpu-example --field-selector=status.phase=Running -o name); do echo "$p: $(kubectl logs -n gpu-test1 $p | grep GPU | tail -1)"; done
```

```shell
pod/dra-gpu-example-6cf7c799b6-6j66p: GPU 0: GRID A100-4C (UUID: GPU-60e223c0-015b-436c-ba00-26da00000000)
pod/dra-gpu-example-6cf7c799b6-mmwp8: GPU 0: GRID A100-4C (UUID: GPU-ba9adcc0-9809-47bb-85a6-b46a00000000)
```

The third claim cannot be satisfied, because both devices are already reserved for other pods — the scheduler refuses to place the pod rather than share an allocated device:

```shell
kubectl get events -n gpu-test1 --field-selector involvedObject.name=dra-gpu-example-6cf7c799b6-zhmhp,reason=FailedScheduling -o custom-columns=REASON:.reason,MESSAGE:.message
```

```shell
REASON             MESSAGE
FailedScheduling   0/3 nodes are available: 1 node(s) had untolerated taint(s), 2 cannot allocate all claims. preemption: 0/3 nodes are available: 3 Preemption is not helpful for scheduling.
```

This demonstrates that on VKS v1.37 each accelerator is allocated exclusively to a single claim/pod, that every pod sees only the device allocated to it, and that a workload without an allocation has no access to accelerator devices on the node.

### Cleanup

```shell
kubectl delete namespace gpu-test1
helm uninstall nvidia-dra-driver-gpu -n nvidia-dra-driver-gpu
kubectl delete namespace nvidia-dra-driver-gpu
kubectl delete crd computedomaincliques.resource.nvidia.com computedomains.resource.nvidia.com
kubectl get resourceslice -o name | xargs -r kubectl delete
```
