package url_mapper

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"

	"errors"
	"fmt"
	"regexp"
)

type UrlMapperExtractor struct {
	pattern     string
	replacement string
}

func (e *UrlMapperExtractor) Setup(config interface{}) error {
	var params map[string]interface{}
	if p, ok := config.(map[string]interface{}); !ok {
		params = make(map[string]interface{})
	} else {
		params = p
	}

	// Validate config
	if pattern, ok := params["pattern"]; !ok {
		return errors.New(fmt.Sprintf("The url mapper extractor must be passed a source url"))
	} else {
		e.pattern = pattern.(string)
	}

	if replacement, ok := params["replacement"]; !ok {
		return errors.New(fmt.Sprintf("The url mapper extractor must be passed a source url"))
	} else {
		e.replacement = replacement.(string)
	}

	return nil
}

func (e *UrlMapperExtractor) Extract(doc document.Document) (interface{}, error) {
	re, err := regexp.Compile(e.pattern)
	if err != nil {
		return nil, err
	}

	return re.ReplaceAllString(doc.URL.String(), e.replacement), nil
}

func init() {
	RegisterPlugin("url_mapper", new(UrlMapperExtractor))
}
