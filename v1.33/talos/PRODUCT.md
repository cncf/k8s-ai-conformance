# Talos Linux — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Sidero Labs |
| **Platform** | Talos Linux |
| **Platform Version** | 1.11.3 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://talos.dev](https://talos.dev) |
| **Documentation** | [Link](https://docs.siderolabs.com/talos/latest/overview/what-is-talos) |

> Talos Linux is a single purpose Linux distribution for running Kubernetes with API management and secure by default configuration.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 6 |
| ⚪ N/A | 3 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | MUST | ✅ |
| Networking | AI Inference | MUST | ✅ |
| Scheduling & Orchestration | Gang Scheduling | MUST | ✅ |
| Scheduling & Orchestration | Cluster Autoscaling | MUST | ⚪ |
| Scheduling & Orchestration | Pod Autoscaling | MUST | ✅ |
| Observability | Accelerator Metrics | MUST | ⚪ |
| Observability | AI Service Metrics | MUST | ⚪ |
| Security | Secure Accelerator Access | MUST | ✅ |
| Operator Support | Robust Controller | MUST | ✅ |

---

## Detailed Requirements

### 🚀 Accelerators

#### ✅ DRA Support

**Level:** 🔴 MUST | **Status:** Implemented

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

**Evidence:**

- [docs.siderolabs.com/kubernetes-guides/advanced-guides/dyn...](https://docs.siderolabs.com/kubernetes-guides/advanced-guides/dynamic-resource-allocation)

**Notes:**

> DRA is supported in Talos Linux and Omni with Kubernetes v1.33 and above and Talos Linux v1.9 and above.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.siderolabs.com/kubernetes-guides/advanced-guides/dep...](https://docs.siderolabs.com/kubernetes-guides/advanced-guides/deploy-traefik)
- [docs.siderolabs.com/kubernetes-guides/cni/deploying-ciliu...](https://docs.siderolabs.com/kubernetes-guides/cni/deploying-cilium#without-kube-proxy)

**Notes:**

> Kubernetes running on Talos Linux and Omni support a variety of Gateway API controllers depending on where you run clusters. For on-prem clusters solutions such as Traefik and Cilium are compatible.

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.siderolabs.com/kubernetes-guides/advanced-guides/kueue](https://docs.siderolabs.com/kubernetes-guides/advanced-guides/kueue)

**Notes:**

> Talos Linux supports gang scheduling that can be deployed on conformant Kubernetes clusters.

#### ⚪ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** N/A

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Notes:**

> Talos Linux works with customer provided machines. In a cloud or virtualized environment autoscalers such as the Kubernetes cluster autoscaler can be used, but it is not required and does not work for edge or on-premesis infrastructure.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.siderolabs.com/kubernetes-guides/advanced-guides/hpa](https://docs.siderolabs.com/kubernetes-guides/advanced-guides/hpa)

**Notes:**

> Talos Linux supports pod autoscalers including vanilla HPA. HPA can be configured to scale based on custom metrics, such as those exposed by NVIDIA DCGM exporter or Neuron Monitor.

### 📊 Observability

#### ⚪ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** N/A

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Notes:**

> Kubernetes cluster managed by Talos linux and Omni support deploying DCGM (NVIDIA), but they do not provide integrated workload deployment, monitoring, or observability.

#### ⚪ AI Service Metrics

**Level:** 🔴 MUST | **Status:** N/A

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Notes:**

> Talos Linux allows for a variety of monitoring tools to be installed in the clusters, but it does not have access to Kubernetes workloads and does not scrape metrics directly.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [see secure_accelerator_access.md](see secure_accelerator_access.md)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.siderolabs.com/kubernetes-guides/advanced-guides/kub...](https://docs.siderolabs.com/kubernetes-guides/advanced-guides/kuberay)
- [github.com/cncf/k8s-conformance/pull/3718](https://github.com/cncf/k8s-conformance/pull/3718)

**Notes:**

> Talos Linux uses a vanilla, conformant Kubernetes cluster. Frameworks such as Ray and Kubeflow have been tested and are used on clusters implemented on Talos Linux and Omni.

---

*Generated from PRODUCT.yaml*
