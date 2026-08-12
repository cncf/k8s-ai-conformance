Original URL: https://blogs.oracle.com/developers/verifying-gpu-accelerator-isolation-oci-kubernetes-engine-oke

# Verifying GPU Accelerator Access Isolation in OCI Kubernetes Engine (OKE)

The industry focus on AI and machine learning has led to the increased use of attached [devices](https://kubernetes.io/docs/reference/glossary/?all=true#term-device), infrastructure resources like hardware accelerators that are directly or indirectly attached to your nodes. In order to improve the user experience related to these devices, the Kubernetes community introduced [Dynamic Resource Allocation (DRA)](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/), which enables users to select, allocate, share, and configure GPUs, NICs and other devices. OCI Kubernetes Engine (OKE) users with clusters running v1.34 and above can make use of DRA in their Kubernetes clusters today. 

## Sharing Device Resources

Simply put, Kubernetes Dynamic Resource Allocation is a feature that lets you request and share attached device resources among pods. Allocating resources with DRA is similar to the [dynamic volume provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/) feature, which allows you to claim storage capacity from storage classes using PersistentVolumeClaims and to request the claimed capacity in your pods. In the case of DRA, device drivers and cluster administrators define device classes that are available for your workloads to claim. Kubernetes will then allocate devices matching claims and schedule pods requesting those claims on nodes with access the allocated devices. The DRA APIs graduated to stable in Kubernetes 1.34 and are now available by default. 

## Verifying Isolation

As with any cluster resource, it is important to prevent unauthorized access or interference between workloads. To maintain a secure posture for your cluster, you may want to verify that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework, in this case, DRA, and the container runtime. The steps below outline the verification process. 

### Setup

1. [Create an OKE cluster](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/create-cluster.htm) of v1.34 or above and deploy a node pool with GPU shapes, for example VM.GPU.A10.2. The DRA APIs are enabled by default on clusters v1.34 or above.

2. Create a dra-helm-values.yaml file which will be used to install the DRA drivers:
```
# Driver root
    nvidiaDriverRoot: "/"
    
    gpuResourcesEnabledOverride: true
    
    resources:
      gpus:
        enabled: true
      computeDomains:
        enabled: false
        
    kubeletPlugin:
      priorityClassName: ""
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
    
    affinity:
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:  
          - matchExpressions:
            # We allow a GPU deployment to be forced by setting the following label to "true"
            - key: "nvidia.com/gpu.present"
              operator: In
              values:
              - "true"
```
3. Using the dra-helm-values.yaml file above, install DRA drivers via helm:
```
helm install nvidia-dra-driver-gpu nvidia/nvidia-dra-driver-gpu --version="25.3.2" --create-namespace --namespace nvidia-dra-driver-gpu -f dra-helm-values.yaml
```
4. Validate that the DRA driver components are running and in a Ready state:
```
$ kubectl get pod -n nvidia-dra-driver-gpu
    
    NAME                                         READY   STATUS    RESTARTS   AGE
    nvidia-dra-driver-gpu-kubelet-plugin-2j5fm   1/1     Running   0          17m
    nvidia-dra-driver-gpu-kubelet-plugin-l7gpq   1/1     Running   0          15m
```
### Test 1
In this test you will deploy a pod to a node with available accelerators without requesting accelerator resources in the pod spec. Execute a command in the Pod to probe for accelerator devices. The command should fail or report that no accelerator devices are found.

1. Create a DRA claim:
```
apiVersion: resource.k8s.io/v1
    kind: ResourceClaimTemplate
    metadata:
      name: gpu-claim-template
    spec:
      spec:
        devices:
          requests:
          - name: single-gpu
            exactly:
              deviceClassName: gpu.nvidia.com
              allocationMode: ExactCount
              count: 1
```
2. Deploy a workload with a GPU claim:
```
apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: dra-gpu-example
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
```
**Note:** Note: The deployment may be stuck at container creation (or terminate if you are trying to delete the deployment) if the DRA driver isn’t used within 30 minutes of installation, caused by https://github.com/kubernetes/kubernetes/issues/133920. If this happens, the workaround is to restart the kubelet manually. 

3. Resource claim results:
```
$ kubectl get resourceclaims
    NAME                                                STATE                AGE
    dra-gpu-example-68f595d7dc-vxqf4-single-gpu-lksr9   allocated,reserved   69s
```
4. NVIDIA DRA driver results:
```
$ kubectl get pods -n nvidia-dra-driver-gpu
    
    NAME                                         READY   STATUS    RESTARTS   AGE
    nvidia-dra-driver-gpu-kubelet-plugin-t48f8   1/1     Running   0          87m
```
5. Successful deployment:
```
$ kubectl get pods -l app=dra-gpu-example
    NAME                               READY   STATUS    RESTARTS   AGE
    dra-gpu-example-68f595d7dc-5vxhm   1/1     Running   0          3m8s
```
6. Demonstrate that the pod can no longer access the accelerator by removing the resource claim from the deployment:
```
apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: dra-gpu-example
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
            #resources:
            #  claims:
            #  - name: single-gpu
          resourceClaims:
          - name: single-gpu
            resourceClaimTemplateName: gpu-claim-template
          tolerations:
          - key: "nvidia.com/gpu"
            operator: "Exists"
            effect: "NoSchedule"
```
7. View the failure to access the GPU:
```
$  kubectl logs dra-gpu-example-5c97694b59-phvgk
    Wed Oct 22 22:31:50 UTC 2025
    bash: line 1: nvidia-smi: command not found
    Waiting...
```

### Test 2
1. Create two pods, each is allocated an accelerator resource.
2. Execute a command in one Pod to attempt to access the other Pod’s accelerator, and should be denied. This can be verified by running this test https://github.com/kubernetes/kubernetes/blob/v1.34.1/test/e2e/dra/dra.go#L180 This is part of the Kubernetes e2e tests specifically to test the DRA functionality.
3. **(IMPORTANT)** Remove the GPU taints first to allow the test to execute:
```
$ kubectl taint nodes <name> nvidia.com/gpu:NoSchedule-
```
4. Clone the Kubernetes GitHub repository:
```
$ git clone https://github.com/kubernetes/kubernetes.git
```
5. Execute the DRA test:
```
$ make WHAT="github.com/onsi/ginkgo/v2/ginkgo k8s.io/kubernetes/test/e2e/e2e.test"
    $ KUBERNETES_PROVIDER="local" hack/ginkgo-e2e.sh --provider=skeleton --ginkgo.focus='must map configs and devices to the right containers'
    Setting up for KUBERNETES_PROVIDER="local".
    Skeleton Provider: prepare-e2e not implemented
    KUBE_MASTER_IP:
    KUBE_MASTER:
      I1022 17:41:11.119534   75247 e2e.go:109] Starting e2e run "3f4a6964-e56d-4229-8dc9-7e16279d1094" on Ginkgo node 1
    Running Suite: Kubernetes e2e suite - /Users/danielberg/development/github/kubernetes/_output/bin
    =================================================================================================
    Random Seed: 1761169270 - will randomize all specs
    
    Will run 1 of 7206 specs
    •
    
    Ran 1 of 7206 Specs in 38.214 seconds
    SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 7205 Skipped
    PASS
    
    Ginkgo ran 1 suite in 38.535845792s
    Test Suite Passed
```
## Conclusion 
DRA is a powerful way to enable your workloads to make use of attached devices including hardware accelerators. To prevent against unauthorized access or interference between workloads, it is important to be able to verify proper isolation of your attached devices. 
