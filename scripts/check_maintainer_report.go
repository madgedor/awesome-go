//go:build ignore
// +build ignore

// check_maintainer_report.go generates a markdown report of unmaintained repositories
// found in the awesome-go list based on last push date.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// generateMaintainerReport builds a markdown-formatted report of unmaintained repos.
func generateMaintainerReport(unmaintained []string, thresholdDays int) string {
	var sb strings.Builder

	sb.WriteString("# Unmaintained Repositories Report\n\n")
	sb.WriteString(fmt.Sprintf("_Generated on %s_\n\n", time.Now().Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf(
		"Repositories with no push activity in the last **%d days** are listed below.\n\n",
		thresholdDays,
	))

	if len(unmaintained) == 0 {
		sb.WriteString("✅ No unmaintained repositories found.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Found **%d** potentially unmaintained repositor", len(unmaintained)))
	if len(unmaintained) == 1 {
		sb.WriteString("y")
	} else {
		sb.WriteString("ies")
	}
	sb.WriteString(":\n\n")
	sb.WriteString("| Repository |\n")
	sb.WriteString("|------------|\n")

	for _, repo := range unmaintained {
		sb.WriteString(fmt.Sprintf("| %s |\n", repo))
	}

	sb.WriteString("\n> Consider replacing or removing these entries from awesome-go.\n")
	return sb.String()
}

// saveMaintainerReport writes the report content to the given file path.
func saveMaintainerReport(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func mainMaintainerReport() {
	const (
		readmePath    = "README.md"
		reportPath    = "reports/unmaintained.md"
		thresholdDays = 365
		token         = ""
	)

	links, err := extractMaintainerLinks(readmePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting links: %v\n", err)
		os.Exit(1)
	}

	unmaintained := findUnmaintained(links, thresholdDays, token)

	if err := os.MkdirAll("reports", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating reports directory: %v\n", err)
		os.Exit(1)
	}

	report := generateMaintainerReport(unmaintained, thresholdDays)
	if err := saveMaintainerReport(reportPath, report); err != nil {
		fmt.Fprintf(os.Stderr, "error saving report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report saved to %s (%d unmaintained repos)\n", reportPath, len(unmaintained))
}
