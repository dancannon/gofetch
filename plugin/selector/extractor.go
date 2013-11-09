package selector

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dancannon/gofetch/document"
)

type Extractor struct {
	selector  string
	attribute string
}

func (e *Extractor) Id() string {
	return "gofetch.selector.extractor"
}

func (e *Extractor) Setup(config map[string]string) error {
	if param, ok := config["attribute"]; ok {
		e.attribute = param
	}
	if param, ok := config["selector"]; ok {
		e.selector = param
	} else {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a CSS selector", e.Id()))
	}

	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
	doc := goquery.NewDocumentFromNode(d.Body)

	n := doc.Find(e.selector)
	if n.Length() == 0 {
		return nil, nil
	}

	if e.attribute == "" {
		return n.First().Text(), nil
	} else {
		attr, _ := n.First().Attr(e.attribute)

		return attr, nil
	}
}
