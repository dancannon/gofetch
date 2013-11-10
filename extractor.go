package gofetch

import (
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
)

type Extractor interface {
	Id() string
	Setup([]config.Value) error
	Extract(*document.Document) (map[string]interface{}, error)
}
