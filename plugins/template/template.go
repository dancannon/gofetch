package selector

import (
	"bytes"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"html/template"

	"errors"
	"fmt"
)

type TemplateExtractor struct {
	template string
}

func (e *TemplateExtractor) Setup(config interface{}) error {
	var params map[string]interface{}
	if p, ok := config.(map[string]interface{}); !ok {
		params = make(map[string]interface{})
	} else {
		params = p
	}

	// Validate config
	if template, ok := params["template"]; !ok {
		return errors.New(fmt.Sprintf("The template extractor must be passed a template"))
	} else {
		e.template = template.(string)
	}

	return nil
}

func (e *TemplateExtractor) Extract(doc document.Document) (res interface{}, err error) {
	t, err := template.New("template").Parse(e.template)
	if err != nil {
		return "", err
	}

	var w = new(bytes.Buffer)
	err = t.Execute(w, doc)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}

func init() {
	RegisterPlugin("template", new(TemplateExtractor))
}
