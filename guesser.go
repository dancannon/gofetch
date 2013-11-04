package gofetch

import (
	"code.google.com/p/go.net/html"
	"errors"
	"regexp"
)

var guessers []Guesser

func init() {
	guessers = []Guesser{
		TextDensityGuesser{},
		ImageCountGuesser{},
	}
}

type Guess struct {
	Type        PageType
	Probability float64
	Result      interface{}
}

type Guesser interface {
	Guess(*Document) Guess
}

// Text Density Guesser
//
// Attempts to guess if a document is a text document based on the amount of
// text per html element.
type TextDensityGuesser struct{}

func (g TextDensityGuesser) Guess(d *Document) Guess {
	guess := Guess{}
	res := g.process(d.Body)
	if res.Length > 0 || res.Elements > 0 {
		guess.Probability = float64(res.Length/res.Elements) / 100
	} else {
		guess.Probability = 0
	}
	guess.Type = Unknown

	if guess.Probability > 0.33 {
		guess.Type = Text
	}

	return guess
}

type textDensityResult struct {
	Length   int
	Elements int
}

func (g *TextDensityGuesser) process(n *html.Node) textDensityResult {
	res := textDensityResult{}

	if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cres := g.process(c)

			if cres.Length == 0 {
				continue
			}

			res.Length += cres.Length
			res.Elements += cres.Elements
		}
	} else if n.Type == html.TextNode {
		data := n.Data
		// Remove any whitespace
		re := regexp.MustCompile("[\t\n\f\r]+|[ ]{2,}")
		data = re.ReplaceAllString(data, "")

		res.Length = len(data)
		res.Elements = 1
	}

	return res
}

// Text Density Guesser
//
// Attempts to guess if a document is a text document based on the amount of
// text per html element.
type ImageCountGuesser struct{}

func (g ImageCountGuesser) Guess(d *Document) Guess {
	images := g.getImages(d.Body)

	// Image elements
	guess := Guess{}
	guess.Result = images
	if len(images) == 1 {
		guess.Probability = 0.5
		guess.Type = Image
	} else if len(images) > 1 {
		guess.Probability = 0.33
		guess.Type = Gallery
	} else {
		guess.Probability = 0.05
		guess.Type = Image
	}

	return guess
}

func (g *ImageCountGuesser) getImages(n *html.Node) []string {
	images := []string{}

	if n.Type == html.ElementNode {
		if n.Data == "img" {
			src, err := g.getSrc(n)
			if err == nil {
				images = append(images, src)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			images = append(images, g.getImages(c)...)
		}
	}

	return images
}
func (g *ImageCountGuesser) getSrc(n *html.Node) (string, error) {
	for _, a := range n.Attr {
		if a.Key == "src" {
			return a.Val, nil
		}
	}

	return "", errors.New("No SRC attribute found")
}
