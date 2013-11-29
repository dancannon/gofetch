package opengraph

import (
	"github.com/dancannon/gofetch/document"
	"strings"
)

type Extractor struct {
}

func (e *Extractor) Id() string {
	return "gofetch.opengraph.extractor"
}

func (e *Extractor) Setup(config map[string]interface{}) error {
	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
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

	return properties, nil
}
