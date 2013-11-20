package article

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
	"regexp"
)

type Extractor struct {
}

func (e *Extractor) Id() string {
	return "gofetch.article.extractor"
}

func (e *Extractor) Setup(_ map[string]interface{}) error {
	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
	blocks := e.parseDocument(d)
	blocks = e.clasifyBlocks(blocks)

	// Get content
	content := ""
	for _, block := range blocks {
		if block.Type == Content {
			content += block.Text + "\n"
		}
	}

	return content, nil
}

func (e *Extractor) parseDocument(d *document.Document) []TextBlock {
	inLink := false
	blocks := []TextBlock{}
	currentBlock := TextBlock{}
	flush := false

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "option", "object", "embed", "applet", "link", "noscript":
				//Ignore
				return
			case "a":
				inLink = true
				fallthrough
			case "strike", "u", "b", "i", "em", "strong", "span", "sup",
				"code", "tt", "sub", "var", "font":
				// Inline
				flush = false
			default:
				// Block
				flush = true
			}
		} else if n.Type == html.TextNode {
			if flush {
				currentBlock.Flush()
				blocks = append(blocks, currentBlock)
				currentBlock = TextBlock{}
				flush = false
			}

			currentBlock.AddText(n.Data, inLink)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				inLink = false
				fallthrough
			case "strike", "u", "b", "i", "em", "strong", "span", "sup",
				"code", "tt", "sub", "var", "font":
				// Inline
				flush = false
			default:
				// Block
				flush = true
			}
		}

		if flush {
			currentBlock.Flush()
			blocks = append(blocks, currentBlock)
			currentBlock = TextBlock{}
		}
	}

	f(d.Body)

	// Clean blocks
	tmp := []TextBlock{}
	re := regexp.MustCompile("^[\t\n\f\r ]+$")
	re1 := regexp.MustCompile("^[\t\n\f\r ]+")
	re2 := regexp.MustCompile("[\t\n\f\r ]+$")

	for _, block := range blocks {
		if block.Text != "" && !re.MatchString(block.Text) {
			// Trim trailing whitespace
			block.Text = re1.ReplaceAllString(block.Text, "")
			block.Text = re2.ReplaceAllString(block.Text, "")
			tmp = append(tmp, block)
		}
	}
	blocks = tmp

	// Calculate block totals
	// for k, block := range e.blocks {
	// }

	return blocks
}

func (e *Extractor) clasifyBlocks(blocks []TextBlock) []TextBlock {
	for k, block := range blocks {
		blockType := Content

		// Get previous and next blocks
		var prev, next TextBlock
		if k == 0 {
			prev = TextBlock{}
		} else if k >= len(blocks)-1 {
			next = TextBlock{}
		} else {
			prev = blocks[k-1]
			next = blocks[k+1]
		}

		if block.LinkDensity <= 0.333333 {
			if prev.LinkDensity <= 0.555555 {
				if block.NumWords <= 16 {
					if next.NumWords <= 15 {
						if prev.NumWords <= 4 {
							blockType = NotContent
						} else {
							blockType = Content
						}
					} else {
						blockType = Content
					}
				} else {
					blockType = Content
				}
			} else {
				if block.NumWords <= 40 {
					if next.NumWords <= 17 {
						blockType = NotContent
					} else {
						blockType = Content
					}
				} else {
					blockType = Content
				}
			}
		} else {

			blockType = NotContent
		}

		blocks[k].Type = blockType
	}

	return blocks
}
