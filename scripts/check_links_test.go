package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestExtractLinks(t *testing.T) {
	f, err := os.CreateTemp("", "readme-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString("# Awesome Go\n")
	_, _ = f.WriteString("- [pkg](https://pkg.go.dev) - some package\n")
	_, _ = f.WriteString("- [github](https://github.com/example/repo) - another\n")
	f.Close()

	links, err := extractLinks(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
}

func TestCheckLink_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	r := checkLink(ts.URL, client)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Status != http.StatusOK {
		t.Fatalf("expected 200, got %d", r.Status)
	}
}

func TestCheckLink_Dead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	r := checkLink(ts.URL, client)
	if r.Status != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", r.Status)
	}
}
