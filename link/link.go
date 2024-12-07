package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, linkBuilder(node))
	}
	return links, nil
}
func linkBuilder(n *html.Node) Link {
	var link Link
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
		}
	}
	link.Text = text(n)
	return link
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var res string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res += text(c)
	}
	return strings.Join(strings.Fields(res), "")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, linkNodes(c)...)
	}
	return nodes
}
