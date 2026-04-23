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

### Can I certify a platform or product that is built on top of an existing Kubernetes distribution?

Yes. Your product does not need to be a standalone Kubernetes distribution to qualify for K8s AI Conformance. It is perfectly acceptable for your product to sit on top of a third-party distribution, provided that the underlying distribution is already Kubernetes Conformant.

Building on a conformant distribution does not automatically make the layered product AI conformant. The layered product must independently demonstrate conformance for its own supported configuration.

#### What must a layered product demonstrate on its own?

<!-- TODO: Open for discussion — define the minimum requirements a layered product
     must satisfy independently of the underlying platform's conformance. -->

The layered product is expected to:
*   Run the conformance tests (or complete the self-assessment checklist) in its own supported and shipped configuration.
*   Show that its packaging, integrations, default settings, and additional components do not break or degrade any conformance requirements.

#### How should a layered product distinguish its conformance claim from the already-conformant base distribution's?

<!-- TODO: Open for discussion — clarify how the certification scope of a layered
     product is distinguished from the already-conformant base distribution. -->

The layered product's submission should clearly articulate:
*   What additional functionality, supported configuration, or validated behavior it contributes beyond the already-conformant base platform.
*   How its claim is distinct from the underlying distribution's claim — i.e., what value or differentiation it adds that warrants a separate certification.
