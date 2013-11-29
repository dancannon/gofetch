package oembed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/dancannon/gofetch/document"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"net/url"
)

type Extractor struct {
	endpoint  string
	extractor string
}

func (e *Extractor) Id() string {
	return "gofetch.oembed.extractor"
}

func (e *Extractor) Setup(config map[string]interface{}) error {
	// Validate config
	if endpoint, ok := config["endpoint"]; !ok {
		return errors.New(fmt.Sprintf("The %s extractor must be passed an endpoint", e.Id()))
	} else {
		e.endpoint = endpoint.(string)
	}

	return nil
}

func (e *Extractor) Extract(d *document.Document) (interface{}, error) {
	url := fmt.Sprintf(e.endpoint, url.QueryEscape(d.Url))

	response, err := http.Get(url)
	spew.Dump(err, url, e.endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Decode result
	var res interface{}
	if e.extractor == "json" || e.extractor == "" {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else {
		decoder := xml.NewDecoder(response.Body)
		err = decoder.Decode(&res)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
