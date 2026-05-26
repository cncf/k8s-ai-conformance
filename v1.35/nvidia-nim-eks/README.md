# NVIDIA NIM on EKS

[NVIDIA NIM](https://developer.nvidia.com/nim) on Amazon EKS is a Kubernetes-based AI inference platform that deploys and manages NVIDIA NIM microservices on EKS with GPU scheduling, autoscaling, and Gateway API integration. NIM microservice lifecycle is managed by the [NIM Operator](https://github.com/NVIDIA/k8s-nim-operator).

## Conformance Submission

- [PRODUCT.yaml](PRODUCT.yaml)

## Layered Product

NIM on EKS is a layered AI inference product on top of an existing Kubernetes-conformant distribution (Amazon EKS). This submission follows the layered-product guidance in [faq.md](../../faq.md):

- **Base platform.** Amazon EKS v1.35 — Kubernetes conformance entry at [cncf/k8s-conformance/v1.35/eks](https://github.com/cncf/k8s-conformance/tree/master/v1.35/eks). EKS provides the prerequisite plain-Kubernetes conformance (control plane, networking, storage, accelerator enablement on managed node groups, the DRA API at GA in v1.35, and the Cluster Autoscaler integration).
- **Layered runtime.** NVIDIA NIM Operator (`1.8.3`), KAI Scheduler, NVIDIA GPU Operator (driver, container toolkit, device plugin, DCGM Exporter, MIG manager, vGPU), NVIDIA DRA driver, kgateway with inference extension CRDs, Prometheus + prometheus-adapter.

Per-requirement layer ownership is recorded in each entry's `notes` field in `PRODUCT.yaml`.

### What this layered product uniquely contributes

| Requirement | Contribution kind | Layered component |
|---|---|---|
| `gang_scheduling` | MUST — meaningfully different variant | KAI Scheduler (vs. Kueue/Volcano; called out in [faq.md](../../faq.md) as an example of a layered product's unique contribution) |
| `ai_inference` | MUST — meaningfully different variant | kgateway with inference extension CRDs (InferencePool, InferenceModelRewrite, InferenceObjective) on top of the base Gateway API |
| `ai_service_metrics` | MUST — meaningfully different variant | NIM inference microservice metrics (token throughput, time-to-first-token, time-per-output-token) on top of generic Prometheus scraping |
| `robust_controller` | MUST — implements | NIM Operator (9 CRDs, admission webhook, NIMService reconcile) — an AI-inference-specific operator beyond the Ray/Kubeflow examples in the spec |
| `accelerator_metrics` | MUST — implements | DCGM Exporter (the layered component that produces per-GPU Prometheus metrics) |
| `pod_autoscaling` | MUST — implements | Prometheus adapter wiring GPU custom metrics into the K8s custom metrics API for HPA |
| `dra_support` | MUST — implements | NVIDIA DRA driver advertising H100 devices via ResourceSlices |
| `secure_accelerator_access` | MUST — implements | GPU Operator + DRA-mediated device isolation |
| `driver_runtime_management` | SHOULD — addition | NVIDIA GPU Operator (full driver/runtime/device-plugin lifecycle) |
| `gpu_sharing` | SHOULD — addition | GPU Operator MIG (A100/H100 hardware partitioning) + time-slicing |
| `virtualized_accelerator` | SHOULD — addition | GPU Operator vGPU deployments |

### Roadmap features

These are on the AI Conformance roadmap and are anticipated for inclusion when the underlying platform mechanisms reach GA:

- **Disaggregated inference** — [NVIDIA Dynamo](https://github.com/ai-dynamo/dynamo) for prefill/decode disaggregation. Called out in [faq.md](../../faq.md) as the canonical example of a meaningfully different variant of disaggregated inference (vs. llm-d). Will be submitted under the same layered-product framework when the requirement formalizes.
- **DRA partitionable devices** — static and dynamic GPU sharing exposed through DRA once partitionable devices is GA. The layered NVIDIA DRA driver is already in place and positions this product to adopt the DRA mechanism without architectural changes.

### Re-test confirmation

Per [faq.md](../../faq.md): a layered product must re-test for full AI Conformance to confirm that the layered components did not break the base platform's existing Kubernetes Conformance behaviors. All 9 MUST requirements were evidenced on the layered cluster (EKS v1.35 with the runtime stack listed above), not on a bare EKS cluster, on 8x NVIDIA H100 80GB HBM3 GPUs running an NVIDIA NIM `Llama 3.2 1B` inference workload via `NIMService` CR. The 3 SHOULD requirements (`driver_runtime_management`, `gpu_sharing`, `virtualized_accelerator`) are implemented through the NVIDIA GPU Operator on the same cluster; `driver_runtime_management` has cluster evidence, and `gpu_sharing` / `virtualized_accelerator` are supported by GPU Operator documentation.

## Evidence

Full per-requirement evidence is hosted in the [AICR repository](https://github.com/NVIDIA/aicr/tree/main/docs/conformance/cncf/v1.35/nim-eks/evidence). Each `PRODUCT.yaml` entry links the corresponding evidence document.

| # | Requirement | Layer / Component | Result | Evidence |
|---|---|---|---|---|
| 1 | `dra_support` | Layered: NVIDIA DRA driver | PASS | [dra-support.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/dra-support.md) |
| 2 | `gang_scheduling` | Layered: KAI Scheduler | PASS | [gang-scheduling.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/gang-scheduling.md) |
| 3 | `secure_accelerator_access` | Layered: GPU Operator + DRA | PASS | [secure-accelerator-access.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/secure-accelerator-access.md) |
| 4 | `accelerator_metrics` | Layered: DCGM Exporter | PASS | [accelerator-metrics.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/accelerator-metrics.md) |
| 5 | `ai_service_metrics` | Layered: NIM inference + Prometheus | PASS | [ai-service-metrics.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/ai-service-metrics.md) |
| 6 | `ai_inference` | Layered: kgateway + inference CRDs | PASS | [inference-gateway.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/inference-gateway.md) |
| 7 | `robust_controller` | Layered: NIM Operator | PASS | [robust-operator.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/robust-operator.md) |
| 8 | `pod_autoscaling` | Layered: Prometheus adapter + HPA | PASS | [pod-autoscaling.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/pod-autoscaling.md) |
| 9 | `cluster_autoscaling` | Base: EKS Cluster Autoscaler | PASS | [cluster-autoscaling.md](https://github.com/NVIDIA/aicr/blob/main/docs/conformance/cncf/v1.35/nim-eks/evidence/cluster-autoscaling.md) |

All 9 MUST conformance requirements are **Implemented**. 3 SHOULD requirements (`driver_runtime_management`, `gpu_sharing`, `virtualized_accelerator`) are also Implemented.
