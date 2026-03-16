//go:build landscape

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// newTestServer creates a simple httptest.Server that returns the given status and body.
func newTestServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
}

func TestParseProductYAML_Valid(t *testing.T) {
	input := []byte(`
metadata:
  kubernetesVersion: v1.34
  platformName: "OpenShift Container Platform"
  platformVersion: "4.21"
  vendorName: "Red Hat"
  websiteUrl: "https://www.redhat.com/en/technologies/cloud-computing/openshift"
  productLogoUrl: "https://www.redhat.com/rhdc/managed-files/Logo-Red_Hat-OpenShift-A-Standard-RGB.svg"
  description: "Red Hat OpenShift Container Platform is an enterprise-ready Kubernetes container platform."
  contactEmailAddress: "test@example.com"
spec:
  accelerators: []
`)
	meta, err := parseProductYAML(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.PlatformName != "OpenShift Container Platform" {
		t.Errorf("PlatformName = %q, want %q", meta.PlatformName, "OpenShift Container Platform")
	}
	if meta.PlatformVersion != "4.21" {
		t.Errorf("PlatformVersion = %q, want %q", meta.PlatformVersion, "4.21")
	}
	if meta.VendorName != "Red Hat" {
		t.Errorf("VendorName = %q, want %q", meta.VendorName, "Red Hat")
	}
	if meta.WebsiteURL != "https://www.redhat.com/en/technologies/cloud-computing/openshift" {
		t.Errorf("WebsiteURL = %q, want correct URL", meta.WebsiteURL)
	}
	if meta.ProductLogoURL != "https://www.redhat.com/rhdc/managed-files/Logo-Red_Hat-OpenShift-A-Standard-RGB.svg" {
		t.Errorf("ProductLogoURL = %q, want correct URL", meta.ProductLogoURL)
	}
	if meta.Description != "Red Hat OpenShift Container Platform is an enterprise-ready Kubernetes container platform." {
		t.Errorf("Description = %q, want correct description", meta.Description)
	}
	if meta.KubernetesVersion != "v1.34" {
		t.Errorf("KubernetesVersion = %q, want %q", meta.KubernetesVersion, "v1.34")
	}
}

func TestParseProductYAML_EmptyPlatformName(t *testing.T) {
	input := []byte(`
metadata:
  kubernetesVersion: v1.34
  platformName: ""
  vendorName: "Red Hat"
spec:
  accelerators: []
`)
	_, err := parseProductYAML(input)
	if err == nil {
		t.Fatal("expected error for empty platformName, got nil")
	}
}

func TestParseProductYAML_MissingPlatformName(t *testing.T) {
	input := []byte(`
metadata:
  kubernetesVersion: v1.34
  vendorName: "Red Hat"
spec:
  accelerators: []
`)
	_, err := parseProductYAML(input)
	if err == nil {
		t.Fatal("expected error for missing platformName, got nil")
	}
}

func TestParseProductYAML_MissingWebsiteUrl(t *testing.T) {
	input := []byte(`
metadata:
  kubernetesVersion: v1.34
  platformName: "CoreWeave Kubernetes Service"
  vendorName: "CoreWeave"
spec:
  accelerators: []
`)
	meta, err := parseProductYAML(input)
	if err != nil {
		t.Fatalf("should not error on missing websiteUrl: %v", err)
	}
	if meta.WebsiteURL != "" {
		t.Errorf("WebsiteURL = %q, want empty", meta.WebsiteURL)
	}
}

func TestParseProductYAML_SnakeCaseFields(t *testing.T) {
	input := []byte(`
metadata:
  kubernetes_version: v1.34
  platform_name: "CoreWeave Kubernetes Service (CKS)"
  platform_version: v1.34
  vendor_name: CoreWeave
  website_url: https://www.coreweave.com/products/coreweave-kubernetes-service
  product_logo_url: https://yoyo.dyne/assets/turbo-encabulator.svg
  description: "CKS is a managed Kubernetes environment."
  contact_email_address: cks@coreweave.com
spec:
  accelerators: []
`)
	meta, err := parseProductYAML(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.PlatformName != "CoreWeave Kubernetes Service (CKS)" {
		t.Errorf("PlatformName = %q, want %q", meta.PlatformName, "CoreWeave Kubernetes Service (CKS)")
	}
	if meta.VendorName != "CoreWeave" {
		t.Errorf("VendorName = %q, want %q", meta.VendorName, "CoreWeave")
	}
	if meta.WebsiteURL != "https://www.coreweave.com/products/coreweave-kubernetes-service" {
		t.Errorf("WebsiteURL = %q, want correct URL", meta.WebsiteURL)
	}
}

// --- Task 2: URL Normalization ---

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://example.com/", "https://example.com"},
		{"https://www.example.com", "https://example.com"},
		{"https://cloud.google.com/kubernetes-engine/", "https://cloud.google.com/kubernetes-engine"},
		{"https://WWW.Example.COM/Path/", "https://example.com/path"},
		{"", ""},
		{"https://example.com", "https://example.com"},
		{"https://www.example.com/", "https://example.com"},
	}
	for _, tc := range tests {
		got := normalizeURL(tc.input)
		if got != tc.want {
			t.Errorf("normalizeURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// --- Task 3: findEntryInLandscape ---

const landscapeFixture = `landscape:
  - category:
    name: Certified Kubernetes - Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Platform
        items:
          - item:
            name: Red Hat OpenShift
            description: OpenShift helps organizations deploy.
            homepage_url: https://www.redhat.com/en/technologies/cloud-computing/openshift
            logo: red-hat-open-shift.svg
            crunchbase: https://www.crunchbase.com/organization/red-hat
          - item:
            name: Google Kubernetes Engine
            description: GKE is a managed Kubernetes service.
            homepage_url: https://cloud.google.com/kubernetes-engine
            logo: google-kubernetes-engine.svg
            crunchbase: https://www.crunchbase.com/organization/google
            second_path:
              - "Platform / Certified Kubernetes - AI Platform"
`

func TestFindEntryInLandscape_Found(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://www.redhat.com/en/technologies/cloud-computing/openshift")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry, got nil")
	}
	if entry.Name != "Red Hat OpenShift" {
		t.Errorf("Name = %q, want %q", entry.Name, "Red Hat OpenShift")
	}
	if entry.HomepageURL != "https://www.redhat.com/en/technologies/cloud-computing/openshift" {
		t.Errorf("HomepageURL = %q", entry.HomepageURL)
	}
	if entry.HasAIPlatformSecondPath {
		t.Error("HasAIPlatformSecondPath should be false for OpenShift")
	}
}

func TestFindEntryInLandscape_FoundWithSecondPath(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://cloud.google.com/kubernetes-engine")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry, got nil")
	}
	if entry.Name != "Google Kubernetes Engine" {
		t.Errorf("Name = %q, want %q", entry.Name, "Google Kubernetes Engine")
	}
	if !entry.HasAIPlatformSecondPath {
		t.Error("HasAIPlatformSecondPath should be true for GKE")
	}
}

func TestFindEntryInLandscape_NotFound(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://nonexistent.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry, got %+v", entry)
	}
}

func TestFindEntryInLandscape_NormalizedMatch(t *testing.T) {
	// Search with trailing slash and www - should still match after normalization
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://cloud.google.com/kubernetes-engine/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry with normalized URL, got nil")
	}
	if entry.Name != "Google Kubernetes Engine" {
		t.Errorf("Name = %q, want %q", entry.Name, "Google Kubernetes Engine")
	}
}

