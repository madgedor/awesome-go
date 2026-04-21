package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func writeMaintainerMarkdown(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "maintainer-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestExtractMaintainerLinks_Found(t *testing.T) {
	content := `## Tools
- [alpha](https://github.com/user/alpha) - A tool.
- [beta](https://github.com/user/beta) - Another tool.
- [gamma](https://example.com/gamma) - Non-GitHub tool.
`
	path := writeMaintainerMarkdown(t, content)
	links, err := extractMaintainerLinks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 GitHub links, got %d", len(links))
	}
	if links[0] != "https://github.com/user/alpha" {
		t.Errorf("expected first link to be alpha, got %s", links[0])
	}
	if links[1] != "https://github.com/user/beta" {
		t.Errorf("expected second link to be beta, got %s", links[1])
	}
}

func TestExtractMaintainerLinks_Dedup(t *testing.T) {
	content := `## Tools
- [alpha](https://github.com/user/alpha) - A tool.
- [alpha-dup](https://github.com/user/alpha) - Duplicate.
`
	path := writeMaintainerMarkdown(t, content)
	links, err := extractMaintainerLinks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 unique link after dedup, got %d", len(links))
	}
}

func TestLastPush_Recent(t *testing.T) {
	// Mock GitHub API returning a recent push date
	recent := time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pushed_at": "` + recent + `"}`))
	}))
	defer server.Close()

	// Override base URL for test by calling lastPush with a full URL
	pushed, err := lastPush(server.URL + "/repos/user/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(pushed) > 60*24*time.Hour {
		t.Errorf("expected recent push, got %v", pushed)
	}
}

func TestLastPush_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not Found"}`))
	}))
	defer server.Close()

	_, err := lastPush(server.URL + "/repos/user/missing")
	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

func TestFindUnmaintained_Empty(t *testing.T) {
	// No links → no unmaintained repos
	results := findUnmaintained([]string{}, 365, "")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestFindUnmaintained_MockStale(t *testing.T) {
	// Mock server returning a very old push date
	old := time.Now().Add(-800 * 24 * time.Hour).Format(time.RFC3339)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pushed_at": "` + old + `"}`))
	}))
	defer server.Close()

	// Use server URL as base; findUnmaintained should flag this repo
	links := []string{"https://github.com/user/oldrepo"}
	results := findUnmaintained(links, 365, server.URL)
	if len(results) != 1 {
		t.Fatalf("expected 1 unmaintained repo, got %d", len(results))
	}
	if results[0] != "https://github.com/user/oldrepo" {
		t.Errorf("unexpected result: %s", results[0])
	}
}
