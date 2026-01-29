package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Requirement represents a single conformance requirement
type Requirement struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Level       string   `yaml:"level"`
	Status      string   `yaml:"status"`
	Evidence    []string `yaml:"evidence"`
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
		// Default: process all YAML files in docs/
		processDocsDirectory()
		return
	}

	for _, inputFile := range os.Args[1:] {
		if err := convertFile(inputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputFile, err)
			os.Exit(1)
		}
	}
}

func processDocsDirectory() {
	// Find the docs directory relative to the script location
	docsDir := "docs"
	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		// Try from the root of the repo
		docsDir = filepath.Join("..", "docs")
	}

	entries, err := os.ReadDir(docsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading docs directory: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			inputFile := filepath.Join(docsDir, entry.Name())
			if err := convertFile(inputFile); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputFile, err)
				os.Exit(1)
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

	// Generate Markdown
	markdown := generateMarkdown(checklist)

	// Write output file
	outputFile := strings.TrimSuffix(inputFile, ".yaml") + ".md"
	if err := os.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputFile)
	return nil
}

func generateMarkdown(checklist ConformanceChecklist) string {
	var sb strings.Builder
	version := checklist.Metadata.KubernetesVersion

	// Header
	sb.WriteString(fmt.Sprintf("# Kubernetes AI Conformance Checklist — %s\n\n", version))

	// Introduction
	sb.WriteString("> This document defines the conformance requirements for certifying a Kubernetes platform\n")
	sb.WriteString("> as capable of reliably running AI and machine learning workloads.\n\n")

	// Build categories
	categories := []Category{
		{
			Name:         "Accelerators",
			Anchor:       "accelerators",
			Icon:         "🚀",
			Requirements: checklist.Spec.Accelerators,
		},
		{
			Name:         "Networking",
			Anchor:       "networking",
			Icon:         "🌐",
			Requirements: checklist.Spec.Networking,
		},
		{
			Name:         "Scheduling & Orchestration",
			Anchor:       "scheduling--orchestration",
			Icon:         "📅",
			Requirements: checklist.Spec.SchedulingOrchestration,
		},
		{
			Name:         "Observability",
			Anchor:       "observability",
			Icon:         "📊",
			Requirements: checklist.Spec.Observability,
		},
		{
			Name:         "Security",
			Anchor:       "security",
			Icon:         "🔒",
			Requirements: checklist.Spec.Security,
		},
		{
			Name:         "Operator Support",
			Anchor:       "operator-support",
			Icon:         "⚙️",
			Requirements: checklist.Spec.Operator,
		},
	}

	// Count requirements
	mustCount, shouldCount := countRequirements(categories)

	// Summary box
	sb.WriteString("---\n\n")
	sb.WriteString("## Overview\n\n")
	sb.WriteString(fmt.Sprintf("| Kubernetes Version | Total Requirements | Mandatory (MUST) | Recommended (SHOULD) |\n"))
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
			sb.WriteString(formatRequirement(req, i+1))
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

func formatRequirement(req Requirement, num int) string {
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