func TestFindEntryInLandscape_LineIndices(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://www.redhat.com/en/technologies/cloud-computing/openshift")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry, got nil")
	}
	// Verify ItemLineIndex points to the "- item:" line
	lines := strings.Split(landscapeFixture, "\n")
	if entry.ItemLineIndex < 0 || entry.ItemLineIndex >= len(lines) {
		t.Fatalf("ItemLineIndex %d out of range", entry.ItemLineIndex)
	}
	itemLine := strings.TrimSpace(lines[entry.ItemLineIndex])
	if itemLine != "- item:" {
		t.Errorf("line at ItemLineIndex = %q, want '- item:'", itemLine)
	}
	// LastFieldLineIndex should be at or after ItemLineIndex
	if entry.LastFieldLineIndex < entry.ItemLineIndex {
		t.Errorf("LastFieldLineIndex %d < ItemLineIndex %d", entry.LastFieldLineIndex, entry.ItemLineIndex)
	}
}

// --- Task 4: insertSecondPath ---

func TestInsertSecondPath_NoExistingSecondPath(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://www.redhat.com/en/technologies/cloud-computing/openshift")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry")
	}

	result := insertSecondPath([]byte(landscapeFixture), entry)
	resultStr := string(result)

	if !strings.Contains(resultStr, "second_path:") {
		t.Error("result should contain 'second_path:'")
	}
	if !strings.Contains(resultStr, `"Platform / Certified Kubernetes - AI Platform"`) {
		t.Error("result should contain AI Platform second_path value")
	}

	// Verify the second_path is correctly indented (12 spaces for key, 14 for list item)
	lines := strings.Split(resultStr, "\n")
	foundKey := false
	for i, line := range lines {
		if strings.TrimSpace(line) == "second_path:" && strings.HasPrefix(line, "            second_path:") {
			foundKey = true
			// Next line should be the list item with 14 spaces
			if i+1 < len(lines) {
				nextLine := lines[i+1]
				if !strings.HasPrefix(nextLine, "              - ") {
					t.Errorf("list item not properly indented: %q", nextLine)
				}
			}
		}
	}
	if !foundKey {
		t.Error("did not find properly indented second_path key")
	}
}

