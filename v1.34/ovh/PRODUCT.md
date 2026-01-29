# OVHcloud Managed Kubernetes Service — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | OVHcloud |
| **Platform** | OVHcloud Managed Kubernetes Service |
| **Platform Version** | 1.0 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://www.ovhcloud.com/en-gb/public-cloud/kubernetes/](https://www.ovhcloud.com/en-gb/public-cloud/kubernetes/) |
| **Documentation** | [Link](https://docs.ovh.com/gb/en/kubernetes/) |

> Benefit from free HA managed Kubernetes service, by hosting your nodes and services on OVH Public Cloud

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

- [./dra.yaml](./dra.yaml)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default. The attached example demonstrates how the NVIDIA driver for DRA can be configured.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [blog.ovhcloud.com/moving-beyond-ingress-why-should-ovhclo...](https://blog.ovhcloud.com/moving-beyond-ingress-why-should-ovhcloud-managed-kubernetes-service-mks-users-start-looking-at-the-gateway-api/)

**Notes:**

> OVHcloud Managed Kubernetes Services support all Gateway API controllers compatible with a vanilla kubernetes deployment.

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [./schedulingOrchestration.md#gang_scheduling](./schedulingOrchestration.md#gang_scheduling)

**Notes:**

> Custom schedulers supporting the gang scheduling algorithm can be deployed on OVHcloud Managed Kubernetes System.

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [github.com/kubernetes/autoscaler/tree/master/cluster-auto...](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/ovhcloud)

**Notes:**

> OVHcloud deploys an instance of the cluster autoscaler when autoscaling is enabled on at least one nodepool. This autoscaler supports pods requesting accelerators using the Dynamic Resource Allocation (DRA) APIs.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [./schedulingOrchestration.md#pod_autoscaling](./schedulingOrchestration.md#pod_autoscaling)

**Notes:**

> A prometheus deployment associated with the prometheus adapter can expose GPU metrics as custom metrics and be used by the HPA.

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [help.ovhcloud.com/csm/fr-public-cloud-kubernetes-monitori...](https://help.ovhcloud.com/csm/fr-public-cloud-kubernetes-monitoring-gpu-application?id=kb_article_view&sysparm_article=KB0055257)

**Notes:**

> The NVIDIA Data Center GPU Manager (DCGM) can be installed on GPU nodes.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [help.ovhcloud.com/csm/fr-public-cloud-kubernetes-monitori...](https://help.ovhcloud.com/csm/fr-public-cloud-kubernetes-monitoring-gpu-application?id=kb_article_view&sysparm_article=KB0055257)

**Notes:**

> Prometheus can be deployed on the OVHcloud Managed Kubernetes Service and collect metrics exposed by NVIDIA DCGM.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [./security.md#secure_accelerator_access](./security.md#secure_accelerator_access)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default.

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [./operator.md](./operator.md)

**Notes:**

> AI operators compatible with a vanilla Kubernetes installation + NVIDIA operators can be installed on OVHcloud Managed Kubernetes Service.

---

*Generated from PRODUCT.yaml*
