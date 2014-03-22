package document

import (
	"bytes"
	"io"
	"net/url"

	"code.google.com/p/go.net/html"
)

type Document struct {
	URL   *url.URL            `json:"url"`
	Title string              `json:"title"`
	Meta  []map[string]string `json:"meta"`
	Doc   *HtmlNode           `json:"doc"`
	Body  *HtmlNode           `json:"body"`
	Raw   string              `json:"raw"`
}

func NewDocument(rawurl string, r io.Reader) (*Document, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		URL:  u,
		Meta: []map[string]string{},
	}

	// Parse the html
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	w := &bytes.Buffer{}
	err = html.Render(w, n)
	if err != nil {
		return nil, err
	}

	doc.Raw = w.String()
	doc.Doc = (*HtmlNode)(n)

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

	return doc, nil
}
