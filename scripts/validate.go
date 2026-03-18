package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const landscapeURL = "https://raw.githubusercontent.com/cncf/landscape/master/landscape.yml"

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// LandscapeData represents the top-level structure of the CNCF landscape YAML
type LandscapeData struct {
	Landscape []LandscapeCategory `yaml:"landscape"`
}

// LandscapeCategory represents a category in the CNCF landscape
type LandscapeCategory struct {
	Name          string                 `yaml:"name"`
	Subcategories []LandscapeSubcategory `yaml:"subcategories"`
}

// LandscapeSubcategory represents a subcategory within a landscape category
type LandscapeSubcategory struct {
	Name  string          `yaml:"name"`
	Items []LandscapeItem `yaml:"items"`
}

// LandscapeItem represents an individual item/entry in the landscape
type LandscapeItem struct {
	Name string `yaml:"name"`
}

// memberSuffixes are the parenthetical suffixes appended to member names in the landscape
var memberSuffixes = []string{" (member)", " (supporter)"}

// fetchCNCFMembers fetches the CNCF landscape YAML and returns a set of member
// names from the "CNCF Members" category. The returned map keys are the member
// names with the trailing "(member)"/"(supporter)" suffix stripped.
func fetchCNCFMembers() (map[string]bool, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := doRequest(client, http.MethodGet, landscapeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CNCF landscape: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch CNCF landscape: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read CNCF landscape response: %v", err)
	}

	var data LandscapeData
	if err := yaml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse CNCF landscape YAML: %v", err)
	}

	members := make(map[string]bool)
	found := false
	for _, category := range data.Landscape {
		if category.Name == "CNCF Members" {
			found = true
			for _, sub := range category.Subcategories {
				for _, item := range sub.Items {
					name := item.Name
					for _, suffix := range memberSuffixes {
						name = strings.TrimSuffix(name, suffix)
					}
					name = strings.TrimSpace(name)
					if name != "" {
						members[name] = true
					}
				}
			}
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("could not find 'CNCF Members' category in landscape data")
	}

	return members, nil
}

// Requirement represents a single checklist item
type Requirement struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Level       string   `yaml:"level"`
	Status      string   `yaml:"status"`
	Evidence    []string `yaml:"evidence"`
	Notes       string   `yaml:"notes"`
}

// ConformanceDoc represents the structure of the YAML files
type ConformanceDoc struct {
	Metadata map[string]interface{}   `yaml:"metadata"`
	Spec     map[string][]Requirement `yaml:"spec"`
}

var validStatuses = map[string]bool{
	"Implemented":           true,
	"Not Implemented":       true,
	"Partially Implemented": true,
	"N/A":                   true,
}

// k8sConformanceURLPattern matches URLs like:
// https://github.com/cncf/k8s-conformance/tree/master/v1.34/gke
// https://github.com/cncf/k8s-conformance/tree/main/v1.34/gke
var k8sConformanceURLPattern = regexp.MustCompile(`^https://github\.com/cncf/k8s-conformance/tree/(master|main)/v\d+\.\d+/[^/]+/?$`)

var metadataFields = map[string]bool{
	"kubernetesVersion":   true,
	"platformName":        true,
	"platformVersion":     true,
	"vendorName":          true,
	"websiteUrl":          true,
	"repoUrl":             false, // Optional
	"documentationUrl":    true,
	"productLogoUrl":      true,
	"description":         true,
	"contactEmailAddress": true,
	"k8sConformanceUrl":   true, // Required: URL to k8s-conformance submission
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/validate.go <path_to_product.yaml> ...")
		os.Exit(1)
	}

	// Fetch CNCF member list once for all validations
	fmt.Println("Fetching CNCF member list from landscape...")
	cncfMembers, err := fetchCNCFMembers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded %d CNCF members from landscape.\n", len(cncfMembers))

	success := true
	for _, path := range os.Args[1:] {
		if !validateProduct(path, cncfMembers) {
			success = false
		}
	}

	if !success {
		os.Exit(1)
	}
}

