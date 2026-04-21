package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewReport(t *testing.T) {
	r := newReport()
	if r == nil {
		t.Fatal("expected non-nil report")
	}
	if r.GeneratedAt.IsZero() {
		t.Error("expected non-zero generated_at")
	}
}

func TestReport_AddEntries(t *testing.T) {
	r := newReport()
	r.AddStale("https://stale.example.com")
	r.AddDead("https://dead.example.com")
	r.AddDuplicate("https://dup.example.com")

	if len(r.StaleLinks) != 1 {
		t.Errorf("expected 1 stale, got %d", len(r.StaleLinks))
	}
	if len(r.DeadLinks) != 1 {
		t.Errorf("expected 1 dead, got %d", len(r.DeadLinks))
	}
	if len(r.Duplicates) != 1 {
		t.Errorf("expected 1 duplicate, got %d", len(r.Duplicates))
	}
}

// TestReport_AddMultipleEntries verifies that adding multiple entries of the
// same category accumulates correctly rather than overwriting.
func TestReport_AddMultipleEntries(t *testing.T) {
	r := newReport()
	r.AddDead("https://dead1.example.com")
	r.AddDead("https://dead2.example.com")
	r.AddDead("https://dead3.example.com")

	if len(r.DeadLinks) != 3 {
		t.Errorf("expected 3 dead links, got %d", len(r.DeadLinks))
	}
}

func TestReport_Save(t *testing.T) {
	r := newReport()
	r.TotalLinks = 42
	r.AddDead("https://dead.example.com")

	tmp, err := os.CreateTemp("", "report-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	t.Cleanup(func() { os.Remove(tmp.Name()) })

	if err := r.Save(tmp.Name()); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var loaded Report
	if err := json.NewDecoder(f).Decode(&loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.TotalLinks != 42 {
		t.Errorf("expected 42, got %d", loaded.TotalLinks)
	}
	if len(loaded.DeadLinks) != 1 {
		t.Errorf("expected 1 dead link, got %d", len(loaded.DeadLinks))
	}
}
