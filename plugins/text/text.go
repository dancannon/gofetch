package text

import (
	htmlutil "html"
	"math"
	"regexp"

	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"

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
	blocks := e.parseNode(doc.Body.Node(), doc)

	// Remove non-content blocks
	blocks = e.getBestBlocks(blocks)

	return blocks.String(e.format == "raw"), nil
}

func (e *TextExtractor) parseNode(n *html.Node, d document.Document) Blocks {
	blocks := Blocks{}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "body":
			return e.parseChildNodes(nil, n, d)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			b := &Block{
				Tag:  n.Data,
				Type: ElementBlock,
			}
			b.Children.Add(&Block{
				Type:   TextBlock,
				Parent: b,
				Data:   e.parseChildNodes(nil, n, d).String(e.format == "raw"),
			})
			return blocks.Add(b)
		case "img":
			return blocks.Add(e.extractImage(n, d))
		case "ol", "ul":
			return blocks.Add(e.extractList(n, d)...)
		case "br":
			return blocks.Add(&Block{
				Type: NewLineBlock,
			})
		case "a":
			b := &Block{
				Tag:   n.Data,
				Type:  ElementBlock,
				Attrs: nodeAttrs(n, d, "href"),
			}
			b.Children = e.parseChildNodes(b, n, d)
			return blocks.Add(b)
		case "article", "aside", "blockquote", "dd", "div", "dl", "fieldset",
			"figcaption", "figure", "footer", "form", "header", "hgroup",
			"output", "p", "pre", "section", "table", "tbody", "tfoot", "tr", "th", "td":
			block := &Block{
				Tag:  n.Data,
				Type: ElementBlock,
			}
			block.Children = e.parseChildNodes(block, n, d)

			return blocks.Add(block)
		default:
			cs := e.parseChildNodes(nil, n, d).String(e.format == "raw")
			if cs != "" {
				return blocks.Add(&Block{
					Type: TextBlock,
					Data: cs,
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
		re = regexp.MustCompile("^[\t\n\f\r ]+$")
		if !re.MatchString(data) {
			return blocks.Add(&Block{
				Type: TextBlock,
				Data: data,
			})
		}
	}

	return blocks
}

func (e *TextExtractor) parseChildNodes(parent *Block, n *html.Node, d document.Document) Blocks {
	if n.Type == html.ElementNode {
		blocks := Blocks{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			for _, b := range e.parseNode(c, d) {
				b.Parent = parent
				blocks = append(blocks, b)
			}
		}
		return blocks
	} else {
		return Blocks{}
	}
}

func (e *TextExtractor) extractImage(n *html.Node, d document.Document) *Block {
	return &Block{
		Tag:   n.Data,
		Type:  SelfClosingBlock,
		Attrs: nodeAttrs(n, d, "src", "width", "height"),
	}
}

func (e *TextExtractor) extractList(n *html.Node, d document.Document) Blocks {
	blocks := Blocks{}

	currBlock := &Block{
		Tag:  n.Data,
		Type: ElementBlock,
	}

	// Collect list items
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if cs := e.parseNode(c, d).String(e.format == "raw"); cs != "" {
				currBlock.Children = append(currBlock.Children, &Block{
					Tag:  "li",
					Type: ElementBlock,
					Children: Blocks{
						&Block{
							Type: TextBlock,
							Data: cs,
						},
					},
				})
			}
		}
	}

	blocks = blocks.Add(currBlock)

	return blocks
}

// Algorithm based on the Goose library
// https://github.com/gravitylabs/goose/
func (e *TextExtractor) getBestBlocks(blocks Blocks) Blocks {
	var textBlocks = Blocks{}
	var parentBlocks = Blocks{}

	startingBoost := float64(0)
	i := 0

	// Get a list of blocks which contain text.
	var f func(*Block)
	f = func(block *Block) {
		if block.Type == ElementBlock {
			switch block.Tag {
			case "article", "aside", "blockquote", "dd", "div", "dl", "fieldset",
				"figcaption", "figure", "footer", "form", "header", "hgroup",
				"output", "p", "pre", "section":
				if block.TextStats.NumStopWords > 2 && block.TextStats.LinkDensity <= 1 {
					textBlocks = append(textBlocks, block)
				}

				for _, c := range block.Children {
					f(c)
				}
			}
		}
	}

	for _, block := range blocks {
		f(block)
	}

	bottomBlocks := float64(len(textBlocks)) * 0.25

	for _, block := range textBlocks {
		boostScore := float64(0)
		if isBlockBoostable(block) {
			if i >= 0 {
				boostScore = ((1.0 / startingBoost) * 50)
				startingBoost += 1
			}
		}
		if len(textBlocks) > 15 {
			if float64(len(textBlocks)-1) <= bottomBlocks {
				booster := bottomBlocks - float64(len(textBlocks)-i)
				boostScore = -math.Pow(booster, 2)
				if math.Abs(boostScore) > 40 {
					boostScore = 5
				}
			}

			upscore := block.TextStats.NumStopWords + int(boostScore)
			if block != nil {
				block.Score += upscore
				if !parentBlocks.Contains(block) {
					parentBlocks = append(parentBlocks, block)
				}

				if block.Parent != nil {
					block.Parent.Score += upscore / 2
					if !parentBlocks.Contains(block.Parent) {
						parentBlocks = append(parentBlocks, block.Parent)
					}

					if block.Parent.Parent != nil {
						block.Parent.Score += upscore / 4
						if !parentBlocks.Contains(block.Parent.Parent) {
							parentBlocks = append(parentBlocks, block.Parent.Parent)
						}
					}
				}

			}

			i++
		}
	}

	// Find the best block
	var bestBlock *Block

	for _, block := range parentBlocks {
		if bestBlock == nil {
			bestBlock = block
		}
		if block.Score >= bestBlock.Score {
			bestBlock = block
		}
	}

	if bestBlock == nil {
		return blocks
	} else {
		return Blocks{bestBlock}
	}
}

func isBlockBoostable(b *Block) bool {
	const MinStopWords int = 5
	const MaxSteps int = 3

	var steps int

	if b.Parent == nil {
		return false
	}

	for _, s := range b.Parent.Children {
		if s == b {
			return false
		}

		switch s.Tag {
		case "article", "aside", "blockquote", "dd", "div", "dl", "fieldset",
			"figcaption", "figure", "footer", "form", "header", "hgroup",
			"output", "p", "pre", "section":

			if steps >= MaxSteps {
				return false
			}
			if s.TextStats.NumStopWords > MinStopWords {
				return true
			}

			steps++
		}
	}

	return false
}

func init() {
	RegisterPlugin("text", new(TextExtractor))
}
