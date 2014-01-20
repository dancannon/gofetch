package text

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"github.com/davecgh/go-spew/spew"

	"regexp"
	"strings"
	"code.google.com/p/go.net/html"
)

type ContentType int

const (
	lineLength = 80

	Content ContentType = iota
	Title
	Tag_Start
	Tag_End
	NotContent
)

type TextBlock struct {
	Type ContentType
	Text string
	Tag  string

	NumChars        int
	NumWords        int
	NumLinkedWords  int
	NumWrappedWords int
	NumLines        int

	TextDensity float64
	LinkDensity float64
}

func (b *TextBlock) AddText(text string, inLink bool) {
	words := strings.Fields(text)

	// Increment counts
	b.NumWords += len(words)
	b.NumChars += len(text)
	if inLink {
		b.NumLinkedWords += len(words)
	}

	b.Text = b.Text + text
}

func (b *TextBlock) Flush() {
	// Count the number of lines
	words := strings.Fields(b.Text)
	currLineChars := 0
	currLineWords := 0

	b.NumLines = 0
	b.NumWrappedWords = 0

	for _, word := range words {
		currLineChars += len(word)
		currLineWords += 1

		if currLineChars > lineLength {
			b.NumLines++
			b.NumWrappedWords = 0

			currLineChars = 0
			currLineWords = 0
		}
	}

	if b.NumLines == 0 {
		b.NumWrappedWords = b.NumWords
		b.NumLines = 1
	} else {
		b.NumWrappedWords = b.NumWords - currLineWords
	}

	b.TextDensity = float64(b.NumWrappedWords) / float64(b.NumLines)
	if b.NumWords == 0 {
		b.LinkDensity = 0
	} else {
		b.LinkDensity = float64(b.NumLinkedWords) / float64(b.NumWords)
	}
}

type TextExtractor struct {
}

func (e *TextExtractor) Setup(_ interface{}) error {
	return nil
}

func (e *TextExtractor) Extract(doc document.Document) (interface{}, error) {
	blocks := e.parseDocument(doc)
	blocks = e.clasifyBlocks(blocks)

	content := ""
	hasContent := false
	spew.Dump(blocks)
	for _, block := range blocks {
		if block.Type == Tag_Start {
			content += "<" + block.Tag + ">"
		} else if block.Type == Tag_End {
			content += "</" + block.Tag + ">"
		} else if block.Type == Content {
			hasContent = true
			content += "<" + block.Tag + ">" + block.Text + "</" + block.Tag + ">"
		}
	}

	if !hasContent {
		return "", nil
	}

	return content, nil
}

func (e *TextExtractor) parseDocument(d document.Document) []TextBlock {
	blocks := []TextBlock{}
	currentBlock := TextBlock{}
	flush := false
	inLink := false
	currTag := ""

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if currTag == "" {
				currTag = n.Data
			}

			switch n.Data {
			case "option", "object", "embed", "applet", "link", "noscript":
				//Ignore
				return
			case "ul", "ol":
				currentBlock.Tag = currTag
				currentBlock.Flush()
				blocks = append(blocks, currentBlock)

				currentBlock = TextBlock{}
				currentBlock.Tag = n.Data
				currentBlock.Type = Tag_End
				currentBlock.Flush()
				blocks = append(blocks, currentBlock)

				currentBlock = TextBlock{}
				currTag = ""
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
				currentBlock.Tag = currTag
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
			case "ul", "ol":
				currentBlock.Tag = currTag
				currentBlock.Flush()
				blocks = append(blocks, currentBlock)

				currentBlock = TextBlock{}
				currentBlock.Tag = n.Data
				currentBlock.Type = Tag_End
				currentBlock.Flush()
				blocks = append(blocks, currentBlock)

				currentBlock = TextBlock{}
				currTag = ""
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
			currentBlock.Tag = currTag
			currentBlock.Flush()

			blocks = append(blocks, currentBlock)
			currentBlock = TextBlock{}
			currTag = ""
		}
	}

	f(d.Body.Node())

	// Clean blocks
	tmp := []TextBlock{}
	re := regexp.MustCompile("^[\t\n\f\r ]+$")
	re1 := regexp.MustCompile("^[\t\n\f\r ]+")
	re2 := regexp.MustCompile("[\t\n\f\r ]+$")

	for _, block := range blocks {
		if block.Type != 0 {
			tmp = append(tmp, block)
		} else if block.Text != "" && !re.MatchString(block.Text) {
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

func (e *TextExtractor) clasifyBlocks(blocks []TextBlock) []TextBlock {
	for k, block := range blocks {
		if block.Type == 0 {
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
					if block.TextDensity <= 10 {
						blocks[k].Type = NotContent
					} else {
						blocks[k].Type = Content
					}
				} else {
					if next.TextDensity <= 10 {
						blocks[k].Type = NotContent
					} else {
						blocks[k].Type = Content
					}
				}
			} else {
				blocks[k].Type = NotContent
			}
		}
	}

	return blocks
}

func init() {
	RegisterPlugin("text", new(TextExtractor))
}
