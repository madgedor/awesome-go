//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var linkRe = regexp.MustCompile(`\(https?://[^)]+\)`)

func extractSectionLinks(filename string) (map[string][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sections := make(map[string][]string)
	currentSection := "unknown"
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
		}
		matches := linkRe.FindAllString(line, -1)
		for _, m := range matches {
			url := m[1 : len(m)-1]
			sections[currentSection] = append(sections[currentSection], url)
		}
	}
	return sections, scanner.Err()
}

func findCrossSectionDuplicates(sections map[string][]string) map[string][]string {
	seen := make(map[string]string)
	duplicates := make(map[string][]string)
	for section, links := range sections {
		for _, link := range links {
			if prev, ok := seen[link]; ok {
				duplicates[link] = append(duplicates[link], prev, section)
			} else {
				seen[link] = section
			}
		}
	}
	return duplicates
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: check_duplicates_links <file.md>")
		os.Exit(1)
	}
	sections, err := extractSectionLinks(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	duplicates := findCrossSectionDuplicates(sections)
	if len(duplicates) == 0 {
		fmt.Println("No cross-section duplicate links found.")
		return
	}
	fmt.Printf("Found %d cross-section duplicate link(s):\n", len(duplicates))
	for link, secs := range duplicates {
		fmt.Printf("  %s (sections: %s)\n", link, strings.Join(secs, ", "))
	}
	os.Exit(1)
}
