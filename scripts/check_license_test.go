package main

import (
	"os"
	"testing"
)

func writeLicenseMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "license_test_*.md")
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

func TestExtractLicenseLinks_HTTPS(t *testing.T) {
	md := "## Section\n- [Foo](https://foo.com) - desc\n- [Bar](https://bar.org) - desc\n"
	path := writeLicenseMarkdown(t, md)
	links, err := extractLicenseLinks(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
}

func TestExtractLicenseLinks_Mixed(t *testing.T) {
	md := "## Section\n- [Foo](http://foo.com) - desc\n- [Bar](https://bar.org) - desc\n"
	path := writeLicenseMarkdown(t, md)
	links, err := extractLicenseLinks(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
}

func TestFindNonHTTPS_None(t *testing.T) {
	links := []LinkLicense{
		{URL: "https://foo.com", Line: 1},
		{URL: "https://bar.com", Line: 2},
	}
	bad := findNonHTTPS(links)
	if len(bad) != 0 {
		t.Fatalf("expected 0 bad links, got %d", len(bad))
	}
}

func TestFindNonHTTPS_Found(t *testing.T) {
	links := []LinkLicense{
		{URL: "http://insecure.com", Line: 3},
		{URL: "https://secure.com", Line: 4},
	}
	bad := findNonHTTPS(links)
	if len(bad) != 1 {
		t.Fatalf("expected 1 bad link, got %d", len(bad))
	}
	if bad[0].URL != "http://insecure.com" {
		t.Errorf("unexpected URL: %s", bad[0].URL)
	}
}
