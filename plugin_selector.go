package gofetch

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type SelectorExtractor struct {
	selector  string
	attribute string
}

func (e *SelectorExtractor) Id() string {
	return "gofetch.selector.extractor"
}

func (e *SelectorExtractor) Setup(config map[string]interface{}) error {
	// Validate config
	if selector, ok := config["selector"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a CSS selector", e.Id()))
	} else {
		e.selector = selector.(string)
	}
	if attribute, ok := config["attribute"]; ok {
		e.attribute = attribute.(string)
	}

	return nil
}

func (e *SelectorExtractor) Extract(d *Document, r *Result) (interface{}, error) {
	doc := goquery.NewDocumentFromNode(d.Body)

	n := doc.Find(e.selector)
	if n.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	if e.attribute == "" {
		return n.First().Text(), nil
	} else {
		attr, _ := n.First().Attr(e.attribute)
		return attr, nil
	}
}
