# How to submit conformance results

## The self assessment

The set of conformance requirements for each Kubernetes version are defined in the `AiConformance-x.yy.yaml` in the [docs](https://github.com/cncf/k8s-ai-conformance/tree/main/docs) folder. 

## Submitting the self assessment 

Prepare a PR to
[https://github.com/cncf/k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance).
Here are [directions](https://help.github.com/en/articles/creating-a-pull-request-from-a-fork) to
prepare a pull request from a fork.
In the descriptions below, `X.Y` refers to the kubernetes major and minor
version, and `$dir` is a short subdirectory name to hold the results for your
product.  Examples would be `gke` or `openshift`.


### Contents of the PR

You must submit the completed self assessment manifest file for the relevant major/minor version of Kubernetes in a subdirectory named after your product (`$dir`), along with test artifacts if applicable:

```
vX.Y/$dir/PRODUCT.yaml: (Required) See below.
vX.Y/$dir/e2e.log: (Recommended for v1.37+) Raw test execution log.
vX.Y/$dir/junit.xml: (Recommended for v1.37+) Machine-readable test report (or results.json).
```

#### PRODUCT.yaml

This file describes your product. It is YAML formatted with the following root-level fields. Please fill in as appropriate.

| Field                   | Description                                                                                                                                                                                                                             |
|-------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `vendorName`            | Name of the legal entity that is certifying. This entity must have a signed participation form on file with the CNCF. This must be an exact match to the organization's name as listed under CNCF Members in the [CNCF Landscape](https://landscape.cncf.io/members). |
| `kubernetesVersion`     | Kubernetes Version to verify                                                                                                                                                                                                            |
| `platformName`          | Name of the product being certified.                                                                                                                                                                                                    |
| `platformVersion`       | The version of the product being certified (not the version of Kubernetes it runs).                                                                                                                                                     |
| `websiteUrl`            | URL to the product information website                                                                                                                                                                                                  |
| `repoUrl`               | If your product is open source, this field is necessary to point to the primary GitHub repo containing the source. It's OK if this is a mirror. OPTIONAL                                                                                |
| `documentationUrl`     | URL to the product documentation                                                                                                                                                                                                        |
| `productLogoUrl`      | URL to the product's logo, (must be in SVG, AI or EPS format -- not a PNG -- and include the product name). OPTIONAL. If not supplied, we'll use your company logo. Please see logo [guidelines](https://github.com/cncf/landscape#logos) |
| `description`           | One sentence description of your offering                                                                                                                                                                                               |
| `contactEmailAddress` | An email address which can be used to contact maintainers regarding the product submitted and updates to the submission process                                                                                                         |
| `k8sConformanceUrl`   | URL to your product's Kubernetes Conformance submission in [cncf/k8s-conformance](https://github.com/cncf/k8s-conformance). Format: `https://github.com/cncf/k8s-conformance/tree/master/vX.Y/product-name`                              |

Examples below are for a fictional Kubernetes implementation called _Turbo
Encabulator_ produced by a company named _Yoyodyne_.

```yaml
metadata:
  kubernetesVersion: v1.34
  platformName: Turbo Encabulator
  platformVersion: v1.7.4
  vendorName: Yoyodyne
  website_url: https://yoyo.dyne/turbo-encabulator
  repo_url: https://github.com/yoyo.dyne/turbo-encabulator
  documentation_url: https://yoyo.dyne/turbo-encabulator/docs
  product_logo_url: https://yoyo.dyne/assets/turbo-encabulator.svg
  description: "The Yoyodyne Turbo Encabulator is a superb Kubernetes distribution for all of your Encabulating needs."
  contact_email_address: yoyodyne@turbo-encabulator.org
  k8s_conformance_url: https://github.com/cncf/k8s-conformance/tree/master/v1.34/turbo-encabulator
```

### Hybrid Verification & Automated Tests (v1.37+)

Starting with Kubernetes v1.37, the AI Conformance program supports a **hybrid verification approach** that combines automated test execution with manual attestation:

- **Automated Tests**: Requirements covered by automated tests (e.g., Secure Accelerator Access, Gang Scheduling, Cluster Autoscaling) are **recommended** to be verified by running the upstream [AI Conformance test suite](https://github.com/kubernetes-sigs/ai-conformance/tree/main/test). The generated test output files (`e2e.log`, `junit.xml`, or `results.json`) should be included in your submission directory (`v1.37/$dir/`).
- **Manual Attestation**: Requirements without automated tests continue to be verified via public documentation or reference URLs in the `evidence` field of your `PRODUCT.yaml`.

#### Running Tests & Generating Test Artifacts

Run the conformance tests against your cluster and generate the necessary artifacts for your submission:

```bash
# 1. Capture full output log (e2e.log)
go test -v ./test | tee e2e.log

# 2. Generate machine-readable JSON (results.json)
go test -v ./test -json > results.json

# 3. Generate JUnit XML report (junit.xml)
# Option A: Using gotestsum (Recommended - runs via go run)
go run gotest.tools/gotestsum@latest --junitfile junit.xml -- ./test
# Option B: Using go-junit-report
go test -v ./test 2>&1 | go run github.com/jstemmer/go-junit-report/v2@latest > junit.xml
```

For test execution details, vendor flags, and prerequisites, refer to the upstream [test documentation](https://github.com/kubernetes-sigs/ai-conformance/tree/main/test).

#### Referencing Test Results in `PRODUCT.yaml`

Link your generated test artifacts under the `evidence` field for the corresponding requirement using relative file paths or the `file://` scheme:

```yaml
    - id: secure_accelerator_access
      description: "Ensure that access to accelerators from within containers is properly isolated..."
      level: MUST
      status: "Implemented"
      evidence:
        - "e2e.log"
        - "junit.xml"
```
or:
```yaml
    - id: secure_accelerator_access
      description: "Ensure that access to accelerators from within containers is properly isolated..."
      level: MUST
      status: "Implemented"
      evidence:
        - "file://v1.37/$dir/junit.xml#TestSecureAcceleratorAccess"
        - "file://v1.37/$dir/e2e.log"
```

#### Opting Out (N/A Requirements)

If a specific requirement is not applicable to your platform (i.e., you answered "N/A" in the conformance checklist), you may skip or opt out of that specific sub-test:
- Use test-specific flags if available (for example, leave `-autoscaler-node-pool-label` unset to skip the cluster autoscaling test).
- Use the `go test -run <regex>` flag to execute only applicable tests.
- When submitting your results, ensure you provide a clear explanation and justification in the `notes` field of `PRODUCT.yaml` for why the requirement is considered N/A.

### Requirements

The self conformance file must be submitted without adjustments or changes in spec to the fields `id, description, level`.
The fields `status, evidence, and notes` need to be filled if the `level` is `MUST`.

To reach conformance, all `MUST` spec fields need to be addressed. Evidence can be publicly reachable documentation URLs or submitted test artifacts (such as `e2e.log` and `junit.xml` in your submission directory). 


## Amendment for Private Review

If you require a private review for an unreleased product, please email a .zip file containing what you would otherwise submit
as a pull request to conformance@cncf.io and the documentation for the evidence. We'll review and confirm that you are ready to be Certified AI Kubernetes Platform
as soon as you open the pull request. We can then often arrange to accept your pull request soon after you make it, at which point you become Certified Kubernetes.

## Review

A reviewer will comment on and/or accept your pull request, typically within 10 business days. If you don't see a response, please contact conformance@cncf.io.

## Issues

If you encounter issues with the conformance program itself during certification (rather than problems with your own implementation), you can file an issue in the [repository](https://github.com/cncf/k8s-ai-conformance).

Questions and comments can also be sent to the [AI Conformance](https://github.com/kubernetes/community/tree/master/wg-ai-conformance) working group, which is the change controller of the conformance definition. You can also reach the community on [Slack](https://kubernetes.slack.com/archives/C09813W8DC2). 
