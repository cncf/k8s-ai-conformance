# JCS for Kubernetes — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | JD Cloud |
| **Platform** | JCS for Kubernetes |
| **Platform Version** | v1.33.3 |
| **Kubernetes Version** | v1.33 |
| **Website** | [https://www.jdcloud.com/en/products/jcs-for-kubernete](https://www.jdcloud.com/en/products/jcs-for-kubernete) |
| **Documentation** | [Link](https://docs.jdcloud.com/en/jcs-for-kubernetes/product-overview) |

> By adopting fully-hosted management node, JCS for Kubernetes offers simple, easy-to-use, high-reliable and powerful container management service to users. Compatible with Kubernetes API, integrating network, storage and other JD Cloud plug-ins.

---

## Compliance Summary

| Status | Count |
|:-------|:-----:|
| ✅ Implemented | 8 |
| **Total** | **8** |

### Requirements at a Glance

| Category | Requirement | Level | Status |
|:---------|:------------|:-----:|:------:|
| Accelerators | DRA Support | SHOULD | ⬜ |
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

#### ⬜ DRA Support

**Level:** 🟡 SHOULD | **Status:** 

> Dynamic Resource Allocation (DRA) APIs enable more flexible and fine-grained resource requests beyond simple counts.

**Evidence:**

- [docs.jdcloud.com/cn/jcs-for-kubernetes/add-gpu](https://docs.jdcloud.com/cn/jcs-for-kubernetes/add-gpu)

**Notes:**

> Through device-plugin integration, JCS for Kubernetes provides end users with the capability to dynamically allocate GPU resources to Pod containers, implementing the Dynamic Resource Allocation (DRA) APIs for flexible and fine-grained resource requests.

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.jdcloud.com/cn/api-gateway/product-overview](https://docs.jdcloud.com/cn/api-gateway/product-overview)
- [docs.jdcloud.com/cn/jdaip/api-key](https://docs.jdcloud.com/cn/jdaip/api-key)
- [docs.jdcloud.com/cn/jdaip/service-invocation-llm](https://docs.jdcloud.com/cn/jdaip/service-invocation-llm)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [docs.jdcloud.com/cn/jdaip/llm-deploy-1](https://docs.jdcloud.com/cn/jdaip/llm-deploy-1)

**Notes:**

> The JD JoyBuild platform which built on JCS for Kubernetes demonstrates gang scheduling conformance by successfully deploying and operating volcano as its primary batch scheduling solution. The platform enables all-or-nothing scheduling for distributed AI workloads, ensuring that multi-pod AI training jobs are either fully scheduled with all required resources or not scheduled at all. This capability is validated through successful execution of AI models using both public container images and user private images, maintaining workload integrity and resource co-location requirements.

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.jdcloud.com/cn/jcs-for-kubernetes/telescopic-nodegroup](https://docs.jdcloud.com/cn/jcs-for-kubernetes/telescopic-nodegroup)
- [docs.jdcloud.com/cn/jcs-for-kubernetes/api/setnodegroupca](https://docs.jdcloud.com/cn/jcs-for-kubernetes/api/setnodegroupca)
- [docs.jdcloud.com/cn/jdaip/create-and-manage-workspace](https://docs.jdcloud.com/cn/jdaip/create-and-manage-workspace)
- [docs.jdcloud.com/cn/jdaip/nodepool](https://docs.jdcloud.com/cn/jdaip/nodepool)
- [docs.jdcloud.com/cn/availability-group/auto-scaling-overview](https://docs.jdcloud.com/cn/availability-group/auto-scaling-overview)

**Notes:**

> JCS for Kubernetes demonstrates gang scheduling conformance by successfully deploying and operating Volcano as its primary batch scheduling solution. The platform enables all-or-nothing scheduling for distributed AI workloads, ensuring that multi-pod AI training jobs are either fully scheduled with all required resources or not scheduled at all. This capability is validated through successful execution of AI models using both public container images and user private images, maintaining workload integrity and resource co-location requirements.

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [docs.jdcloud.com/cn/jcs-for-kubernetes/cronhpa](https://docs.jdcloud.com/cn/jcs-for-kubernetes/cronhpa)
- [docs.jdcloud.com/cn/jdaip/create-queue](https://docs.jdcloud.com/cn/jdaip/create-queue)

**Notes:**

> JCS for Kubernetes demonstrates HorizontalPodAutoscaler functionality for accelerator-utilizing pods through its CronHPA and AI resource queue mechanisms. The platform enables dynamic scaling of AI/ML workloads based on custom metrics relevant to machine learning operations, such as GPU utilization, training progress, or inference latency. This ensures that pods with accelerators can automatically scale up during high-demand periods and scale down when resources are underutilized, optimizing both performance and cost for AI workloads.

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.jdcloud.com/cn/monitoring/learning](https://docs.jdcloud.com/cn/monitoring/learning)
- [docs.jdcloud.com/cn/jdaip/monitor](https://docs.jdcloud.com/cn/jdaip/monitor)
- [docs.jdcloud.com/cn/jdaip/model-observations](https://docs.jdcloud.com/cn/jdaip/model-observations)

**Notes:**

> JCS for Kubernetes implements robust accelerator metrics collection using monitor agent, NVIDIA DCGM, and vendor-provided monitoring tools. The platform exposes comprehensive performance metrics through standardized endpoints, covering core requirements like per-GPU utilization and memory usage, plus additional metrics such as thermal data, power metrics, and interconnect bandwidth when supported by hardware. All metrics follow OpenTelemetry conventions for machine-readable format and cross-platform interoperability.

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [docs.jdcloud.com/cn/jcs-for-kubernetes/Kubernetes-install...](https://docs.jdcloud.com/cn/jcs-for-kubernetes/Kubernetes-install-jdmon)
- [docs.jdcloud.com/cn/jcs-for-kubernetes/prometheus-instanc...](https://docs.jdcloud.com/cn/jcs-for-kubernetes/prometheus-instance-grafana-dashboard)
- [docs.jdcloud.com/cn/jcs-for-kubernetes/custom-pod-metric-...](https://docs.jdcloud.com/cn/jcs-for-kubernetes/custom-pod-metric-notifications)

**Notes:**

> JCS for Kubernetes  provides a fully integrated monitoring system based on Prometheus, which automatically discovers and scrapes metrics endpoints exposed by workloads in the standard Prometheus exposition format, ensuring seamless integration for collecting and displaying key metrics from common AI frameworks and servers.

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [docs.jdcloud.com/cn/jdaip/permission-control](https://docs.jdcloud.com/cn/jdaip/permission-control)
- [docs.jdcloud.com/cn/virtual-machines/install-GPU](https://docs.jdcloud.com/cn/virtual-machines/install-GPU)
- [docs.jdcloud.com/cn/jcs-for-kubernetes/custom-gpu-driver](https://docs.jdcloud.com/cn/jcs-for-kubernetes/custom-gpu-driver)
- [docs.jdcloud.com/cn/gcs/loginInstance](https://docs.jdcloud.com/cn/gcs/loginInstance)

**Notes:**

> JCS for Kubernetes ensures proper accelerator isolation and mediation through permission control and its device plugin implementation, which enforces strict resource boundaries within the Kubernetes resource management framework. Only allocated GPUs are accessible within workload containers, preventing unauthorized access or interference between different workloads. The device plugin works in conjunction with the container runtime to mediate all accelerator access, ensuring that each container can only utilize the specific GPU resources assigned to it, thereby maintaining workload isolation and security.

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.jdcloud.com/cn/jdaip/create-trainjob](https://docs.jdcloud.com/cn/jdaip/create-trainjob)

**Notes:**

> The JD JoyBuild platform, built on JCS for Kubernetes, demonstrates full AI operator support by successfully installing and operating PyTorch and Ray training operators. These complex AI operators with their respective CRDs are fully functional, with all operator pods running correctly, webhooks operational, and custom resources being properly reconciled to support distributed machine learning workloads.

---

*Generated from PRODUCT.yaml*
