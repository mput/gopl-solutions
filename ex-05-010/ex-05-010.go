package main

import (
	"fmt"
)

// Rewrite topoSort using map.

var prereqs = map[string]map[string]bool{
	"alogrithms": {"data structures": true},
	"calculus":   {"linear algebra": true},
	"compilers": {
		"data structures":       true,
		"formal languages":      true,
		"computer organization": true,
	},
	"data structures":       {"discret math": true},
	"databases":             {"data structures": true},
	"discret math":          {"intro to programing": true},
	"formal languages":      {"discret math": true},
	"networks":              {"operating systems": true},
	"operating systems":     {"data structures": true, "computer organization": true},
	"programming languages": {"data structures": true, "computer organization": true},
}

func main() {
	for i, elm := range topoSort(prereqs) {
		fmt.Printf("%2d -> %s\n", i+1, elm)
	}
}

func topoSort(p map[string]map[string]bool) []string {
	res := []string{}
	var sortAll func(keys []string)
	seen := make(map[string]bool)

	sortAll = func(keys []string) {
		for _, key := range keys {
			if !seen[key] {
				seen[key] = true
				preKeys := []string{}
				for k, _ := range p[key] {
					preKeys = append(preKeys, k)
				}
				sortAll(preKeys)
				res = append(res, key)
			}
		}
	}

	keys := []string{}
	for k, _ := range p {
		keys = append(keys, k)
	}
	sortAll(keys)
	return res
}
