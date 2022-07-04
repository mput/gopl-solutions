package main

import (
	"log"
	"net/http"

	"mput.me/gopl/links"
)

func breadthFirst(f func(item string) []string, worklist []string) {
}

func processPage(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	links.Extact()
	return []string{}
}

func main() {
	breadthFirst(processPage, []string{"https://go.dev/"})
}
