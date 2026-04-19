//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LinkLicense struct {
	URL  string
	Line int
}

var licensePattern = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]+)\)`)

func extractLicenseLinks(filename string) ([]LinkLicense, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var links []LinkLicense
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := licensePattern.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			links = append(links, LinkLicense{URL: m[2], Line: lineNum})
		}
	}
	return links, scanner.Err()
}

func findNonHTTPS(links []LinkLicense) []LinkLicense {
	var bad []LinkLicense
	for _, l := range links {
		if strings.HasPrefix(l.URL, "http://") {
			bad = append(bad, l)
		}
	}
	return bad
}

func main() {
	filename := "README.md"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	links, err := extractLicenseLinks(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	bad := findNonHTTPS(links)
	if len(bad) == 0 {
		fmt.Println("All links use HTTPS.")
		return
	}

	fmt.Printf("Found %d non-HTTPS link(s):\n", len(bad))
	for _, l := range bad {
		fmt.Printf("  line %d: %s\n", l.Line, l.URL)
	}
	os.Exit(1)
}
