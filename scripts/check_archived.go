//go:build ignore

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var githubRepoRe = regexp.MustCompile(`https://github\.com/([^/]+/[^/)\s]+)`)

func extractGithubRepos(path string) ([]string, error) {
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
		matches := githubRepoRe.FindAllString(line, -1)
		for _, m := range matches {
			m = strings.TrimRight(m, "/")
			if !seen[m] {
				seen[m] = true
				repos = append(repos, m)
			}
		}
	}
	return repos, scanner.Err()
}

func isArchived(repoURL, token string) (bool, error) {
	parts := strings.TrimPrefix(repoURL, "https://github.com/")
	apiURL := "https://api.github.com/repos/" + parts
	req, _ := http.NewRequest("GET", apiURL, nil)
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var data struct {
		Archived bool `json:"archived"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, err
	}
	return data.Archived, nil
}

func findArchivedRepos(repos []string, token string) []string {
	var archived []string
	for _, r := range repos {
		ok, err := isArchived(r, token)
		if err == nil && ok {
			archived = append(archived, r)
		}
	}
	return archived
}

func main() {
	path := "README.md"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	token := os.Getenv("GITHUB_TOKEN")
	repos, err := extractGithubRepos(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	archived := findArchivedRepos(repos, token)
	if len(archived) > 0("Archived repositories:")
		.Println("}
		os.Exit(1)
	}
	fmt.Println("No archived repositories found.")
}
