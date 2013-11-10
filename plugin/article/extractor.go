package article

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	"regexp"
	"strings"
)

type Extractor struct {
	blocks       []string
	currentBlock string
	flush        bool
}

func (e *Extractor) Id() string {
	return "gofetch.article.extractor"
}

func (e *Extractor) Setup(_ []config.Value) error {
	e.blocks = []string{}
	e.currentBlock = ""
	e.flush = false

	return nil
}

func (e *Extractor) Extract(d *document.Document) (map[string]interface{}, error) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "option", "object", "embed", "applet", "link", "noscript":
				//Ignore
				return
			case "strike", "a", "u", "b", "i", "em", "strong", "span", "sup",
				"code", "tt", "sub", "var", "font":
				// Inline
				e.flush = false
			default:
				// Block
				e.flush = true
			}
		} else if n.Type == html.TextNode {
			if e.flush {
				e.flushBlock()
				e.flush = false
			}

			e.currentBlock = e.currentBlock + n.Data
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "strike", "a", "u", "b", "i", "em", "strong", "span", "sup",
				"code", "tt", "sub", "var", "font":
				// Inline
				e.flush = false
			default:
				// Block
				e.flush = true
			}
		}

		if e.flush {
			e.flushBlock()
		}
	}

	f(d.Body)

	// Clean blocks
	blocks := []string{}
	re := regexp.MustCompile("^[\t\n\f\r ]+$")
	re1 := regexp.MustCompile("^[\t\n\f\r ]+")
	re2 := regexp.MustCompile("[\t\n\f\r ]+$")

	for _, block := range e.blocks {
		if block != "" && !re.MatchString(block) {
			// Trim trailing whitespace
			block = re1.ReplaceAllString(block, "")
			block = re2.ReplaceAllString(block, "")
			blocks = append(blocks, block)
		}
	}

	return map[string]interface{}{
		"type": "text",
		"text": strings.Join(blocks, "\n"),
	}, nil
}

func (e *Extractor) flushBlock() {
	e.blocks = append(e.blocks, e.currentBlock)
	e.currentBlock = ""
}
