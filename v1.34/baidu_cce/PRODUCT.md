# CCE（Cloud Container Engine） — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Baidu Cloud |
| **Platform** | CCE（Cloud Container Engine） |
| **Platform Version** | 1.34 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://cloud.baidu.com/](https://cloud.baidu.com/) |
| **Documentation** | [Link](https://cloud.baidu.com/doc/CCE/index.html) |

> Cloud Container Engine (CCE) is a highly scalable, high-performance container management service. It allows you to easily run applications on hosted cloud server instance clusters. With CCE, there's no need to install, operate, and scale cluster management infrastructure. You can start and stop Docker applications, check the complete status of the cluster, and access various cloud services with simple API calls. Containers can be deployed in your cluster based on your resource and availability requirements, meeting the specific needs of your service or application.

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

- [intl.cloud.baidu.com/en/doc/CCE/s/xmhddquk7](https://intl.cloud.baidu.com/en/doc/CCE/s/xmhddquk7)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/kmhddqumv](https://intl.cloud.baidu.com/en/doc/CCE/s/kmhddqumv)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/Umhddquu8](https://intl.cloud.baidu.com/en/doc/CCE/s/Umhddquu8)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/Fmhddqus0](https://intl.cloud.baidu.com/en/doc/CCE/s/Fmhddqus0)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/amhddquwk](https://intl.cloud.baidu.com/en/doc/CCE/s/amhddquwk)

**Notes:**

> Implemented

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/Ulps6uxwe-intl-en](https://intl.cloud.baidu.com/en/doc/CCE/s/Ulps6uxwe-intl-en)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/slps718vs-intl-en](https://intl.cloud.baidu.com/en/doc/CCE/s/slps718vs-intl-en)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/Lmhddquze](https://intl.cloud.baidu.com/en/doc/CCE/s/Lmhddquze)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [intl.cloud.baidu.com/en/doc/CCE/s/Jmhddqupe](https://intl.cloud.baidu.com/en/doc/CCE/s/Jmhddqupe)

---

*Generated from PRODUCT.yaml*
