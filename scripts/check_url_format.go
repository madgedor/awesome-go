//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var linkLineRe = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]+)\)`)

type URLEntry struct {
	Line int
	URL  string
	Raw  string
}

func extractURLEntries(path string) ([]URLEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []URLEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := linkLineRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			entries = append(entries, URLEntry{Line: lineNum, URL: m[2], Raw: m[0]})
		}
	}
	return entries, scanner.Err()
}

func findMalformedURLs(entries []URLEntry) []URLEntry {
	var bad []URLEntry
	for _, e := range entries {
		u, err := url.Parse(e.URL)
		if err != nil || u.Host == "" || strings.Contains(e.URL, " ") {
			bad = append(bad, e)
		}
	}
	return bad
}

func main() {
	path := "README.md"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	entries, err := extractURLEntries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	bad := findMalformedURLs(entries)
	if len(bad) == 0 {
		fmt.Println("All URLs are well-formed.")
		return
	}
	fmt.Printf("Found %d malformed URL(s):\n", len(bad))
	for _, e := range bad {
		fmt.Printf("  line %d: %s\n", e.Line, e.URL)
	}
	os.Exit(1)
}
