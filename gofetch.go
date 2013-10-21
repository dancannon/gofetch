package gofetch

import (
	"code.google.com/p/go.net/html"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

var scs spew.ConfigState = spew.ConfigState{Indent: "\t"}
var scs2 spew.ConfigState = spew.ConfigState{Indent: "\t", MaxDepth: 2}

func Fetch(url string) (Result, error) {
	// Make request
	res, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}

	// Check the returned MIME type
	if isContentTypeParsable(res) {
		// If the page was HTML then parse the HTMl otherwise return the plain
		// text
		if isContentTypeHtml(res) {
			return Result{
				Url:      res.Request.URL.String(),
				PageType: PlainText,
				Body:     res.Body,
			}, nil
		} else {
			return Result{
				Url:      res.Request.URL.String(),
				PageType: PlainText,
				Body:     res.Body,
			}, nil
		}
	} else {
		// If the content cannot be parsed then guess the page type based on the
		// Content-Type header
		return Result{
			Url:      res.Request.URL.String(),
			PageType: guessPageTypeFromMime(res),
		}, nil
	}
}

type parseState struct {
	state string
}

func newParseState() *parseState {
	return &parseState{
		state: "inline",
	}
}

func prepareDocument(r Result) *Document {
	doc := NewDocument()
	s := newParseState()

	n, err := html.Parse(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	prepareNode(s, n, doc)

	return doc
}

func prepareNode(s *parseState, n *html.Node, d *Document) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			s.state = "title"
			scs2.Dump(n)
		case "meta":
			for _, a := range n.Attr {
				d.Meta[a.Key] = a.Val
			}
		}
	} else if n.Type == html.DocumentNode {
	} else if n.Type == html.CommentNode {
	} else if n.Type == html.DoctypeNode {
	} else if n.Type == html.TextNode {
		switch s.state {
		case "title":
			d.Title = n.Data
		}
	} else {
		scs.Dump(n)
		panic("Unknown node type")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		prepareNode(s, c, d)
	}

	// Run after the end tag
	s.state = "inline"

}
