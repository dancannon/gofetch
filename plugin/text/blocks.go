package text

import (
	"strings"
)

type ContentType int

const (
	lineLength = 80

	Content ContentType = iota
	Title
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
