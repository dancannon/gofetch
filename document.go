package gofetch

import (
	"bytes"
	"code.google.com/p/go.net/html"
)

type Document struct {
	Url     string
	Title   string
	Mime    string
	Meta    map[string]interface{}
	Body    *html.Node
	Content bytes.Buffer
}

func NewDocument(url string) *Document {
	doc := &Document{
		Url: url,
	}
	doc.Meta = map[string]interface{}{}

	return doc
}