func TestInsertSecondPath_HasExistingSecondPathNotAI(t *testing.T) {
	// Create a fixture where the entry already has second_path but NOT AI Platform
	fixture := `landscape:
  - category:
    name: Certified Kubernetes - Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Platform
        items:
          - item:
            name: Red Hat OpenShift
            description: OpenShift helps organizations deploy.
            homepage_url: https://www.redhat.com/en/technologies/cloud-computing/openshift
            logo: red-hat-open-shift.svg
            crunchbase: https://www.crunchbase.com/organization/red-hat
            second_path:
              - "Platform / Certified Kubernetes - Distribution"
`
	entry, err := findEntryInLandscape([]byte(fixture),
		"https://www.redhat.com/en/technologies/cloud-computing/openshift")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry")
	}
	if entry.HasAIPlatformSecondPath {
		t.Fatal("should not have AI Platform second path yet")
	}

	result := insertSecondPath([]byte(fixture), entry)
	resultStr := string(result)

	// Should have both the existing and new second_path items
	if !strings.Contains(resultStr, "Certified Kubernetes - Distribution") {
		t.Error("should still contain existing second_path value")
	}
	if !strings.Contains(resultStr, "Certified Kubernetes - AI Platform") {
		t.Error("should contain new AI Platform second_path value")
	}

	// Should only have one second_path key (not two)
	count := strings.Count(resultStr, "second_path:")
	if count != 1 {
		t.Errorf("expected 1 second_path key, got %d", count)
	}
}

func TestInsertSecondPath_AlreadyHasAIPlatform(t *testing.T) {
	entry, err := findEntryInLandscape([]byte(landscapeFixture),
		"https://cloud.google.com/kubernetes-engine")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find GKE entry")
	}
	if !entry.HasAIPlatformSecondPath {
		t.Fatal("GKE should already have AI Platform second path")
	}

	result := insertSecondPath([]byte(landscapeFixture), entry)
	// Should be unchanged
	if string(result) != landscapeFixture {
		t.Error("result should be unchanged when AI Platform already present")
	}
}

