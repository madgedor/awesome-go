//go:build ignore
// +build ignore

// check_maintainer.go checks GitHub repositories in awesome-go README.md
// to identify projects that appear unmaintained (no commits in over 2 years).
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	// staleThreshold defines how long without a commit before a repo is considered unmaintained
	staleThreshold = 2 * 365 * 24 * time.Hour
	githubAPIBase  = "https://api.github.com/repos"
)

var githubRepoRe = regexp.MustCompile(`https://github\.com/([^/]+/[^/)\s]+)`)

// repoInfo holds the last push date returned by the GitHub API.
type repoInfo struct {
	PushedAt time.Time `json:"pushed_at"`
}

// extractMaintainerLinks parses markdown and returns unique GitHub repo slugs.
func extractMaintainerLinks(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := map[string]bool{}
	var repos []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		matches := githubRepoRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			slug := strings.TrimSuffix(m[1], ".git")
			if !seen[slug] {
				seen[slug] = true
				repos = append(repos, slug)
			}
		}
	}
	return repos, scanner.Err()
}

// lastPush fetches the most recent push timestamp for a GitHub repo slug.
// It honours the GITHUB_TOKEN environment variable for authenticated requests.
func lastPush(slug string) (time.Time, error) {
	url := fmt.Sprintf("%s/%s", githubAPIBase, slug)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return time.Time{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("GitHub API returned %d for %s", resp.StatusCode, slug)
	}

	var info repoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return time.Time{}, err
	}
	return info.PushedAt, nil
}

// findUnmaintained returns slugs whose last push is older than staleThreshold.
func findUnmaintained(slugs []string) []string {
	var stale []string
	cutoff := time.Now().Add(-staleThreshold)
	for _, slug := range slugs {
		pushed, err := lastPush(slug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not fetch %s: %v\n", slug, err)
			continue
		}
		if pushed.Before(cutoff) {
			stale = append(stale, fmt.Sprintf("%s (last push: %s)", slug, pushed.Format("2006-01-02")))
		}
	}
	return stale
}

func main() {
	readme := "README.md"
	if len(os.Args) > 1 {
		readme = os.Args[1]
	}

	slugs, err := extractMaintainerLinks(readme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", readme, err)
		os.Exit(1)
	}

	unmaintained := findUnmaintained(slugs)
	if len(unmaintained) == 0 {
		fmt.Println("All repositories appear to be actively maintained.")
		return
	}

	fmt.Printf("Found %d potentially unmaintained repositories (no push in 2+ years):\n", len(unmaintained))
	for _, r := range unmaintained {
		fmt.Println(" -", r)
	}
	os.Exit(1)
}
