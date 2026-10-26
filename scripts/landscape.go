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

// nonProductName matches landscape cards that represent companies or
// certifications (KCSP, KCNTP, KTP, membership) rather than products.
var nonProductName = regexp.MustCompile(`\((KCSP|KCNTP|KTP|member)\)`)

// findEntryInLandscape searches a landscape.yml byte slice for a product entry
// whose homepage_url matches the given URL (after normalization).
func findEntryInLandscape(data []byte, targetURL string) (*LandscapeEntry, error) {
	entries, err := collectEntries(data)
	if err != nil {
		return nil, err
	}
	normalizedTarget := normalizeURL(targetURL)
	for _, e := range entries {
		if nonProductName.MatchString(e.Name) {
			continue
		}
		if normalizeURL(e.HomepageURL) == normalizedTarget {
			return e, nil
		}
	}
	return nil, nil
}

// findEntryByName searches a landscape.yml byte slice for a product entry
// whose name matches the given product name. Exact token matches win;
// otherwise the best-scoring containment/prefix match wins, where entries
// mentioning the vendor break ties (e.g. Huawei's vs Baidu's
// "Cloud Container Engine").
func findEntryByName(data []byte, vendorName, platformName string) (*LandscapeEntry, error) {
	entries, err := collectEntries(data)
	if err != nil {
		return nil, err
	}
	vendorTokens := nameTokens(vendorName)
	var best *LandscapeEntry
	bestScore := 0
	for _, e := range entries {
		if nonProductName.MatchString(e.Name) {
			continue
		}
		level, overlap := nameMatchLevel(platformName, e.Name)
		if level == nameMatchNone {
			level, overlap = vendorQualifiedMatch(vendorName, platformName, e.Name)
		}
		if level == nameMatchExact {
			return e, nil
		}
		if level != nameMatchContained {
			continue
		}
		score := overlap * 2
		if mentionsAny(e.Name, vendorTokens) {
			score++
		}
		if score > bestScore {
			best = e
			bestScore = score
		}
	}
	return best, nil
}

// vendorQualifiedMatch handles entries that prepend the vendor to the product
// name, e.g. entry "Baidu Cloud Container Engine" for vendor "Baidu Cloud" +
// product "CCE（Cloud Container Engine）". The entry tokens (>= 3) must appear
// as an ordered subsequence of vendor+product tokens.
func vendorQualifiedMatch(vendorName, platformName, entryName string) (int, int) {
	et := nameTokens(entryName)
	if len(et) < 3 {
		return nameMatchNone, 0
	}
	vpt := nameTokens(vendorName + " " + platformName)
	i := 0
	for _, tok := range vpt {
		if tok == et[i] {
			i++
			if i == len(et) {
				return nameMatchContained, len(et)
			}
		}
	}
	return nameMatchNone, 0
}

func mentionsAny(name string, tokens []string) bool {
	nt := nameTokens(name)
	for _, t := range tokens {
		for _, n := range nt {
			if t == n {
				return true
			}
		}
	}
	return false
}

// collectEntries parses landscape YAML and returns all item entries in
// document order.
func collectEntries(data []byte) ([]*LandscapeEntry, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing landscape YAML: %w", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("unexpected YAML structure: expected document node")
	}
	var entries []*LandscapeEntry
	walkEntries(root.Content[0], &entries)
	return entries, nil
}

// walkEntries recursively walks the YAML node tree collecting item mappings.
func walkEntries(node *yaml.Node, entries *[]*LandscapeEntry) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		if entry := entryFromMapping(node); entry != nil {
			*entries = append(*entries, entry)
		}
		for i := 1; i < len(node.Content); i += 2 {
			walkEntries(node.Content[i], entries)
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			walkEntries(child, entries)
		}
	}
}

