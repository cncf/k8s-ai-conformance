//go:build validate

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSupportsHybridVerification(t *testing.T) {
	cases := map[string]bool{
		"1.33": false,
		"1.36": false,
		"1.37": true,
		"1.38": true,
		"2.0":  true,
		"bad":  false,
		"":     false,
	}
	for version, want := range cases {
		if got := supportsHybridVerification(version); got != want {
			t.Errorf("supportsHybridVerification(%q) = %v, want %v", version, got, want)
		}
	}
}

func TestResolveEvidencePath(t *testing.T) {
	productDir := filepath.Join("v1.37", "gke")
	cases := []struct {
		name         string
		link         string
		wantPath     string
		wantFragment string
		wantErr      bool
	}{
		{name: "bare relative", link: "junit.xml", wantPath: filepath.Join(productDir, "junit.xml")},
		{name: "bare with fragment", link: "junit.xml#TestGangScheduling", wantPath: filepath.Join(productDir, "junit.xml"), wantFragment: "TestGangScheduling"},
		{name: "file scheme product relative", link: "file://e2e.log", wantPath: filepath.Join(productDir, "e2e.log")},
		{name: "file scheme repo-root style", link: "file://v1.37/gke/junit.xml#TestSecureAcceleratorAccess", wantPath: filepath.Join(productDir, "junit.xml"), wantFragment: "TestSecureAcceleratorAccess"},
		{name: "file scheme triple slash", link: "file:///v1.37/gke/results.json", wantPath: filepath.Join(productDir, "results.json")},
		{name: "subdirectory", link: "artifacts/junit.xml", wantPath: filepath.Join(productDir, "artifacts", "junit.xml")},
		{name: "other product dir rejected", link: "file://v1.37/other/junit.xml", wantErr: true},
		{name: "other version same product rejected", link: "v1.36/gke/junit.xml", wantErr: true},
		{name: "parent traversal rejected", link: "../other/junit.xml", wantErr: true},
		{name: "file scheme traversal rejected", link: "file://../../secret", wantErr: true},
		{name: "empty path", link: "file://#TestX", wantErr: true},
		{name: "self reference rejected", link: ".", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := resolveEvidencePath(productDir, tc.link)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", ref)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ref.Path != tc.wantPath {
				t.Errorf("Path = %q, want %q", ref.Path, tc.wantPath)
			}
			if ref.Fragment != tc.wantFragment {
				t.Errorf("Fragment = %q, want %q", ref.Fragment, tc.wantFragment)
			}
		})
	}
}

func TestClassifyArtifact(t *testing.T) {
	cases := map[string]artifactKind{
		"junit.xml":           artifactJUnit,
		"artifacts/JUnit.XML": artifactJUnit,
		"results.json":        artifactGoTestJSON,
		"e2e.log":             artifactE2ELog,
		"v1.37/acme/e2e.log":  artifactE2ELog,
		"install.log":         artifactNone,
		"driver-setup.log":    artifactNone,
		"cluster-config.json": artifactNone,
		"crd.json":            artifactNone,
		"report.xml":          artifactNone,
		"README.md":           artifactNone,
	}
	for name, want := range cases {
		if got := classifyArtifact(name); got != want {
			t.Errorf("classifyArtifact(%q) = %v, want %v", name, got, want)
		}
	}
}

const junitGotestsum = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="4" failures="0" errors="0" time="120.5">
  <testsuite tests="4" failures="0" time="120.5" name="sigs.k8s.io/ai-conformance/test">
    <testcase classname="test" name="TestSecureAcceleratorAccess" time="30.1"></testcase>
    <testcase classname="test" name="TestGangScheduling" time="40.2"></testcase>
    <testcase classname="test" name="TestGangScheduling/all_or_nothing" time="20.0"></testcase>
    <testcase classname="test" name="TestAcceleratorClusterAutoscaling" time="0.0">
      <skipped message="-autoscaler-node-pool-label not set"></skipped>
    </testcase>
  </testsuite>
