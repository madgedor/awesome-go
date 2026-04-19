//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	TotalLinks  int       `json:"total_links"`
	StaleLinks  []string  `json:"stale_links"`
	DeadLinks   []string  `json:"dead_links"`
	Duplicates  []string  `json:"duplicates"`
}

func newReport() *Report {
	return &Report{
		GeneratedAt: time.Now().UTC(),
	}
}

func (r *Report) AddStale(url string) {
	r.StaleLinks = append(r.StaleLinks, url)
}

func (r *Report) AddDead(url string) {
	r.DeadLinks = append(r.DeadLinks, url)
}

func (r *Report) AddDuplicate(url string) {
	r.Duplicates = append(r.Duplicates, url)
}

func (r *Report) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (r *Report) Print() {
	fmt.Printf("Report generated at: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Printf("Total links: %d\n", r.TotalLinks)
	fmt.Printf("Stale links: %d\n", len(r.StaleLinks))
	fmt.Printf("Dead links:  %d\n", len(r.DeadLinks))
	fmt.Printf("Duplicates:  %d\n", len(r.Duplicates))
}
