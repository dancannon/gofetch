package document

import (
	"bytes"
	"io"

	"code.google.com/p/go.net/html"
)

type HtmlNode html.Node

func (n *HtmlNode) MarshalJSON() ([]byte, error) {
	w := bytes.Buffer{}
	err := html.Render(&w, n.Node())

	return w.Bytes(), err
}

func (n *HtmlNode) Node() *html.Node {
	return (*html.Node)(n)
}

type Document struct {
	Url   string              `json:"url"`
	Title string              `json:"title"`
	Meta  []map[string]string `json:"meta"`
	Doc   *HtmlNode           `json:"doc"`
	Body  *HtmlNode           `json:"body"`
}

func NewDocument(url string, r io.Reader) *Document {
	doc := &Document{
		Url:  url,
		Meta: []map[string]string{},
	}

	// Parse the html
	n, err := html.Parse(r)
	if err != nil {
		panic("Error parsing the html")
	}

	doc.Doc = (*HtmlNode)(&(*n))

	// Process the document html to extract the title/meta tags
	var processNode func(*html.Node)
	var processTitleNode func(*html.Node)

	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				processTitleNode(c)
			}
			return
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			attrs := map[string]string{}
			for _, a := range n.Attr {
				attrs[a.Key] = a.Val
			}

			doc.Meta = append(doc.Meta, attrs)
		} else if n.Type == html.ElementNode && n.Data == "body" {
			doc.Body = (*HtmlNode)(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	processTitleNode = func(n *html.Node) {
		if n.Type == html.TextNode {
			doc.Title += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processTitleNode(c)
		}
	}
	processNode(n)

	return doc
}
