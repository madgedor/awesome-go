//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type starsReport struct {
	GeneratedAt string     `json:"generated_at"`
	Threshold   int        `json:"threshold"`
	LowRepos    []repoInfo `json:"low_repos"`
}

func generateStarsReport(links []string, threshold int) starsReport {
	low := findLowStarRepos(links, threshold)
	return starsReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Threshold:   threshold,
		LowRepos:    low,
	}
}

func saveStarsReport(r starsReport, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func mainStarsReport() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: check_github_stars_report <README.md> <output.json>")
		os.Exit(1)
	}
	links, err := extractGithubLinks(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	report := generateStarsReport(links, 10)
	if err := saveStarsReport(report, os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Report saved to %s (%d low-star repos)\n", os.Args[2], len(report.LowRepos))
}
