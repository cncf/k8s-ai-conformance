# Secure Accelerator Access on Talos Linux

**MUST**: Ensure that access to accelerators from within containers is
properly isolated and mediated by the Kubernetes resource management framework
(device plugin or DRA) and container runtime, preventing unauthorized access
or interference between workloads.

## Automated verification

This requirement is verified end-to-end by the upstream
[AI Conformance test suite](https://github.com/kubernetes-sigs/ai-conformance)
in the same submission directory:

- **`junit.xml#TestSecureAcceleratorAccess`** — machine-readable JUnit results
- **`e2e.log`** — full human-readable log
- **`results.json`** — Go-test JSON events for CI ingestion

The test suite auto-detects the allocation mode (DRA vs device-plugin) and
runs both positive and negative sub-tests:

| Sub-test | What it proves |
| --- | --- |
| `TestSecureAcceleratorAccess/PositiveAccessTest` | A pod that requests one accelerator via the platform's resource-management framework (DRA `ResourceClaim` on this cluster) is scheduled and observes exactly one accelerator device from inside the container. |
| `TestSecureAcceleratorAccess/NegativeIsolationTest` | A pod that does **not** request an accelerator is scheduled but observes **zero** accelerator devices — no device files leak in, no driver access is granted. |

Both sub-tests passed against the two-node Talos Linux v1.14.1 / Kubernetes
v1.37.0 cluster described in this submission.

## Cluster configuration verified

The submission was collected on:

- **Control plane**: Raspberry Pi 5 (arm64, no accelerator)
- **Worker**: NVIDIA DGX Spark GB10 (arm64, single GPU exposed as
  `gpu.nvidia.com` via the [NVIDIA DRA driver](https://github.com/NVIDIA/k8s-dra-driver-gpu))

Isolation on Talos Linux is enforced by three cooperating layers:

1. **The signed `nvidia-container-toolkit-production` system extension**
   provides the runtime hooks that mount `/dev/nvidia*` devices only when a
   Container Device Interface (CDI) claim exists for the container.
2. **The NVIDIA DRA driver** allocates individual GPU devices to pods that
   hold a `ResourceClaim` and rejects overlapping claims from unrelated pods.
3. **Talos Linux itself** is immutable and configured entirely through
   declarative machine config; there is no host-level SSH or shell to bypass
   the container runtime and manually attach a GPU to a pod.

Together these make it impossible for a pod without a DRA claim to see or
touch a GPU device on this cluster, which is what the automated negative
isolation test verifies.
