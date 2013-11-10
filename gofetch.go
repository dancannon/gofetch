package gofetch

import (
	"fmt"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"regexp"
	"sort"
)

var scs spew.ConfigState = spew.ConfigState{Indent: "\t"}
var scs2 spew.ConfigState = spew.ConfigState{Indent: "\t", MaxDepth: 2}
var c = config.LoadConfig("config.xml")

func Fetch(url string) (Result, error) {
	// Load all the rules
	for _, pc := range c.RuleProviders {
		provider, err := loadProvider(pc.Id, pc.Parameters)
		if err != nil {
			continue
		}

		c.Rules = append(c.Rules, provider.Provide()...)
	}

	// Sort the rules
	sort.Sort(config.RuleSlice(c.Rules))

	// Make request
	res, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}

	// Check the returned MIME type
	if isContentTypeParsable(res) {
		// If the page was HTML then parse the HTMl otherwise return the plain
		// text
		if isContentTypeHtml(res) {
			return parseHtml(Result{
				Url:  res.Request.URL.String(),
				Body: res.Body,
			}), nil
		} else {
			return Result{
				Url:      res.Request.URL.String(),
				PageType: PlainText,
				Body:     res.Body,
			}, nil
		}
	} else {
		// If the content cannot be parsed then guess the page type based on the
		// Content-Type header
		return Result{
			Url:      res.Request.URL.String(),
			PageType: guessPageTypeFromMime(res),
		}, nil
	}
}

func parseHtml(res Result) Result {
	doc := document.NewDocument(res.Url, res.Body)
	cleanDocument(doc)

	// Iterate through all registered extractors and find one that can be used
	for _, rule := range c.Rules {
		for _, url := range rule.Urls {
			re := regexp.MustCompile(url)
			if re.MatchString(doc.Url) {
				if extractor, ok := extractors[rule.Extractor]; ok {
					err := extractor.Setup(rule.Values)
					if err != nil {
						panic(err.Error())
					}

					content, err := extractor.Extract(doc)
					if err != nil {
						panic(err.Error())
					}
					res.Content = content
					return res
				} else {
					panic(fmt.Sprintf("Extractor %s not found", rule.Extractor))
				}
			}
		}
	}

	return res
}
