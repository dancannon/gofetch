package message

import (
	"github.com/dancannon/gofetch/document"
)

type ExtractMessage struct {
	PageType string
	Value    interface{}

	Document *document.Document
}
