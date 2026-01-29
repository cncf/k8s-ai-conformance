# Linode Kubernetes Engine (LKE) — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Akamai |
| **Platform** | Linode Kubernetes Engine (LKE) |
| **Platform Version** | v1.34 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://www.akamai.com/cloud](https://www.akamai.com/cloud) |
| **Documentation** | [Link](https://techdocs.akamai.com/cloud-computing/docs/linode-kubernetes-engine) |

> Linode Kubernetes Engine (LKE) is Akamai’s managed container orchestration engine built on top of Kubernetes. Through LKE, you can quickly deploy and manage your containerized applications without needing to build and maintain your own Kubernetes cluster. This enables you to utilize Kubernetes without the specialized knowledge, added complexity, and additional overhead typically associated with manual deployments.

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

- [github.com/linode/lke-ai-conformance/blob/main/README.md](https://github.com/linode/lke-ai-conformance/blob/main/README.md)
- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#dra-support)

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#ai-inference)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#gang-scheduling)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [techdocs.akamai.com/cloud-computing/docs/manage-nodes-and...](https://techdocs.akamai.com/cloud-computing/docs/manage-nodes-and-node-pools#autoscale-automatically-resize-node-pools)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#pod-autoscaling)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#accelerator-metrics)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [www.linode.com/docs/guides/deploy-prometheus-operator-wit...](https://www.linode.com/docs/guides/deploy-prometheus-operator-with-grafana-on-lke/)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#secure-accelerator-access)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [github.com/linode/lke-ai-conformance/blob/main/README.md#...](https://github.com/linode/lke-ai-conformance/blob/main/README.md#robust-controller)

---

*Generated from PRODUCT.yaml*
