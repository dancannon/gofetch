package selector

import (
	"errors"
	"fmt"
	"github.com/dancannon/gofetch/document"
	"regexp"
)

type Extractor struct {
	pattern     string
	replacement string
}

func (e *Extractor) Id() string {
	return "gofetch.url_mapper.extractor"
}

func (e *Extractor) Setup(values []Value) error {
	// Validate config
	if pattern, ok := config["pattern"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a source url", e.Id()))
	} else {
		e.pattern = pattern.(string)
	}

	if replacement, ok := config["replacement"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a source url", e.Id()))
	} else {
		e.replacement = replacement.(string)
	}

	return nil
}

func (e *Extractor) Extract(d *document.Document) (map[string]interface{}, error) {
	re := regexp.Compile(e.pattern)
	return re.ReplaceAllString(d.Url, e.replacement)
}
