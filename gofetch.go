package gofetch

import (
	"code.google.com/p/go.net/html"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	doc := NewDocument(r.Url)
	s := newParseState()

	n, err := html.Parse(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	prepareNode(s, n, doc)

	return doc
}

var (
	ignorableIdentifiers = []string{
		"comment", "extra", "foot", "head", "topbar", "nav", "menu", "sidebar", "page",
		"breadcrumb", "hide", "hidden", "no-?display", "\\bad\\b", "advert", "promo",
		"featured", "toolbox", "toolbar", "tools", "actions", "buttons", "related",
		"share", "social", "pop",
	}
)

func prepareNode(s *parseState, n *html.Node, d *Document) {
	if n.Data == "script" {
		s.state = "inline"
		n.Parent.RemoveChild(n)
		return
	}
	if n.Type == html.ElementNode {
		// Ensure that the body tag is added to the result document
		if n.Data == "body" {
			d.Body = n
		} else {
			tmpAttrs := []html.Attribute{}
			for _, a := range n.Attr {
				if a.Key == "id" || a.Key == "class" {
					for _, ident := range ignorableIdentifiers {
						matched, _ := regexp.MatchString(ident, strings.ToLower(a.Val))
						if matched {
							n.Parent.RemoveChild(n)
							s.state = "inline"
							return
						}
					}
				} else if a.Key == "href" || a.Key == "src" {
					// Attempt to fix URLs
					urlb, err := url.Parse(d.Url)
					if err != nil {
						continue
					}
					urlr, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					a.Val = urlb.ResolveReference(urlr).String()
				}

				tmpAttrs = append(tmpAttrs, a)
			}
			n.Attr = tmpAttrs

			switch n.Data {
			case "title":
				s.state = "title"
			case "meta":
				for _, a := range n.Attr {
					d.Meta[a.Key] = a.Val
				}
			// Remove un-needed tags
			case "script", "style", "link", "noscript":
				n.Parent.RemoveChild(n)
				s.state = "inline"
				return
			}
		}
	} else if n.Type == html.DocumentNode {
	} else if n.Type == html.CommentNode {
		n.Parent.RemoveChild(n)
	} else if n.Type == html.DoctypeNode {
	} else if n.Type == html.TextNode {
		switch s.state {
		case "title":
			d.Title = n.Data
		}
	} else {
		panic("Unknown node type")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		prepareNode(s, c, d)
	}

	// Run after the end tag
	s.state = "inline"
}
