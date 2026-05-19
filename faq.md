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

Yes. Your product does not need to be a standalone Kubernetes distribution to qualify for K8s AI Conformance. It is acceptable for your product to sit on top of a third-party distribution, provided that the underlying distribution is already Kubernetes Conformant.

AI Conformance applies to a Kubernetes platform as a whole, the combination of infrastructure, Kubernetes, and runtime/add-ons that together meet the required (must-have) features. Building on a Kubernetes AI Conformant distribution does not automatically make the layered platform AI Conformant: the combined platform must still demonstrate that every required feature is met by some layer.

#### How are requirements satisfied across layers?

The required features are grouped into categories that map to the layers of a typical platform:

*   **Infrastructure and/or base Kubernetes** — networking, storage, accelerator enablement, and the APIs and behaviors covered by the prerequisite Kubernetes Conformance (i.e., plain Kubernetes Conformance, not Kubernetes AI Conformance).
*   **Runtime / add-ons** — schedulers, operators, device plugins, and the inference/training stacks layered on top.

For each required feature, the submission attaches a reference, an evidence link, or "doesn't apply", depending on which layer owns the feature in that platform.

> NOTE: A layered product needs to re-test for full AI Conformance because layered components can change AI-conformant behaviors that the base platform's standard Kubernetes Conformance tests do not exercise. The re-test confirms the layered product's configurations did not break the base platform's existing AI Conformance. The submission references the base distribution's existing Kubernetes AI Conformance certification, when applicable. A layered product must uniquely contribute at least one feature — a MUST replacement, a meaningfully different variant of an existing MUST or SHOULD, a SHOULD addition, or a roadmap feature — to be considered AI Conformant.

#### What are the common patterns for layering on an existing Kubernetes platform?

A product can layer on an existing Kubernetes platform, whether or not the base is itself AI Conformant, in any of these ways:

1.  **Layered components that replace or supplement an existing feature.** Covers MUST-level replacements as well as meaningfully different variants of an existing MUST or SHOULD feature (for example, a more advanced implementation of a capability the base already ships). *Examples as of May 2026:* `kai-scheduler` instead of `Volcano` for gang scheduling; `Dynamo` instead of `llm-d` for disaggregated inference.
2.  **Layered components that add a feature the base does not provide.** Covers optional (SHOULD-level) features and roadmap features the layered product introduces beyond the base; sometimes anticipating capabilities that may become required in a future revision.

In all cases, the submission identifies the layer and component that owns each requirement and provides the corresponding evidence for full AI conformance of the entire stack.

#### How should a layered product distinguish its conformance claim from the base platform's?

The submission should make clear, on a per-requirement basis, which layer owns each required feature and whether the layered product's contribution replaces, matches, or extends what the base platform provides. Useful things to call out include:

*   Required (MUST) features the layered product implements independently of the base platform.
*   MUST or SHOULD features that both the layered product and the base platform implement, but the layered product offers a meaningfully different variant (for example, a more advanced implementation of the same capability).
*   Optional (SHOULD) features the layered product implements that the base platform does not.
*   Features the layered product includes that are on the AI Conformance roadmap but are not yet required.

Distinguishing should describe what the layered product does and contributes, not what the base platform lacks, to keep submissions focused on conformance rather than competitive comparison.
