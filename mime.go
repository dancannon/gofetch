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
	TypeMap = map[PageType][]string{
		Image: []string{"image"},
		Audio: []string{"audio"},
		Video: []string{"video"},
		Flash: []string{"flash"},
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

func guessPageTypeFromMime(res *http.Response) PageType {
	for typ, strs := range TypeMap {
		for _, str := range strs {
			if strings.Contains(res.Header.Get("Content-Type"), str) {
				return typ
			}
		}
	}

	return Unknown
}
