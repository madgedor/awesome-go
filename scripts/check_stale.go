//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

var linkRe = regexp.MustCompile(`\(https?://[^)]+\)`)

type StaleResult struct {
	URL        string
	LastCommit string
	Stale      bool
}

func extractLinks(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var links []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matches := linkRe.FindAllString(scanner.Text(), -1)
		for _, m := range matches {
			links = append(links, m[1:len(m)-1])
		}
	}
	return links, scanner.Err()
}

func checkStale(url string, threshold time.Duration) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	lastMod := resp.Header.Get("Last-Modified")
	if lastMod == "" {
		return false, nil
	}
	t, err := http.ParseTime(lastMod)
	if err != nil {
		return false, nil
	}
	return time.Since(t) > threshold, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: check_stale <README.md>")
		os.Exit(1)
	}

	links, err := extractLinks(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	threshold := 365 * 24 * time.Hour
	staleCount := 0
	for _, link := range links {
		stale, err := checkStale(link, threshold)
		if err != nil {
			continue
		}
		if stale {
			fmt.Printf("[STALE] %s\n", link)
			staleCount++
		}
	}
	if staleCount > 0 {
		os.Exit(1)
	}
}
