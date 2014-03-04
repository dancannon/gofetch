package base

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	. "github.com/dancannon/gofetch/plugins/oembed"
	. "github.com/dancannon/gofetch/plugins/opengraph"
	. "github.com/dancannon/gofetch/plugins/text"
	. "github.com/dancannon/gofetch/plugins/title"
)

type BaseExtractor struct {
	params map[string]interface{}

	oembedExtractor    OEmbedExtractor
	opengraphExtractor OpengraphExtractor
	textExtractor      TextExtractor
	titleExtractor     TitleExtractor
}

func (e *BaseExtractor) Setup(config interface{}) error {
	if p, ok := config.(map[string]interface{}); !ok {
		e.params = make(map[string]interface{})
	} else {
		e.params = p
	}

	return nil
}

func (e *BaseExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	if e.oembedExtractor.Supports(doc) {
		// Attempt the OEmbed extractor first as it has the best accuracy
		if err := e.oembedExtractor.Setup(e.getExtractorParams("oembed")); err != nil {
			return nil, "", err
		}

		return e.oembedExtractor.ExtractValues(doc)
	} else if e.opengraphExtractor.Supports(doc) {
		// Next try the opengraph extractor
		if err := e.opengraphExtractor.Setup(e.getExtractorParams("opengraph")); err != nil {
			return nil, "", err
		}

		return e.opengraphExtractor.ExtractValues(doc)
	} else {
		// If the previous two extractors did not work then manually create a rule
		// using the text and title extractors
		if err := e.titleExtractor.Setup(e.getExtractorParams("title")); err != nil {
			return nil, "", err
		}
		if err := e.textExtractor.Setup(e.getExtractorParams("text")); err != nil {
			return nil, "", err
		}

		title, err := e.titleExtractor.Extract(doc)
		if err != nil {
			return nil, "", err
		}

		text, err := e.textExtractor.Extract(doc)
		if err != nil {
			return nil, "", err
		}

		return map[string]interface{}{
			"title": title,
			"text":  text,
		}, "text", nil
	}
}

func (e *BaseExtractor) getExtractorParams(id string) interface{} {
	if v, ok := e.params[id]; ok {
		return v
	}

	return map[string]interface{}{}
}

func init() {
	RegisterPlugin("base", new(BaseExtractor))
}
