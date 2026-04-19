//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var categoryRe = regexp.MustCompile(`^#{1,3} (.+)`)
var linkRe = regexp.MustCompile(`^- \[`)

type CategoryStats struct {
	Name  string
	Links int
}

func extractCategories(filename string) ([]CategoryStats, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var stats []CategoryStats
	var current *CategoryStats

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if m := categoryRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				stats = append(stats, *current)
			}
			current = &CategoryStats{Name: strings.TrimSpace(m[1])}
		} else if linkRe.MatchString(line) && current != nil {
			current.Links++
		}
	}
	if current != nil {
		stats = append(stats, *current)
	}
	return stats, scanner.Err()
}

func findEmptyCategories(stats []CategoryStats) []string {
	var empty []string
	for _, s := range stats {
		if s.Links == 0 {
			empty = append(empty, s.Name)
		}
	}
	return empty
}

func main() {
	filename := "README.md"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	stats, err := extractCategories(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	empty := findEmptyCategories(stats)
	if len(empty) > 0 {
", e)
		tos.Exit(1)
	}
	fmt.Println("All categories have at least one link.")
}
