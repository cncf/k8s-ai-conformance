# Gardener — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | NeoNephos Foundation |
| **Platform** | Gardener |
| **Platform Version** | v1.130.0 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://gardener.cloud/](https://gardener.cloud/) |
| **Documentation** | [Link](https://gardener.cloud/docs/) |

> Gardener is an open-source project that provides a Kubernetes-native way to manage Kubernetes clusters as a service. It enables users to provision, manage, and operate conformant Kubernetes clusters across various cloud and on-premise infrastructures.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 8 |
| 🟡 Partially Implemented | 1 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | SHOULD | 🟡 |
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

#### 🟡 DRA Support

**Level:** 🟡 SHOULD | **Status:** Partially Implemented

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/dra_support/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/dra_support/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/dra_support/test_result.log)

**Notes:**

> Verified that the resource.k8s.io API group is available with all required DRA resource types (deviceclasses, resourceclaims, resourceclaimtemplates, resourceslices) at v1beta1. DRA v1 APIs (required by spec) are GA in Kubernetes v1.34+, hence the partial implementation status for v1.33.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_inference/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_inference/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_inference/test_result.log)

**Notes:**

> Installed Traefik gateway controller, created GatewayClass, Gateway, and HTTPRoute with weighted traffic splitting (70/30) and header-based routing. Verified all resources accepted and functional.

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/gang_scheduling/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/gang_scheduling/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/gang_scheduling/test_result.log)

**Notes:**

> Installed Kueue v0.14.2 gang scheduling solution, configured resource quotas, and submitted a multi-pod job requiring 3 pods to run in parallel. Verified Kueue admitted the job and all 3 pods were scheduled atomically (all-or-nothing).

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/cluster_autoscaling/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/cluster_autoscaling/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/cluster_autoscaling/test_result.log)

**Notes:**

> Started with 1 GPU node, deployed 2 pods each requesting 1 GPU (exceeding capacity). Verified autoscaler scaled up to 2 nodes so both pods could run. Then deleted the workload and verified autoscaler scaled back down to 1 node.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/pod_autoscaling/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/pod_autoscaling/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/pod_autoscaling/test_result.log)

**Notes:**

> Deployed Prometheus stack with DCGM exporter integration, created a custom GPU utilization metric (pod_gpu_utilization) via PrometheusRule, deployed prometheus-adapter to expose it via Custom Metrics API, then created an HPA targeting a GPU workload. Verified HPA scaled up when GPU load exceeded threshold and scaled down when load was removed.

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/accelerator_metrics/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/accelerator_metrics/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/accelerator_metrics/test_result.log)

**Notes:**

> Verified that NVIDIA DCGM Exporter (pre-installed via GPU Operator) exposes GPU metrics at http://nvidia-dcgm-exporter.gpu-operator.svc:9400/metrics in Prometheus format, including per-accelerator utilization, memory, temperature, and power metrics.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_service_metrics/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_service_metrics/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/ai_service_metrics/test_result.log)

**Notes:**

> Deployed a test AI application (podinfo) exposing Prometheus metrics, then deployed our own Prometheus stack with pod annotation-based discovery. Generated traffic and verified metrics were successfully scraped and queryable.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/secure_accelerator_access/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/secure_accelerator_access/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/secure_accelerator_access/test_result.log)

**Notes:**

> Deployed a pod without GPU resource requests to a GPU node. Verified it cannot access GPU devices (/dev/nvidia* not present). Also, deployed 2 pods each requesting 1 GPU. Verified each pod received a different GPU (different UUIDs), could only see exactly 1 GPU via nvidia-smi, and could not access unauthorized GPU device files.

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/robust_controller/README.md)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/robust_controller/test_procedure.sh)
- [github.com/gardener/gardener-ai-conformance/blob/main/v1....](https://github.com/gardener/gardener-ai-conformance/blob/main/v1.33/robust_controller/test_result.log)

**Notes:**

> Installed KubeRay operator v1.3.0, verified CRDs were registered (RayCluster, RayJob, RayService), tested webhook validation by submitting an invalid RayCluster spec (correctly rejected), created a valid RayCluster that reconciled to ready state, and executed distributed Ray tasks to confirm functionality.

---

*Generated from PRODUCT.yaml*
