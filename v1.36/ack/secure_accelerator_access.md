## Description

Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

#### Prerequisites

- An ACK managed cluster running Kubernetes 1.34 or later, with `kubectl` configured to access the cluster.
- A GPU-accelerated node pool provisioned in the cluster.
- The `ack.node.gpu.schedule: disabled` label applied to GPU nodes so the default device plugin does not report hardware, preventing duplicate allocation when DRA takes over.

References:

- https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/scheduling-gpu-using-dra

#### Install the NVIDIA DRA Driver

Add the NVIDIA Helm repository to initiate the driver installation:

```shell
helm repo add nvidia https://helm.ngc.nvidia.com/nvidia && helm repo update
```

Install the driver:

```shell
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu --version="25.3.2" --create-namespace --namespace nvidia-dra-driver-gpu \
    --set gpuResourcesEnabledOverride=true \
    --set controller.affinity=null \
    --set "kubeletPlugin.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0].key=ack.node.gpu.schedule" \
    --set "kubeletPlugin.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0].operator=In" \
    --set "kubeletPlugin.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0].values[0]=disabled"
```

```shell
NAME: nvidia-dra-driver-gpu
LAST DEPLOYED: Wed Aug  5 19:30:27 2026
NAMESPACE: nvidia-dra-driver-gpu
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

The installation creates the `nvidia-dra-driver-gpu` namespace. Verify that the DRA controller and kubelet-plugin pods are running:

```shell
kubectl get pods -n nvidia-dra-driver-gpu
```

```shell
NAME                                                READY   STATUS    RESTARTS   AGE
nvidia-dra-driver-gpu-controller-7d67d6566c-pw4bl   1/1     Running   0          44m
nvidia-dra-driver-gpu-kubelet-plugin-jv9ms          2/2     Running   0          33m
```

Confirm that the driver has registered its `DeviceClass` and `ResourceSlice` objects:

```shell
kubectl get deviceclass,resourceslice
```

```shell
NAME                                                                    AGE
deviceclass.resource.k8s.io/compute-domain-daemon.nvidia.com            46m
deviceclass.resource.k8s.io/compute-domain-default-channel.nvidia.com   46m
deviceclass.resource.k8s.io/gpu.nvidia.com                              46m
deviceclass.resource.k8s.io/mig.nvidia.com                              46m

NAME                                                                                       NODE                         DRIVER                      POOL                         AGE
resourceslice.resource.k8s.io/cn-shanghai.10.xxx.xxx.xxx-compute-domain.nvidia.com-nfngn   cn-shanghai.10.xxx.xxx.xxx   compute-domain.nvidia.com   cn-shanghai.10.xxx.xxx.xxx   33m
resourceslice.resource.k8s.io/cn-shanghai.10.xxx.xxx.xxx-gpu.nvidia.com-kzrd8              cn-shanghai.10.xxx.xxx.xxx   gpu.nvidia.com              cn-shanghai.10.xxx.xxx.xxx   33m
```

#### Run the E2E Test

Clone the ai-conformance repository and check out the `main` branch:

```shell
mkdir workspace && cd workspace && git clone https://github.com/kubernetes-sigs/ai-conformance
cd ai-conformance
git checkout main

git rev-parse HEAD
# c0b0b164ecc3fb16540506aa74795e6e7b63e344
```

Run the secure accelerator access test:

```shell
$ go test -v ./test -run 'Test(SecureAcceleratorAccess)$'

=== RUN   TestSecureAcceleratorAccess
    util_test.go:296: Checking environment: Found ResourceSlice: cn-shanghai.10.xxx.xxx.xxx-gpu.nvidia.com-kzrd8 (Node: cn-shanghai.10.xxx.xxx.xxx, Driver: gpu.nvidia.com, Devices: 1)
    util_test.go:184: Auto-detected allocation mode: dra
    util_test.go:151: Secure accelerator access test running with allocation mode: dra (accelerator node: cn-shanghai.10.xxx.xxx.xxx)
=== RUN   TestSecureAcceleratorAccess/PositiveAccessTest
    util_test.go:528: Waiting for Pod pos-pod to be running...
    util_test.go:708: Waiting to see if Pod pos-pod/prober logs contain 'RESULT: ACCELERATOR_FOUND'...
    util_test.go:722: PASS: prober isolation/access verified.
=== RUN   TestSecureAcceleratorAccess/NegativeIsolationTest
    util_test.go:528: Waiting for Pod neg-pod to be running...
    util_test.go:708: Waiting to see if Pod neg-pod/prober logs contain 'RESULT: ACCELERATOR_MISSING'...
    util_test.go:722: PASS: prober isolation/access verified.
=== RUN   TestSecureAcceleratorAccess/MultiContainerIsolationTest
    util_test.go:528: Waiting for Pod multi-container-pod to be running...
    util_test.go:708: Waiting to see if Pod multi-container-pod/authorized logs contain 'RESULT: ACCELERATOR_FOUND'...
    util_test.go:722: PASS: authorized isolation/access verified.
    util_test.go:708: Waiting to see if Pod multi-container-pod/unauthorized logs contain 'RESULT: ACCELERATOR_MISSING'...
    util_test.go:722: PASS: unauthorized isolation/access verified.
=== NAME  TestSecureAcceleratorAccess
    ai_conformance_test.go:44: Cleaning up namespace ai-conformance-c5ll6...
--- PASS: TestSecureAcceleratorAccess (115.43s)
    --- PASS: TestSecureAcceleratorAccess/PositiveAccessTest (36.34s)
    --- PASS: TestSecureAcceleratorAccess/NegativeIsolationTest (34.32s)
    --- PASS: TestSecureAcceleratorAccess/MultiContainerIsolationTest (36.36s)
PASS
ok  	github.com/kubernetes-sigs/ai-conformance/test	116.378s
```
