package text

import (
	"bytes"
	"fmt"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"github.com/davecgh/go-spew/spew"
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
	blocks := e.parseNode(doc.Body.Node())
	// e.parseNode(doc.Body.Node())

	// spew.Dump(blocks.String(false))
	fmt.Println(blocks.String(true))

	content := ""
	return content, nil
}

func (e *TextExtractor) parseNode(n *html.Node) Blocks {
	blocks := Blocks{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		addEndTag := false

		if n.Type == html.ElementNode {
			fmt.Println(n.Data)
			switch n.Data {
			case "body":
			case "h1", "h2", "h3", "h4", "h5", "h6":
				blocks = append(blocks, Block{
					Tag:     n.Data,
					TagType: ElementTag,
					Data:    e.getNodeText(n).String(false),
				})
				return
			case "img":
				blocks = append(blocks, e.extractImage(n))
				return
			case "ol", "ul":
				blocks = append(blocks, e.extractList(n)...)
				return
			case "br":
				blocks = append(blocks, Block{
					TagType: NewLineTag,
				})
			case "table", "tbody", "tfoot", "tr", "th", "td":
				buf := &bytes.Buffer{}
				html.Render(buf, n)

				blocks = append(blocks, Block{
					TagType: RawTag,
					Data:    buf.String(),
				})
				return
			case "a":
			default:
				data := e.getNodeText(n).String(false)
				if data != "" {
					blocks = append(blocks, Block{
						TagType: TextTag,
						Data:    data,
					})
				}
				return

			// Block elements
			case "article", "aside", "blockquote", "dd", "div", "dl", "fieldset",
				"figcaption", "figure", "footer", "form", "header", "hgroup",
				"output", "p", "pre", "section":
				// Attempt to check if the element only contains one non-empty child
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					spew.Dump(c)
				}

				addEndTag = true
				blocks = append(blocks, Block{
					Tag:     n.Data,
					TagType: StartTag,
				})
			}
		} else {
			data := e.getNodeText(n).String(false)
			if data != "" {
				blocks = append(blocks, Block{
					TagType: TextTag,
					Data:    data,
				})
			}
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if addEndTag {
			blocks = append(blocks, Block{
				Tag:     "p",
				TagType: EndTag,
			})
		}
	}

	f(n)

	return blocks
}

func (e *TextExtractor) extractImage(n *html.Node) Block {
	attrs := map[string]string{}

	// Collect attributes
	for _, a := range n.Attr {
		switch a.Key {
		case "src", "width", "height":
			attrs[a.Key] = a.Val
		}
	}

	return Block{
		Tag:     n.Data,
		TagType: SelfClosingTag,
		Attrs:   attrs,
	}
}

func (e *TextExtractor) extractList(n *html.Node) Blocks {
	blocks := Blocks{}

	blocks = append(blocks, Block{
		Tag:     n.Data,
		TagType: StartTag,
	})

	// Collect list items
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			blocks = append(blocks, Block{
				Tag:     "li",
				TagType: ElementTag,
				Data:    e.getNodeText(c).String(false),
			})
		}
	}

	blocks = append(blocks, Block{
		Tag:     n.Data,
		TagType: EndTag,
	})

	return blocks
}

func (e *TextExtractor) getNodeText(n *html.Node) Blocks {
	if n.Type == html.TextNode {
		var re *regexp.Regexp

		data := n.Data
		re = regexp.MustCompile("[\t\r\n]+")
		data = re.ReplaceAllString(data, "")
		re = regexp.MustCompile(" {2,}")
		data = re.ReplaceAllString(data, " ")
		data = htmlutil.EscapeString(data)

		if data != "" {
			return Blocks{
				Block{
					TagType: TextTag,
					Data:    data,
				},
			}
		} else {
			return Blocks{}
		}
	} else if n.Type == html.ElementNode {
		switch n.Data {
		case "article", "aside", "blockquote", "div", "fieldset", "p", "pre", "td", "section":
			return e.parseNode(n)
		default:
			blocks := Blocks{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				blocks = append(blocks, e.parseNode(c)...)
			}
			return blocks
		}
	} else {
		return Blocks{}
	}
}

func (e *TextExtractor) clasifyBlocks(blocks []TextBlock) []TextBlock {
	for k, block := range blocks {
		bp := &block
		bp.Classify()
		blocks[k] = *bp
	}

	return blocks
}

func init() {
	RegisterPlugin("text", new(TextExtractor))
}
