package gofetch

import (
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

func (e *OEmbedExtractor) Id() string {
	return "gofetch.oembed.extractor"
}

func (e *OEmbedExtractor) Setup(config map[string]interface{}) error {
	// Validate config
	if endpoint, ok := config["endpoint"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed an endpoint", e.Id()))
	} else {
		e.endpoint = endpoint.(string)
	}

	if format, ok := config["format"]; ok {
		e.format = format.(string)
	}

	return nil
}

func (e *OEmbedExtractor) Extract(d *Document, r *Result) (interface{}, error) {
	url := fmt.Sprintf(e.endpoint, url.QueryEscape(d.Url))

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Decode result
	var res map[string]interface{}

	if e.format == "json" || e.format == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		decoder := xml.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	// Override the result page type
	r.PageType = res["type"].(string)

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

	return interface{}(res), nil
}
