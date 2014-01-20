package plugins

import "github.com/dancannon/gofetch/document"

var (
	plugins = make(map[string]Plugin)
)

type Plugin interface {
	Setup(config interface{}) error
}

type Supported interface {
	Supports(doc document.Document) bool
}

type Extractor interface {
	Plugin

	Extract(doc document.Document) (interface{}, error)
}

type MultiExtractor interface {
	Plugin

	ExtractValues(doc document.Document) (values interface{}, pagetype string, err error)
}

func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}

func GetExtractor(name string) Extractor {
	if plugin, ok := plugins[name]; ok {
		if extractor, ok := plugin.(Extractor); ok {
			return extractor
		}
	}

	return nil
}

func GetMultiExtractor(name string) MultiExtractor {
	if plugin, ok := plugins[name]; ok {
		if extractor, ok := plugin.(MultiExtractor); ok {
			return extractor
		}
	}

	return nil
}
