# Amazon Elastic Kubernetes Service (EKS) — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Amazon Web Services (AWS) |
| **Platform** | Amazon Elastic Kubernetes Service (EKS) |
| **Platform Version** | 1.34.1-eks.4 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://aws.amazon.com/eks/](https://aws.amazon.com/eks/) |
| **Documentation** | [Link](https://docs.aws.amazon.com/eks/latest/userguide/what-is-eks.html) |

> Amazon EKS is a managed Kubernetes service that you can use to build, run, and scale production-ready Kubernetes applications across any environment

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 9 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | MUST | ✅ |
| Networking | AI Inference | MUST | ✅ |
| Scheduling & Orchestration | Gang Scheduling | MUST | ✅ |
| Scheduling & Orchestration | Cluster Autoscaling | MUST | ✅ |
| Scheduling & Orchestration | Pod Autoscaling | MUST | ✅ |
| Observability | Accelerator Metrics | MUST | ✅ |
| Observability | AI Service Metrics | MUST | ✅ |
| Security | Secure Accelerator Access | MUST | ✅ |
| Operator Support | Robust Controller | MUST | ✅ |

---

## Detailed Requirements

### 🚀 Accelerators

#### ✅ DRA Support

**Level:** 🔴 MUST | **Status:** Implemented

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

**Evidence:**

- [docs.aws.amazon.com/eks/latest/best-practices/aiml-comput...](https://docs.aws.amazon.com/eks/latest/best-practices/aiml-compute.html#aiml-dra)

**Notes:**

> DRA is supported in EKS for Kubernetes versions 1.33 and above. EKS launched support for Kubernetes 1.33 in May 2025 and support for Kubernetes 1.34 in October 2025. In EKS, DRA is enabled by default on the EKS control plane and in the EKS-optimized AL2023 and Bottlerocket AMIs.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [kubernetes-sigs.github.io/aws-load-balancer-controller/v2...](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.14/guide/gateway/l7gateway/)
- [docs.aws.amazon.com/elasticloadbalancing/latest/applicati...](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html)

**Notes:**

> EKS supports Kubernetes Gateway API and Ingress through the AWS Load Balancer Controller, which manages AWS Application Load Balancers (L7) and Network Load Balancers (L4).

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [aws.amazon.com/blogs/hpc/scaling-your-llm-inference-workl...](https://aws.amazon.com/blogs/hpc/scaling-your-llm-inference-workloads-multi-node-deployment-with-tensorrt-llm-and-triton-on-amazon-eks/)
- [awslabs.github.io/ai-on-eks/docs/infra/trainium#volcano-s...](https://awslabs.github.io/ai-on-eks/docs/infra/trainium#volcano-scheduler)
- [builder.aws.com/content/2zsUDYed1Df7zHHZFc11MRctb6M/fract...](https://builder.aws.com/content/2zsUDYed1Df7zHHZFc11MRctb6M/fractional-gpu-sharing-on-amazon-eks-using-nvidia-kai-scheduler-and-accelerated-ec2-instances)
- [aws.amazon.com/blogs/big-data/deploy-apache-yunikorn-batc...](https://aws.amazon.com/blogs/big-data/deploy-apache-yunikorn-batch-scheduler-for-amazon-emr-on-eks/)
- [aws.amazon.com/blogs/hpc/gang-scheduling-pods-on-amazon-e...](https://aws.amazon.com/blogs/hpc/gang-scheduling-pods-on-amazon-eks-using-aws-batch-multi-node-processing-jobs/)

**Notes:**

> EKS supports using custom schedulers, LeaderWorkerSet, as well as integrations with AWS Batch and Amazon Sagemaker for gang scheduling.

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.aws.amazon.com/eks/latest/userguide/auto-accelerated...](https://docs.aws.amazon.com/eks/latest/userguide/auto-accelerated.html)
- [aws.amazon.com/blogs/containers/scaling-a-large-language-...](https://aws.amazon.com/blogs/containers/scaling-a-large-language-model-with-nvidia-nim-on-amazon-eks-with-karpenter/)

**Notes:**

> EKS supports Karpenter, which automatically scales nodes up/down for pending pods based on node selection constraints. When the Karpenter node selection constraints in a NodePool define instances with accelerators, Karpenter scales nodes up/down based on the pending pods requesting accelerators. Karpenter is used in EKS Auto Mode, an EKS feature where the cluster’s nodes, networking, storage, and load balancing are fully managed and scaled by EKS. EKS also supports Cluster Autoscaler.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [github.com/awslabs/ai-on-eks/tree/main/ai-conformance/1.3...](https://github.com/awslabs/ai-on-eks/tree/main/ai-conformance/1.34/pod_autoscaling)

**Notes:**

> EKS supports pod autoscalers including vanilla HPA and KEDA with HPA. HPA can be configured to scale based on custom metrics, such as those exposed by NVIDIA DCGM exporter or Neuron Monitor. The example linked in the evidence shows how to use HPA with the DCGM_FI_DEV_GPU_UTIL metric from DCGM on EKS.

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.aws.amazon.com/eks/latest/best-practices/aiml-observ...](https://docs.aws.amazon.com/eks/latest/best-practices/aiml-observability.html)
- [aws.amazon.com/blogs/containers/part-1-introduction-to-ob...](https://aws.amazon.com/blogs/containers/part-1-introduction-to-observing-machine-learning-workloads-on-amazon-eks/)

**Notes:**

> EKS supports collecting accelerator metrics through DCGM (NVIDIA) and Nueron Monitor (Nueron), and integrates these solutions with CloudWatch Observability Agent and AWS Distro for Open Telemetry for export and CloudWatch and Grafana with Prometheus for visualization. With these solutions, users can measure GPU, network, storage, and memory utilization and can inspect these metrics at varying Kubernetes resource levels.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.aws.amazon.com/eks/latest/best-practices/aiml-observ...](https://docs.aws.amazon.com/eks/latest/best-practices/aiml-observability.html)
- [aws.amazon.com/blogs/containers/part-1-introduction-to-ob...](https://aws.amazon.com/blogs/containers/part-1-introduction-to-observing-machine-learning-workloads-on-amazon-eks/)

**Notes:**

> EKS integrates with Amazon CloudWatch, open source Prometheus, and Amazon Managed Prometheus. The CloudWatch Observability Agent can be easily installed in EKS clusters via EKS add-ons to scrape node and application-level metrics and export them to CloudWatch metrics dashboards. CloudWatch Container Insights can be used with these metrics for out-of-the-box deep analysis. The AWS Distro for Open Telemetry add-on can be easily installed in EKS clusters via EKS add-ons to scrape node and application-level metrics and export them to open source Prometheus or Amazon Managed Prometheus.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [github.com/awslabs/ai-on-eks/tree/main/ai-conformance/1.3...](https://github.com/awslabs/ai-on-eks/tree/main/ai-conformance/1.34/secure_accelerator_access)

**Notes:**

> See EKS-Conformance-Test-1.md and EKS-Conformance-Test-2.md in the linked resource for steps.

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [awslabs.github.io/ai-on-eks/docs/blueprints/inference](https://awslabs.github.io/ai-on-eks/docs/blueprints/inference)

**Notes:**

> EKS supports a breadth of inference serving frameworks including Ray, vLLM, KubeFlow, AIBrix, NVIDIA Triton. The linked resource contains infrastructure-as-code templates to deploy these solutions on EKS.

---

*Generated from PRODUCT.yaml*
