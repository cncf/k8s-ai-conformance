package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Evidence can be either a string or a list of strings in YAML
type Evidence []string

func (e *Evidence) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as a string first
	var str string
	if err := value.Decode(&str); err == nil {
		if str != "" && str != "N/A" {
			*e = []string{str}
		} else {
			*e = []string{}
		}
		return nil
	}

	// Otherwise unmarshal as a list
	var list []string
	if err := value.Decode(&list); err != nil {
		return err
	}
	*e = list
	return nil
}

// Requirement represents a single conformance requirement
type Requirement struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Level       string   `yaml:"level"`
	Status      string   `yaml:"status"`
	Evidence    Evidence `yaml:"evidence"`
	Notes       string   `yaml:"notes"`
}

// Metadata represents the platform metadata
type Metadata struct {
	KubernetesVersion   string `yaml:"kubernetesVersion"`
	PlatformName        string `yaml:"platformName"`
	PlatformVersion     string `yaml:"platformVersion"`
	VendorName          string `yaml:"vendorName"`
	WebsiteURL          string `yaml:"websiteUrl"`
	RepoURL             string `yaml:"repoUrl"`
	DocumentationURL    string `yaml:"documentationUrl"`
	ProductLogoURL      string `yaml:"productLogoUrl"`
	Description         string `yaml:"description"`
	ContactEmailAddress string `yaml:"contactEmailAddress"`
	K8sConformanceURL   string `yaml:"k8sConformanceUrl"`
}

// Spec represents the conformance specification
type Spec struct {
	Accelerators            []Requirement `yaml:"accelerators"`
	Networking              []Requirement `yaml:"networking"`
	SchedulingOrchestration []Requirement `yaml:"schedulingOrchestration"`
	Observability           []Requirement `yaml:"observability"`
	Security                []Requirement `yaml:"security"`
	Operator                []Requirement `yaml:"operator"`
}

// ConformanceChecklist represents the full YAML structure
type ConformanceChecklist struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec"`
}

// Category holds display information for a requirement category
type Category struct {
	Name         string
	Anchor       string
	Icon         string
	Requirements []Requirement
}

func main() {
	if len(os.Args) < 2 {
		// Default: process all files
		processAllFiles()
		return
	}

	for _, inputFile := range os.Args[1:] {
		if err := convertFile(inputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputFile, err)
			os.Exit(1)
		}
	}
}

func processAllFiles() {
	// Process docs/ directory (templates)
	processDirectory("docs", "AIConformance-*.yaml")

	// Process version directories (v1.33, v1.34, etc.)
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading current directory: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "v") {
			processVersionDirectory(entry.Name())
		}
	}
}

func processDirectory(dir, pattern string) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error globbing %s: %v\n", dir, err)
		return
	}

	for _, match := range matches {
		if err := convertFile(match); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", match, err)
		}
	}
}

func processVersionDirectory(versionDir string) {
	entries, err := os.ReadDir(versionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", versionDir, err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			productFile := filepath.Join(versionDir, entry.Name(), "PRODUCT.yaml")
			if _, err := os.Stat(productFile); err == nil {
				if err := convertFile(productFile); err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", productFile, err)
				}
			}
		}
	}
}

func convertFile(inputFile string) error {
	// Read YAML file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Parse YAML
	var checklist ConformanceChecklist
	if err := yaml.Unmarshal(data, &checklist); err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	// Determine if this is a template or a product submission
	isTemplate := strings.Contains(inputFile, "docs/") || isPlaceholder(checklist.Metadata.PlatformName)

	// Generate Markdown
	var markdown string
	if isTemplate {
		markdown = generateTemplateMarkdown(checklist)
	} else {
		markdown = generateProductMarkdown(checklist)
	}

	// Write output file
	outputFile := strings.TrimSuffix(inputFile, ".yaml") + ".md"
	if err := os.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputFile)
	return nil
}

func isPlaceholder(value string) bool {
	return strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]")
}

