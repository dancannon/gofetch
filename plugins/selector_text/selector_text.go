package selector_text

import (
	"bytes"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	. "github.com/dancannon/gofetch/plugins/text"
	"runtime"
	"strings"

	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type SelectorTextExtractor struct {
	selector      string
	restype       string
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
	if restype, ok := params["restype"]; ok {
		e.restype = restype.(string)
	}

	// Setup text SelectorTextExtractor
	e.textextractor = &TextExtractor{}
	err := e.textextractor.Setup(config)
	if err != nil {
		return err
	}

	return nil
}

func (e *SelectorTextExtractor) Extract(doc document.Document) (value interface{}, err error) {
	// GoQuery panics so we need to catch the errors
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if v, ok := r.(string); ok {
				err = errors.New(v)
			} else {
				err = r.(error)
			}
		}
	}()

	qdoc, err := goquery.NewDocumentFromReader(strings.NewReader(doc.Raw))
	if err != nil {
		return nil, err
	}

	n := qdoc.Find(e.selector)
	if n.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	switch e.restype {
	case "first":
		doc.Body = (*document.HtmlNode)(n.Get(0))
		value, err = e.textextractor.Extract(doc)
		if err != nil {
			return
		}
	case "all":
		s := []interface{}{}
		var v interface{}

		for i := 0; i < n.Length(); i++ {
			doc.Body = (*document.HtmlNode)(n.Get(i))
			v, err = e.textextractor.Extract(doc)
			if err != nil {
				return
			}
			s = append(s, v)
		}
		value = s
	case "merge":
		fallthrough
	default:
		buf := bytes.Buffer{}

		var v interface{}
		for i, _ := range n.Nodes {
			doc.Body = (*document.HtmlNode)(n.Nodes[i])
			v, err = e.textextractor.Extract(doc)
			if err != nil {
				return
			}

			buf.WriteString(v.(string))
			if i < len(n.Nodes)-1 {
				// buf.WriteString("\n")
			}

		}

		value = buf.String()
	}
	return
}

func init() {
	RegisterPlugin("selector_text", new(SelectorTextExtractor))
}
