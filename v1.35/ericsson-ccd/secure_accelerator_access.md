# Secure Accelerator Access Validation

## Overview

This guide covers the validation of secure accelerator access using the NVIDIA Dynamic Resource Allocation (DRA) plugin on Ericsson Cloud Container Distribution (ECCD) 2.34.0.

DRA is a Kubernetes feature that enables fine-grained, claim-based allocation of hardware resources such as GPUs to workloads. This guide validates that resource isolation is enforced correctly — ensuring that pods cannot access accelerator devices that have not been explicitly allocated to them.


## Step 1: Install the required drivers, GPU Operator and DRA Plugin

See Step 1 in DRA_plugin.md



## Tests

### Test 1: Verify Pod Cannot Access Unallocated GPU Resources

#### 1. Create the testing namespace and a DRA ResourceClaimTemplate with Accelerator device

```bash
cat <<EOF | envsubst | kubectl apply -f -
---
apiVersion: v1
kind: Namespace
metadata:
  name: testing
---
apiVersion: resource.k8s.io/v1
kind: ResourceClaimTemplate
metadata:
  name: single-gpu
  namespace: testing
spec:
  spec:
    devices:
      requests:
      - name: req-0
        exactly:
          deviceClassName: gpu.nvidia.com
          allocationMode: ExactCount
          count: 1
EOF
```

#### 2. Create a test application claiming the resource

```bash
export IMAGE_TAG=$(curl -s -q https://registry.eccd.local:5000/v2/httpd/tags/list | jq '.tags[]' -r | grep -v 'sha256')

cat <<EOF | envsubst | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: testing
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: test-app
          image: registry.eccd.local:5000/httpd:${IMAGE_TAG}
          command: ["/bin/sh", "-c"]
          args:
            - |
              while true; do
                nvidia-smi -L
                sleep 30
              done
          resources:
            claims:
              - name: single-gpu
      resourceClaims:
        - name: single-gpu
          resourceClaimTemplateName: single-gpu
EOF
```

#### 3. Check the resource claim

```bash
kubectl get resourceclaim -n testing
NAME                                         STATE                AGE
test-app-545fdccb78-96w7r-single-gpu-nfmsd   allocated,reserved   6m23s
```

#### 4. Check that the Accelerator is accessible in the pod.

```bash
#Logs show the nvidia-smi tool listing the device
kubectl logs -n testing test-app-545fdccb78-96w7r
GPU 0: Tesla T4 (UUID: GPU-6d5c0ac5-74d7-3130-46f1-e93ed57f469a)

#Devices are also mounted in the pod
kubectl exec -it -n testing test-app-545fdccb78-96w7r  -- ls -la /dev/nvidia*
crw-rw-rw- 1 root root 195, 254 Feb 27 16:36 /dev/nvidia-modeset
crw-rw-rw- 1 root root 507,   0 Feb 27 16:36 /dev/nvidia-uvm
crw-rw-rw- 1 root root 507,   1 Feb 27 16:36 /dev/nvidia-uvm-tools
crw-rw-rw- 1 root root 195,   0 Feb 27 16:36 /dev/nvidia0
crw-rw-rw- 1 root root 195, 255 Feb 27 16:36 /dev/nvidiactl

/dev/nvidia-caps:
total 0
drwxr-xr-x 2 root root     80 Feb 27 16:36 .
drwxr-xr-x 7 root root    500 Feb 27 16:36 ..
cr-------- 1 root root 510, 1 Feb 27 16:36 nvidia-cap1
cr--r--r-- 1 root root 510, 2 Feb 27 16:36 nvidia-cap2
```

#### 5. Delete the deployment and apply it without the resource claim

```bash
#Delete the deployment
kubectl delete deployment -n testing test-app

#Create the deployment without resource claims
cat <<EOF | envsubst | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
  namespace: testing
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
        - name: test-app
          image: registry.eccd.local:5000/httpd:${IMAGE_TAG}
          command: ["/bin/sh", "-c"]
          args:
            - |
              while true; do
                nvidia-smi -L
                sleep 30
              done
         #resources:
         #   claims:
         #     - name: single-gpu
      resourceClaims:
        - name: single-gpu
          resourceClaimTemplateName: single-gpu
EOF

deployment.apps/test-app configured
```

#### 6. Verify that the pod cannot access the accelerator

```bash
#Get the pod
kubectl get pods -n testing
NAME                       READY   STATUS    RESTARTS   AGE
test-app-86d7dbc84-bzh6v   1/1     Running   0          3s

#Check the pod logs
kubectl logs -n testing test-app-86d7dbc84-bzh6v
/bin/sh: line 1: nvidia-smi: command not found

#Ensure that the container cannot access the GPU
kubectl exec -it -n testing test-app-86d7dbc84-bzh6v  -- ls -la /dev/nvidia*
ls: cannot access '/dev/nvidia0': No such file or directory
ls: cannot access '/dev/nvidia-caps': No such file or directory
ls: cannot access '/dev/nvidiactl': No such file or directory
ls: cannot access '/dev/nvidia-modeset': No such file or directory
ls: cannot access '/dev/nvidia-uvm': No such file or directory
ls: cannot access '/dev/nvidia-uvm-tools': No such file or directory
command terminated with exit code 2
```
The logs show that both the nvidia-smi tool and the GPU devices are not mounted and are not accessible from the container.



### Test 2: Verify Pod Cannot Access GPU Devices Allocated to Another Pod

This scenario is covered by the Kubernetes end-to-end tests for DRA at (https://github.com/kubernetes/kubernetes/blob/v1.35.0/test/e2e/dra/dra.go#L331). Successful execution of this test on the platform confirms that GPU isolation between pods is enforced correctly.

#### 1. Clone kubernetes repo and checkout 1.35.0

```bash
#Install needed packages
sudo zypper in git go

#Clone the kubernetes repo and checkout v1.35.0
git clone https://github.com/kubernetes/kubernetes
cd kubernetes
git checkout v1.35.0
```

#### 2. Execute the e2e test suite for DRA and ensure that it passes

```bash
#Set kubeconfig to the local cluster
export KUBECONFIG=~/.kube/config

#Run the e2e test for DRA

go run github.com/onsi/ginkgo/v2/ginkgo -focus='must map configs and devices to the right containers' ./test/e2e
  I0227 17:30:19.018465 2453969 test_context.go:564] The --provider flag is not set. Continuing as if --provider=skeleton had been used.
  I0227 17:30:19.018631 2453969 e2e.go:109] Starting e2e run "5682ad31-968e-4c10-ac7e-758010f4c5d6" on Ginkgo node 1
Running Suite: Kubernetes e2e suite - /home/eccd/echekjo/kubernetes/test/e2e
============================================================================
Random Seed: 1772213405 - will randomize all specs

Will run 1 of 7348 specs
SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS
<truncated>

Ran 1 of 7348 Specs in 19.903 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 7347 Skipped
PASS

Ginkgo ran 1 suite in 34.006434924s
Test Suite Passed
```