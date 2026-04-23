package main

import (
	"os"
	"testing"
)

func writeAlphaMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "alpha-*.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestExtractSections_Sorted(t *testing.T) {
	md := `## Web
- [Alpha](https://alpha.io) - A.
- [Beta](https://beta.io) - B.
- [Gamma](https://gamma.io) - G.
`
	path := writeAlphaMarkdown(t, md)
	sections, err := extractSections(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}
	if len(sections[0].Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(sections[0].Items))
	}
}

func TestFindNonAlphabetical_Sorted(t *testing.T) {
	sections := []sectionItems{
		{Header: "## Web", Items: []string{"alpha", "beta", "gamma"}},
	}
	issues := findNonAlphabetical(sections)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestFindNonAlphabetical_Unsorted(t *testing.T) {
	sections := []sectionItems{
		{Header: "## Web", Items: []string{"alpha", "gamma", "beta"}},
	}
	issues := findNonAlphabetical(sections)
	if len(issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(issues))
	}
}

func TestFindNonAlphabetical_MultiSection(t *testing.T) {
	sections := []sectionItems{
		{Header: "## A", Items: []string{"apple", "banana"}},
		{Header: "## B", Items: []string{"zebra", "ant"}},
	}
	issues := findNonAlphabetical(sections)
	if len(issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(issues))
	}
}

// TestFindNonAlphabetical_EmptySection verifies that sections with no items
// or only one item do not produce false positives.
func TestFindNonAlphabetical_EmptySection(t *testing.T) {
	sections := []sectionItems{
		{Header: "## Empty", Items: []string{}},
		{Header: "## Single", Items: []string{"only"}},
	}
	issues := findNonAlphabetical(sections)
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty/single-item sections, got %v", issues)
	}
}
