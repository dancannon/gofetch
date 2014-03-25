package oembed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"code.google.com/p/go.net/html"
	"github.com/clbanning/mxj"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	"github.com/dancannon/gofetch/util"
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
					e.endpoint = href + "&maxwidth=480&maxheight=360"

					return true
				} else if typ == "text/xml+oembed" {
					e.format = "xml"
					e.endpoint = href + "&maxwidth=480&maxheight=360"

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
		endpoint = fmt.Sprintf(e.endpointFormat, doc.URL.String())
	}

	// Resolve absolute endpoint url
	eu, err := url.Parse(endpoint)
	if err != nil {
		return nil, "", err
	}
	endpoint = doc.URL.ResolveReference(eu).String()

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	// Ensure that the endpoint was found
	if response.StatusCode != 200 {
		return nil, "", errors.New("There was an error fetching result from the OEmbed endpoint")
	}

	// Decode result
	var props map[string]interface{}

	if e.format == "json" || e.format == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&props)
		if err != nil {
			return nil, "", err
		}
	} else {
		props, err = mxj.NewMapXmlReader(response.Body)
		if err != nil {
			return nil, "", err
		}

		if _, ok := props["oembed"]; !ok {
			return nil, "", errors.New("OEmbed data in incorrect format")
		}
		if _, ok := props["oembed"].(map[string]interface{}); !ok {
			return nil, "", errors.New("OEmbed data in incorrect format")
		}

		props = props["oembed"].(map[string]interface{})
	}

	var resptype string
	if t, ok := props["type"].(string); ok {
		resptype = t
	} else {
		resptype = "unknown"
	}

	switch resptype {
	case "photo", "image":
		res := util.CreateMapFromProps(props, map[string]string{
			"title":      "title",
			"html":       "html",
			"width:int":  "width",
			"height:int": "height",
			"url":        "url",
		})
		if _, ok := props["author_url"]; ok {
			res["author"] = util.CreateMapFromProps(props, map[string]string{
				"name": "author_name",
				"url":  "author_url",
			})
		}
		if _, ok := props["thumbnail_url"]; ok {
			res["thumbnail"] = util.CreateMapFromProps(props, map[string]string{
				"url":        "thumbnail_url",
				"width:int":  "thumbnail_width",
				"height:int": "thumbnail_height",
			})
		}

		return res, "image", nil
	case "video":
		res := util.CreateMapFromProps(props, map[string]string{
			"title": "title",
			"html":  "html",
		})
		if _, ok := props["author_url"]; ok {
			res["author"] = util.CreateMapFromProps(props, map[string]string{
				"name": "author_name",
				"url":  "author_url",
			})
		}
		if _, ok := props["thumbnail_url"]; ok {
			res["thumbnail"] = util.CreateMapFromProps(props, map[string]string{
				"url":        "thumbnail_url",
				"width:int":  "thumbnail_width",
				"height:int": "thumbnail_height",
			})
		}

		return res, "video", nil
	case "rich":
		res := util.CreateMapFromProps(props, map[string]string{
			"title": "title",
			"html":  "html",
		})

		return res, "rich", nil
	}

	return props, "unknown", nil
}

func init() {
	RegisterPlugin("oembed", new(OEmbedExtractor))
}
