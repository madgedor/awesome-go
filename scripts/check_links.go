package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

var linkRegex = regexp.MustCompile(`https?://[^\s\)\]]+`)

type Result struct {
	URL    string
	Status int
	Err    error
}

func extractLinks(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var links []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matches := linkRegex.FindAllString(scanner.Text(), -1)
		links = append(links, matches...)
	}
	return links, scanner.Err()
}

func checkLink(url string, client *http.Client) Result {
	resp, err := client.Get(url)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	defer resp.Body.Close()
	return Result{URL: url, Status: resp.StatusCode}
}

func main() {
	links, err := extractLinks("README.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading README.md: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	results := make(chan Result, len(links))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for _, link := range links {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results <- checkLink(u, client)
		}(link)
	}

	wg.Wait()
	close(results)

	failed := 0
	for r := range results {
		if r.Err != nil || r.Status >= 400 {
			fmt.Printf("DEAD [%d] %s\n", r.Status, r.URL)
			failed++
		}
	}
	if failed > 0 {
		fmt.Printf("\n%d dead link(s) found\n", failed)
		os.Exit(1)
	}
	fmt.Println("All links OK")
}
