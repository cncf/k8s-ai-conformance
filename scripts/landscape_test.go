//go:build landscape

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	err := downloadLogo("http://127.0.0.1:1/nonexistent", filepath.Join(t.TempDir(), "test-logo-bad.svg"))
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

func TestDownloadLogo_InvalidScheme(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"file scheme", "file:///etc/passwd"},
		{"ftp scheme", "ftp://example.com/logo.svg"},
		{"no scheme", "://example.com/logo.svg"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := downloadLogo(tc.url, filepath.Join(t.TempDir(), "logo.svg"))
			if err == nil {
				t.Fatal("expected error for invalid scheme")
			}
			if !strings.Contains(err.Error(), "scheme") {
				t.Errorf("error should mention scheme, got: %v", err)
			}
		})
	}
}

func TestDownloadLogo_OversizedFile(t *testing.T) {
	// Create a test server that returns a response larger than maxLogoSize
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write slightly more than maxLogoSize (10MB + 1KB)
		w.WriteHeader(200)
		buf := make([]byte, 1024)
		for i := range buf {
			buf[i] = 'x'
		}
		for written := 0; written <= maxLogoSize; written += len(buf) {
			w.Write(buf)
		}
	}))
	defer ts.Close()

	destPath := filepath.Join(t.TempDir(), "oversized.svg")
	err := downloadLogo(ts.URL, destPath)
	if err == nil {
		t.Fatal("expected error for oversized file")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("error should mention size limit, got: %v", err)
	}
	// File should have been cleaned up
	if _, err := os.Stat(destPath); err == nil {
		t.Error("oversized file should have been removed")
	}
}

// --- Edge Case Tests ---