func buildCategories(spec Spec) []Category {
	return []Category{
		{
			Name:         "Accelerators",
			Anchor:       "accelerators",
			Icon:         "🚀",
			Requirements: spec.Accelerators,
		},
		{
			Name:         "Networking",
			Anchor:       "networking",
			Icon:         "🌐",
			Requirements: spec.Networking,
		},
		{
			Name:         "Scheduling & Orchestration",
			Anchor:       "scheduling--orchestration",
			Icon:         "📅",
			Requirements: spec.SchedulingOrchestration,
		},
		{
			Name:         "Observability",
			Anchor:       "observability",
			Icon:         "📊",
			Requirements: spec.Observability,
		},
		{
			Name:         "Security",
			Anchor:       "security",
			Icon:         "🔒",
			Requirements: spec.Security,
		},
		{
			Name:         "Operator Support",
			Anchor:       "operator-support",
			Icon:         "⚙️",
			Requirements: spec.Operator,
		},
	}
}

// generateTemplateMarkdown generates markdown for template files (docs/)
func generateTemplateMarkdown(checklist ConformanceChecklist) string {
	var sb strings.Builder
	version := checklist.Metadata.KubernetesVersion

	// Header
	sb.WriteString(fmt.Sprintf("# Kubernetes AI Conformance Checklist — %s\n\n", version))

	// Introduction
	sb.WriteString("> This document defines the conformance requirements for certifying a Kubernetes platform\n")
	sb.WriteString("> as capable of reliably running AI and machine learning workloads.\n\n")

	categories := buildCategories(checklist.Spec)

	// Count requirements
	mustCount, shouldCount := countRequirements(categories)

	// Summary box
	sb.WriteString("---\n\n")
	sb.WriteString("## Overview\n\n")
	sb.WriteString("| Kubernetes Version | Total Requirements | Mandatory (MUST) | Recommended (SHOULD) |\n")
	sb.WriteString("|:------------------:|:------------------:|:----------------:|:--------------------:|\n")
	sb.WriteString(fmt.Sprintf("| **%s** | %d | %d | %d |\n\n",
		version, mustCount+shouldCount, mustCount, shouldCount))

	// Requirement levels explanation
	sb.WriteString("### Requirement Levels\n\n")
	sb.WriteString("| Level | Meaning |\n")
	sb.WriteString("|:-----:|:--------|\n")
	sb.WriteString("| **MUST** | Mandatory for conformance. Platform cannot be certified without implementing this requirement. |\n")
	sb.WriteString("| **SHOULD** | Recommended but not mandatory. Platforms are encouraged to implement for better AI/ML support. |\n\n")

	// Table of Contents
	sb.WriteString("---\n\n")
	sb.WriteString("## Table of Contents\n\n")
	for _, cat := range categories {
		if len(cat.Requirements) > 0 {
			must, should := countCategoryRequirements(cat.Requirements)
			sb.WriteString(fmt.Sprintf("- [%s %s](#%s)", cat.Icon, cat.Name, cat.Anchor))
			sb.WriteString(fmt.Sprintf(" — %s\n", formatRequirementCount(must, should)))
		}
	}
	sb.WriteString("\n---\n\n")

	// Each category
	for _, cat := range categories {
		if len(cat.Requirements) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s %s\n\n", cat.Icon, cat.Name))

		// Category summary
		must, should := countCategoryRequirements(cat.Requirements)
		sb.WriteString(fmt.Sprintf("*%s*\n\n", formatRequirementCount(must, should)))

		// Requirements
		for i, req := range cat.Requirements {
			sb.WriteString(formatTemplateRequirement(req, i+1))
		}

		sb.WriteString("---\n\n")
	}

	// Footer
	sb.WriteString("## Submission Instructions\n\n")
	sb.WriteString("To submit your platform for conformance certification:\n\n")
	sb.WriteString("1. Copy this checklist template to `PRODUCT.yaml`\n")
	sb.WriteString("2. Fill in the `metadata` section with your platform details\n")
	sb.WriteString("3. For each requirement, set the `status` field to one of:\n")
	sb.WriteString("   - `Implemented` — Requirement is fully supported\n")
	sb.WriteString("   - `Not Implemented` — Requirement is not supported\n")
	sb.WriteString("   - `Partially Implemented` — Requirement is partially supported\n")
	sb.WriteString("   - `N/A` — Requirement does not apply (must provide justification in `notes`)\n")
	sb.WriteString("4. Provide `evidence` URLs linking to documentation or test results\n")
	sb.WriteString("5. Submit a pull request to the [k8s-ai-conformance](https://github.com/cncf/k8s-ai-conformance) repository\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("*Generated from AIConformance-%s.yaml*\n", strings.TrimPrefix(version, "v")))

	return sb.String()
}

// generateProductMarkdown generates markdown for product submissions (v*/)
func generateProductMarkdown(checklist ConformanceChecklist) string {
	var sb strings.Builder
	meta := checklist.Metadata

	// Header with platform name
	sb.WriteString(fmt.Sprintf("# %s — AI Conformance Report\n\n", meta.PlatformName))

	// Platform info card
	sb.WriteString("## Platform Information\n\n")
	sb.WriteString("| | |\n")
	sb.WriteString("|:--|:--|\n")
	sb.WriteString(fmt.Sprintf("| **Vendor** | %s |\n", meta.VendorName))
	sb.WriteString(fmt.Sprintf("| **Platform** | %s |\n", meta.PlatformName))
	sb.WriteString(fmt.Sprintf("| **Platform Version** | %s |\n", meta.PlatformVersion))
	sb.WriteString(fmt.Sprintf("| **Kubernetes Version** | %s |\n", meta.KubernetesVersion))
	if meta.WebsiteURL != "" && !isPlaceholder(meta.WebsiteURL) {
		sb.WriteString(fmt.Sprintf("| **Website** | [%s](%s) |\n", meta.WebsiteURL, meta.WebsiteURL))
	}
	if meta.DocumentationURL != "" && !isPlaceholder(meta.DocumentationURL) {
		sb.WriteString(fmt.Sprintf("| **Documentation** | [Link](%s) |\n", meta.DocumentationURL))
	}
	sb.WriteString("\n")

	// Description
	if meta.Description != "" && !isPlaceholder(meta.Description) {
		sb.WriteString(fmt.Sprintf("> %s\n\n", meta.Description))
	}

	categories := buildCategories(checklist.Spec)

	// Compliance summary
	sb.WriteString("---\n\n")
	sb.WriteString("## Compliance Summary\n\n")

	implemented, notImplemented, partial, na := countStatuses(categories)
	total := implemented + notImplemented + partial + na

	sb.WriteString("| Status | Count |\n")
	sb.WriteString("|:-------|:-----:|\n")
	sb.WriteString(fmt.Sprintf("| ✅ Implemented | %d |\n", implemented))
	if partial > 0 {
		sb.WriteString(fmt.Sprintf("| 🟡 Partially Implemented | %d |\n", partial))
	}
	if notImplemented > 0 {
		sb.WriteString(fmt.Sprintf("| ❌ Not Implemented | %d |\n", notImplemented))
	}
	if na > 0 {
		sb.WriteString(fmt.Sprintf("| ⚪ N/A | %d |\n", na))
	}
	sb.WriteString(fmt.Sprintf("| **Total** | **%d** |\n\n", total))

	// Quick status table
	sb.WriteString("### Requirements at a Glance\n\n")
	sb.WriteString("| Category | Requirement | Level | Status |\n")
	sb.WriteString("|:---------|:------------|:-----:|:------:|\n")
	for _, cat := range categories {
		for _, req := range cat.Requirements {
			statusIcon := getStatusIcon(req.Status)
			levelBadge := "MUST"
			if req.Level == "SHOULD" {
				levelBadge = "SHOULD"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				cat.Name, formatID(req.ID), levelBadge, statusIcon))
		}
	}
	sb.WriteString("\n")

	// Detailed requirements
	sb.WriteString("---\n\n")
	sb.WriteString("## Detailed Requirements\n\n")

	for _, cat := range categories {
		if len(cat.Requirements) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("### %s %s\n\n", cat.Icon, cat.Name))

		for _, req := range cat.Requirements {
			sb.WriteString(formatProductRequirement(req))
		}
	}

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString("*Generated from PRODUCT.yaml*\n")

	return sb.String()
}

