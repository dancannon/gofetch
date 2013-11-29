package gofetch

import (
	"github.com/dancannon/gofetch/plugin/oembed"
	"github.com/dancannon/gofetch/plugin/opengraph"
	"github.com/dancannon/gofetch/plugin/selector"
	"github.com/dancannon/gofetch/plugin/selector_text"
	"github.com/dancannon/gofetch/plugin/text"
	"github.com/dancannon/gofetch/plugin/title"
)

var (
	extractors map[string]Extractor = map[string]Extractor{}
)

func RegisterExtractor(extractor Extractor) {
	extractors[extractor.Id()] = extractor
}

func init() {
	// Register all plugins
	RegisterExtractor(new(oembed.Extractor))
	RegisterExtractor(new(opengraph.Extractor))
	RegisterExtractor(new(selector.Extractor))
	RegisterExtractor(new(selector_text.Extractor))
	RegisterExtractor(new(text.Extractor))
	RegisterExtractor(new(title.Extractor))
}
