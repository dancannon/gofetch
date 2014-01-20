package oembed

import (
	"code.google.com/p/go.net/html"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"

	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type OEmbedExtractor struct {
	endpoint string
	format   string
}

func (e *OEmbedExtractor) Setup(config interface{}) error {
	params := config.(map[string]interface{})

	// Validate config
	if endpoint, ok := params["endpoint"]; !ok {
		return errors.New(fmt.Sprintf("The oembed extractor must be passed an endpoint"))
	} else {
		e.endpoint = endpoint.(string)
	}

	if format, ok := params["format"]; ok {
		e.format = format.(string)
	}

	return nil
}

func (e *OEmbedExtractor) Supports(doc document.Document) bool {
	// Look for an oembed like tag
	var findTag func(*html.Node) bool
	findTag = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			if n.Data == "like" {
				// Load values
				var typ, href string
				for _, attr := range n.Attr {
					if attr.Key == "type" {
						typ = attr.Val
					} else if attr.Key == "href" {
						href = attr.Val
					}

					// Check if type is a valid oembed discovery type
					if typ == "application/json+oembed" {
						e.format = "json"
						e.endpoint = href

						return true
					} else if typ == "text/xml+oembed" {
						e.format = "xml"
						e.endpoint = href

						return true
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findTag(c) {
				return true
			}
		}

		return false
	}

	return findTag(doc.Doc.Node())
}

func (e *OEmbedExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	url := fmt.Sprintf(e.endpoint, url.QueryEscape(doc.Url))

	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	// Decode result
	var res map[string]interface{}

	if e.format == "json" || e.format == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, "", err
		}
	} else {
		decoder := xml.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, "", err
		}
	}

	switch res["type"] {
	case "photo":
		res = map[string]interface{}{
			"title": res["title"],
			"author": map[string]interface{}{
				"name": res["author_name"],
				"url":  res["author_url"],
			},
			"thumbnail": map[string]interface{}{
				"url":    res["thumbnail_url"],
				"width":  res["thumbnail_width"],
				"height": res["thumbnail_height"],
			},
			"url":    res["url"],
			"width":  res["width"],
			"height": res["height"],
		}
	case "video":
		res = map[string]interface{}{
			"title": res["title"],
			"author": map[string]interface{}{
				"name": res["author_name"],
				"url":  res["author_url"],
			},
			"thumbnail": map[string]interface{}{
				"url":    res["thumbnail_url"],
				"width":  res["thumbnail_width"],
				"height": res["thumbnail_height"],
			},
			"html":   res["html"],
			"width":  res["width"],
			"height": res["height"],
		}
	}

	return res, res["type"].(string), nil
}

func init() {
	RegisterPlugin("oembed", new(OEmbedExtractor))
}
