# Frequently Asked Questions (FAQ)

This answers common questions about the CNCF Kubernetes AI Conformance program.

### What is CNCF Kubernetes AI Conformance?

The CNCF Kubernetes AI Conformance defines a set of capabilities, APIs, and configurations that a Kubernetes cluster must offer to reliably run AI/ML workloads. A Kubernetes platform or distribution must be certified as Kubernetes conformant before it can be certified as AI conformant.

### What are the goals of the AI Conformance program?

The primary goals of the program are to:
*   Simplify AI/ML on Kubernetes and accelerate adoption.
*   Guarantee interoperability and portability for AI workloads.
*   Enable ecosystem growth for AI tools on an industry standard foundation.

### Will there be AI conformance tests?

Automated conformance tests are planned for 2026. Currently, certification is based on self-assessment using a conformance checklist.

### Why are some requirements conditional?

The context of the requirements doesn’t always apply to all platforms. For example, autoscalers are not always applicable to on-prem clusters. In such cases, platforms can still be AI conformant by answering “Not Applicable” to the status and providing a justification.

### Does a "Not Applicable" (N/A) status mean platforms can skip requirements?

No. The justification for an N/A status cannot be "we don’t support this feature," but rather an explanation of why the requirement's context doesn't apply to the platform's environment.

### Is self-certification per product or per company?

AI conformance is certified per product and per configuration. For example, a cloud deployment is a different configuration from an air-gapped one.

### What use cases is this conformance for?

This effort is intended to cover all popular AI/ML workloads, including but not limited to training, inference, and agentic workloads. It is important to note that while not every workload type will use every single conformant feature, they all benefit from running on a standardized and capable platform. This ensures a consistent and interoperable environment for the entire AI/ML lifecycle.

### Why a unified AI conformance instead of separate ones for training and serving?

A unified approach is taken to avoid fragmentation and ensure workload portability across platforms. This reflects the increasing convergence of training and serving workloads.

### What is the expected revision cycle for conformance?

The revision cycle will be aligned with the release cycles of Kubernetes.
