# RKE2 — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | SUSE |
| **Platform** | RKE2 |
| **Platform Version** | v1.33 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://www.rancher.com/products/secure-kubernetes-distribution](https://www.rancher.com/products/secure-kubernetes-distribution) |
| **Documentation** | [Link](https://documentation.suse.com/cloudnative/rke2/) |

> RKE2 is an enterprise-grade conformant Kubernetes distribution that is a foundational component of SUSE AI and SUSE Rancher Prime.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 9 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | SHOULD | ✅ |
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

**Level:** 🟡 SHOULD | **Status:** Implemented

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#support-dynamic-res...](https://docs.rke2.io/reference/ai_conformance#support-dynamic-resource-allocation-dra])

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#support-the-gateway...](https://docs.rke2.io/reference/ai_conformance#support-the-gateway-api, https://raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/main/docs/cncf-ai-conformance/v1.33/SUSE-AI/specs/ai_inference.md)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#gang-scheduling](https://docs.rke2.io/reference/ai_conformance#gang-scheduling)
- [raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/m...](https://raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/main/docs/cncf-ai-conformance/v1.33/SUSE-AI/specs/gang_scheduling.md)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#cluster-autoscaler](https://docs.rke2.io/reference/ai_conformance#cluster-autoscaler)

**Notes:**

> The platform does not provide an autoscaler. It can be used in and out of environments where autoscaling it possible. The platform is tested to work with Kubernetes autoscaler.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#horizontal-pod-auto...](https://docs.rke2.io/reference/ai_conformance#horizontal-pod-autoscaler, https://raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/main/docs/cncf-ai-conformance/v1.33/SUSE-AI/specs/pod_autoscaling.md)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [documentation.suse.com/suse-ai/1.0/html/AI-deployment-int...](https://documentation.suse.com/suse-ai/1.0/html/AI-deployment-intro/index.html#observability-installing)
- [documentation.suse.com/suse-ai/1.0/html/AI-monitoring/ind...](https://documentation.suse.com/suse-ai/1.0/html/AI-monitoring/index.html#ai-monitoring-gpu)

**Notes:**

> Accelerator metrics are available through SUSE AI.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [documentation.suse.com/suse-ai/1.0/html/AI-deployment-int...](https://documentation.suse.com/suse-ai/1.0/html/AI-deployment-intro/index.html#observability-installing,https://documentation.suse.com/suse-ai/1.0/html/AI-monitoring/index.html)

**Notes:**

> Service metrics are available through SUSE AI.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#secure-accelerator-...](https://docs.rke2.io/reference/ai_conformance#secure-accelerator-access)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.rke2.io/reference/ai_conformance#robust-crd-and-cont...](https://docs.rke2.io/reference/ai_conformance#robust-crd-and-controller-operation)
- [raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/m...](https://raw.githubusercontent.com/SUSE/suse-ai-stack/refs/heads/main/docs/cncf-ai-conformance/v1.33/SUSE-AI/specs/robust_controller.md)

---

*Generated from PRODUCT.yaml*
