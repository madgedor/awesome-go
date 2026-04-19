package main

import (
	"os"
	"testing"
)

func writeURLMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "urltest-*.md")
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

func TestExtractURLEntries_Valid(t *testing.T) {
	md := "## Section\n- [Foo](https://example.com) - desc\n- [Bar](https://go.dev) - desc\n"
	path := writeURLMarkdown(t, md)
	entries, err := extractURLEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestExtractURLEntries_Empty(t *testing.T) {
	path := writeURLMarkdown(t, "## No links here\n")
	entries, err := extractURLEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestFindMalformedURLs_None(t *testing.T) {
	entries := []URLEntry{
		{Line: 1, URL: "https://example.com"},
		{Line: 2, URL: "https://go.dev/pkg"},
	}
	bad := findMalformedURLs(entries)
	if len(bad) != 0 {
		t.Fatalf("expected no bad URLs, got %d", len(bad))
	}
}

func TestFindMalformedURLs_Found(t *testing.T) {
	entries := []URLEntry{
		{Line: 1, URL: "https://good.com"},
		{Line: 2, URL: "https://bad url.com"},
		{Line: 3, URL: "not-a-url"},
	}
	bad := findMalformedURLs(entries)
	if len(bad) != 2 {
		t.Fatalf("expected 2 bad URLs, got %d", len(bad))
	}
}
