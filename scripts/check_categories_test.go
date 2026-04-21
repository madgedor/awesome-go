package main

import (
	"os"
	"testing"
)

func writeCategoryMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "categories-*.md")
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

func TestExtractCategories_WithLinks(t *testing.T) {
	md := `## Audio
- [beep](https://github.com/faiface/beep) - Sound library.
- [flac](https://github.com/mewkiz/flac) - FLAC support.
## Video
- [gmf](https://github.com/3d0c/gmf) - Go bindings for FFmpeg.
`
	path := writeCategoryMarkdown(t, md)
	stats, err := extractCategories(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(stats))
	}
	if stats[0].Links != 2 {
		t.Errorf("expected 2 links in Audio, got %d", stats[0].Links)
	}
	if stats[1].Links != 1 {
		t.Errorf("expected 1 link in Video, got %d", stats[1].Links)
	}
}

func TestExtractCategories_Empty(t *testing.T) {
	md := `## Audio
- [beep](https://github.com/faiface/beep) - Sound library.
## EmptySection
## Video
- [gmf](https://github.com/3d0c/gmf) - Bindings.
`
	path := writeCategoryMarkdown(t, md)
	stats, err := extractCategories(path)
	if err != nil {
		t.Fatal(err)
	}
	empty := findEmptyCategories(stats)
	if len(empty) != 1 || empty[0] != "EmptySection" {
		t.Errorf("expected [EmptySection], got %v", empty)
	}
}

func TestFindEmptyCategories_None(t *testing.T) {
	stats := []CategoryStats{
		{Name: "Audio", Links: 3},
		{Name: "Video", Links: 1},
	}
	empty := findEmptyCategories(stats)
	if len(empty) != 0 {
		t.Errorf("expected no empty categories, got %v", empty)
	}
}

// TestFindEmptyCategories_AllEmpty verifies that all categories are reported
// as empty when none of them contain any links.
func TestFindEmptyCategories_AllEmpty(t *testing.T) {
	stats := []CategoryStats{
		{Name: "Audio", Links: 0},
		{Name: "Video", Links: 0},
		{Name: "Databases", Links: 0},
	}
	empty := findEmptyCategories(stats)
	if len(empty) != 3 {
		t.Errorf("expected 3 empty categories, got %v", empty)
	}
}
