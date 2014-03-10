package text

import (
	"bytes"
	"fmt"
	"strings"
)

type TagType uint32

const (
	ElementTag TagType = iota
	TextTag
	StartTag
	EndTag
	NewLineTag
	SelfClosingTag
	RawTag
)

type Blocks []Block

func (blocks Blocks) String(html bool) string {
	buf := bytes.Buffer{}
	if html {
		for _, block := range blocks {
			switch block.TagType {
			case StartTag:
				buf.WriteString(fmt.Sprintf("<%s%s>\n", block.Tag, block.AttrString()))
			case EndTag:
				buf.WriteString(fmt.Sprintf("</%s>\n", block.Tag))
			case SelfClosingTag:
				buf.WriteString(fmt.Sprintf("<%s%s />", block.Tag, block.AttrString()))
			case NewLineTag:
				buf.WriteString("<br />")
			case ElementTag:
				buf.WriteString(fmt.Sprintf("<%s%s>%s</%s>\n", block.Tag, block.AttrString(), block.Data, block.Tag))
			case TextTag, RawTag:
				buf.WriteString(block.Data)
			}
		}
	} else {
		for _, block := range blocks {
			switch block.TagType {
			case StartTag:
				buf.WriteString("\n")
			case EndTag:
				buf.WriteString("\n")
			case NewLineTag:
				buf.WriteString("\n")
			case ElementTag:
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
					buf.WriteString(fmt.Sprintf("%s\n", block.Data))
				}
			case TextTag:
				buf.WriteString(block.Data)
			}
		}
	}

	return buf.String()
}

type Block struct {
	IsContent bool
	Tag       string
	TagType   TagType
	Attrs     map[string]string
	Data      string
}

func (b Block) AttrString() string {
	buf := bytes.Buffer{}
	for k, v := range b.Attrs {
		buf.WriteString(fmt.Sprintf(" %s=%s", k, v))
	}
	return buf.String()
}

type TextBlock struct {
	Type ContentType

	Tag  string
	Text string

	NumChars        int
	NumWords        int
	NumLinkedWords  int
	NumWrappedWords int
	NumLines        int

	TextDensity float64
	LineDensity float64
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

	if b.NumWords == 0 {
		b.TextDensity = 0
	} else {
		b.TextDensity = (float64(countDistinctWords(b.Text)) / float64(b.NumWords)) * 100
	}
	if b.NumWrappedWords == 0 {
		b.LineDensity = 0
	} else {
		b.LineDensity = (float64(b.NumLines) / float64(b.NumWrappedWords)) * 100
	}
	if b.NumWords == 0 {
		b.LinkDensity = 0
	} else {
		b.LinkDensity = 100 - ((float64(b.NumLinkedWords) / float64(b.NumWords)) * 100)
	}
}

// The algorithm used for classifying text is based on the boilerpipe library
// https://code.google.com/p/boilerpipe/
func (b *TextBlock) Classify() {
	// If the block is all mark as not content
	if b.LinkDensity == 0 {
		b.Type = NotContent
	} else {
		score := ((b.TextDensity * 0.5) + ((b.LinkDensity) * 0.3) + (b.LineDensity*0.2)*100)
		if score >= 0.45 {
			b.Type = Content
		} else {
			b.Type = NotContent
		}
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
