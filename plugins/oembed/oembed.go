package oembed

import (
	. "github.com/dancannon/gofetch/message"
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

func (e *OEmbedExtractor) Extract(msg *ExtractMessage) error {
	url := fmt.Sprintf(e.endpoint, url.QueryEscape(msg.Document.Url))

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// Decode result
	var res map[string]interface{}

	if e.format == "json" || e.format == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return err
		}
	} else {
		decoder := xml.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return err
		}
	}

	// Override the result page type
	msg.PageType = res["type"].(string)

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

	msg.Value = res

	return nil
}

func init() {
	RegisterPlugin("oembed", new(OEmbedExtractor))
}
