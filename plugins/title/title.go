package title

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"sort"
	"strings"
)

type title struct {
	text  string
	level int
}

type titleSlice []title

func (s titleSlice) Len() int           { return len(s) }
func (s titleSlice) Less(i, j int) bool { return s[i].level < s[j].level }
func (s titleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type TitleExtractor struct {
}

func (e *TitleExtractor) Setup(_ interface{}) error {
	return nil
}

func (e *TitleExtractor) Extract(doc document.Document) (interface{}, error) {
	var currTitle title

	titles := titleSlice{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				currTitle = title{}

				switch n.Data {
				case "h1":
					currTitle.level = 1
				case "h2":
					currTitle.level = 2
				case "h3":
					currTitle.level = 3
				case "h4":
					currTitle.level = 4
				case "h5":
					currTitle.level = 5
				case "h6":
					currTitle.level = 6
				}
			}
		} else if n.Type == html.TextNode {
			currTitle.text += n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				titles = append(titles, currTitle)
			}
		}
	}

	f(doc.Body.Node())

	sort.Sort(titles)
	for _, t := range titles {
		if strings.Contains(doc.Title, t.text) {
			return t.text, nil
		}
	}

	if len(titles) > 1 {
		return titles[0], nil
	}

	return doc.Title, nil
}

func init() {
	RegisterPlugin("title", new(TitleExtractor))
}
