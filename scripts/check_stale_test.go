package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func writeMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.md")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestExtractLinks_Stale(t *testing.T) {
	path := writeMarkdown(t, "- [foo](https://example.com/foo) - desc\n- [bar](https://go.dev) - desc\n")
	links, err := extractLinks(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
}

func TestCheckStale_NotStale(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	stale, err := checkStale(ts.URL, 365*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if stale {
		t.Error("expected not stale")
	}
}

func TestCheckStale_Stale(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		w.Header().Set("Last-Modified", old.Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	stale, err := checkStale(ts.URL, 365*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if !stale {
		t.Error("expected stale")
	}
}

func TestCheckStale_NoHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	stale, err := checkStale(ts.URL, 365*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if stale {
		t.Error("expected not stale when no header")
	}
}
