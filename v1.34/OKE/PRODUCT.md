# OCI Kubernetes Engine (OKE) — AI Conformance Report

## Platform Information

| | |
|:--|:--|
| **Vendor** | Oracle |
| **Platform** | OCI Kubernetes Engine (OKE) |
| **Platform Version** | v1.34 |
| **Kubernetes Version** | v1.34 |
| **Website** | [https://www.oracle.com/cloud/cloud-native/kubernetes-engine/](https://www.oracle.com/cloud/cloud-native/kubernetes-engine/) |
| **Documentation** | [Link](https://docs.oracle.com/en-us/iaas/Content/ContEng/home.htm) |

> Oracle Cloud Infrastructure Kubernetes Engine (OKE) is a fully-managed, scalable, and highly available service that you can use to deploy your containerized applications to the cloud.

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

- [docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/conte...](https://docs.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengaboutk8sversions.htm#supportedk8sversions)

**Notes:**

> DRA v1 APIs are enabled in 1.34 by default

### 🌐 Networking

#### ✅ AI Inference

**Level:** 🔴 MUST | **Status:** Implemented

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

**Evidence:**

- [docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/conteng-...](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/conteng-using-istio.htm#:~:text=Using%20the%20Kubernetes%20Gateway%20API)

### 📅 Scheduling & Orchestration

#### ✅ Gang Scheduling

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

**Evidence:**

- [github.com/oracle-quickstart/oci-hpc-oke/blob/main/docs/u...](https://github.com/oracle-quickstart/oci-hpc-oke/blob/main/docs/using-rdma-network-locality-when-running-workloads-on-oke.md#using-kueue-with-topology-aware-scheduling)
- [www.oracle.com/in/a/ocom/docs/cloud/accelerate-ai-with-oc...](https://www.oracle.com/in/a/ocom/docs/cloud/accelerate-ai-with-oci-supercluster.pdf)

#### ✅ Cluster Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

**Evidence:**

- [docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengu...](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengusingclusterautoscaler.htm#:~:text=In%20addition%20to%20scaling%20based%20on%20CPU%2C%20memory%2C%20and%20other%20properties%2C%20the%20Kubernetes%20Cluster%20Autoscaler%20also%20supports%20scaling%20based%20on%20requests%20for%20GPU%20accelerators.%20For%20example%2C%20the%20Kubernetes%20Cluster%20Autoscaler%20will%20scale%20up%20a%20node%20pool%20when%20a%20pod%20requesting%20an%20accelerator%20cannot%20be%20scheduled.)

#### ✅ Pod Autoscaling

**Level:** 🔴 MUST | **Status:** Implemented

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

**Evidence:**

- [blogs.oracle.com/developers/post/autoscaling-gpu-workload...](https://blogs.oracle.com/developers/post/autoscaling-gpu-workloads-oci-kubernetes-engine-oke)
- [docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengu...](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengusinghorizontalpodautoscaler.htm)

### 📊 Observability

#### ✅ Accelerator Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

**Evidence:**

- [docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengi...](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengintroducingclusteraddons.htm#contengintroducingclusteraddons__section-optional-addons:~:text=Add%2Don.-,NVIDIA%20GPU%20Plugin%3A,-The%20optional%20NVIDIA)
- [docs.oracle.com/en-us/iaas/Content/Compute/Tasks/transiti...](https://docs.oracle.com/en-us/iaas/Content/Compute/Tasks/transitioning-to-hpc-plugin.htm)
- [github.com/oracle-quickstart/oci-gpu-scanner](https://github.com/oracle-quickstart/oci-gpu-scanner)

#### ✅ AI Service Metrics

**Level:** 🔴 MUST | **Status:** Implemented

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

**Evidence:**

- [github.com/oci-hpc/oci-hpc-oke/blob/manual-monitoring-dep...](https://github.com/oci-hpc/oci-hpc-oke/blob/manual-monitoring-deployment/docs/deploying-monitoring-stack-manually.md)
- [github.com/oracle-quickstart/oci-gpu-scanner](https://github.com/oracle-quickstart/oci-gpu-scanner)

### 🔒 Security

#### ✅ Secure Accelerator Access

**Level:** 🔴 MUST | **Status:** Implemented

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

**Evidence:**

- [blogs.oracle.com/developers/post/verifying-gpu-accelerato...](https://blogs.oracle.com/developers/post/verifying-gpu-accelerator-isolation-oci-kubernetes-engine-oke)

### ⚙️ Operator Support

#### ✅ Robust Controller

**Level:** 🔴 MUST | **Status:** Implemented

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

**Evidence:**

- [docs.anyscale.com/admin/cloud/kubernetes](https://docs.anyscale.com/admin/cloud/kubernetes)
- [www.oracle.com/artificial-intelligence/machine-learning-w...](https://www.oracle.com/artificial-intelligence/machine-learning-with-kubeflow/)

---

*Generated from PRODUCT.yaml*