func formatTemplateRequirement(req Requirement, num int) string {
	var sb strings.Builder

	// Requirement header with level badge
	levelBadge := "🔴 **MUST**"
	if req.Level == "SHOULD" {
		levelBadge = "🟡 **SHOULD**"
	}

	// Format ID for display
	displayID := formatID(req.ID)

	sb.WriteString(fmt.Sprintf("### %d. %s\n\n", num, displayID))
	sb.WriteString(fmt.Sprintf("**Level:** %s\n\n", levelBadge))

	// Description with proper wrapping
	sb.WriteString("**Description:**\n\n")
	sb.WriteString(fmt.Sprintf("> %s\n\n", req.Description))

	// Compliance fields (for template reference)
	sb.WriteString("<details>\n")
	sb.WriteString("<summary><strong>Compliance Fields</strong></summary>\n\n")
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|:------|:------|\n")
	sb.WriteString(fmt.Sprintf("| `id` | `%s` |\n", req.ID))
	sb.WriteString("| `status` | `Implemented` \\| `Not Implemented` \\| `Partially Implemented` \\| `N/A` |\n")
	sb.WriteString("| `evidence` | List of URLs to documentation/test results |\n")
	sb.WriteString("| `notes` | Additional context (required if status is `N/A`) |\n\n")
	sb.WriteString("</details>\n\n")

	return sb.String()
}