func TestIdempotentSecondPath(t *testing.T) {
	// Start with a landscape entry WITHOUT second_path
	fixture := `landscape:
  - category:
    name: Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Distribution
        items:
          - item:
            name: TestPlatform
            description: A test platform.
            homepage_url: https://test.example.com
            logo: test.svg
            crunchbase: https://www.crunchbase.com/organization/test
`
	// Step 1: Find the entry, confirm no AI Platform second_path
	entry, err := findEntryInLandscape([]byte(fixture), "https://test.example.com")
	if err != nil {
		t.Fatalf("step 1: unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("step 1: expected to find entry, got nil")
	}
	if entry.HasAIPlatformSecondPath {
		t.Fatal("step 1: HasAIPlatformSecondPath should be false before insertion")
	}

	// Step 2: Insert second_path
	modified := insertSecondPath([]byte(fixture), entry)
	modifiedStr := string(modified)
	if !strings.Contains(modifiedStr, "second_path:") {
		t.Fatal("step 2: modified data should contain 'second_path:'")
	}
	if !strings.Contains(modifiedStr, "Certified Kubernetes - AI Platform") {
		t.Fatal("step 2: modified data should contain AI Platform path")
	}

	// Step 3: Find the entry again in the MODIFIED data
	entry2, err := findEntryInLandscape(modified, "https://test.example.com")
	if err != nil {
		t.Fatalf("step 3: unexpected error: %v", err)
	}
	if entry2 == nil {
		t.Fatal("step 3: expected to find entry in modified data, got nil")
	}
	if !entry2.HasAIPlatformSecondPath {
		t.Fatal("step 3: HasAIPlatformSecondPath should be true after insertion")
	}

	// Step 4: Insert again — should be idempotent (unchanged)
	result := insertSecondPath(modified, entry2)
	if string(result) != string(modified) {
		t.Error("step 4: second insertSecondPath should return data unchanged (idempotent)")
	}
}

func TestURLNormalizationEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"URL with port", "https://example.com:8080/path/", "https://example.com:8080/path"},
		{"URL with query params", "https://example.com/path?foo=bar", "https://example.com/path?foo=bar"},
		{"URL with fragment", "https://example.com/path#section", "https://example.com/path#section"},
		{"URL without scheme", "example.com", "example.com"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeURL(tc.input)
			if got != tc.want {
				t.Errorf("normalizeURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestMissingMetadataFields(t *testing.T) {
	t.Run("missing vendorName succeeds", func(t *testing.T) {
		input := []byte(`
metadata:
  platformName: "TestPlatform"
  description: "A test platform."
spec:
  accelerators: []
`)
		meta, err := parseProductYAML(input)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if meta.VendorName != "" {
			t.Errorf("VendorName = %q, want empty", meta.VendorName)
		}
		if meta.PlatformName != "TestPlatform" {
			t.Errorf("PlatformName = %q, want %q", meta.PlatformName, "TestPlatform")
		}
	})

	t.Run("missing description succeeds", func(t *testing.T) {
		input := []byte(`
metadata:
  platformName: "TestPlatform"
  vendorName: "TestVendor"
spec:
  accelerators: []
`)
		meta, err := parseProductYAML(input)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if meta.Description != "" {
			t.Errorf("Description = %q, want empty", meta.Description)
		}
	})

	t.Run("missing productLogoUrl succeeds", func(t *testing.T) {
		input := []byte(`
metadata:
  platformName: "TestPlatform"
  vendorName: "TestVendor"
  description: "A test."
spec:
  accelerators: []
`)
		meta, err := parseProductYAML(input)
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}
		if meta.ProductLogoURL != "" {
			t.Errorf("ProductLogoURL = %q, want empty", meta.ProductLogoURL)
		}
	})

	t.Run("empty metadata section errors", func(t *testing.T) {
		input := []byte(`
metadata:
spec:
  accelerators: []
`)
		_, err := parseProductYAML(input)
		if err == nil {
			t.Fatal("expected error for empty metadata section, got nil")
		}
	})

	t.Run("invalid YAML errors", func(t *testing.T) {
		input := []byte(`not: valid: yaml: [[[`)
		_, err := parseProductYAML(input)
		if err == nil {
			t.Fatal("expected error for invalid YAML, got nil")
		}
	})
}

func TestFindEntryInLandscapeWithRealFormat(t *testing.T) {
	realFixture := `landscape:
  - category:
    name: Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Distribution
        items:
          - item:
            name: Red Hat OpenShift
            description: OpenShift® helps organizations focus on building and scaling their business with fully supported enterprise Kubernetes by Red Hat®.
            homepage_url: https://www.redhat.com/en/technologies/cloud-computing/openshift
            repo_url: https://github.com/openshift/kubernetes
            logo: red-hat-open-shift.svg
            twitter: https://twitter.com/openshift
            crunchbase: https://www.crunchbase.com/organization/red-hat
            second_path:
              - "Platform / Certified Kubernetes - AI Platform"
          - item:
            name: RKE Government
            description: A Kubernetes distribution focused on enabling Federal government compliance-based use cases.
            homepage_url: https://docs.rke2.io/
            repo_url: https://github.com/rancher/rke2
            logo: rke-government.svg
            crunchbase: https://www.crunchbase.com/organization/suse
      - subcategory:
        name: Certified Kubernetes - AI Platform
        items: []
`

	t.Run("find OpenShift with AI Platform second_path", func(t *testing.T) {
		entry, err := findEntryInLandscape([]byte(realFixture),
			"https://www.redhat.com/en/technologies/cloud-computing/openshift")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry == nil {
			t.Fatal("expected to find OpenShift entry, got nil")
		}
		if entry.Name != "Red Hat OpenShift" {
			t.Errorf("Name = %q, want %q", entry.Name, "Red Hat OpenShift")
		}
		if !entry.HasAIPlatformSecondPath {
			t.Error("OpenShift should have HasAIPlatformSecondPath=true")
		}
	})

	t.Run("find RKE Government without AI Platform second_path", func(t *testing.T) {
		entry, err := findEntryInLandscape([]byte(realFixture),
			"https://docs.rke2.io/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry == nil {
			t.Fatal("expected to find RKE Government entry, got nil")
		}
		if entry.Name != "RKE Government" {
			t.Errorf("Name = %q, want %q", entry.Name, "RKE Government")
		}
		if entry.HasAIPlatformSecondPath {
			t.Error("RKE Government should have HasAIPlatformSecondPath=false")
		}
	})

	t.Run("insertSecondPath on RKE Government produces valid output", func(t *testing.T) {
		entry, err := findEntryInLandscape([]byte(realFixture),
			"https://docs.rke2.io/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry == nil {
			t.Fatal("expected to find RKE Government entry")
		}

		result := insertSecondPath([]byte(realFixture), entry)
		resultStr := string(result)

		// Should contain second_path with proper indentation
		if !strings.Contains(resultStr, "second_path:") {
			t.Error("result should contain 'second_path:'")
		}
		if !strings.Contains(resultStr, `"Platform / Certified Kubernetes - AI Platform"`) {
			t.Error("result should contain AI Platform second_path value")
		}

		// Verify proper indentation: 12 spaces for second_path key
		lines := strings.Split(resultStr, "\n")
		foundKey := false
		for i, line := range lines {
			// Find the NEW second_path (for RKE Government, not OpenShift's existing one)
			if strings.TrimSpace(line) == "second_path:" {
				// Check if this follows the RKE Government crunchbase line
				if i > 0 && strings.Contains(lines[i-1], "suse") {
					foundKey = true
					if !strings.HasPrefix(line, "            second_path:") {
						t.Errorf("second_path key not properly indented: %q", line)
					}
					if i+1 < len(lines) {
						nextLine := lines[i+1]
						if !strings.HasPrefix(nextLine, "              - ") {
							t.Errorf("second_path list item not properly indented: %q", nextLine)
						}
					}
				}
			}
		}
		if !foundKey {
			t.Error("did not find properly placed second_path for RKE Government")
		}
	})
}

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal name",
			input: "OpenShift Container Platform",
			want:  "openshift-container-platform",
		},
		{
			name:  "special characters",
			input: "Platform (v2.0) #1",
			want:  "platform-v2-0-1",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeBranchName(tc.input)
			if got != tc.want {
				t.Errorf("sanitizeBranchName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}

	// Long name test: verify truncation to max 50 chars and no trailing dash
	t.Run("long name truncated to 50 chars max", func(t *testing.T) {
		longName := "This Is A Very Long Platform Name That Should Definitely Be Truncated To Fifty Characters"
		got := sanitizeBranchName(longName)
		if len(got) > 50 {
			t.Errorf("sanitizeBranchName(%q) length = %d, want <= 50", longName, len(got))
		}
		if strings.HasSuffix(got, "-") {
			t.Errorf("sanitizeBranchName(%q) = %q, should not end with '-'", longName, got)
		}
		if got == "" {
			t.Error("sanitizeBranchName of long name should not be empty")
		}
	})
}

func TestNameMatchLevel(t *testing.T) {
	tests := []struct {
		product string
		entry   string
		want    int
	}{
		{"Azure Kubernetes Service (AKS)", "Azure Kubernetes Service (AKS)", nameMatchExact},
		{"RKE2", "RKE2", nameMatchExact},
		{"DaoCloud Enterprise", "DaoCloud Enterprise", nameMatchExact},
		{"Alibaba Cloud Container Service for Kubernetes", "Alibaba Cloud Container Service for Kubernetes", nameMatchExact},
		{"Linode Kubernetes Engine (LKE)", "Linode Kubernetes Engine", nameMatchContained},
		{"vSphere Kubernetes Service", "VMware vSphere Kubernetes Service", nameMatchContained},
		{"Alauda Container Platform", "Alauda Container Platform (ACP)", nameMatchContained},
		{"Cloud Container Distribution", "Ericsson Cloud Container Distribution", nameMatchContained},
		{"BKS", "Breqwatr BKS", nameMatchContained},
		{"OVHcloud Managed Kubernetes Service", "OVH Managed Kubernetes Service", nameMatchContained},
		{"DaoCloud Enterprise", "DaoCloud (KCSP)", nameMatchNone},
		{"Amazon Elastic Kubernetes Service (EKS)", "Amazon Elastic Container Service for Kubernetes (EKS)", nameMatchNone},
		{"Azure Kubernetes Service (AKS)", "Azure Kubernetes Service (AKS) (Wasm)", nameMatchContained},
		{"k0s", "k0smotron", nameMatchNone},
		{"JDOS", "JD Cloud (KCSP)", nameMatchNone},
		{"NVIDIA NIM on EKS", "NVIDIA (member)", nameMatchNone},
		{"vSphere Kubernetes Service", "Kubernetes", nameMatchNone},
		{"JCS for Kubernetes", "Kubernetes", nameMatchNone},
		{"Cloud Container Distribution", "Distribution", nameMatchNone},
		{"Cloud Container Engine", "Huawei Cloud Container Engine (CCE)", nameMatchContained},
	}
	for _, tt := range tests {
		if got, _ := nameMatchLevel(tt.product, tt.entry); got != tt.want {
			t.Errorf("nameMatchLevel(%q, %q) = %d, want %d", tt.product, tt.entry, got, tt.want)
		}
	}
}

func TestFindEntryByName_PrefersExactOverContained(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Certified Kubernetes - Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Platform
        items:
          - item:
            name: Azure Kubernetes Service (AKS) (Wasm)
            homepage_url: https://example.com/wasm
            logo: aks-wasm.svg
          - item:
            name: Azure Kubernetes Service (AKS)
            homepage_url: https://example.com/aks
            logo: aks.svg
`
	entry, err := findEntryByName([]byte(fixture), "Microsoft", "Azure Kubernetes Service (AKS)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry, got nil")
	}
	if entry.Name != "Azure Kubernetes Service (AKS)" {
		t.Errorf("Name = %q, want exact match preferred", entry.Name)
	}
}

func TestFindEntryByName_NotFound(t *testing.T) {
	entry, err := findEntryByName([]byte(landscapeFixture), "Nobody", "Totally Unrelated Product")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry, got %+v", entry)
	}
}

func TestFindEntryByName_PrefersLargestOverlap(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Certified Kubernetes - Platform
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Platform
        items:
          - item:
            name: OVH (KCSP)
            homepage_url: https://example.com/kcsp
            logo: ovh-kcsp.svg
          - item:
            name: VMware vSphere Kubernetes Service
            homepage_url: https://example.com/vks
            logo: vks.svg
`
	entry, err := findEntryByName([]byte(fixture), "Broadcom", "vSphere Kubernetes Service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected to find entry, got nil")
	}
	if entry.Name != "VMware vSphere Kubernetes Service" {
		t.Errorf("Name = %q, want VMware entry", entry.Name)
	}
}

func TestFindEntryInLandscape_SkipsNonProductCards(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Kubernetes Certified Service Provider
    subcategories:
      - subcategory:
        name: Kubernetes Certified Service Provider
        items:
          - item:
            name: Huawei (KCSP)
            homepage_url: https://huawei.example.com
            logo: huawei-kcsp.svg
`
	entry, err := findEntryInLandscape([]byte(fixture), "https://huawei.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected KCSP card to be skipped, got %+v", entry)
	}
}

func TestFindEntryByName_VendorQualified(t *testing.T) {
	fixture := `landscape:
  - category:
    name: Certified Kubernetes - Hosted
    subcategories:
      - subcategory:
        name: Certified Kubernetes - Hosted
        items:
          - item:
            name: Baidu Cloud Container Engine
            homepage_url: https://example.com/baidu-cce
            logo: baidu.svg
          - item:
            name: Huawei Cloud Container Engine (CCE)
            homepage_url: https://example.com/huawei-cce
            logo: huawei.svg
`
	entry, err := findEntryByName([]byte(fixture), "Baidu Cloud", "CCE（Cloud Container Engine）")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil || entry.Name != "Baidu Cloud Container Engine" {
		t.Fatalf("want Baidu entry, got %+v", entry)
	}

	entry, err = findEntryByName([]byte(fixture), "Huawei", "Cloud Container Engine")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil || entry.Name != "Huawei Cloud Container Engine (CCE)" {
		t.Fatalf("want Huawei entry (vendor tiebreak), got %+v", entry)
	}
}
