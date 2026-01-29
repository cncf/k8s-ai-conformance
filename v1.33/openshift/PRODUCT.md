# OpenShift Container Platform — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Red Hat |
| **Platform** | OpenShift Container Platform |
| **Platform Version** | 4.20 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://www.redhat.com/en/technologies/cloud-computing/openshift](https://www.redhat.com/en/technologies/cloud-computing/openshift) |
| **Documentation** | [Link](https://docs.openshift.com/container-platform/4.20/) |

> Red Hat OpenShift Container Platform is an enterprise-ready Kubernetes container platform with full-stack automated operations to manage hybrid cloud, multicloud, and edge deployments.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 8 |
| ❌ Not Implemented | 1 |
| **Total** | **9** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | SHOULD | ❌ |
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

#### ❌ DRA Support

**Level:** 🟡 SHOULD | **Status:** Not Implemented

> Dynamic Resource Allocation (DRA) APIs enable more flexible and fine-grained resource requests beyond simple counts.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/network_apis/gateway-gateway-networking-k8s-io-v1)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/ai_workloads/red-hat-build-of-kueue#gangscheduling)

**Notes:**

> Red Hat build of Kueue enables gang admission

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/hardware_accelerators/about-hardware-accelerators)
- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/machine_management/applying-autoscaling)
- [www.redhat.com/en/blog/autoscaling-nvidia-gpus-on-red-hat...](https://www.redhat.com/en/blog/autoscaling-nvidia-gpus-on-red-hat-openshift)

**Notes:**

> The OpenShift cluster autoscaler implementation satisfies this requirement. We have tested with several different models of GPU and users are able to control how their workloads are matched with specific hardware needs.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/nodes/automatically-scaling-pods-with-the-custom-metrics-autoscaler-operator#nodes-cma-autoscaling-custom-trigger-prom-gpu_nodes-cma-autoscaling-custom-trigger)
- [developers.redhat.com/articles/2025/08/12/boost-ai-effici...](https://developers.redhat.com/articles/2025/08/12/boost-ai-efficiency-gpu-autoscaling-openshift#custom_metrics_autoscaler__keda__and_prometheus)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.nvidia.com/datacenter/cloud-native/openshift/latest/...](https://docs.nvidia.com/datacenter/cloud-native/openshift/latest/enable-gpu-monitoring-dashboard.html)
- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/hardware_accelerators/nvidia-gpu-architecture)
- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/hardware_accelerators/amd-gpu-operator)
- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/red_hat_build_of_opentelemetry/index)

**Notes:**

> As part of the OpenShift observability solution, OpenShift provides comprehensive support for AI accelerators (e.g. NVIDIA, AMD) through dedicated GPU operators that enable standardized metrics collection and monitoring. NVIDIA GPU Operator integrates DCGM-based monitoring, exposing GPU utilization, power consumption (watts), temperature (Celsius), utilization (percent), and memory metrics. AMD GPU Operator with ROCm integration provides equivalent AI accelerator monitoring capabilities. GPU telemetry is exposed via DCGM Exporter for Prometheus consumption through /metrics endpoints. OpenShift observability solution also provides native integration with OpenTelemetry standards via the Red Hat build of OpenTelemetry.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.redhat.com/en/documentation/openshift_container_plat...](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/pdf/monitoring/OpenShift_Container_Platform-4.20-Monitoring-en-US.pdf)

**Notes:**

> OpenShift provides a fully integrated monitoring system based on Prometheus, which automatically discovers and scrapes metrics endpoints exposed by workloads in the standard Prometheus exposition format, ensuring seamless integration for collecting and displaying key metrics from common AI frameworks and servers.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [docs.google.com/document/d/14t54X8N6Xg8P0p6-MKZxySLrSuJQm...](https://docs.google.com/document/d/14t54X8N6Xg8P0p6-MKZxySLrSuJQm15M/edit?usp=sharing&ouid=104474006979808168216&rtpof=true&sd=true)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [developers.redhat.com/articles/2025/04/22/fine-tune-llms-...](https://developers.redhat.com/articles/2025/04/22/fine-tune-llms-kubeflow-trainer-openshift-ai)
- [docs.google.com/document/d/1XKqETt-sXbznYwrX5toJ1pVhcYaoq...](https://docs.google.com/document/d/1XKqETt-sXbznYwrX5toJ1pVhcYaoq6PL/edit?usp=sharing&ouid=104474006979808168216&rtpof=true&sd=true)

**Notes:**

> Kubeflow Trainer operator V1 with Multiple CRDs (TrainJob, TrainingRuntime, ClusterTrainingRuntime), admission webhooks, and controller reconciliation. Supports distributed training for PyTorch, DeepSpeed, MLX, MPI frameworks. Supported for self-managed OpenShift and Red Hat OpenShift AI Cloud Service.

---

*Generated from PRODUCT.yaml*
