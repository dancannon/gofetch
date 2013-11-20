package gofetch

import (
	"fmt"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var scs spew.ConfigState = spew.ConfigState{Indent: "\t"}
var scs2 spew.ConfigState = spew.ConfigState{Indent: "\t", MaxDepth: 2}
var c = config.LoadConfig("config.json")

func Fetch(url string) (Result, error) {
	// Load all the rules
	for _, pc := range c.RuleProviders {
		provider, err := loadProvider(pc)
		if err != nil {
			continue
		}

		c.Rules = append(c.Rules, provider.Provide()...)
	}

	// Sort the rules
	sort.Sort(sort.Reverse(config.RuleSlice(c.Rules)))

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

	// // Iterate through all registered extractors and find one that can be used
	// for _, rule := range c.Rules {
	// 	for _, url := range rule.Urls {
	// 		re := regexp.MustCompile(url)
	// 		if re.MatchString(doc.Url) {
	// 			if extractor, ok := extractors[rule.Extractor]; ok {
	// 				err := extractor.Setup(rule.Values)
	// 				if err != nil {
	// 					panic(err.Error())
	// 				}

	// 				content, err := extractor.Extract(doc)
	// 				if err != nil {
	// 					panic(err.Error())
	// 				}
	// 				res.Content = content
	// 				return res
	// 			} else {
	// 				panic(fmt.Sprintf("Extractor %s not found", rule.Extractor))
	// 			}
	// 		}
	// 	}
	// }
	res.Content = make(map[string]interface{})

	// Iterate through all registered rules and find one that can be used
	for _, rule := range c.Rules {
		for _, url := range rule.Urls {
			re := regexp.MustCompile(url)
			if re.MatchString(doc.Url) {
				res.Content = loadValues(rule.Values, doc)
				return res
			}
		}
	}

	return res
}

func loadValues(values map[string]interface{}, doc *document.Document) interface{} {
	m := map[string]interface{}{}

	for key, val := range values {
		// If value is an extractor reference then run the extractor and merge
		// the result. Extractors start with a '@'
		if strings.Index(key, "@") == 0 {
			// Ensure that the extractor config is in the right format
			ec, ok := val.(map[string]interface{})
			if !ok {
				panic("The extractor configuration is invalid")
			}

			return runExtractor(ec, doc)
		} else {
			switch val := val.(type) {
			case map[string]interface{}:
				// m[key] = make(map[string]interface{})
				m[key] = loadValues(val, doc)
			default:
				m[key] = val
			}
		}
	}

	return m
}

func runExtractor(config map[string]interface{}, doc *document.Document) interface{} {
	// Validate extractor config
	id, ok := config["id"].(string)
	if !ok {
		panic("The extractor configuration is invalid")
	}
	params, ok := config["params"].(map[string]interface{})
	if !ok {
		params = make(map[string]interface{})
	}

	// Load and execute the extractor
	if extractor, ok := extractors[id]; ok {
		err := extractor.Setup(params)
		if err != nil {
			panic(err.Error())
		}

		eres, err := extractor.Extract(doc)
		if err != nil {
			panic(err.Error())
		}

		return eres
	} else {
		panic(fmt.Sprintf("Extractor %s not found", id))
	}
}
