package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

var URL = "https://scrape-me.dreamsofcode.io/"
var hmap = map[string]bool{}

func main() {
	resp, err := http.Get(URL)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	links, err := extractLinks(resp.Body)
	if err != nil {
		panic(err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
}

func extractLinks(body io.Reader) ([]string, error) {
	var links []string
	// Parse into Tree Node
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML %w", err)
	}
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					var link string
					if attr.Val[:4] != "http" {
						link = URL[:len(URL)-1] + attr.Val
					} else {
						link = attr.Val
					}
					if !hmap[link] {
						links = append(links, link)
						hmap[link] = true
						break
					} else {
						continue
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return links, nil
}
