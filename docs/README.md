# Conformance Checklists per Release

## Summary

Each `AIConformance-1.xx.yaml` file contains the conformance requirements for that Kubernetes release. Vendors use these as templates to demonstrate their platform meets the requirements.

- All `MUST` requirements must be fulfilled with documented evidence to achieve conformance
- All `SHOULD` requirements are recommended but optional

## Release Freeze

The `AIConformance-1.xx.yaml` files are finalized at Kubernetes Release Freeze and won't change after that point.

Once frozen, a new `AIConformance-1.xx.yaml` is added to this repo for the upcoming release so vendors can submit their conformance results for CNCF review and certification.

The conformance requirements are defined by the [Kubernetes AI Conformance](https://github.com/kubernetes-sigs/ai-conformance/blob/main/RELEASE.md) project. Join the conversation on [Slack](https://kubernetes.slack.com/archives/C09813W8DC2).