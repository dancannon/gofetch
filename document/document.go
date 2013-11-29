package document

import (
	"code.google.com/p/go.net/html"
	"io"
)

type Document struct {
	Url   string
	Title string
	Meta  []map[string]string
	Doc   *html.Node
	Body  *html.Node
	Raw   io.ReadCloser
}

func NewDocument(url string, r io.ReadCloser) *Document {
	doc := &Document{
		Url:  url,
		Raw:  r,
		Meta: []map[string]string{},
	}

	// Parse the html
	n, err := html.Parse(r)
	if err != nil {
		panic("Error parsing the html")
	}

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
			doc.Doc = n
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
