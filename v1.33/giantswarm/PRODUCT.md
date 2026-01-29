# Giant Swarm Platform — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Giant Swarm |
| **Platform** | Giant Swarm Platform |
| **Platform Version** | 1.33.0 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://www.giantswarm.io/](https://www.giantswarm.io/) |
| **Documentation** | [Link](https://docs.giantswarm.io/) |

> Giant Swarm Platform is an enterprise-grade managed Kubernetes platform for containerized applications, including stateful and stateless, AI and ML, Linux and Windows, complex and simple web apps, API, and backend services.

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

- [docs.giantswarm.io/tutorials/fleet-management/cluster-man...](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/dynamic-resource-allocation/)

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.giantswarm.io/tutorials/connectivity/gateway-api/](https://docs.giantswarm.io/tutorials/connectivity/gateway-api/)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.giantswarm.io/tutorials/fleet-management/job-managem...](https://docs.giantswarm.io/tutorials/fleet-management/job-management/kueue/)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.giantswarm.io/tutorials/fleet-management/cluster-man...](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/aws-cluster-scaling/)
- [docs.giantswarm.io/tutorials/fleet-management/cluster-man...](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/cluster-autoscaler/)
- [karpenter.sh/docs/concepts/scheduling/#acceleratorsgpu-re...](https://karpenter.sh/docs/concepts/scheduling/#acceleratorsgpu-resources)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.giantswarm.io/tutorials/fleet-management/scaling-wor...](https://docs.giantswarm.io/tutorials/fleet-management/scaling-workloads/scaling-based-on-custom-metrics)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.giantswarm.io/tutorials/fleet-management/cluster-man...](https://docs.giantswarm.io/tutorials/fleet-management/cluster-management/gpu/#monitoring)
- [docs.giantswarm.io/overview/observability/configuration/](https://docs.giantswarm.io/overview/observability/configuration/)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.giantswarm.io/getting-started/observe-your-clusters-...](https://docs.giantswarm.io/getting-started/observe-your-clusters-and-apps/)
- [docs.giantswarm.io/overview/observability/data-management...](https://docs.giantswarm.io/overview/observability/data-management/data-ingestion/)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [secure_accelerator_access_tests.md](secure_accelerator_access_tests.md)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.giantswarm.io/tutorials/fleet-management/job-managem...](https://docs.giantswarm.io/tutorials/fleet-management/job-management/kuberay)

---

*Generated from PRODUCT.yaml*
