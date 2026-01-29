# Alibaba Cloud Container Service for Kubernetes — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Alibaba Cloud |
| **Platform** | Alibaba Cloud Container Service for Kubernetes |
| **Platform Version** | 1.34.1-aliyun.1 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://www.alibabacloud.com/en/product/kubernetes](https://www.alibabacloud.com/en/product/kubernetes) |
| **Documentation** | [Link](https://www.alibabacloud.com/help/en/ack/) |

> Alibaba Cloud Container Service for Kubernetes (ACK) provides high-performance management services for containerized applications. You can manage enterprise-level containerized applications throughout the application lifecycle. This service allows you to run containerized applications in the cloud in an efficient manner.

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

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/scheduling-gpu-using-dra)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [www.alibabacloud.com/help/en/cs/user-guide/gateway-with-i...](https://www.alibabacloud.com/help/en/cs/user-guide/gateway-with-inference-extension-overview)
- [www.alibabacloud.com/help/en/cs/user-guide/generative-ai-...](https://www.alibabacloud.com/help/en/cs/user-guide/generative-ai-service-enhancement/)
- [www.alibabacloud.com/help/en/asm/sidecar/use-gateway-api-...](https://www.alibabacloud.com/help/en/asm/sidecar/use-gateway-api-to-define-routing-rules)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/work-with-gang-scheduling)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/overview-of-node-scaling/)
- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/configure-automatic-node-scaling-for-gpu-applications)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/enable-auto-scaling-based-on-gpu-metrics)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/observability-overview)
- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/gpu-monitoring/)
- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/use-alibaba-cloud-prometheus-service-to-monitor-an-ack-cluster)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/observability-overview)
- [www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/us...](https://www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/user-guide/configure-monitoring-for-llm-inference-services)
- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/user-guide/use-alibaba-cloud-prometheus-service-to-monitor-an-ack-cluster)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [github.com/AliyunContainerService/ai-models-on-ack/tree/m...](https://github.com/AliyunContainerService/ai-models-on-ack/tree/main/ai-conformance/v1.34/secure_accelerator_access)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/us...](https://www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/use-cases/ray-cluster-best-practices)
- [www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedi...](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/use-cases/run-apache-spark-workloads-on-ack)
- [www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/us...](https://www.alibabacloud.com/help/en/ack/cloud-native-ai-suite/user-guide/overview-of-ai-jobs)

---

*Generated from PRODUCT.yaml*