func validateProduct(path string, cncfMembers map[string]bool) bool {
	fmt.Printf("Validating %s...\n", path)

	// Extract version
	re := regexp.MustCompile(`(?:^|/)(v\d+\.\d+)/`)
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		fmt.Printf("Error: Could not determine version from path %s\n", path)
		return false
	}
	version := matches[1]
	schemaVersion := version[1:] // 1.33

	// Find schema
	schemaPath := filepath.Join("docs", fmt.Sprintf("AIConformance-%s.yaml", schemaVersion))
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		fmt.Printf("Error: Schema not found at %s\n", schemaPath)
		return false
	}

	// Load Schema
	schema, err := loadYaml(schemaPath)
	if err != nil {
		fmt.Printf("Error loading schema: %v\n", err)
		return false
	}

	// Load Product
	product, err := loadYaml(path)
	if err != nil {
		fmt.Printf("Error loading product: %v\n", err)
		return false
	}

	errors := []string{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	addError := func(msg string) {
		mu.Lock()
		errors = append(errors, msg)
		mu.Unlock()
	}

	// Validate Metadata
	if product.Metadata == nil {
		addError("Missing 'metadata' section")
	} else {
		for field, required := range metadataFields {
			snakeField := toSnakeCase(field)

			val, ok := product.Metadata[field]
			if !ok {
				val, ok = product.Metadata[snakeField]
			}

			if !ok || val == nil {
				if required {
					addError(fmt.Sprintf("Missing metadata field: %s (or %s)", field, snakeField))
				}
				continue
			}

			strVal, isStr := val.(string)
			if !isStr {
				addError(fmt.Sprintf("Metadata field %s is not a string", field))
				continue
			}

			if strVal == "" {
				if required {
					addError(fmt.Sprintf("Metadata field %s is empty", field))
				}
			} else if strings.HasPrefix(strVal, "[") && strings.HasSuffix(strVal, "]") {
				addError(fmt.Sprintf("Metadata field %s has placeholder value: %s", field, strVal))
			} else if field == "k8sConformanceUrl" || snakeField == "k8s_conformance_url" {
				// Special validation for k8s-conformance URL
				if !k8sConformanceURLPattern.MatchString(strVal) {
					addError(fmt.Sprintf("Invalid k8sConformanceUrl format: %s. Expected format: https://github.com/cncf/k8s-conformance/tree/master/v{version}/{product}", strVal))
				} else {
					// Also validate that the URL is accessible
					wg.Add(1)
					go func(url, fName string) {
						defer wg.Done()
						if err := validateURL(url); err != nil {
							addError(fmt.Sprintf("k8sConformanceUrl is not accessible: %s (%v)", url, err))
						}
					}(strVal, field)
				}
			} else if strings.HasPrefix(strVal, "http") {
				wg.Add(1)
				go func(url, fName string) {
					defer wg.Done()
					if err := validateURL(url); err != nil {
						addError(fmt.Sprintf("Invalid URL in metadata %s: %s (%v)", fName, url, err))
					}
				}(strVal, field)
			}
		}
	}

	// Validate CNCF Membership
	if product.Metadata != nil {
		vendorName := ""
		if v, ok := product.Metadata["vendorName"]; ok {
			vendorName, _ = v.(string)
		} else if v, ok := product.Metadata["vendor_name"]; ok {
			vendorName, _ = v.(string)
		}
		vendorName = strings.TrimSpace(vendorName)
		if vendorName != "" && !cncfMembers[vendorName] {
			addError(fmt.Sprintf("vendorName '%s' does not match any CNCF member in the CNCF Landscape. The vendorName must exactly match the organization name listed at https://landscape.cncf.io/members", vendorName))
		}
	}

	// Validate Spec
	if product.Spec == nil {
		addError("Missing 'spec' section")
	} else {
		for category, schemaReqs := range schema.Spec {
			prodReqs, ok := product.Spec[category]
			if !ok {
				addError(fmt.Sprintf("Missing spec category: %s", category))
				continue
			}

			prodReqMap := make(map[string]Requirement)
			for _, r := range prodReqs {
				prodReqMap[r.ID] = r
			}

			for _, sReq := range schemaReqs {
				pReq, exists := prodReqMap[sReq.ID]
				if !exists {
					addError(fmt.Sprintf("Missing requirement '%s' in category '%s'", sReq.ID, category))
					continue
				}

				// Special case for SHOULD level with empty status and evidence
				if sReq.Level == "SHOULD" && pReq.Status == "" && len(pReq.Evidence) == 0 {
					continue
				}

				// Check Status
				if !validStatuses[pReq.Status] {
					addError(fmt.Sprintf("Invalid status '%s' for '%s'. Must be one of %v", pReq.Status, sReq.ID, keys(validStatuses)))
				}

				// Check MUST level
				if sReq.Level == "MUST" {
					if pReq.Status != "Implemented" {
						addError(fmt.Sprintf("Requirement '%s' is MUST level but status is '%s'. It must be 'Implemented'.", sReq.ID, pReq.Status))
					}
				}

				// Check N/A notes
				if pReq.Status == "N/A" && pReq.Notes == "" {
					addError(fmt.Sprintf("Notes required for '%s' when status is N/A", sReq.ID))
				}

				// Validate Evidence Links
				productDir := filepath.Dir(path)
				for _, link := range pReq.Evidence {
					if link == "" {
						continue
					}
					wg.Add(1)
					go func(url, reqID string) {
						defer wg.Done()
						if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
							if err := validateURL(url); err != nil {
								addError(fmt.Sprintf("Invalid evidence URL for '%s': %s (%v)", reqID, url, err))
							}
						} else {
							// Check local file
							fullPath := filepath.Join(productDir, url)
							if _, err := os.Stat(fullPath); os.IsNotExist(err) {
								addError(fmt.Sprintf("Invalid evidence file for '%s': %s (not found at %s)", reqID, url, fullPath))
							}
						}
					}(link, sReq.ID)
				}
			}
		}
	}

	wg.Wait()

	if len(errors) > 0 {
		fmt.Println("Validation failed:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return false
	}

	fmt.Println("Validation successful!")
	return true
}

func loadYaml(path string) (*ConformanceDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc ConformanceDoc
	err = yaml.Unmarshal(data, &doc)
	return &doc, err
}

func validateURL(urlStr string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// Try HEAD first
	resp, err := doRequest(client, http.MethodHead, urlStr)
	if err == nil && resp.StatusCode < 400 {
		resp.Body.Close()
		return nil
	}

	// If HEAD fails or returns error, try GET
	resp, err = doRequest(client, http.MethodGet, urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

// doRequest creates an HTTP request with browser-like headers to avoid being
// blocked by WAFs/bot-protection that reject Go's default User-Agent.
func doRequest(client *http.Client, method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	return client.Do(req)
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}
