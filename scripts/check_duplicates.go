package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var linkLineRe = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]+)\)`)

// extractAllLinks scans a markdown file and returns a map of url -> list of line numbers.
func extractAllLinks(filename string) (map[string][]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	urlLines := make(map[string][]int)
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := linkLineRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			url := strings.TrimSpace(m[2])
			urlLines[url] = append(urlLines[url], lineNum)
		}
	}
	return urlLines, scanner.Err()
}

// findDuplicates returns urls that appear more than once.
func findDuplicates(urlLines map[string][]int) map[string][]int {
	dupes := make(map[string][]int)
	for url, lines := range urlLines {
		if len(lines) > 1 {
			dupes[url] = lines
		}
	}
	return dupes
}

func main() {
	filename := "README.md"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	urlLines, err := extractAllLinks(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	dupes := findDuplicates(urlLines)
	if len(dupes) == 0 {
		fmt.Println("No duplicate links found.")
		return
	}

	fmt.Printf("Found %d duplicate link(s):\n", len(dupes))
	for url, lines := range dupes {
		fmt.Printf("  %s (lines: %v)\n", url, lines)
	}
	os.Exit(1)
}
