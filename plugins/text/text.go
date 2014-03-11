package text

import (
	"bytes"
	"fmt"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	htmlutil "html"
	"regexp"

	"code.google.com/p/go.net/html"
)

type ContentType int

const (
	lineLength = 80

	Content ContentType = iota
	Title
	Tag_Start
	Tag_End
	Tag
	NotContent
)

type TextExtractor struct {
	format string
}

func (e *TextExtractor) Setup(config interface{}) error {
	params := config.(map[string]interface{})

	// Validate config
	if format, ok := params["format"]; !ok {
		e.format = "raw"
	} else {
		e.format = format.(string)
	}

	return nil
}

func (e *TextExtractor) Extract(doc document.Document) (interface{}, error) {
	var tmp Blocks
	blocks := e.parseNode(doc.Body.Node())

	// Attempt to remove non-content blocks
	// Based on the Goose library - https://github.com/GravityLabs/goose
	for i := 0; i < len(blocks); i++ {
		var tb TextBlock
		var bb Blocks

		if blocks[i].TagType == StartTag && blocks[i].EndBlock != nil {
			// Create text block
			lastJ := 0
			for j := i; j < len(blocks) && blocks[i].EndBlock != blocks[j]; j++ {
				tb.AddText(blocks[j].Data, blocks[j].Tag == "a")
				bb = append(bb, blocks[j])
				lastJ = j
			}

			i = lastJ
		} else {
			tb.AddText(blocks[i].Data, blocks[i].Tag == "a")
			bb = append(bb, blocks[i])
		}

		tb.Flush()
		if stopWordCount(tb.Data) > 2 && tb.LinkDensity < 1 {
			tmp = append(tmp, bb...)
		}
	}
	// for _, block := range blocks {

	// }

	// spew.Dump(blocks.String(false))
	fmt.Println(blocks.String(true))

	// return blocks.String(e.format == "raw"), nil
	return "", nil
}

func (e *TextExtractor) parseNode(n *html.Node) Blocks {
	blocks := Blocks{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		isLink := false
		addEndTag := false

		if n.Type == html.ElementNode {
			switch n.Data {
			case "body":
			case "h1", "h2", "h3", "h4", "h5", "h6":
				blocks = append(blocks, &Block{
					Tag:     n.Data,
					TagType: ElementTag,
					Data:    e.parseChildNodes(n).String(false),
				})
				return
			case "img":
				blocks = append(blocks, e.extractImage(n))
				return
			case "ol", "ul":
				blocks = append(blocks, e.extractList(n)...)
				return
			case "br":
				blocks = append(blocks, &Block{
					TagType: NewLineTag,
				})
			case "table", "tbody", "tfoot", "tr", "th", "td":
				buf := &bytes.Buffer{}
				html.Render(buf, n)

				blocks = append(blocks, &Block{
					TagType: RawTag,
					Data:    buf.String(),
				})
				return
			default:
				data := e.parseChildNodes(n).String(false)
				if data != "" {
					blocks = append(blocks, &Block{
						TagType: TextTag,
						Data:    data,
					})
				}
				return
			case "a":
				isLink = true
				addEndTag = true
				blocks = blocks.AddStartBlock(&Block{
					Tag:     n.Data,
					TagType: StartTag,
					Attrs:   nodeAttrs(n, "href"),
				})
			// Block elements
			case "article", "aside", "blockquote", "dd", "div", "dl", "fieldset",
				"figcaption", "figure", "footer", "form", "header", "hgroup",
				"output", "p", "pre", "section":
				// Attempt to check if the element only contains one non-empty child
				count := 0
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Data != "" {
						count++
					}
				}

				if count > 0 {
					addEndTag = true
					blocks = blocks.AddStartBlock(&Block{
						Tag:     n.Data,
						TagType: StartTag,
					})
				}

			}
		} else if n.Type == html.TextNode {
			var re *regexp.Regexp

			data := n.Data
			re = regexp.MustCompile("[\t\r\n]+")
			data = re.ReplaceAllString(data, "")
			re = regexp.MustCompile(" {2,}")
			data = re.ReplaceAllString(data, " ")
			data = htmlutil.EscapeString(data)

			if data != "" {
				blocks = append(blocks, &Block{
					TagType: TextTag,
					Data:    data,
				})
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if isLink {
			isLink = false
		}
		if addEndTag {
			blocks = blocks.AddEndBlock(&Block{
				Tag:     n.Data,
				TagType: EndTag,
			})
		}
	}

	f(n)

	return blocks
}

func (e *TextExtractor) parseChildNodes(n *html.Node) Blocks {
	if n.Type == html.ElementNode {
		blocks := Blocks{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			blocks = append(blocks, e.parseNode(c)...)
		}
		return blocks
	} else {
		return Blocks{}
	}
}

func (e *TextExtractor) extractImage(n *html.Node) *Block {
	attrs := map[string]string{}

	// Collect attributes
	for _, a := range n.Attr {
		switch a.Key {
		case "src", "width", "height":
			attrs[a.Key] = a.Val
		}
	}

	return &Block{
		Tag:     n.Data,
		TagType: SelfClosingTag,
		Attrs:   nodeAttrs(n, "src", "width", "height"),
	}
}

func (e *TextExtractor) extractList(n *html.Node) Blocks {
	blocks := Blocks{}

	blocks = blocks.AddStartBlock(&Block{
		Tag:     n.Data,
		TagType: StartTag,
	})

	// Collect list items
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			blocks = append(blocks, &Block{
				Tag:     "li",
				TagType: ElementTag,
				Data:    e.parseNode(c).String(false),
			})
		}
	}

	blocks = blocks.AddEndBlock(&Block{
		Tag:     n.Data,
		TagType: EndTag,
	})

	return blocks
}

func init() {
	RegisterPlugin("text", new(TextExtractor))
}
