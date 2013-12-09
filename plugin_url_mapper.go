package gofetch

import (
	"errors"
	"fmt"
	"regexp"
)

type UrlMapperExtractor struct {
	pattern     string
	replacement string
}

func (e *UrlMapperExtractor) Id() string {
	return "gofetch.url_mapper.extractor"
}

func (e *UrlMapperExtractor) Setup(config map[string]interface{}) error {
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

func (e *UrlMapperExtractor) Extract(d *Document, r *Result) (interface{}, error) {
	re, err := regexp.Compile(e.pattern)
	if err != nil {
		return nil, err
	}

	return re.ReplaceAllString(d.Url, e.replacement), nil
}
