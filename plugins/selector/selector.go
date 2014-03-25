package selector

import (
	"bytes"
	"runtime"
	"strings"

	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"github.com/dancannon/gofetch/util"
	"github.com/davecgh/go-spew/spew"

	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type SelectorExtractor struct {
	selector  string
	attribute string
	restype   string
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
	if restype, ok := params["restype"]; ok {
		e.restype = restype.(string)
	}

	return nil
}

func (e *SelectorExtractor) Extract(doc document.Document) (res interface{}, err error) {
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
	spew.Config.MaxDepth = 2
	spew.Dump(e.selector)
	if n.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("Selector '%s' not found", e.selector))
	}

	var value interface{}
	if e.attribute == "" {
		switch e.restype {
		case "first":
			value = util.SelectionToString(n.First())
		case "all":
			value = n.Map(func(index int, n *goquery.Selection) string {
				return util.SelectionToString(n)
			})
		case "html":
			w := &bytes.Buffer{}
			html.Render(w, n.Nodes[0])
			value = w.String()
		case "merge":
			fallthrough
		default:
			value = util.SelectionToString(n)
		}
	} else {
		switch e.restype {
		case "first":
			value, _ = n.First().Attr(e.attribute)
		case "all", "merge":
			fallthrough
		default:
			value = n.Map(func(index int, n *goquery.Selection) string {
				res, _ := n.Attr(e.attribute)
				return res
			})
		}
		value, _ = n.First().Attr(e.attribute)
	}

	return value, nil
}

func init() {
	RegisterPlugin("selector", new(SelectorExtractor))
}
