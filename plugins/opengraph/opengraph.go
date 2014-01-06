package opengraph

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"

	"strings"
)

type OpengraphExtractor struct {
}

func (e *OpengraphExtractor) Setup(_ interface{}) error {
	return nil
}

func (e *OpengraphExtractor) Extract(doc document.Document) (interface{}, error) {
	properties := map[string]interface{}{}

	for _, meta := range doc.Meta {
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

func init() {
	RegisterPlugin("opengraph", new(OpengraphExtractor))
}
