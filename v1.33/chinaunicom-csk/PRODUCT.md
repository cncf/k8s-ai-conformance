# CSK — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Chinaunicom Cloud |
| **Platform** | CSK |
| **Platform Version** | v1.33 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://www.cucloud.cn/product/csk.html](https://www.cucloud.cn/product/csk.html) |
| **Documentation** | [Link](https://support.cucloud.cn/document/127/581/900.html?id=900&arcid=1756) |

> Chinaunicom container service product based on Kubernetes. It provides high-performance and scalable container application management capabilities

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

**Notes:**

> DRA APIs are disabled in 1.33 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5492)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5352)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5485)

**Notes:**

> The platform does not provide an autoscaler. It can be used in and out of environments where autoscaling it possible. The platform is tested to work with Kubernetes autoscaler.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5486)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5487)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=6052)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5489)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [support.cucloud.cn/document/127/581/2541.html?id=2541&arc...](https://support.cucloud.cn/document/127/581/2541.html?id=2541&arcid=5490)

---

*Generated from PRODUCT.yaml*
