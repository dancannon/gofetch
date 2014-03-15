package gofetch

import (
	"net/url"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
)

var (
	ignorableIdentifiers = "comment|extra|foot|topbar|nav|menu|sidebar|breadcrumb|hide|hidden|no-?display|\\bad\\b|advert|promo|featured|toolbox|toolbar|tools|actions|buttons|related|share|social|facebook|twitter|google|pop|links|meta$|scroll|shoutbox|sponsor|contact|form|community|subscribe"
	ignorableRegex       = regexp.MustCompile(ignorableIdentifiers)
)

func cleanDocument(d *document.Document) {
	cleanNode(d.Body.Node(), d)
}

func cleanNode(n *html.Node, d *document.Document) {
	if n.Type == html.ElementNode {
		switch n.Data {
		// Remove un-needed tags
		case "script", "style", "link", "noscript":
			deleteNode(n)
			return
		}

		if n.Data != "body" {
			// Ensure that the body tag is added to the result document
			tmpAttrs := []html.Attribute{}
			for _, a := range n.Attr {
				if a.Key == "id" || a.Key == "class" || a.Key == "name" {
					if ignorableRegex.MatchString(strings.ToLower(a.Val)) {
						deleteNode(n)
						return
					}
				} else if a.Key == "href" || a.Key == "src" {
					// Attempt to fix URLs
					urlr, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					a.Val = d.URL.ResolveReference(urlr).String()
				}

				tmpAttrs = append(tmpAttrs, a)
			}
			n.Attr = tmpAttrs
		} else if n.Type == html.CommentNode {
			deleteNode(n)
		}
	}

	// Build the list of children node before iterating. This is needed because we will be
	// deleting nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cleanNode(c, d)
	}
}

func deleteNode(n *html.Node) {
	if n.Parent != nil {

	}
	if n.Parent.FirstChild == n {
		n.FirstChild = n.NextSibling
	}
	if n.NextSibling != nil {
		n.NextSibling.PrevSibling = n.PrevSibling
	}
	if n.Parent.LastChild == n {
		n.LastChild = n.PrevSibling
	}
	if n.PrevSibling != nil {
		n.PrevSibling.NextSibling = n.NextSibling
	}
	n.Parent = nil
}
