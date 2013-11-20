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

func (e *Extractor) Setup(config map[string]interface{}) error {
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

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
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

// func (e *Extractor) extractValues(doc *goquery.Document, values []value) map[string]interface{} {
// 	m := map[string]interface{}{}

// 	for _, v := range values {
// 		if v.value == "" {
// 			if len(v.children) > 0 {
// 				m[v.name] = e.extractValues(doc, v.children)
// 			} else {
// 				n := doc.Find(v.selector)
// 				if n.Length() == 0 {
// 					continue
// 				}

// 				if v.attribute == "" {
// 					m[v.name] = n.First().Text()
// 				} else {
// 					attr, _ := n.First().Attr(v.attribute)

// 					m[v.name] = attr
// 				}
// 			}
// 		} else {
// 			m[v.name] = v.value
// 		}
// 	}

// 	return m
// }