</testsuites>
`

const junitFailing = `<testsuite tests="2" failures="1" name="test">
  <testcase name="TestSecureAcceleratorAccess" time="1.0">
    <failure message="Failed">expected isolation</failure>
  </testcase>
  <testcase name="TestGangScheduling" time="1.0"></testcase>
</testsuite>
`

const junitErrored = `<testsuites>
  <testsuite name="test">
    <testcase name="TestGangScheduling">
      <error message="panic">runtime error</error>
    </testcase>
  </testsuite>
</testsuites>
`

func TestParseJUnit(t *testing.T) {
	report, err := parseJUnit(strings.NewReader(junitGotestsum))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 0 {
		t.Errorf("Failures = %v, want none", report.Failures)
	}
	want := map[string]testOutcome{
		"TestSecureAcceleratorAccess":       outcomePass,
		"TestGangScheduling":                outcomePass,
		"TestGangScheduling/all_or_nothing": outcomePass,
		"TestAcceleratorClusterAutoscaling": outcomeSkip,
		"TestMissing":                       outcomeUnknown,
	}
	for name, o := range want {
		if got := report.outcomeFor(name); got != o {
			t.Errorf("outcomeFor(%q) = %v, want %v", name, got, o)
		}
	}
}

func TestParseJUnit_FailureAndError(t *testing.T) {
	report, err := parseJUnit(strings.NewReader(junitFailing))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 1 || report.Failures[0] != "TestSecureAcceleratorAccess" {
		t.Errorf("Failures = %v, want [TestSecureAcceleratorAccess]", report.Failures)
	}
	if report.outcomeFor("TestGangScheduling") != outcomePass {
		t.Errorf("TestGangScheduling should pass")
	}

	report, err = parseJUnit(strings.NewReader(junitErrored))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.outcomeFor("TestGangScheduling") != outcomeFail {
		t.Errorf("<error> should be treated as failure")
	}
}

func TestParseJUnit_Invalid(t *testing.T) {
	if _, err := parseJUnit(strings.NewReader("not xml")); err == nil {
		t.Error("expected error for non-XML input")
	}
	if _, err := parseJUnit(strings.NewReader("<testsuites></testsuites>")); err == nil {
		t.Error("expected error for XML with no testcases")
	}
}

const goTestJSONPassing = `{"Time":"2026-09-01T00:00:00Z","Action":"start","Package":"sigs.k8s.io/ai-conformance/test"}
{"Time":"2026-09-01T00:00:01Z","Action":"run","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestSecureAcceleratorAccess"}
{"Time":"2026-09-01T00:00:02Z","Action":"output","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestSecureAcceleratorAccess","Output":"=== RUN   TestSecureAcceleratorAccess\n"}
{"Time":"2026-09-01T00:00:03Z","Action":"pass","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestSecureAcceleratorAccess","Elapsed":2}
{"Time":"2026-09-01T00:00:04Z","Action":"run","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestGangScheduling"}
{"Time":"2026-09-01T00:00:05Z","Action":"pass","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestGangScheduling","Elapsed":1}
{"Time":"2026-09-01T00:00:06Z","Action":"run","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestAcceleratorClusterAutoscaling"}
{"Time":"2026-09-01T00:00:07Z","Action":"skip","Package":"sigs.k8s.io/ai-conformance/test","Test":"TestAcceleratorClusterAutoscaling","Elapsed":0}
{"Time":"2026-09-01T00:00:08Z","Action":"pass","Package":"sigs.k8s.io/ai-conformance/test","Elapsed":8}
`

const goTestJSONFailing = `{"Action":"run","Package":"p","Test":"TestGangScheduling"}
{"Action":"fail","Package":"p","Test":"TestGangScheduling","Elapsed":1}
{"Action":"fail","Package":"p","Elapsed":1}
`

const goTestJSONBuildFailure = `{"Action":"output","Package":"p","Output":"# p\n./x.go:1:1: syntax error\n"}
{"Action":"fail","Package":"p","Elapsed":0}
`

func TestParseGoTestJSON(t *testing.T) {
	report, err := parseGoTestJSON(strings.NewReader(goTestJSONPassing))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 0 {
		t.Errorf("Failures = %v, want none", report.Failures)
	}
	if report.outcomeFor("TestSecureAcceleratorAccess") != outcomePass {
		t.Error("TestSecureAcceleratorAccess should pass")
	}
	if report.outcomeFor("TestAcceleratorClusterAutoscaling") != outcomeSkip {
		t.Error("TestAcceleratorClusterAutoscaling should be skipped")
	}

	report, err = parseGoTestJSON(strings.NewReader(goTestJSONFailing))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 1 || report.Failures[0] != "TestGangScheduling" {
		t.Errorf("Failures = %v, want [TestGangScheduling] (package-level fail must not duplicate)", report.Failures)
	}

	report, err = parseGoTestJSON(strings.NewReader(goTestJSONBuildFailure))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 1 || report.Failures[0] != "package p" {
		t.Errorf("Failures = %v, want [package p]", report.Failures)
	}
}

func TestParseGoTestJSON_Invalid(t *testing.T) {
	if _, err := parseGoTestJSON(strings.NewReader("{not json")); err == nil {
		t.Error("expected error for malformed JSON")
	}
	if _, err := parseGoTestJSON(strings.NewReader(`{"Action":"output","Package":"p"}` + "\n")); err == nil {
		t.Error("expected error when no test events present")
	}
}

const e2eLogPassing = `=== RUN   TestSecureAcceleratorAccess
    device_probe_test.go:42: probing accelerators
