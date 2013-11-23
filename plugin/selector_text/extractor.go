package selector_text

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dancannon/gofetch/document"
	"github.com/dancannon/gofetch/plugin/text"
)

type Extractor struct {
	selector      string
	textextractor *text.Extractor
}

func (e *Extractor) Id() string {
	return "gofetch.selector_text.extractor"
}

func (e *Extractor) Setup(config map[string]interface{}) error {
	// Validate config
	if selector, ok := config["selector"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a CSS selector", e.Id()))
	} else {
		e.selector = selector.(string)
	}

	// Setup text extractor
	e.textextractor = &text.Extractor{}
	err := e.textextractor.Setup(config)
	if err != nil {
		return err
	}

	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
	doc := goquery.NewDocumentFromNode(d.Body)

	n := doc.Find(e.selector)
	if n.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	// Create a new document using the selected node
	d.Body = n.Get(0)

	// Run the new document through the text selector
	return e.textextractor.Extract(d)
}
