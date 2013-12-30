package plugins

import (
	"github.com/dancannon/gofetch/message"
	"github.com/davecgh/go-spew/spew"
)

var (
	Extractors = make(map[string]Extractor)
)

type Plugin interface {
	Setup(config interface{}) error
}

type Extractor interface {
	Plugin

	Extract(msg *message.ExtractMessage) error
}

func RegisterPlugin(name string, plugin Plugin) {
	// Check if plugin is an extractor
	spew.Dump(plugin.(Extractor))
	if extractor, ok := plugin.(Extractor); ok {
		Extractors[name] = extractor
	}
}