--- PASS: TestSecureAcceleratorAccess (30.12s)
=== RUN   TestGangScheduling
=== RUN   TestGangScheduling/all_or_nothing
    --- PASS: TestGangScheduling/all_or_nothing (10.00s)
--- PASS: TestGangScheduling (40.00s)
=== RUN   TestAcceleratorClusterAutoscaling
    cluster_autoscaling_test.go:20: -autoscaler-node-pool-label not set, skipping
--- SKIP: TestAcceleratorClusterAutoscaling (0.00s)
PASS
ok  	sigs.k8s.io/ai-conformance/test	70.120s
`

const e2eLogFailing = `=== RUN   TestGangScheduling
--- FAIL: TestGangScheduling (5.00s)
FAIL
FAIL	sigs.k8s.io/ai-conformance/test	5.000s
`

func TestParseE2ELog(t *testing.T) {
	report, err := parseE2ELog(strings.NewReader(e2eLogPassing))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 0 {
		t.Errorf("Failures = %v, want none", report.Failures)
	}
	if report.outcomeFor("TestGangScheduling/all_or_nothing") != outcomePass {
		t.Error("indented subtest result line should be parsed")
	}
	if report.outcomeFor("TestAcceleratorClusterAutoscaling") != outcomeSkip {
		t.Error("TestAcceleratorClusterAutoscaling should be skipped")
	}

	report, err = parseE2ELog(strings.NewReader(e2eLogFailing))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Failures) != 1 || report.Failures[0] != "TestGangScheduling" {
		t.Errorf("Failures = %v, want [TestGangScheduling]", report.Failures)
	}

	if _, err := parseE2ELog(strings.NewReader("no results here\n")); err == nil {
		t.Error("expected error for log without result lines")
	}
}

func TestArtifactReport_SubtestAggregation(t *testing.T) {
	r := newArtifactReport()
	r.record("TestX/a", outcomePass)
	r.record("TestX/b", outcomeSkip)
	if got := r.outcomeFor("TestX"); got != outcomePass {
		t.Errorf("mixed pass/skip subtests = %v, want pass", got)
	}

	r = newArtifactReport()
	r.record("TestY/a", outcomeSkip)
	r.record("TestY/b", outcomeSkip)
	if got := r.outcomeFor("TestY"); got != outcomeSkip {
		t.Errorf("all-skip subtests = %v, want skip", got)
	}

	r = newArtifactReport()
	r.record("TestZ/a", outcomePass)
	r.record("TestZ/b", outcomeFail)
	if got := r.outcomeFor("TestZ"); got != outcomeFail {
		t.Errorf("any-fail subtests = %v, want fail", got)
	}

	r = newArtifactReport()
	r.record("TestW", outcomeFail)
	r.record("TestW", outcomePass)
	if got := r.outcomeFor("TestW"); got != outcomeFail {
		t.Errorf("fail must not be downgraded to pass, got %v", got)
	}
}

func writeArtifact(t *testing.T, dir, name, content string) evidenceRef {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return evidenceRef{Path: p}
}

func TestCheckArtifactEvidence_Fragment(t *testing.T) {
	dir := t.TempDir()
	ref := writeArtifact(t, dir, "junit.xml", junitGotestsum)
	reports := map[string]*artifactReport{}

	ref.Fragment = "TestSecureAcceleratorAccess"
	errs := checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestSecureAcceleratorAccess", "secure_accelerator_access", "Implemented", reports)
	if len(errs) != 0 {
		t.Errorf("passing fragment: errors=%v, want none", errs)
	}

	ref.Fragment = "TestDoesNotExist"
	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestDoesNotExist", "gang_scheduling", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "was not found in junit.xml") {
		t.Errorf("missing fragment: errors=%v, want one 'not found' error naming the artifact without fragment", errs)
	}

	ref.Fragment = "TestAcceleratorClusterAutoscaling"
	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestAcceleratorClusterAutoscaling", "cluster_autoscaling", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "was skipped in junit.xml") {
		t.Errorf("skipped+Implemented: errors=%v, want one 'was skipped' error", errs)
	}

	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestAcceleratorClusterAutoscaling", "cluster_autoscaling", "N/A", reports)
	if len(errs) != 0 {
		t.Errorf("skipped+N/A: errors=%v, want none", errs)
	}
}

func TestCheckArtifactEvidence_FailuresReportedOnce(t *testing.T) {
	dir := t.TempDir()
	ref := writeArtifact(t, dir, "junit.xml", junitFailing)
	reports := map[string]*artifactReport{}

	ref.Fragment = "TestGangScheduling"
	errs := checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestGangScheduling", "gang_scheduling", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "Test artifact junit.xml contains failing tests: TestSecureAcceleratorAccess") {
		t.Errorf("first reference: errors=%v, want one 'contains failing tests' error with fragment stripped from artifact name", errs)
	}

	ref.Fragment = ""
	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml", "secure_accelerator_access", "Implemented", reports)
	for _, e := range errs {
		if strings.Contains(e, "contains failing tests") {
			t.Errorf("artifact-level failure reported twice: %v", errs)
		}
	}
	if len(errs) != 1 || !strings.Contains(errs[0], "failed in") {
		t.Errorf("implicit failed test: errors=%v, want one 'failed in' error", errs)
	}

	ref.Fragment = "TestSecureAcceleratorAccess"
	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml#TestSecureAcceleratorAccess", "secure_accelerator_access", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "failed in") {
		t.Errorf("failed fragment: errors=%v, want one 'failed in' error", errs)
	}
}

func TestCheckArtifactEvidence_ImplicitTest(t *testing.T) {
	dir := t.TempDir()
	ref := writeArtifact(t, dir, "results.json", goTestJSONPassing)
	reports := map[string]*artifactReport{}

	errs := checkArtifactEvidence(ref, artifactGoTestJSON, "results.json", "secure_accelerator_access", "Implemented", reports)
	if len(errs) != 0 {
		t.Errorf("present implicit test: errors=%v, want none", errs)
	}

	errs = checkArtifactEvidence(ref, artifactGoTestJSON, "results.json", "cluster_autoscaling", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "was skipped") {
		t.Errorf("skipped implicit test + Implemented: errors=%v, want one skip error", errs)
	}

	errs = checkArtifactEvidence(ref, artifactGoTestJSON, "results.json", "cluster_autoscaling", "N/A", reports)
	if len(errs) != 0 {
		t.Errorf("skipped implicit test + N/A: errors=%v, want none", errs)
	}

	errs = checkArtifactEvidence(ref, artifactGoTestJSON, "results.json", "ai_inference", "Implemented", reports)
	if len(errs) != 0 {
		t.Errorf("non-auto-tested requirement: errors=%v, want none", errs)
	}

	ref = writeArtifact(t, dir, "e2e.log", e2eLogFailing)
	errs = checkArtifactEvidence(ref, artifactE2ELog, "e2e.log", "secure_accelerator_access", "Implemented", reports)
	if len(errs) != 2 || !strings.Contains(errs[0], "contains failing tests") || !strings.Contains(errs[1], "was not found") {
		t.Errorf("failing log + missing implicit test: errors=%v, want failing-tests error then not-found error", errs)
	}
}

func TestCheckArtifactEvidence_ParseErrorReportedOnce(t *testing.T) {
	dir := t.TempDir()
	ref := writeArtifact(t, dir, "junit.xml", "garbage")
	reports := map[string]*artifactReport{}

	errs := checkArtifactEvidence(ref, artifactJUnit, "junit.xml", "gang_scheduling", "Implemented", reports)
	if len(errs) != 1 || !strings.Contains(errs[0], "Invalid test artifact") {
		t.Errorf("first reference: errors=%v, want one parse error", errs)
	}
	errs = checkArtifactEvidence(ref, artifactJUnit, "junit.xml", "secure_accelerator_access", "Implemented", reports)
	if len(errs) != 0 {
		t.Errorf("second reference to unparseable artifact should be silent, got %v", errs)
	}
}

func TestCheckRequirementStatus(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		level, status string
		notes         string
		wantErrs      int
		wantWarnings  int
	}{
		{"MUST implemented", "secure_accelerator_access", "MUST", "Implemented", "", 0, 0},
		{"conditional MUST N/A with notes", "cluster_autoscaling", "MUST", "N/A", "ships no cluster autoscaler", 0, 1},
		{"conditional MUST N/A without notes", "cluster_autoscaling", "MUST", "N/A", "", 1, 1},
		{"conditional MUST N/A with whitespace notes", "cluster_autoscaling", "MUST", "N/A", "  \n", 1, 1},
		{"unconditional MUST N/A with notes", "secure_accelerator_access", "MUST", "N/A", "not applicable", 1, 0},
		{"MUST not implemented", "secure_accelerator_access", "MUST", "Not Implemented", "", 1, 0},
		{"MUST partially implemented", "secure_accelerator_access", "MUST", "Partially Implemented", "", 1, 0},
		{"conditional MUST not implemented", "cluster_autoscaling", "MUST", "Not Implemented", "", 1, 0},
		{"MUST invalid status", "secure_accelerator_access", "MUST", "Done", "", 2, 0},
		{"SHOULD not implemented", "dra_support", "SHOULD", "Not Implemented", "", 0, 0},
		{"SHOULD N/A with notes", "dra_support", "SHOULD", "N/A", "no such hardware", 0, 0},
		{"SHOULD N/A without notes", "dra_support", "SHOULD", "N/A", "", 1, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs, warnings := checkRequirementStatus(tc.id, tc.level, tc.status, tc.notes)
			if len(errs) != tc.wantErrs {
				t.Errorf("errors=%v, want %d", errs, tc.wantErrs)
			}
			if len(warnings) != tc.wantWarnings {
				t.Errorf("warnings=%v, want %d", warnings, tc.wantWarnings)
			}
		})
	}
}
