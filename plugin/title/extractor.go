package title

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
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

type Extractor struct {
}

func (e *Extractor) Id() string {
	return "gofetch.title.extractor"
}

func (e *Extractor) Setup(config map[string]interface{}) error {
	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
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

	f(d.Body)

	sort.Sort(titles)
	for _, t := range titles {
		if strings.Contains(d.Title, t.text) {
			return t.text, nil
		}
	}

	return d.Title, nil
}
