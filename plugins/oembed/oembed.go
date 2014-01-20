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
	endpoint       string
	endpointFormat string
	format         string
}

func (e *OEmbedExtractor) Setup(config interface{}) error {
	var params map[string]interface{}
	if p, ok := config.(map[string]interface{}); !ok {
		params = make(map[string]interface{})
	} else {
		params = p
	}

	// Validate config
	if e.endpoint == "" {
		if endpointFormat, ok := params["endpoint"]; !ok {
			return errors.New(fmt.Sprintf("The oembed extractor must be passed an endpoint"))
		} else {
			e.endpointFormat = endpointFormat.(string)
		}
	}

	if format, ok := params["format"]; ok {
		e.format = format.(string)
	}

	return nil
}

func (e *OEmbedExtractor) Supports(doc document.Document) bool {
	// Look for an oembed like tag
	var findOEmbedTag func(*html.Node) bool
	findOEmbedTag = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			if n.Data == "link" {
				// Load values
				var typ, href string
				for _, attr := range n.Attr {
					if attr.Key == "type" {
						typ = attr.Val
					} else if attr.Key == "href" {
						href = attr.Val
					}
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

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findOEmbedTag(c) {
				return true
			}
		}

		return false
	}

	return findOEmbedTag(doc.Doc.Node())
}

func (e *OEmbedExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	var endpoint string
	if e.endpoint != "" {
		endpoint = e.endpoint
	} else {
		endpoint = fmt.Sprintf(e.endpoint, url.QueryEscape(doc.Url))
	}

	// Resolve absolute endpoint url
	hu, err := url.Parse(doc.Url)
	if err != nil {
		return nil, "", err
	}
	eu, err := url.Parse(endpoint)
	if err != nil {
		return nil, "", err
	}
	endpoint = hu.ResolveReference(eu).String()

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	// Decode result
	var resp map[string]interface{}

	if e.format == "json" || e.format == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&resp)
		if err != nil {
			return nil, "", err
		}
	} else {
		decoder := xml.NewDecoder(response.Body)
		err = decoder.Decode(&resp)
		if err != nil {
			return nil, "", err
		}
	}

	var resptype string
	if t, ok := resp["type"].(string); ok {
		resptype = t
	} else {
		resptype = "unknown"
	}

	switch resptype {
	case "photo":
		resp = map[string]interface{}{
			"title": resp["title"],
			"author": map[string]interface{}{
				"name": resp["author_name"],
				"url":  resp["author_url"],
			},
			"thumbnail": map[string]interface{}{
				"url":    resp["thumbnail_url"],
				"width":  resp["thumbnail_width"],
				"height": resp["thumbnail_height"],
			},
			"url":    resp["url"],
			"width":  resp["width"],
			"height": resp["height"],
		}
	case "video":
		resp = map[string]interface{}{
			"title": resp["title"],
			"author": map[string]interface{}{
				"name": resp["author_name"],
				"url":  resp["author_url"],
			},
			"thumbnail": map[string]interface{}{
				"url":    resp["thumbnail_url"],
				"width":  resp["thumbnail_width"],
				"height": resp["thumbnail_height"],
			},
			"html":   resp["html"],
			"width":  resp["width"],
			"height": resp["height"],
		}
	}

	return resp, resptype, nil
}

func init() {
	RegisterPlugin("oembed", new(OEmbedExtractor))
}
