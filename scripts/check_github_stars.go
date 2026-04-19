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

var githubRe = regexp.MustCompile(`https://github\.com/([^/]+/[^/)\s]+)`)

type repoInfo struct {
	URL   string
	Stars int
}

func extractGithubLinks(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var links []string
	seen := map[string]bool{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matches := githubRe.FindAllString(scanner.Text(), -1)
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				links = append(links, m)
			}
		}
	}
	return links, scanner.Err()
}

func fetchStars(repoURL string) (int, error) {
	parts := strings.TrimPrefix(repoURL, "https://github.com/")
	apiURL := "https://api.github.com/repos/" + parts
	req, _ := http.NewRequest("GET", apiURL, nil)
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var data struct {
		Stars int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Stars, nil
}

func findLowStarRepos(links []string, threshold int) []repoInfo {
	var low []repoInfo
	for _, link := range links {
		stars, err := fetchStars(link)
		if err != nil {
			continue
		}
		if stars < threshold {
			low = append(low, repoInfo{URL: link, Stars: stars})
		}
	}
	return low
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: check_github_stars <README.md>")
		os.Exit(1)
	}
	links, err := extractGithubLinks(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	low := findLowStarRepos(links, 10)
	if len(low) > 0 {
		fmt.Println("Repos with fewer than 10 stars:")
		for _, r := range low {
			fmt.Printf("  %s (%d stars)\n", r.URL, r.Stars)
		}
		os.Exit(1)
	}
	fmt.Println("All repos meet the star threshold.")
}
