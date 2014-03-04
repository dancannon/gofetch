package text

import (
	"github.com/davecgh/go-spew/spew"
	"strings"
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
		spew.Dump(b)
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
