//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type URLFormatReport struct {
	GeneratedAt string     `json:"generated_at"`
	File        string     `json:"file"`
	Total       int        `json:"total_urls"`
	Malformed   []URLEntry `json:"malformed"`
}

func generateURLReport(path string) (*URLFormatReport, error) {
	entries, err := extractURLEntries(path)
	if err != nil {
		return nil, err
	}
	bad := findMalformedURLs(entries)
	return &URLFormatReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		File:        path,
		Total:       len(entries),
		Malformed:   bad,
	}, nil
}

func saveURLReport(r *URLFormatReport, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func mainReport() {
	// Default input file; override by passing a path as the first argument.
	path := "README.md"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Default output file; override by passing a path as the second argument.
	output := "url_format_report.json"
	if len(os.Args) > 2 {
		output = os.Args[2]
	}

	report, err := generateURLReport(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := saveURLReport(report, output); err != nil {
		fmt.Fprintf(os.Stderr, "error saving report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Report saved to %s. Total: %d, Malformed: %d\n", output, report.Total, len(report.Malformed))
}
