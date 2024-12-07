package main

import (
	"cligame/link"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Loc struct {
	Value string `xml:"loc"`
}
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Str     string   `xml:"xlmns,attr"`
	Urls    []Loc    `xml:"url"`
}

func main() {
	addr := flag.String("addr", "https://gophercises.com", "enter web address")
	maxDepth := flag.Int("depth", 10, "provide depth of search")
	flag.Parse()

	pages := bfs(*addr, *maxDepth)
	var toXml Urlset
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, Loc{page})
	}
	toXml.Str = "http://www.sitemaps.org/schemas/sitemap/0.9"
	output, _ := os.Create("xmlF")
	defer output.Close()

	output.WriteString(xml.Header)

	enc := xml.NewEncoder(output)
	enc.Indent("", "  ")
	err := enc.Encode(toXml)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func bfs(urlStr string, depth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i < depth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	var ret []string
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	initialLinks, _ := link.Parse(r)
	var ret []string
	for _, l := range initialLinks {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// another way to crawl a website

func crawlWebsite(addr string, links []link.Link, visitedLinks []string) ([]link.Link, []string) {
	var allLinks []link.Link
	visitedSet := make(map[string]bool)
	for _, v := range visitedLinks {
		visitedSet[v] = true
	}
	for _, l := range links {
		if strings.HasPrefix(l.Href, "/") {
			l.Href = addr + l.Href
		}
		if strings.HasPrefix(l.Href, addr) && !visitedSet[l.Href] {
			visitedSet[l.Href] = true
			visitedLinks = append(visitedLinks, l.Href)

			resp, err := http.Get(l.Href)
			if err != nil {
				log.Println("error fetching:", l.Href, err)
				continue
			}
			defer resp.Body.Close()

			newLinks, err := link.Parse(resp.Body)
			if err != nil {
				log.Println("error parsing", l.Href, err)
				continue
			}

			allLinks = append(allLinks, newLinks...)
			allLinks, visitedLinks = crawlWebsite(addr, allLinks, visitedLinks)
		}

		allLinks = append(allLinks, l)
	}

	return allLinks, visitedLinks
}

func removeDuplicates(links []string) []string {
	uniqueUrls := make(map[string]bool)
	var result []string

	for _, link := range links {
		if !uniqueUrls[link] {
			uniqueUrls[link] = true
			result = append(result, link)
		}
	}
	return result
}
