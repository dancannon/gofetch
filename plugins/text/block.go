package text

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/dancannon/gofetch/document"

	"code.google.com/p/go.net/html"
)

type BlockType uint32

const (
	ElementBlock BlockType = iota
	TextBlock
	NewLineBlock
	SelfClosingBlock
	RawBlock
)

type Blocks []*Block

var startBlock *Block

func (b Blocks) Add(blocks ...*Block) Blocks {
	for _, block := range blocks {
		block.updateTextStats()
	}
	return append(b, blocks...)
}

func (blocks Blocks) String(html bool) string {
	buf := bytes.Buffer{}
	if html {
		buf.WriteString(blocks.htmlString(""))
	} else {
		for _, block := range blocks {
			switch block.Type {
			case NewLineBlock:
				buf.WriteString("\n")
			case ElementBlock:
				switch block.Tag {
				case "li":
					buf.WriteString(fmt.Sprintf("  - %s\n", block.Data))
				case "h1":
					buf.WriteString(fmt.Sprintf("#%s\n", block.Data))
				case "h2":
					buf.WriteString(fmt.Sprintf("##%s\n", block.Data))
				case "h3":
					buf.WriteString(fmt.Sprintf("###%s\n", block.Data))
				case "h4":
					buf.WriteString(fmt.Sprintf("####%s\n", block.Data))
				case "h5":
					buf.WriteString(fmt.Sprintf("#####%s\n", block.Data))
				case "h6":
					buf.WriteString(fmt.Sprintf("######%s\n", block.Data))
				default:
					if cs := block.Children.String(html); cs != "" {
						buf.WriteString(fmt.Sprintf("%s\n", cs))
					}
				}
			case TextBlock:
				buf.WriteString(block.Data + " ")
			}
		}
	}

	return buf.String()
}

func (blocks Blocks) htmlString(prefix string) string {
	buf := bytes.Buffer{}

	for _, block := range blocks {
		switch block.Type {
		case SelfClosingBlock:
			buf.WriteString(fmt.Sprintf("<%s%s />", block.Tag, block.AttrString()))
		case NewLineBlock:
			buf.WriteString("<br />")
		case ElementBlock:
			if cs := block.Children.htmlString(prefix + "\t"); cs != "" {
				buf.WriteString(fmt.Sprintf("<%s%s>%s</%s>", block.Tag, block.AttrString(), cs, block.Tag))
			}
		case TextBlock, RawBlock:
			buf.WriteString(block.Data + " ")
		}
	}

	return buf.String()
}

func (blocks Blocks) Contains(block *Block) bool {
	for _, b := range blocks {
		if b == block {
			return true
		}
	}

	return false
}

type Block struct {
	TextStats BlockTextStats
	Score     int

	IsContent bool

	Tag   string
	Type  BlockType
	Attrs map[string]string
	Data  string

	Parent   *Block
	Children Blocks
}

func (b Block) AttrString() string {
	buf := bytes.Buffer{}
	for k, v := range b.Attrs {
		buf.WriteString(fmt.Sprintf(" %s=%s", k, v))
	}
	return buf.String()
}

func (b *Block) updateTextStats() {
	var f func(*Block, bool)
	f = func(block *Block, inLink bool) {
		if block.Type == ElementBlock {
			if inLink {
				for _, c := range block.Children {
					f(c, true)
				}
			} else {
				for _, c := range block.Children {
					f(c, block.Tag == "a")
				}
			}
		} else if block.Type == TextBlock {
			b.TextStats.AddText(block.Data, inLink)
		}
	}

	f(b, false)

	b.TextStats.Flush()
}

type BlockTextStats struct {
	Text string

	NumWords       int
	NumLinks       int
	NumLinkedWords int
	NumStopWords   int

	LinkDensity float64
}

func (b *BlockTextStats) AddText(text string, inLink bool) {
	words := strings.Fields(text)

	// Increment counts
	b.NumWords += len(words)
	if inLink {
		b.NumLinks++
		b.NumLinkedWords += len(words)
	}

	b.Text = b.Text + text
}

func (b *BlockTextStats) Flush() {
	b.NumStopWords = stopWordCount(b.Text)
	if b.NumWords == 0 {
		b.LinkDensity = 0
	} else {
		b.LinkDensity = (float64(b.NumLinkedWords) / float64(b.NumWords))
	}
}

func countDistinctWords(text string) int {
	m := map[string]struct{}{}
	count := 0
	for _, word := range strings.Fields(text) {
		if _, ok := m[word]; !ok {
			m[word] = struct{}{}
			count++
		}
	}

	return count
}

func nodeAttrs(n *html.Node, d document.Document, keys ...string) map[string]string {
	attrs := map[string]string{}

	// Collect attributes
	for _, a := range n.Attr {
		if len(keys) == 0 || containsKey(keys, a.Key) {
			switch a.Key {
			case "src":
				// Attempt to fix URLs
				urlr, err := url.Parse(a.Val)
				if err != nil {
					continue
				}
				attrs[a.Key] = d.URL.ResolveReference(urlr).String()
			default:
				attrs[a.Key] = a.Val
			}
		}
	}

	return attrs
}

func containsKey(s []string, k string) bool {
	for _, e := range s {
		if e == k {
			return true
		}
	}
	return false
}
