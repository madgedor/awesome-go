//go:build ignore

package main

import (
	"os"
	"testing"
)

func writeSectionMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.md")
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

func TestExtractSectionLinks(t *testing.T) {
	md := `## Web
- [Gin](https://github.com/gin-gonic/gin)
## CLI
- [Cobra](https://github.com/spf13/cobra)
`
	path := writeSectionMarkdown(t, md)
	sections, err := extractSectionLinks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sections["Web"]) != 1 {
		t.Errorf("expected 1 link in Web, got %d", len(sections["Web"]))
	}
	if len(sections["CLI"]) != 1 {
		t.Errorf("expected 1 link in CLI, got %d", len(sections["CLI"]))
	}
}

func TestFindCrossSectionDuplicates_Found(t *testing.T) {
	sections := map[string][]string{
		"Web": {"https://github.com/gin-gonic/gin"},
		"CLI": {"https://github.com/gin-gonic/gin"},
	}
	duplicates := findCrossSectionDuplicates(sections)
	if len(duplicates) != 1 {
		t.Errorf("expected 1 duplicate, got %d", len(duplicates))
	}
}

func TestFindCrossSectionDuplicates_None(t *testing.T) {
	sections := map[string][]string{
		"Web": {"https://github.com/gin-gonic/gin"},
		"CLI": {"https://github.com/spf13/cobra"},
	}
	duplicates := findCrossSectionDuplicates(sections)
	if len(duplicates) != 0 {
		t.Errorf("expected 0 duplicates, got %d", len(duplicates))
	}
}

// TestFindCrossSectionDuplicates_MultipleURLs checks that multiple duplicate
// URLs across sections are all detected correctly.
func TestFindCrossSectionDuplicates_MultipleURLs(t *testing.T) {
	sections := map[string][]string{
		"Web":     {"https://github.com/gin-gonic/gin", "https://github.com/spf13/cobra"},
		"CLI":     {"https://github.com/spf13/cobra"},
		"Routers": {"https://github.com/gin-gonic/gin"},
	}
	duplicates := findCrossSectionDuplicates(sections)
	if len(duplicates) != 2 {
		t.Errorf("expected 2 duplicates, got %d", len(duplicates))
	}
}
