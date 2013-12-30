package url_mapper

import (
	. "github.com/dancannon/gofetch/message"
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
	params := config.(map[string]interface{})

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

func (e *UrlMapperExtractor) Extract(msg *ExtractMessage) error {
	re, err := regexp.Compile(e.pattern)
	if err != nil {
		return err
	}

	msg.Value = re.ReplaceAllString(msg.Document.Url, e.replacement)

	return nil
}

func init() {
	RegisterPlugin("url_mapper", new(UrlMapperExtractor))
}
