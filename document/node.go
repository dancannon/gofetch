package document

import (
	"bytes"
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
