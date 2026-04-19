package main

import (
	"os"
	"testing"
)

func writeTempMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "awesome-*.md")
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

func TestExtractAllLinks_NoDuplicates(t *testing.T) {
	md := "- [Foo](https://foo.com)\n- [Bar](https://bar.com)\n"
	file := writeTempMarkdown(t, md)

	urlLines, err := extractAllLinks(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(urlLines) != 2 {
		t.Fatalf("expected 2 links, got %d", len(urlLines))
	}
}

func TestExtractAllLinks_WithDuplicates(t *testing.T) {
	md := "- [Foo](https://foo.com)\n- [Foo again](https://foo.com)\n"
	file := writeTempMarkdown(t, md)

	urlLines, err := extractAllLinks(file)
	if err != nil {
		t.Fatal(err)
	}
	lines, ok := urlLines["https://foo.com"]
	if !ok {
		t.Fatal("expected https://foo.com to be present")
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 occurrences, got %d", len(lines))
	}
}

func TestFindDuplicates(t *testing.T) {
	urlLines := map[string][]int{
		"https://foo.com": {1, 5},
		"https://bar.com": {2},
	}
	dupes := findDuplicates(urlLines)
	if len(dupes) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(dupes))
	}
	if _, ok := dupes["https://foo.com"]; !ok {
		t.Error("expected https://foo.com to be a duplicate")
	}
}

func TestFindDuplicates_None(t *testing.T) {
	urlLines := map[string][]int{
		"https://a.com": {1},
		"https://b.com": {3},
	}
	dupes := findDuplicates(urlLines)
	if len(dupes) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(dupes))
	}
}
