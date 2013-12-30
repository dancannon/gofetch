package selector

import (
	. "github.com/dancannon/gofetch/message"
	. "github.com/dancannon/gofetch/plugins"

	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type SelectorExtractor struct {
	selector  string
	attribute string
}

func (e *SelectorExtractor) Setup(config interface{}) error {
	params := config.(map[string]interface{})

	// Validate config
	if selector, ok := params["selector"]; !ok {
		return errors.New(fmt.Sprintf("The selector extractor must be passed a CSS selector"))
	} else {
		e.selector = selector.(string)
	}
	if attribute, ok := params["attribute"]; ok {
		e.attribute = attribute.(string)
	}

	return nil
}

func (e *SelectorExtractor) Extract(msg *ExtractMessage) error {
	doc := goquery.NewDocumentFromNode(msg.Document.Body.Node())

	n := doc.Find(e.selector)
	if n.Length() == 0 {
		return errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	if e.attribute == "" {
		msg.Value = n.First().Text()

		return nil
	} else {
		msg.Value, _ = n.First().Attr(e.attribute)
		return nil
	}
}

func init() {
	RegisterPlugin("selector", new(SelectorExtractor))
}
