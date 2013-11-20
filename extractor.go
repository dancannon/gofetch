package gofetch

import (
	// "github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
)

type Extractor interface {
	Id() string
	Setup(map[string]interface{}) error
	Extract(*document.Document) (interface{}, error)
}
