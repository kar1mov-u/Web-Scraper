package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var URL = "https://scrape-me.dreamsofcode.io/"
var hmap = map[string]bool{}
var dead = []string{}

func main() {
	// Start crawling with a depth of 5
	links := helper(URL, 3)
	fmt.Println("Links found:", links)
}

func helper(link string, depth int) []string {
	if depth == 0 || hmap[link] {
		return nil
	}

	// Mark link as visited
	hmap[link] = true
	fmt.Println("Crawling:", link, "at depth:", depth)

	// Fetch the page
	resp, err := http.Get(link)

	if err != nil {
		fmt.Println("Error fetching:", link, err)
		return nil
	}
	if resp.StatusCode >= 400 {
		fmt.Printf("Dead link is found %s (status %d)\n", link, resp.StatusCode)
		dead = append(dead, link)
	}
	defer resp.Body.Close()

	// Extract links
	links, err := extractLinks(resp.Body, link)
	if err != nil {
		fmt.Println("Error extracting links from:", link, err)
		return nil
	}

	// Recursive crawling
	res := []string{link} // Include the current link
	for _, l := range links {
		res = append(res, helper(l, depth-1)...)
	}
	return res
}

func extractLinks(body io.Reader, base string) ([]string, error) {
	var links []string

	// Parse the HTML document
	doc, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Traverse the HTML tree
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// Resolve relative URLs
					link, err := resolveURL(attr.Val, base)
					if err != nil || len(link) < 5 {
						continue
					}
					if !hmap[link] {
						links = append(links, link)
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

func resolveURL(href, base string) (string, error) {
	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(u).String(), nil
}
