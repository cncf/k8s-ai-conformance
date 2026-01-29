# CoreWeave Kubernetes Service (CKS) — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | CoreWeave |
| **Platform** | CoreWeave Kubernetes Service (CKS) |
| **Platform Version** | v1.33 |
| **Kubernetes Version** | v1.33 |

> CKS is a managed Kubernetes environment purpose-built for building, training, and deploying AI applications.

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

- [docs.coreweave.com/docs/products/cks/clusters/scheduling/...](https://docs.coreweave.com/docs/products/cks/clusters/scheduling/imex-dra-scheduling)
- [docs.coreweave.com/docs/products/cks/clusters/scheduling/...](https://docs.coreweave.com/docs/products/cks/clusters/scheduling/dynamic-reasource-featuregates)

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/tutorials/deploy-vll...](https://docs.coreweave.com/docs/products/cks/tutorials/deploy-vllm-inference)

**Notes:**

> Users on CKS can deploy and manage any known Gateway Controller which implements the Gateway API including AI Gateways like K-Gateway.

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/clusters/coreweave-c...](https://docs.coreweave.com/docs/products/cks/clusters/coreweave-charts/kueue)
- [docs.coreweave.com/docs/products/sunk#slurm-on-kubernetes-1](https://docs.coreweave.com/docs/products/sunk#slurm-on-kubernetes-1)

**Notes:**

> We support co-schedulers, the SUNK / slurm scheduler and many others as well.

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/nodes/autoscaling#co...](https://docs.coreweave.com/docs/products/cks/nodes/autoscaling#configure-autoscaling)

**Notes:**

> We support standard cluster autoscaling using our NodePool concepts.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [github.com/kubernetes-sigs/metrics-server](https://github.com/kubernetes-sigs/metrics-server)
- [docs.coreweave.com/docs/products/cks/reference/cluster-co...](https://docs.coreweave.com/docs/products/cks/reference/cluster-components)
- [docs.coreweave.com/docs/products/cks/clusters/frameworks/...](https://docs.coreweave.com/docs/products/cks/clusters/frameworks/kubeflow)

**Notes:**

> We support standard HPA installation and use of metrics-server from upstream.

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/reference/cluster-co...](https://docs.coreweave.com/docs/products/cks/reference/cluster-components)
- [docs.coreweave.com/docs/observability/managed-grafana](https://docs.coreweave.com/docs/observability/managed-grafana)

**Notes:**

> We expose and manage all known device metrics from accelerators and bare metal devices.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/reference/cluster-co...](https://docs.coreweave.com/docs/products/cks/reference/cluster-components)
- [docs.coreweave.com/docs/observability/managed-grafana](https://docs.coreweave.com/docs/observability/managed-grafana)

**Notes:**

> We expose and manage all known device metrics from accelerators and bare metal devices.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/reference/cluster-co...](https://docs.coreweave.com/docs/products/cks/reference/cluster-components)

**Notes:**

> Current devices are managed with device plugins, with ongoing efforts to transition to DRA once vendor support is more mature.

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.coreweave.com/docs/products/cks/clusters/frameworks/...](https://docs.coreweave.com/docs/products/cks/clusters/frameworks/introduction)

**Notes:**

> We support all known AI Frameworks and Controllers on CKS.

---

*Generated from PRODUCT.yaml*
