package gofetch

import (
	"net/http"
)

type Result struct {
	Url      string
	PageType string
	Content  interface{}
	Response http.Response `json:"-"`
}
