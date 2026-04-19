//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var linkLineRe = regexp.MustCompile(`^\s*-\s+\[([^\]]+)\]\(([^)]+)\)\s*-?s*(.*)`)

type LinkEntry struct {
	Line    int
	Name    string
	URL     string
	Desc    string
}

func extractLinkEntries(path string) ([]LinkEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []LinkEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		m := linkLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		entries = append(entries, LinkEntry{
			Line: lineNum,
			Name: m[1],
			URL:  m[2],
			Desc: strings.TrimSpace(m[3]),
		})
	}
	return entries, scanner.Err()
}

func findMissingDescriptions(entries []LinkEntry) []LinkEntry {
	var missing []LinkEntry
	for _, e := range entries {
		if e.Desc == "" {
			missing = append(missing, e)
		}
	}
	return missing
}

func main() {
	path := "README.md"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	entries, err := extractLinkEntries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	missing := findMissingDescriptions(entries)
	if len(missing) == 0 {
		fmt.Println("All entries have descriptions.")
		return
	}
	for _, e := range missing {
		fmt.Printf("line %d: [%s](%s) — missing description\n", e.Line, e.Name, e.URL)
	}
	os.Exit(1)
}
