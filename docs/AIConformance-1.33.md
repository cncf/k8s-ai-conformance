# Kubernetes AI Conformance Checklist — v1.33

> This document defines the conformance requirements for certifying a Kubernetes platform
> as capable of reliably running AI and machine learning workloads.

---

## Overview

| Kubernetes Version | Total Requirements | Mandatory (MUST) | Recommended (SHOULD) |
|:------------------:|:------------------:|:----------------:|:--------------------:|
| **v1.33** | 9 | 8 | 1 |

### Requirement Levels

| Level | Meaning |
|:-----:|:--------|
| **MUST** | Mandatory for conformance. Platform cannot be certified without implementing this requirement. |
| **SHOULD** | Recommended but not mandatory. Platforms are encouraged to implement for better AI/ML support. |

---

## Table of Contents

- [🚀 Accelerators](#accelerators) — 1 SHOULD
- [🌐 Networking](#networking) — 1 MUST
- [📅 Scheduling & Orchestration](#scheduling--orchestration) — 3 MUST
- [📊 Observability](#observability) — 2 MUST
- [🔒 Security](#security) — 1 MUST
- [⚙️ Operator Support](#operator-support) — 1 MUST

---

## 🚀 Accelerators

*1 SHOULD*

### 1. DRA Support

**Level:** 🟡 **SHOULD**

**Description:**

> Support Dynamic Resource Allocation (DRA) APIs to enable more flexible and fine-grained resource requests beyond simple counts.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `dra_support` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## 🌐 Networking

*1 MUST*

### 1. AI Inference

**Level:** 🔴 **MUST**

**Description:**

> Support the Kubernetes Gateway API with an implementation for advanced traffic management for inference services, which enables capabilities like weighted traffic splitting, header-based routing (for OpenAI protocol headers), and optional integration with service meshes.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `ai_inference` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## 📅 Scheduling & Orchestration

*3 MUST*

### 1. Gang Scheduling

**Level:** 🔴 **MUST**

**Description:**

> The platform must allow for the installation and successful operation of at least one gang scheduling solution that ensures all-or-nothing scheduling for distributed AI workloads (e.g. Kueue, Volcano, etc.) To be conformant, the vendor must demonstrate that their platform can successfully run at least one such solution.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `gang_scheduling` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

### 2. Cluster Autoscaling

**Level:** 🔴 **MUST**

**Description:**

> If the platform provides a cluster autoscaler or an equivalent mechanism, it must be able to scale up/down node groups containing specific accelerator types based on pending pods requesting those accelerators.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `cluster_autoscaling` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

### 3. Pod Autoscaling

**Level:** 🔴 **MUST**

**Description:**

> If the platform supports the HorizontalPodAutoscaler, it must function correctly for pods utilizing accelerators. This includes the ability to scale these Pods based on custom metrics relevant to AI/ML workloads.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `pod_autoscaling` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## 📊 Observability

*2 MUST*

### 1. Accelerator Metrics

**Level:** 🔴 **MUST**

**Description:**

> For supported accelerator types, the platform must allow for the installation and successful operation of at least one accelerator metrics solution that exposes fine-grained performance metrics via a standardized, machine-readable metrics endpoint. This must include a core set of metrics for per-accelerator utilization and memory usage. Additionally, other relevant metrics such as temperature, power draw, and interconnect bandwidth should be exposed if the underlying hardware or virtualization layer makes them available. The list of metrics should align with emerging standards, such as OpenTelemetry metrics, to ensure interoperability. The platform may provide a managed solution, but this is not required for conformance.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `accelerator_metrics` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

### 2. AI Service Metrics

**Level:** 🔴 **MUST**

**Description:**

> Provide a monitoring system capable of discovering and collecting metrics from workloads that expose them in a standard format (e.g. Prometheus exposition format). This ensures easy integration for collecting key metrics from common AI frameworks and servers.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `ai_service_metrics` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## 🔒 Security

*1 MUST*

### 1. Secure Accelerator Access

**Level:** 🔴 **MUST**

**Description:**

> Ensure that access to accelerators from within containers is properly isolated and mediated by the Kubernetes resource management framework (device plugin or DRA) and container runtime, preventing unauthorized access or interference between workloads.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `secure_accelerator_access` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## ⚙️ Operator Support

*1 MUST*

### 1. Robust Controller

**Level:** 🔴 **MUST**

**Description:**

> The platform must prove that at least one complex AI operator with a CRD (e.g., Ray, Kubeflow) can be installed and functions reliably. This includes verifying that the operator's pods run correctly, its webhooks are operational, and its custom resources can be reconciled.

<details>
<summary><strong>Compliance Fields</strong></summary>

| Field | Value |
|:------|:------|
| `id` | `robust_controller` |
| `status` | `Implemented` \| `Not Implemented` \| `Partially Implemented` \| `N/A` |
| `evidence` | List of URLs to documentation/test results |
| `notes` | Additional context (required if status is `N/A`) |

</details>

---

## Submission Instructions

To submit your platform for conformance certification:

1. Copy this checklist template to `PRODUCT.yaml`
2. Fill in the `metadata` section with your platform details
3. For each requirement, set the `status` field to one of:
   - `Implemented` — Requirement is fully supported
   - `Not Implemented` — Requirement is not supported
   - `Partially Implemented` — Requirement is partially supported
   - `N/A` — Requirement does not apply (must provide justification in `notes`)
4. Provide `evidence` URLs linking to documentation or test results
5. Submit a pull request to the [k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance) repository

---

*Generated from AIConformance-1.33.yaml*
