# Conformance Checklists per Release

## Summary

Each `AIConformance-1.xx.yaml` file contains the conformance requirements for that Kubernetes release. Vendors use these as templates to demonstrate their platform meets the requirements.

- All `MUST` requirements must be fulfilled with documented evidence to achieve conformance
- All `SHOULD` requirements are recommended but optional

## Hybrid Verification Submission (v1.37+)

Starting with v1.37, we support a hybrid verification approach that combines automated test results with manual attestation:

1.  **Automated Tests**: For requirements with automated tests (e.g., Secure Accelerator Access, Gang Scheduling, Cluster Autoscaling), it is **recommended** that you run the upstream [AI Conformance test suite](https://github.com/kubernetes-sigs/ai-conformance/tree/main/test) and include the output files (`e2e.log`, `junit.xml`, or `results.json`) in your submission PR.
2.  **Manual Attestation**: For requirements without automated tests, continue to provide documentation or reference URLs in the `evidence` field of the checklist YAML.

Reference your automated test results in the checklist YAML using the `file://` scheme or relative paths to your submitted test artifacts. See [instructions.md](../instructions.md#hybrid-verification--automated-tests-v137) for details on running tests and generating evidence.

## Release Freeze

The `AIConformance-1.xx.yaml` files are finalized at Kubernetes Release Freeze and won't change after that point.

Once frozen, a new `AIConformance-1.xx.yaml` is added to this repo for the upcoming release so vendors can submit their conformance results for CNCF review and certification.

The conformance requirements are defined by the [Kubernetes AI Conformance](https://github.com/kubernetes-sigs/ai-conformance/blob/main/RELEASE.md) project. Join the conversation on [Slack](https://kubernetes.slack.com/archives/C09813W8DC2).