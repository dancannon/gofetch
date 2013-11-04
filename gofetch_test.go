package gofetch

import (
	"code.google.com/p/go.net/html"
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	// res, err := Fetch("http://getbootstrap.com/examples/starter-template/")
	// res, err := Fetch("http://getbootstrap.com/examples/jumbotron/")
	// res, err := Fetch("http://getbootstrap.com/examples/carousel/")
	// res, err := Fetch("http://www.theguardian.com/technology/2013/nov/01/caa-easa-electronic-devices-flight-take-off-landing")
	// res, err := Fetch("http://www.birmingham.ac.uk/index.aspx")
	// res, err := Fetch("http://www.birmingham.ac.uk/university/index.aspx")
	// res, err := Fetch("http://www.youtube.com/watch?v=C0DPdy98e4c")
	// res, err := Fetch("https://www.google.co.uk/?gws_rd=cr&ei=IMtzUuLkI-Hb0QX-woD4CA#q=test")
	res, err := Fetch("http://imgur.com")
	// res, err := Fetch("http://imgur.com/rXmjOMe")
	if err != nil {
		t.Errorf("Error was returned(%s)", err)
	}

	doc := prepareDocument(res)
	f, err := os.Create("test.html")
	defer f.Close()

	html.Render(f, doc.Body)

	// Execute all guesses and reduce probabilities to a single value for each
	// page type
	guesses := make(map[PageType]Guess)
	for _, guesser := range guessers {
		guess := guesser.Guess(doc)

		if g, ok := guesses[guess.Type]; ok {
			g.Probability = g.Probability * guess.Probability
		} else {
			guesses[guess.Type] = guess
		}
	}

	// Find the most likely guess
	var highest Guess
	for _, guess := range guesses {
		if guess.Probability >= highest.Probability {
			highest = guess
		}
	}

	spew.Printf("Guess: %s(%d)", highest.Type, highest.Probability)

	switch highest.Type {
	case Text:
		extractor := TextExtractor{}
		spew.Dump(extractor.Extract(doc))
	case Image, Gallery:
		spew.Dump(highest.Result)
	default:
		spew.Dump("Unknown page type")
	}
}
