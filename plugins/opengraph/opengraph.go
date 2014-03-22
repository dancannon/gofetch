package opengraph

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"github.com/dancannon/gofetch/util"

	"strings"
)

type OpengraphExtractor struct {
}

func (e *OpengraphExtractor) Setup(_ interface{}) error {
	return nil
}

func (e *OpengraphExtractor) Supports(doc document.Document) bool {
	// Look for a property starting with the string "og:"
	for _, meta := range doc.Meta {
		for key, val := range meta {
			if key == "property" || key == "name" {
				if strings.HasPrefix(val, "og:") {
					return true
				}
			}
		}
	}

	return false
}

func (e *OpengraphExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	props := map[string]interface{}{}

	// Load opengraph properties into a map
	for _, meta := range doc.Meta {
		var property, content string

		for key, val := range meta {
			if key == "property" || key == "name" {
				property = val
			} else if key == "content" {
				content = val
			}
		}

		if property != "" && content != "" && strings.HasPrefix(property, "og:") {
			props[property[3:]] = content
		}
	}

	var pagetype string
	if _, ok := props["type"]; !ok {
		return props, "unknown", nil
	}

	var values map[string]interface{}

	// Not an official type but some sites use it (Flickr)
	if strings.Contains(props["type"].(string), "image") ||
		strings.Contains(props["type"].(string), "photo") {
		pagetype = "image"
		values = util.CreateMapFromProps(props, map[string]string{
			"title":  "title",
			"url":    "image",
			"width":  "image:width",
			"height": "image:height",
		})
	} else {
		pagetype = "general"
		values = util.CreateMapFromProps(props, map[string]string{
			"title":   "title",
			"content": "description",
		})
	}

	return values, pagetype, nil
}

func init() {
	RegisterPlugin("opengraph", new(OpengraphExtractor))
}