// --- Task 5: sanitizeLogoName, insertNewEntry, downloadLogo ---

func TestSanitizeLogoName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"OpenShift Container Platform", "openshift-container-platform.svg"},
		{"Google Kubernetes Engine", "google-kubernetes-engine.svg"},
		{"CoreWeave Kubernetes Service (CKS)", "coreweave-kubernetes-service-cks.svg"},
		{"Simple", "simple.svg"},
		{"Already-Dashed", "already-dashed.svg"},
		{"  Spaces & Symbols!  ", "spaces-symbols.svg"},
	}
	for _, tc := range tests {
		got := sanitizeLogoName(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeLogoName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestInsertNewEntry_EmptyItems(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - AI Platform
        items: []
`
	meta := &ProductMeta{
		PlatformName: "TestPlatform",
		Description:  "A test platform for AI.",
		WebsiteURL:   "https://test.example.com",
	}
	result, err := insertNewEntry([]byte(fixture), meta, "test-platform.svg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := string(result)

	if strings.Contains(resultStr, "items: []") {
		t.Error("should have replaced 'items: []'")
	}
	if !strings.Contains(resultStr, "name: TestPlatform") {
		t.Error("should contain the platform name")
	}
	if !strings.Contains(resultStr, "homepage_url: https://test.example.com") {
		t.Error("should contain homepage_url")
	}
	if !strings.Contains(resultStr, "logo: test-platform.svg") {
		t.Error("should contain logo filename")
	}
	if !strings.Contains(resultStr, "A test platform for AI.") {
		t.Error("should contain description")
	}
}

func TestInsertNewEntry_ExistingItems(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - AI Platform
        items:
          - item:
            name: Existing Platform
            description: Already here.
            homepage_url: https://existing.example.com
            logo: existing.svg
`
	meta := &ProductMeta{
		PlatformName: "NewPlatform",
		Description:  "A new platform.",
		WebsiteURL:   "https://new.example.com",
	}
	result, err := insertNewEntry([]byte(fixture), meta, "new-platform.svg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := string(result)

	// Should still have the existing entry
	if !strings.Contains(resultStr, "name: Existing Platform") {
		t.Error("should still contain existing platform")
	}
	// Should have the new entry
	if !strings.Contains(resultStr, "name: NewPlatform") {
		t.Error("should contain new platform name")
	}
	if !strings.Contains(resultStr, "homepage_url: https://new.example.com") {
		t.Error("should contain new homepage_url")
	}
}

func TestInsertNewEntry_SubcategoryNotFound(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Platform
    subcategories:
      - subcategory:
        name: Something Else
        items: []
`
	meta := &ProductMeta{
		PlatformName: "TestPlatform",
		Description:  "A test platform.",
		WebsiteURL:   "https://test.example.com",
	}
	_, err := insertNewEntry([]byte(fixture), meta, "test.svg")
	if err == nil {
		t.Fatal("expected error when subcategory not found")
	}
}

func TestDownloadLogo_BadURL(t *testing.T) {
	err := downloadLogo("http://127.0.0.1:1/nonexistent", "/tmp/test-logo-bad.svg")
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestDownloadLogo_HTTPError(t *testing.T) {
	// Use httptest for a 404 response
	ts := newTestServer(404, "not found")
	defer ts.Close()

	err := downloadLogo(ts.URL, t.TempDir()+"/logo.svg")
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error should mention 404, got: %v", err)
	}
}

func TestDownloadLogo_Success(t *testing.T) {
	ts := newTestServer(200, "<svg>test</svg>")
	defer ts.Close()

	destPath := t.TempDir() + "/logo.svg"
	err := downloadLogo(ts.URL, destPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("could not read downloaded file: %v", err)
	}
	if string(data) != "<svg>test</svg>" {
		t.Errorf("file content = %q, want %q", string(data), "<svg>test</svg>")
	}
}
