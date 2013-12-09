package gofetch

import (
	"strings"
)

type OpengraphExtractor struct {
}

func (e *OpengraphExtractor) Id() string {
	return "gofetch.opengraph.extractor"
}

func (e *OpengraphExtractor) Setup(config map[string]interface{}) error {
	return nil
}

func (e *OpengraphExtractor) Extract(d *Document, r *Result) (interface{}, error) {
	properties := map[string]interface{}{}

	for _, meta := range d.Meta {
		var property, content string

		for key, val := range meta {
			if key == "property" {
				property = val
			} else if key == "content" {
				content = val
			}
		}

		if property != "" && content != "" && strings.HasPrefix(property, "og:") {
			properties[property[3:]] = content
		}
	}

	r.PageType = "misc"

	return properties, nil
}
