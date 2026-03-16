//go:build landscape

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ProductMeta holds metadata extracted from a PRODUCT.yaml file.
type ProductMeta struct {
	PlatformName      string
	PlatformVersion   string
	VendorName        string
	WebsiteURL        string
	ProductLogoURL    string
	Description       string
	KubernetesVersion string
}

// productFile is the top-level structure for unmarshalling PRODUCT.yaml.
type productFile struct {
	Metadata map[string]interface{} `yaml:"metadata"`
}

// parseProductYAML parses a PRODUCT.yaml byte slice and extracts ProductMeta.
// It supports both camelCase and snake_case field names.
// Returns an error if platformName is empty or missing.
func parseProductYAML(data []byte) (*ProductMeta, error) {
	var pf productFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parsing PRODUCT.yaml: %w", err)
	}
	if pf.Metadata == nil {
		return nil, fmt.Errorf("PRODUCT.yaml missing metadata section")
	}

	get := func(camel, snake string) string {
		if v, ok := pf.Metadata[camel]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		if v, ok := pf.Metadata[snake]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}

	meta := &ProductMeta{
		PlatformName:      get("platformName", "platform_name"),
		PlatformVersion:   get("platformVersion", "platform_version"),
		VendorName:        get("vendorName", "vendor_name"),
		WebsiteURL:        get("websiteUrl", "website_url"),
		ProductLogoURL:    get("productLogoUrl", "product_logo_url"),
		Description:       get("description", "description"),
		KubernetesVersion: get("kubernetesVersion", "kubernetes_version"),
	}

	if meta.PlatformName == "" {
		return nil, fmt.Errorf("PRODUCT.yaml: platformName is required and must not be empty")
	}

	return meta, nil
}

// LandscapeEntry represents a found entry in landscape.yml.
type LandscapeEntry struct {
	Name                    string
	HomepageURL             string
	HasAIPlatformSecondPath bool
	ItemLineIndex           int // 0-indexed line of the "- item:" line
	LastFieldLineIndex      int // 0-indexed line of the last field in the entry
}

// normalizeURL normalizes a URL for matching: lowercase, strip trailing /, strip www. from host.
func normalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(strings.ToLower(raw))
	if err != nil {
		return strings.TrimRight(strings.ToLower(raw), "/")
	}

	parsed.Host = strings.TrimPrefix(parsed.Host, "www.")
	parsed.Path = strings.TrimRight(parsed.Path, "/")

	return parsed.String()
}

// findEntryInLandscape searches a landscape.yml byte slice for an entry
// whose homepage_url matches the given URL (after normalization).
func findEntryInLandscape(data []byte, targetURL string) (*LandscapeEntry, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing landscape YAML: %w", err)
	}

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("unexpected YAML structure: expected document node")
	}

	normalizedTarget := normalizeURL(targetURL)
	return walkForEntry(root.Content[0], normalizedTarget)
}

// walkForEntry recursively walks the YAML node tree looking for a mapping node
// that has a homepage_url matching the target.
func walkForEntry(node *yaml.Node, normalizedTarget string) (*LandscapeEntry, error) {
	if node == nil {
		return nil, nil
	}

	switch node.Kind {
	case yaml.MappingNode:
		// Check if this mapping has homepage_url that matches
		entry := checkMappingForEntry(node, normalizedTarget)
		if entry != nil {
			return entry, nil
		}
		// Recurse into mapping values
		for i := 1; i < len(node.Content); i += 2 {
			result, err := walkForEntry(node.Content[i], normalizedTarget)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			result, err := walkForEntry(child, normalizedTarget)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}
	}

	return nil, nil
}

// checkMappingForEntry checks whether a YAML mapping node represents a landscape
// item with a homepage_url matching the target. Returns nil if not a match.
func checkMappingForEntry(node *yaml.Node, normalizedTarget string) *LandscapeEntry {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	var name, homepageURL string
	var hasItem bool
	var hasAIPlatform bool
	var secondPathNode *yaml.Node
	maxLine := 0 // track the last line in this mapping (1-indexed from yaml.Node)
	itemLine := 0

	for i := 0; i < len(node.Content)-1; i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]

		switch key.Value {
		case "item":
			hasItem = true
			// The "- item:" line is the item key line; but the actual sequence entry
			// starts at the key's line. We need to subtract 1 since the "- " prefix
			// is on the same line as the key.
			itemLine = key.Line
		case "name":
			name = val.Value
		case "homepage_url":
			homepageURL = val.Value
		case "second_path":
			secondPathNode = val
		}

		// Track the maximum line number for this mapping
		lastLine := lastNodeLine(val)
		if lastLine > maxLine {
			maxLine = lastLine
		}
		if key.Line > maxLine {
			maxLine = key.Line
		}
	}

	if !hasItem || homepageURL == "" {
		return nil
	}

	if normalizeURL(homepageURL) != normalizedTarget {
		return nil
	}

	// Check if second_path already contains AI Platform
	if secondPathNode != nil && secondPathNode.Kind == yaml.SequenceNode {
		for _, item := range secondPathNode.Content {
			if strings.Contains(item.Value, "Certified Kubernetes - AI Platform") {
				hasAIPlatform = true
				break
			}
		}
	}

	return &LandscapeEntry{
		Name:                    name,
		HomepageURL:             homepageURL,
		HasAIPlatformSecondPath: hasAIPlatform,
		ItemLineIndex:           itemLine - 1, // convert 1-indexed to 0-indexed
		LastFieldLineIndex:      maxLine - 1,  // convert 1-indexed to 0-indexed
	}
}

