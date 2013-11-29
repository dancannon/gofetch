package gofetch

import (
	"github.com/dancannon/gofetch/document"
)

type Result struct {
	Url      string
	PageType PageType
	Document *document.Document
	Content  interface{}
}
