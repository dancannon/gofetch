package parser

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"fmt"
	"strings"
)

type Document struct {
	title  string
	blocks []string
}

type Parser struct {
	doc   *Document
	flush bool
}

func NewParser() *Parser {
	return &Parser{
		doc: &Document{},
	}
}

func (p *Parser) Parse(s string) error {
	r := bytes.NewBufferString(s)
	z := html.NewTokenizer(r)

	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return z.Err()
		case html.TextToken:
			text := strings.TrimSpace(string(z.Text()))
			if depth > 0 && text != "" {
				fmt.Println(text)
			}
		case html.StartTagToken:
			depth++
		case html.EndTagToken:
			depth--
		}
	}
}
