package gofetch

import (
	"io"
)

type Result struct {
	Url      string
	PageType PageType
	Body     io.ReadCloser
	Content  interface{}
}