// entryFromMapping converts a YAML mapping node into a LandscapeEntry if it
// represents a landscape item. Returns nil otherwise.
func entryFromMapping(node *yaml.Node) *LandscapeEntry {
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

const (
	nameMatchNone = iota
	nameMatchContained
	nameMatchExact
)

var nameTokenSplit = regexp.MustCompile(`[^a-z0-9]+`)

func nameTokens(s string) []string {
	var tokens []string
	for _, t := range nameTokenSplit.Split(strings.ToLower(s), -1) {
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

// nameMatchLevel compares a product name against a landscape entry name using
// token-based heuristics, returning the match level and overlap size in
// tokens. Beyond exact equality, a match is declared when:
//   - the product tokens appear contiguously in the entry name, e.g.
//     "BKS" in "Breqwatr BKS" (single-token products need >= 3 chars);
//   - the entry tokens appear contiguously in the product name, e.g.
//     "Linode Kubernetes Engine" in "Linode Kubernetes Engine (LKE)" —
//     entries need >= 2 tokens so generic cards like "Kubernetes" never match;
//   - token counts are equal and each pair is equal or a prefix of the other,
//     e.g. "OVHcloud ..." vs "OVH ...".
func nameMatchLevel(productName, entryName string) (int, int) {
	pt := nameTokens(productName)
	et := nameTokens(entryName)
	if len(pt) == 0 || len(et) == 0 {
		return nameMatchNone, 0
	}
	if tokensEqual(pt, et) {
		return nameMatchExact, len(pt)
	}
	if len(pt) < len(et) && (len(pt) >= 2 || len(pt[0]) >= 3) && tokensContain(et, pt) {
		return nameMatchContained, len(pt)
	}
	if len(et) < len(pt) && len(et) >= 2 && tokensContain(pt, et) {
		return nameMatchContained, len(et)
	}
	if len(pt) == len(et) && len(pt) >= 2 && tokensPrefixEqual(pt, et) {
		return nameMatchContained, len(pt)
	}
	return nameMatchNone, 0
}

func tokensEqual(a, b []string) bool {
	return len(a) == len(b) && equalAll(a, b)
}

func equalAll(a, b []string) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// tokensContain reports whether short appears as a contiguous subsequence of
// long. Single-token needles must be at least 3 characters to avoid spurious
// matches.
func tokensContain(long, short []string) bool {
	if len(short) == 1 && len(short[0]) < 3 {
		return false
	}
	for i := 0; i+len(short) <= len(long); i++ {
		if equalAll(long[i:i+len(short)], short) {
			return true
		}
	}
	return false
}

// tokensPrefixEqual reports whether same-length token slices match pairwise,
// where each pair is equal or one token is a prefix (>= 3 chars) of the other.
func tokensPrefixEqual(a, b []string) bool {
	for i := range a {
		if a[i] == b[i] {
			continue
		}
		x, y := a[i], b[i]
		if len(x) > len(y) {
			x, y = y, x
		}
		if len(x) < 3 || !strings.HasPrefix(y, x) {
			return false
		}
	}
	return true
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

// maxLogoSize is the maximum allowed logo download size (10 MB).
const maxLogoSize = 10 << 20

// downloadLogo fetches a logo from a URL and writes it to destPath.
// Only http and https schemes are allowed. Downloads are capped at maxLogoSize.
// Returns an error on invalid schemes, HTTP 4xx/5xx responses, or oversized files.
func downloadLogo(logoURL, destPath string) error {
	parsed, err := url.Parse(logoURL)
	if err != nil {
		return fmt.Errorf("invalid logo URL %q: %w", logoURL, err)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return fmt.Errorf("logo URL must use http or https scheme, got %q", parsed.Scheme)
	}

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

	limited := io.LimitReader(resp.Body, maxLogoSize+1)
	n, err := io.Copy(f, limited)
	if err != nil {
		f.Close()
		os.Remove(destPath)
		return fmt.Errorf("writing logo to %s: %w", destPath, err)
	}
	if n > maxLogoSize {
		f.Close()
		os.Remove(destPath)
		return fmt.Errorf("logo from %s exceeds maximum size of %d bytes", logoURL, maxLogoSize)
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

// productResult records what was done for one product in a run.
type productResult struct {
	meta   *ProductMeta
	action string
}

// processProduct applies the landscape edit for a single product against the
// landscape.yml at landscapePath. It returns a human-readable action string,
// or "" if no change was needed.
func processProduct(tmpDir, landscapePath string, meta *ProductMeta) (string, error) {
	landscapeData, err := os.ReadFile(landscapePath)
	if err != nil {
		return "", fmt.Errorf("reading landscape.yml: %w", err)
	}

	entry, err := findEntryInLandscape(landscapeData, meta.WebsiteURL)
	if err != nil {
		return "", fmt.Errorf("searching landscape: %w", err)
	}
	if entry == nil {
		entry, err = findEntryByName(landscapeData, meta.VendorName, meta.PlatformName)
		if err != nil {
			return "", fmt.Errorf("searching landscape by name: %w", err)
		}
		if entry != nil {
			log.Printf("Matched %q to existing entry %q by name", meta.PlatformName, entry.Name)
		}
	}

	if entry != nil {
		if entry.HasAIPlatformSecondPath {
			log.Printf("Entry %q already has Certified Kubernetes - AI Platform second_path. Nothing to do.", entry.Name)
			return "", nil
		}
		log.Printf("Found existing entry %q, adding AI Platform second_path...", entry.Name)
		modified := insertSecondPath(landscapeData, entry)
		if err := os.WriteFile(landscapePath, modified, 0644); err != nil {
			return "", fmt.Errorf("writing modified landscape.yml: %w", err)
		}
		return "added second_path to existing entry", nil
	}

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
		return "", fmt.Errorf("inserting new entry: %w", err)
	}
	if err := os.WriteFile(landscapePath, modified, 0644); err != nil {
		return "", fmt.Errorf("writing modified landscape.yml: %w", err)
	}
	return "added new entry", nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run -tags landscape scripts/landscape.go <PRODUCT.yaml path> [<PRODUCT.yaml path>...] [--pr-url <url>] [--branch <name>]")
	}

	var productPaths []string
	var prURL, branchOverride string
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--pr-url":
			if i+1 < len(os.Args) {
				prURL = os.Args[i+1]
				i++
			}
		case "--branch":
			if i+1 < len(os.Args) {
				branchOverride = os.Args[i+1]
				i++
			}
		default:
			productPaths = append(productPaths, os.Args[i])
		}
	}
	if len(productPaths) == 0 {
		log.Fatal("At least one PRODUCT.yaml path is required")
	}
	batch := len(productPaths) > 1

	// 1. Read and parse all PRODUCT.yaml files
	var metas []*ProductMeta
	for _, productPath := range productPaths {
		data, err := os.ReadFile(productPath)
		if err != nil {
			if batch {
				log.Printf("WARNING: skipping %s: %v", productPath, err)
				continue
			}
			log.Fatalf("Reading PRODUCT.yaml: %v", err)
		}
		meta, err := parseProductYAML(data)
		if err != nil {
			if batch {
				log.Printf("WARNING: skipping %s: %v", productPath, err)
				continue
			}
			log.Fatalf("Parsing PRODUCT.yaml: %v", err)
		}
		if meta.WebsiteURL == "" {
			if batch {
				log.Printf("WARNING: skipping %s: websiteUrl is required", productPath)
				continue
			}
			log.Fatal("PRODUCT.yaml: websiteUrl is required for landscape integration")
		}
		log.Printf("Parsed product: %s by %s (k8s %s)", meta.PlatformName, meta.VendorName, meta.KubernetesVersion)
		metas = append(metas, meta)
	}
	if len(metas) == 0 {
		log.Fatal("No valid PRODUCT.yaml files to process")
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

	// 3. Apply changes per product
	landscapePath := filepath.Join(tmpDir, "landscape.yml")
	var results []productResult
	for _, meta := range metas {
		action, err := processProduct(tmpDir, landscapePath, meta)
		if err != nil {
			if batch {
				log.Printf("WARNING: skipping %s: %v", meta.PlatformName, err)
				continue
			}
			log.Fatalf("Processing %s: %v", meta.PlatformName, err)
		}
		if action == "" {
			continue
		}
		results = append(results, productResult{meta: meta, action: action})
	}
	if len(results) == 0 {
		log.Println("No landscape changes needed. Nothing to do.")
		return
	}

	// 4. Create branch, commit, push
	var branchName, commitMsg, prTitle string
	if batch {
		branchName = branchOverride
		if branchName == "" {
			branchName = "ai-conformance/catch-up-" + time.Now().UTC().Format("2006-01-02")
		}
		commitMsg = "Add certified AI Platform entries to Certified Kubernetes - AI Platform"
		prTitle = fmt.Sprintf("Add %d certified platforms to Certified Kubernetes - AI Platform", len(results))
	} else {
		branchName = branchOverride
		if branchName == "" {
			branchName = "ai-conformance/" + sanitizeBranchName(results[0].meta.PlatformName)
		}
		commitMsg = fmt.Sprintf("Add %s to Certified Kubernetes - AI Platform", results[0].meta.PlatformName)
		prTitle = commitMsg
	}
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
	if err := runCmdInDir(tmpDir, "git", "commit", "--signoff", "-m", commitMsg); err != nil {
		log.Fatalf("Committing changes: %v", err)
	}
	// Use --force in case the branch exists from a previous failed run
	if err := runCmdInDir(tmpDir, "git", "push", "--force", "-u", "origin", branchName); err != nil {
		log.Fatalf("Pushing branch: %v", err)
	}

	// 5. Open PR with gh CLI
	submissionLine := ""
	if prURL != "" {
		submissionLine = fmt.Sprintf("**Conformance Submission:** %s\n", prURL)
	}

	var prBody string
	if batch {
		var sb strings.Builder
		sb.WriteString("## AI Conformance Certification\n\n")
		sb.WriteString("This PR reconciles the landscape with certified AI conformant platforms:\n\n")
		sb.WriteString("| Product | Vendor | Kubernetes | Change |\n|---|---|---|---|\n")
		for _, r := range results {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				r.meta.PlatformName, r.meta.VendorName, r.meta.KubernetesVersion, r.action))
		}
		sb.WriteString("\n" + submissionLine)
		sb.WriteString("\nAutomated by [k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance).")
		prBody = sb.String()
	} else {
		prBody = fmt.Sprintf(`## AI Conformance Certification

**Product:** %s
**Vendor:** %s
**Kubernetes Version:** %s
%s
This PR %s.

Automated by [k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance).`,
			results[0].meta.PlatformName,
			results[0].meta.VendorName,
			results[0].meta.KubernetesVersion,
			submissionLine,
			results[0].action,
		)
	}

	log.Println("Creating PR on cncf/landscape...")
	if err := runCmdInDir(tmpDir, "gh", "pr", "create",
		"--repo", "cncf/landscape",
		"--base", "master",
		"--head", branchName,
		"--title", prTitle,
		"--body", prBody,
		"--reviewer", "taylorwaggoner",
	); err != nil {
		log.Fatalf("Creating PR: %v", err)
	}

	log.Println("Done! PR created successfully.")
}
