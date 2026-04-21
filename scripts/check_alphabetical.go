//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var sectionHeaderRe = regexp.MustCompile(`^#{1,3} `)
var listItemRe = regexp.MustCompile(`^[-*] `)

type sectionItems struct {
	Header string
	Items  []string
}

func extractSections(Items, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sections []sectionItems
	var current *sectionItems

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if sectionHeaderRe.MatchString(line) {
			if current != nil {
				sections = append(sections, *current)
			}
			current = &sectionItems{Header: line}
		} else if current != nil {
			if m := listItemRe.FindStringSubmatch(line); m != nil {
				// Strip leading articles ("a ", "an ", "the ") before comparing,
				// so e.g. "The Foo" sorts alongside "Foo" rather than at the end.
				name := strings.ToLower(m[1])
				name = stripArticle(name)
				current.Items = append(current.Items, name)
			}
		}
	}
	if current != nil {
		sections = append(sections, *current)
	}
	return sections, scanner.Err()
}

// stripArticle removes a leading "a ", "an ", or "the " from s.
func stripArticle(s string) string {
	for _, article := range []string{"the ", "an ", "a "} {
		if strings.HasPrefix(s, article) {
			return strings.TrimPrefix(s, article)
		}
	}
	return s
}

func findNonAlphabetical(sections []sectionItems) []string {
	var issues []string
	for _, sec := range sections {
		for i := 1; i < len(sec.Items); i++ {
			if sec.Items[i] < sec.Items[i-1] {
				issues = append(issues, fmt.Sprintf("section %q: %q should come before %q", sec.Header, sec.Items[i], sec.Items[i-1]))
			}
		}
	}
	return issues
}

func main() {
	filename := "README.md"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	sections, err := extractSections(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	issues := findNonAlphabetical(sections)
	if len(issues) > 0 {
		for _, iss := range issues {
			fmt.Println(iss)
		}
		os.Exit(1)
	}
	fmt.Println("All sections are alphabetically sorted.")
}
