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

func (e *OpengraphExtractor) Supports(doc document.Document) (bool, error) {
	// Look for a property starting with the string "og:"
	for _, meta := range doc.Meta {
		for key, val := range meta {
			if key == "property" {
				if strings.HasPrefix(val, "og:") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (e *OpengraphExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	props := map[string]interface{}{}

	// Load opengraph properties into a map
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
			props[property[3:]] = content
		}
	}

	var pagetype string

	// Not an official type but some sites use it (Flickr)
	if strings.Contains(props["type"].(string), "image") {
		pagetype = "image"
		props = map[string]interface{}{
			"title":  props["title"],
			"url":    props["image"],
			"width":  props["image:width"],
			"height": props["image:height"],
		}
	} else if strings.Contains(props["type"].(string), "video") {
		pagetype = "video"
		props = map[string]interface{}{
			"title":       props["title"],
			"description": props["description"],
		}
	} else if strings.Contains(props["type"].(string), "article") {
		pagetype = "text"
		props = map[string]interface{}{
			"title": props["title"],
			"text":  props["description"],
		}
	} else {
		pagetype = "unknown"
	}

	return props, pagetype, nil
}

func init() {
	RegisterPlugin("opengraph", new(OpengraphExtractor))
}
