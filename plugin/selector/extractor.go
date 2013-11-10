package selector

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
)

type Extractor struct {
	values []value
}

type value struct {
	name      string
	value     string
	attribute string
	selector  string
	children  []value
}

func (e *Extractor) Id() string {
	return "gofetch.selector.extractor"
}

func (e *Extractor) Setup(values []config.Value) error {
	var err error
	e.values, err = e.validateValueNodes(values)
	return err
}

func (e *Extractor) validateValueNodes(values []config.Value) ([]value, error) {
	var err error

	vs := []value{}
	for _, cv := range values {
		// Validate node
		v := value{}
		v.name = cv.Name
		v.value = cv.Value

		if cv.Name == "" {
			return vs, errors.New(fmt.Sprintf("Each value must have have at least a name"))
		}

		// If the node already has a value then skip validation
		if cv.Value == "" {
			if len(cv.Children) > 0 {
				// Validate children nodes
				v.children, err = e.validateValueNodes(cv.Children)
				if err != nil {
					return vs, err
				}
				vs = append(vs, v)
			} else {
				m := config.ParameterSlice(cv.Parameters).ToMap()
				if param, ok := m["attribute"]; ok {
					v.attribute = param
				}
				if param, ok := m["selector"]; ok {
					v.selector = param
				} else {
					return vs, errors.New(fmt.Sprintf("The %s extractor must be passed a CSS selector", e.Id()))
				}
			}
		}

		vs = append(vs, v)
	}

	return vs, nil
}

func (e *Extractor) Extract(d *document.Document) (map[string]interface{}, error) {
	doc := goquery.NewDocumentFromNode(d.Body)

	return e.extractValues(doc, e.values), nil
}

func (e *Extractor) extractValues(doc *goquery.Document, values []value) map[string]interface{} {
	m := map[string]interface{}{}

	for _, v := range values {
		if v.value == "" {
			if len(v.children) > 0 {
				m[v.name] = e.extractValues(doc, v.children)
			} else {
				n := doc.Find(v.selector)
				if n.Length() == 0 {
					continue
				}

				if v.attribute == "" {
					m[v.name] = n.First().Text()
				} else {
					attr, _ := n.First().Attr(v.attribute)

					m[v.name] = attr
				}
			}
		} else {
			m[v.name] = v.value
		}
	}

	return m
}