// lastNodeLine returns the last line number (1-indexed) used by a yaml.Node,
// accounting for sequences and mappings.
func lastNodeLine(node *yaml.Node) int {
	if node == nil {
		return 0
	}
	max := node.Line
	for _, child := range node.Content {
		cl := lastNodeLine(child)
		if cl > max {
			max = cl
		}
	}
	return max
}

// insertSecondPath inserts the AI Platform second_path into an existing landscape entry.
// If the entry already has a second_path block, it appends the new item.
// If not, it inserts both the second_path key and the list item.
func insertSecondPath(data []byte, entry *LandscapeEntry) []byte {
	lines := strings.Split(string(data), "\n")
	insertAfter := entry.LastFieldLineIndex

	var newLines []string
	if entry.HasAIPlatformSecondPath {
		// Already has it, nothing to do
		return data
	}

	// Determine if the entry already has a second_path key by scanning entry lines
	hasSecondPath := false
	for i := entry.ItemLineIndex; i <= entry.LastFieldLineIndex && i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "second_path:" {
			hasSecondPath = true
			break
		}
	}

	if hasSecondPath {
		// Append just the list item after the last line of the entry
		newLines = make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:insertAfter+1]...)
		newLines = append(newLines, `              - "Platform / Certified Kubernetes - AI Platform"`)
		newLines = append(newLines, lines[insertAfter+1:]...)
	} else {
		// Insert both second_path key and list item
		newLines = make([]string, 0, len(lines)+2)
		newLines = append(newLines, lines[:insertAfter+1]...)
		newLines = append(newLines, `            second_path:`)
		newLines = append(newLines, `              - "Platform / Certified Kubernetes - AI Platform"`)
		newLines = append(newLines, lines[insertAfter+1:]...)
	}

	return []byte(strings.Join(newLines, "\n"))
}

// sanitizeLogoName converts a platform name to a safe logo filename.
// Lowercase, replace non-alphanumeric with -, append .svg.
var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func sanitizeLogoName(platformName string) string {
	s := strings.ToLower(platformName)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s + ".svg"
}

// insertNewEntry inserts a new landscape entry into the Certified Kubernetes - AI Platform
// subcategory. It handles both empty (items: []) and populated item lists.
func insertNewEntry(data []byte, meta *ProductMeta, logoFilename string) ([]byte, error) {
	content := string(data)

	// Build the entry block
	homepageURL := meta.WebsiteURL

	// Sanitize description: collapse to single line, escape for YAML
	desc := strings.ReplaceAll(meta.Description, "\n", " ")
	desc = strings.Join(strings.Fields(desc), " ")

	entryBlock := fmt.Sprintf("          - item:\n"+
		"            name: %s\n"+
		"            description: >-\n"+
		"              %s\n"+
		"            homepage_url: %s\n"+
		"            logo: %s", meta.PlatformName, desc, homepageURL, logoFilename)

	// Look for "Certified Kubernetes - AI Platform" subcategory
	lines := strings.Split(content, "\n")
	subcatIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "Certified Kubernetes - AI Platform") {
			// Make sure this is a subcategory/category name, not a second_path reference
			trimmed := strings.TrimSpace(line)
			isNameLine := strings.HasPrefix(trimmed, "name:") || strings.HasPrefix(trimmed, "- name:")
			if isNameLine && strings.Contains(trimmed, "Certified Kubernetes - AI Platform") {
				subcatIdx = i
				break
			}
		}
	}
	if subcatIdx == -1 {
		return nil, fmt.Errorf("subcategory 'Certified Kubernetes - AI Platform' not found in landscape data")
	}

	// Find the items line for this subcategory
	for i := subcatIdx + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "items: []" {
			// Replace empty items with our entry
			newLines := make([]string, 0, len(lines)+6)
			newLines = append(newLines, lines[:i]...)
			newLines = append(newLines, "        items:")
			newLines = append(newLines, entryBlock)
			newLines = append(newLines, lines[i+1:]...)
			return []byte(strings.Join(newLines, "\n")), nil
		}
		if trimmed == "items:" {
			// Find end of existing items and append
			// Items start at i, entries follow
			j := i + 1
			for j < len(lines) {
				lt := strings.TrimSpace(lines[j])
				if lt == "" {
					j++
					continue
				}
				// Check if we've left the items section (next subcategory or category)
				if !strings.HasPrefix(lines[j], "          ") && lt != "" {
					break
				}
				j++
			}
			// Insert before j
			newLines := make([]string, 0, len(lines)+6)
			newLines = append(newLines, lines[:j]...)
			newLines = append(newLines, entryBlock)
			newLines = append(newLines, lines[j:]...)
			return []byte(strings.Join(newLines, "\n")), nil
		}

		// If we hit the next subcategory or category before finding items, break
		if strings.HasPrefix(trimmed, "- name:") {
			break
		}
	}

	return nil, fmt.Errorf("could not find items list for 'Certified Kubernetes - AI Platform' subcategory")
}

