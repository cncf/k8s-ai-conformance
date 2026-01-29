# Google Kubernetes Engine — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Google |
| **Platform** | Google Kubernetes Engine |
| **Platform Version** | 1.34.0-gke.1662000 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://cloud.google.com/kubernetes-engine/](https://cloud.google.com/kubernetes-engine/) |
| **Documentation** | [Link](https://cloud.google.com/kubernetes-engine/docs/) |

> GKE is an enterprise-grade platform for containerized applications, including stateful and stateless, AI and ML, Linux and Windows, complex and simple web apps, API, and backend services.

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

- [cloud.google.com/kubernetes-engine/docs/how-to/set-up-dra](https://cloud.google.com/kubernetes-engine/docs/how-to/set-up-dra)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/concepts/gateway-api](https://cloud.google.com/kubernetes-engine/docs/concepts/gateway-api)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/tutorials/kueue-i...](https://cloud.google.com/kubernetes-engine/docs/tutorials/kueue-intro)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/concepts/cluster-...](https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-autoscaler)
- [cloud.google.com/kubernetes-engine/docs/concepts/about-cu...](https://cloud.google.com/kubernetes-engine/docs/concepts/about-custom-compute-classes#gpu-rule)
- [cloud.google.com/kubernetes-engine/docs/concepts/about-cu...](https://cloud.google.com/kubernetes-engine/docs/concepts/about-custom-compute-classes#tpu_configuration)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/how-to/machine-le...](https://cloud.google.com/kubernetes-engine/docs/how-to/machine-learning/inference/autoscaling)
- [cloud.google.com/kubernetes-engine/docs/how-to/machine-le...](https://cloud.google.com/kubernetes-engine/docs/how-to/machine-learning/inference/autoscaling-tpu)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/concepts/gpus#mon...](https://cloud.google.com/kubernetes-engine/docs/concepts/gpus#monitoring)
- [cloud.google.com/kubernetes-engine/docs/how-to/tpus#metrics](https://cloud.google.com/kubernetes-engine/docs/how-to/tpus#metrics)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [cloud.google.com/kubernetes-engine/docs/concepts/observab...](https://cloud.google.com/kubernetes-engine/docs/concepts/observability)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [docs.google.com/document/d/1nx-J0oakpOU85LLtLIs-hWNSPnGv-...](https://docs.google.com/document/d/1nx-J0oakpOU85LLtLIs-hWNSPnGv-qya2x_E4dHp4HM/edit?tab=t.0)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [cloud.google.com/blog/products/ai-machine-learning/build-...](https://cloud.google.com/blog/products/ai-machine-learning/build-a-ml-platform-with-kubeflow-and-ray-on-gke)

---

*Generated from PRODUCT.yaml*
