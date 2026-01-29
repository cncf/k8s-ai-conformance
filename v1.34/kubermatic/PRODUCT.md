# Kubermatic Kubernetes Platform — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Kubermatic |
| **Platform** | Kubermatic Kubernetes Platform |
| **Platform Version** | v2.29 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://www.kubermatic.com/](https://www.kubermatic.com/) |
| **Documentation** | [Link](https://docs.kubermatic.com/) |

> Kubermatic Kubernetes Platform is in an open source project to centrally manage the global automation of thousands of Kubernetes clusters across multicloud, on-prem and edge.

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

- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/dyn...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/dynamic-resource-allocation/)

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/kub...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/kubelb/)
- [docs.kubermatic.com/kubelb/v1.2/tutorials/aigateway/](https://docs.kubermatic.com/kubelb/v1.2/tutorials/aigateway/)
- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/net...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/networking/ai-inference-routing/)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/architecture/concept...](https://docs.kubermatic.com/kubermatic/v2.29/architecture/concept/kkp-concepts/applications/default-applications-catalog/kueue/)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/kkp...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/kkp-autoscaler/cluster-autoscaler/)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/hpa...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/hpa-with-custom-gpu-metrics/)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/architecture/monitor...](https://docs.kubermatic.com/kubermatic/v2.29/architecture/monitoring-logging-alerting/user-cluster/)
- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/mon...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/monitoring-logging-alerting/user-cluster/user-guide/)
- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/mon...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/monitoring-logging-alerting/user-cluster/health-assessment/)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/architecture/monitor...](https://docs.kubermatic.com/kubermatic/v2.29/architecture/monitoring-logging-alerting/user-cluster/)
- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/mon...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/monitoring-logging-alerting/user-cluster/user-guide/)
- [docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/mon...](https://docs.kubermatic.com/kubermatic/v2.29/tutorials-howtos/monitoring-logging-alerting/user-cluster/health-assessment/)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [secure_accelerator_access.md](secure_accelerator_access.md)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.kubermatic.com/kubermatic/v2.29/architecture/concept...](https://docs.kubermatic.com/kubermatic/v2.29/architecture/concept/kkp-concepts/addons/kubeflow/)

---

*Generated from PRODUCT.yaml*
