package gofetch

import (
	"net/http"
	"strings"
)

var (
	HtmlTypes = []string{
		"html",
		"xml",
	}
	ParseableTypes = []string{
		"text",
		"html",
		"xml",
	}
)

func isContentTypeParsable(res *http.Response) bool {
	for _, typ := range ParseableTypes {
		if strings.Contains(res.Header.Get("Content-Type"), typ) {
			return true
		}
	}

	return false
}

func isContentTypeHtml(res *http.Response) bool {
	for _, typ := range HtmlTypes {
		if strings.Contains(res.Header.Get("Content-Type"), typ) {
			return true
		}
	}

	return false
}
