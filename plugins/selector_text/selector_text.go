package selector_text

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/message"
	. "github.com/dancannon/gofetch/plugins"
	. "github.com/dancannon/gofetch/plugins/text"

	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type SelectorTextExtractor struct {
	selector      string
	textextractor *TextExtractor
}

func (e *SelectorTextExtractor) Setup(config interface{}) error {
	params := config.(map[string]interface{})

	// Validate config
	if selector, ok := params["selector"]; !ok {
		return errors.New(fmt.Sprintf("The selector extractor must be passed a CSS selector"))
	} else {
		e.selector = selector.(string)
	}

	// Setup text SelectorTextExtractor
	e.textextractor = &TextExtractor{}
	err := e.textextractor.Setup(config)
	if err != nil {
		return err
	}

	return nil
}

func (e *SelectorTextExtractor) Extract(msg *ExtractMessage) error {
	doc := goquery.NewDocumentFromNode(msg.Document.Body.Node())

	n := doc.Find(e.selector)
	if n.Length() == 0 {
		return errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	// Create a new document using the selected node
	msg.Document.Body = (*document.HtmlNode)(n.Get(0))

	// Run the new document through the text selector
	return e.textextractor.Extract(msg)
}

func init() {
	RegisterPlugin("selector_text", new(SelectorTextExtractor))
}
