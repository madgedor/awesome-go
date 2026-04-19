package main

import (
	"os"
	"testing"
)

func writeDescMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "desc-*.md")
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

func TestExtractLinkEntries_WithDesc(t *testing.T) {
	md := `## Section
- [Foo](https://foo.com) - A foo library
- [Bar](https://bar.com) - A bar library
`
	path := writeDescMarkdown(t, md)
	entries, err := extractLinkEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Desc == "" {
		t.Error("expected description for Foo")
	}
}

func TestExtractLinkEntries_MissingDesc(t *testing.T) {
	md := `## Section
- [Foo](https://foo.com) - A foo library
- [Bar](https://bar.com)
`
	path := writeDescMarkdown(t, md)
	entries, err := extractLinkEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].Desc != "" {
		t.Errorf("expected empty description for Bar, got %q", entries[1].Desc)
	}
}

func TestFindMissingDescriptions_None(t *testing.T) {
	entries := []LinkEntry{
		{Name: "Foo", URL: "https://foo.com", Desc: "A library"},
	}
	missing := findMissingDescriptions(entries)
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %d", len(missing))
	}
}

func TestFindMissingDescriptions_Found(t *testing.T) {
	entries := []LinkEntry{
		{Name: "Foo", URL: "https://foo.com", Desc: ""},
		{Name: "Bar", URL: "https://bar.com", Desc: "ok"},
	}
	missing := findMissingDescriptions(entries)
	if len(missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(missing))
	}
	if missing[0].Name != "Foo" {
		t.Errorf("expected Foo, got %s", missing[0].Name)
	}
}
