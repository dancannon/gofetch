package selector

import (
	"errors"
	"fmt"
	"github.com/dancannon/gofetch/document"
	"regexp"
)

type Extractor struct {
	values
}

func (e *Extractor) Id() string {
	return "gofetch.url_mapper.extractor"
}

func (e *Extractor) Setup(values []Value) error {
	e.values = values

	return nil
}

func (e *Extractor) Extract(d *document.Document) (map[string]interface{}, error) {
	// re := regexp.Compile(e.pattern)
	// return re.ReplaceAllString(d.Url, e.replacement)
	return map[string]interface{}{}, nil
}
