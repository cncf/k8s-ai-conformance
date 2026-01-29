# DaoCloud Enterprise — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | DaoCloud |
| **Platform** | DaoCloud Enterprise |
| **Platform Version** | v5.0 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://docs.daocloud.io/](https://docs.daocloud.io/) |
| **Documentation** | [Link](https://download.daocloud.io/DaoCloud_Enterprise/DaoCloud_Enterprise) |

> Daocloud helps you provide a reliable and consistent basic support environment to meet the high SLA requirements of enterprise critical applications.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 7 |
| ⚪ N/A | 2 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | SHOULD | ⚪ |
| Networking | AI Inference | MUST | ✅ |
| Scheduling & Orchestration | Gang Scheduling | MUST | ✅ |
| Scheduling & Orchestration | Cluster Autoscaling | MUST | ⚪ |
| Scheduling & Orchestration | Pod Autoscaling | MUST | ✅ |
| Observability | Accelerator Metrics | MUST | ✅ |
| Observability | AI Service Metrics | MUST | ✅ |
| Security | Secure Accelerator Access | MUST | ✅ |
| Operator Support | Robust Controller | MUST | ✅ |

---

## Detailed Requirements

### 🚀 Accelerators

#### ⚪ DRA Support

**Level:** 🟡 SHOULD | **Status:** N/A

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

**Notes:**

> DRA APIs are disabled in 1.33 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)
- [docs.daocloud.io/en/hydra/intro/deploy-ws/#create-service...](https://docs.daocloud.io/en/hydra/intro/deploy-ws/#create-service-mesh-istio-gateway-api)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/networking.go#L16

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)
- [docs.daocloud.io/en/baize/best-practice/train-with-tas/](https://docs.daocloud.io/en/baize/best-practice/train-with-tas/)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/scheduling_orchestration.go#L41

#### ⚪ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** N/A

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Notes:**

> DaoCloud Enterprise run Kubernetes on-premises and does not provide any cluster autoscaler like karpenter.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)
- [docs.daocloud.io/en/kpanda/user-guide/scale/custom-hpa/](https://docs.daocloud.io/en/kpanda/user-guide/scale/custom-hpa/)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/scheduling_orchestration.go#L270

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)
- [docs.daocloud.io/en/insight/quickstart/otel/operator/](https://docs.daocloud.io/en/insight/quickstart/otel/operator/)
- [docs.daocloud.io/en/kpanda/user-guide/gpu/nvidia/gpu-moni...](https://docs.daocloud.io/en/kpanda/user-guide/gpu/nvidia/gpu-monitoring-alarm/gpu-metrics)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/observability.go#L26

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)
- [docs.daocloud.io/en/insight/user-guide/collection-manag/c...](https://docs.daocloud.io/en/insight/user-guide/collection-manag/collection-manag/)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/observability.go#L101

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/security.go#L23

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [e2e.log](e2e.log)
- [junit_01.xml](junit_01.xml)

**Notes:**

> test code: https://github.com/carlory/ai-conformance/blob/c4b99e98160525e78e475165a72b0c920501f57c/e2e/ai/operator.go#L41

---

*Generated from PRODUCT.yaml*
