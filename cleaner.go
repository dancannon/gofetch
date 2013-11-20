package gofetch

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
	"net/url"
	"regexp"
	"strings"
)

var (
	ignorableIdentifiers = "comment|extra|foot|topbar|nav|menu|sidebar|breadcrumb|hide|hidden|no-?display|\\bad\\b|advert|promo|featured|toolbox|toolbar|tools|actions|buttons|related|share|social|facebook|twitter|google|pop|links|meta$|scroll|shoutbox|sponsor|contact|form|community|subscribe"
	ignorableRegex       = regexp.MustCompile(ignorableIdentifiers)
)

func cleanDocument(d *document.Document) {
	cleanNode(d.Doc, d)
}

func cleanNode(n *html.Node, d *document.Document) {
	if n.Type == html.ElementNode {
		// Ensure that the body tag is added to the result document
		if n.Data == "body" {
			d.Body = n
		} else {
			tmpAttrs := []html.Attribute{}
			for _, a := range n.Attr {
				if a.Key == "id" || a.Key == "class" || a.Key == "name" {
					if ignorableRegex.MatchString(strings.ToLower(a.Val)) {
						n.Parent.RemoveChild(n)
						return
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
			// Remove un-needed tags
			case "script", "style", "link", "noscript":
				n.Parent.RemoveChild(n)
				return
			}
		}
	} else if n.Type == html.CommentNode {
		n.Parent.RemoveChild(n)
	}

	// Build the list of children node before iterating. This is needed because we will be
	// deleting nodes
	childNodes := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childNodes = append(childNodes, c)
	}

	for _, c := range childNodes {
		if c != nil {
			cleanNode(c, d)
		}
	}
}
