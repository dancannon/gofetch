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

func (e *Extractor) Setup(config map[string]string) error {
	if param, ok := config["pattern"]; ok {
		e.pattern = param
	} else {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a regular expression pattern", e.Id()))
	}
	if param, ok := config["replacement"]; ok {
		e.replacement = param
	} else {
		return errors.New(fmt.Sprintf("The %s extractor must be passed a replacement string", e.Id()))
	}

	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
	re := regexp.Compile(e.pattern)
	return re.ReplaceAllString(d.Url, e.replacement)
}
