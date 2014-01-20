package selector

import (
	"github.com/dancannon/gofetch/document"
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
	var params map[string]interface{}
	if p, ok := config.(map[string]interface{}); !ok {
		params = make(map[string]interface{})
	} else {
		params = p
	}

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

func (e *SelectorExtractor) Extract(doc document.Document) (interface{}, error) {
	qdoc := goquery.NewDocumentFromNode(doc.Body.Node())

	n := qdoc.Find(e.selector)
	if n.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	var value interface{}
	if e.attribute == "" {
		value = n.First().Text()
	} else {
		value, _ = n.First().Attr(e.attribute)
	}

	return value, nil
}

func init() {
	RegisterPlugin("selector", new(SelectorExtractor))
}
