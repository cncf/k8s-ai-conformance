## Description

Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

## Evidence

### Setup and Documentation

The VKE cluster is running Kubernetes v1.36. The NVIDIA GPU DRA (Dynamic Resource Allocation) driver is installed and configured according to the official Volcengine documentation: [Dynamic Resource Allocation (DRA) - Volcengine Documentation](https://www.volcengine.com/docs/6460/2353389?lang=zh).

### Secure Accelerator Access Test

The AI Conformance test suite validates that accelerator devices are correctly mapped only to the containers that request them, ensuring proper isolation and secure access.

#### Test Execution

The test is executed using the `ai-conformance` test suite:

```shell
# Clone the AI Conformance repository
git clone https://github.com/kubernetes-sigs/ai-conformance.git
cd ai-conformance

# Run the Secure Accelerator Access test
go test -v ./test -run 'Test(SecureAcceleratorAccess)$'
```

#### Test Results

Execution output from the VKE v1.36 validation environment:

```text
➜  ai-conformance git:(main) ✗ go test -v ./test -run 'Test(SecureAcceleratorAccess)$'
=== RUN   TestSecureAcceleratorAccess
    util_test.go:325: Checking environment: Found ResourceSlice: 192.*.*.*-gpu.nvidia.com-ncfkn (Node: 192.*.*.*, Driver: gpu.nvidia.com, Devices: 1)
    util_test.go:213: Auto-detected allocation mode: dra
    util_test.go:180: Secure accelerator access test running with allocation mode: dra (accelerator node: 192.*.*.*)
=== RUN   TestSecureAcceleratorAccess/PositiveAccessTest
    util_test.go:557: Waiting for Pod pos-pod to be running...
    util_test.go:737: Waiting to see if Pod pos-pod/prober logs contain 'RESULT: ACCELERATOR_FOUND'...
    util_test.go:751: PASS: prober isolation/access verified.
=== RUN   TestSecureAcceleratorAccess/NegativeIsolationTest
    util_test.go:557: Waiting for Pod neg-pod to be running...
    util_test.go:737: Waiting to see if Pod neg-pod/prober logs contain 'RESULT: ACCELERATOR_MISSING'...
    util_test.go:751: PASS: prober isolation/access verified.
=== RUN   TestSecureAcceleratorAccess/MultiContainerIsolationTest
    util_test.go:557: Waiting for Pod multi-container-pod to be running...
    util_test.go:737: Waiting to see if Pod multi-container-pod/authorized logs contain 'RESULT: ACCELERATOR_FOUND'...
    util_test.go:751: PASS: authorized isolation/access verified.
    util_test.go:737: Waiting to see if Pod multi-container-pod/unauthorized logs contain 'RESULT: ACCELERATOR_MISSING'...
    util_test.go:751: PASS: unauthorized isolation/access verified.
=== NAME  TestSecureAcceleratorAccess
    ai_conformance_test.go:44: Cleaning up namespace ai-conformance-548fw...
--- PASS: TestSecureAcceleratorAccess (149.61s)
    --- PASS: TestSecureAcceleratorAccess/PositiveAccessTest (74.34s)
    --- PASS: TestSecureAcceleratorAccess/NegativeIsolationTest (34.34s)
    --- PASS: TestSecureAcceleratorAccess/MultiContainerIsolationTest (34.57s)
PASS
ok      github.com/kubernetes-sigs/ai-conformance/test  152.738s
```

This test result demonstrates that accelerator access in VKE is mediated by the Kubernetes DRA framework and the container runtime, preventing unauthorized access between workloads.