// downloadLogo fetches a logo from a URL and writes it to destPath.
// Returns an error on HTTP 4xx/5xx responses.
func downloadLogo(logoURL, destPath string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(logoURL)
	if err != nil {
		return fmt.Errorf("downloading logo from %s: %w", logoURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("downloading logo from %s: HTTP %d", logoURL, resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating logo file %s: %w", destPath, err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(destPath)
		return fmt.Errorf("writing logo to %s: %w", destPath, err)
	}

	return f.Close()
}

// sanitizeBranchName converts a name to a safe git branch suffix.
// Lowercase, non-alphanumeric characters replaced with dashes, max 50 chars.
func sanitizeBranchName(name string) string {
	s := strings.ToLower(name)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 50 {
		s = s[:50]
		s = strings.TrimRight(s, "-")
	}
	return s
}

// runCmd runs an external command, piping stdout and stderr to os.Stdout/os.Stderr.
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCmdInDir runs an external command in the specified directory,
// piping stdout and stderr to os.Stdout/os.Stderr.
func runCmdInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// checkExistingPR checks if there's already an open PR for the given branch on cncf/landscape.
// Returns the PR URL if one exists, empty string otherwise.
func checkExistingPR(repoDir, branchName string) string {
	cmd := exec.Command("gh", "pr", "list",
		"--repo", "cncf/landscape",
		"--head", branchName,
		"--state", "open",
		"--json", "url",
		"--jq", ".[0].url // empty",
	)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "" // If gh fails, proceed anyway
	}
	return strings.TrimSpace(string(out))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run -tags landscape scripts/landscape.go <PRODUCT.yaml path> [--pr-url <url>]")
	}

	productPath := os.Args[1]
	var prURL string
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--pr-url" && i+1 < len(os.Args) {
			prURL = os.Args[i+1]
			i++
		}
	}

	// 1. Read and parse PRODUCT.yaml
	data, err := os.ReadFile(productPath)
	if err != nil {
		log.Fatalf("Reading PRODUCT.yaml: %v", err)
	}

	meta, err := parseProductYAML(data)
	if err != nil {
		log.Fatalf("Parsing PRODUCT.yaml: %v", err)
	}
	log.Printf("Parsed product: %s by %s (k8s %s)", meta.PlatformName, meta.VendorName, meta.KubernetesVersion)

	if meta.WebsiteURL == "" {
		log.Fatal("PRODUCT.yaml: websiteUrl is required for landscape integration")
	}

	// 2. Clone cncf/landscape repo (shallow)
	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		log.Fatal("GH_TOKEN environment variable is required")
	}

	tmpDir, err := os.MkdirTemp("", "landscape-*")
	if err != nil {
		log.Fatalf("Creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log.Println("Cloning cncf/landscape (shallow)...")
	if err := runCmd("git", "clone", "--depth", "1", "https://github.com/cncf/landscape.git", tmpDir); err != nil {
		log.Fatalf("Cloning landscape repo: %v", err)
	}
	// Set authenticated remote URL for push (avoids token in clone command/process list)
	authURL := fmt.Sprintf("https://x-access-token:%s@github.com/cncf/landscape.git", ghToken)
	if err := runCmdInDir(tmpDir, "git", "remote", "set-url", "origin", authURL); err != nil {
		log.Fatalf("Setting authenticated remote URL: %v", err)
	}

	// 3. Read landscape.yml
	landscapePath := filepath.Join(tmpDir, "landscape.yml")
	landscapeData, err := os.ReadFile(landscapePath)
	if err != nil {
		log.Fatalf("Reading landscape.yml: %v", err)
	}

	// 4. Search for existing entry
	entry, err := findEntryInLandscape(landscapeData, meta.WebsiteURL)
	if err != nil {
		log.Fatalf("Searching landscape: %v", err)
	}

	var prBodyAction string

	if entry != nil {
		if entry.HasAIPlatformSecondPath {
			// Already has AI Platform second_path — nothing to do
			log.Printf("Entry %q already has Certified Kubernetes - AI Platform second_path. Nothing to do.", entry.Name)
			return
		}
		// Found but no AI Platform — insert second_path
		log.Printf("Found existing entry %q, adding AI Platform second_path...", entry.Name)
		modified := insertSecondPath(landscapeData, entry)
		if err := os.WriteFile(landscapePath, modified, 0644); err != nil {
			log.Fatalf("Writing modified landscape.yml: %v", err)
		}
		prBodyAction = `This PR adds the "Certified Kubernetes - AI Platform" designation to the existing landscape entry via ` + "`second_path`."
	} else {
		// Not found — download logo and insert new entry
		log.Printf("No existing entry found for %q, creating new entry...", meta.PlatformName)

		logoFilename := sanitizeLogoName(meta.PlatformName)
		if meta.ProductLogoURL != "" {
			logoDestPath := filepath.Join(tmpDir, "hosted_logos", logoFilename)
			if err := downloadLogo(meta.ProductLogoURL, logoDestPath); err != nil {
				log.Printf("WARNING: Failed to download logo: %v (continuing without logo)", err)
			} else {
				log.Printf("Downloaded logo to %s", logoDestPath)
			}
		} else {
			log.Println("WARNING: No productLogoUrl provided, skipping logo download")
		}

		modified, err := insertNewEntry(landscapeData, meta, logoFilename)
		if err != nil {
			log.Fatalf("Inserting new entry: %v", err)
		}
		if err := os.WriteFile(landscapePath, modified, 0644); err != nil {
			log.Fatalf("Writing modified landscape.yml: %v", err)
		}
		prBodyAction = `This PR adds a new entry to the "Certified Kubernetes - AI Platform" subcategory.`
	}

	// 5. Create branch, commit, push
	branchName := "ai-conformance/" + sanitizeBranchName(meta.PlatformName)
	log.Printf("Creating branch %s...", branchName)

	// Check if a PR already exists for this branch
	existingPR := checkExistingPR(tmpDir, branchName)
	if existingPR != "" {
		log.Printf("An open PR already exists for branch %s: %s", branchName, existingPR)
		log.Println("Skipping — delete the existing PR/branch to re-run.")
		return
	}

	if err := runCmdInDir(tmpDir, "git", "checkout", "-b", branchName); err != nil {
		log.Fatalf("Creating branch: %v", err)
	}
	if err := runCmdInDir(tmpDir, "git", "add", "-A"); err != nil {
		log.Fatalf("Staging changes: %v", err)
	}
	commitMsg := fmt.Sprintf("Add %s to Certified Kubernetes - AI Platform", meta.PlatformName)
	if err := runCmdInDir(tmpDir, "git", "commit", "-m", commitMsg); err != nil {
		log.Fatalf("Committing changes: %v", err)
	}
	// Use --force in case the branch exists from a previous failed run
	if err := runCmdInDir(tmpDir, "git", "push", "--force", "-u", "origin", branchName); err != nil {
		log.Fatalf("Pushing branch: %v", err)
	}

	// 6. Open PR with gh CLI
	prTitle := fmt.Sprintf("Add %s to Certified Kubernetes - AI Platform", meta.PlatformName)

	submissionLine := ""
	if prURL != "" {
		submissionLine = fmt.Sprintf("**Conformance Submission:** %s\n", prURL)
	}

	prBody := fmt.Sprintf(`## AI Conformance Certification

**Product:** %s
**Vendor:** %s
**Kubernetes Version:** %s
%s
%s

Automated by [k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance).`,
		meta.PlatformName,
		meta.VendorName,
		meta.KubernetesVersion,
		submissionLine,
		prBodyAction,
	)

	log.Println("Creating PR on cncf/landscape...")
	if err := runCmdInDir(tmpDir, "gh", "pr", "create",
		"--title", prTitle,
		"--body", prBody,
		"--reviewer", "taylorwaggoner",
	); err != nil {
		log.Fatalf("Creating PR: %v", err)
	}

	log.Println("Done! PR created successfully.")
}
