package gofetch

import (
	"github.com/dancannon/gofetch/document"
)

type Extractor interface {
	Id() string
	Setup(map[string]string) error
	Extract(*document.Document) (interface{}, error)
}
