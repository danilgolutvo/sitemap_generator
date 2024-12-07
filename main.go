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
	"strings"
)

type URL struct {
	Loc string `xml:"loc"`
}
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

func main() {
	addr := flag.String("addr", "https://gophercises.com", "enter web address")
	flag.Parse()
	resp, err := http.Get(*addr)
	if err != nil {
		log.Fatal(err)
	}

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	pages := hrefs(resp.Body, base)
	for _, l := range pages {
		fmt.Println(l)
	}
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

//allLinks, visited := crawlWebsite(*addr, initialLinks, nil)
////log.Println("All Links:", allLinks)
//_ = allLinks
//fmt.Printf(generateSitemap(removeDuplicates(visited)))

// visit all urls from gotten list of links

// need to figure out a way to determine if a link goes to the same domain or a different one

// when all links gotten, this data has to be outputted in XML format

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

func generateSitemap(links []string) string {
	var urls []URL
	for _, link := range links {
		urls = append(urls, URL{Loc: link})
	}
	sitemap := Sitemap{URLs: urls}

	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		log.Fatal("Error generating sitemap:", err)
	}
	return string(output)
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
