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
	if strings.Contains(props["type"].(string), "image") {
		pagetype = "image"
		values = createMapFromProps(props, map[string]string{
			"title":  "title",
			"url":    "image",
			"width":  "image:width",
			"height": "image:height",
		})
	} else if strings.Contains(props["type"].(string), "video") {
		pagetype = "video"
		values = createMapFromProps(props, map[string]string{
			"title":       "title",
			"description": "description",
		})
	} else {
		pagetype = "text"
		values = createMapFromProps(props, map[string]string{
			"title": "title",
			"text":  "description",
		})
	}

	return values, pagetype, nil
}

func createMapFromProps(props map[string]interface{}, keys map[string]string) map[string]interface{} {
	m := make(map[string]interface{})

	for mapKey, propKey := range keys {
		if val, ok := props[propKey]; ok {
			m[mapKey] = val
		}
	}

	return m
}

func init() {
	RegisterPlugin("opengraph", new(OpengraphExtractor))
}