func formatProductRequirement(req Requirement) string {
	var sb strings.Builder

	// Status icon and level
	statusIcon := getStatusIcon(req.Status)
	levelBadge := "🔴 MUST"
	if req.Level == "SHOULD" {
		levelBadge = "🟡 SHOULD"
	}

	displayID := formatID(req.ID)

	sb.WriteString(fmt.Sprintf("#### %s %s\n\n", statusIcon, displayID))
	sb.WriteString(fmt.Sprintf("**Level:** %s | **Status:** %s\n\n", levelBadge, req.Status))

	// Description
	sb.WriteString(fmt.Sprintf("> %s\n\n", req.Description))

	// Evidence
	if len(req.Evidence) > 0 && hasValidEvidence(req.Evidence) {
		sb.WriteString("**Evidence:**\n\n")
		for _, e := range req.Evidence {
			if e != "" {
				sb.WriteString(fmt.Sprintf("- [%s](%s)\n", truncateURL(e), e))
			}
		}
		sb.WriteString("\n")
	}

	// Notes
	if req.Notes != "" {
		sb.WriteString("**Notes:**\n\n")
		sb.WriteString(fmt.Sprintf("> %s\n\n", req.Notes))
	}

	return sb.String()
}

func getStatusIcon(status string) string {
	switch strings.ToLower(status) {
	case "implemented":
		return "✅"
	case "not implemented":
		return "❌"
	case "partially implemented":
		return "🟡"
	case "n/a":
		return "⚪"
	default:
		return "⬜"
	}
}

func hasValidEvidence(evidence []string) bool {
	for _, e := range evidence {
		if e != "" {
			return true
		}
	}
	return false
}

func truncateURL(url string) string {
	// Remove protocol for display
	display := strings.TrimPrefix(url, "https://")
	display = strings.TrimPrefix(display, "http://")

	// Truncate if too long
	if len(display) > 60 {
		return display[:57] + "..."
	}
	return display
}

func formatID(id string) string {
	// Convert snake_case to Title Case
	words := strings.Split(id, "_")
	for i, word := range words {
		if len(word) > 0 {
			// Handle common acronyms
			upper := strings.ToUpper(word)
			if upper == "DRA" || upper == "API" || upper == "GPU" || upper == "AI" || upper == "URL" || upper == "CRD" {
				words[i] = upper
			} else {
				words[i] = strings.ToUpper(word[:1]) + word[1:]
			}
		}
	}
	return strings.Join(words, " ")
}

func countRequirements(categories []Category) (must, should int) {
	for _, cat := range categories {
		for _, req := range cat.Requirements {
			if req.Level == "MUST" {
				must++
			} else {
				should++
			}
		}
	}
	return
}

func countCategoryRequirements(reqs []Requirement) (must, should int) {
	for _, req := range reqs {
		if req.Level == "MUST" {
			must++
		} else {
			should++
		}
	}
	return
}

func countStatuses(categories []Category) (implemented, notImplemented, partial, na int) {
	for _, cat := range categories {
		for _, req := range cat.Requirements {
			switch strings.ToLower(req.Status) {
			case "implemented":
				implemented++
			case "not implemented":
				notImplemented++
			case "partially implemented":
				partial++
			case "n/a":
				na++
			}
		}
	}
	return
}

func formatRequirementCount(must, should int) string {
	parts := []string{}
	if must > 0 {
		parts = append(parts, fmt.Sprintf("%d MUST", must))
	}
	if should > 0 {
		parts = append(parts, fmt.Sprintf("%d SHOULD", should))
	}
	if len(parts) == 0 {
		return "No requirements"
	}
	return strings.Join(parts, ", ")
}
