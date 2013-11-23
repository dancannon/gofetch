package gofetch

import (
	"fmt"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	"github.com/imdario/mergo"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type Fetcher struct {
	Config config.Config
}

func NewFetcher(config config.Config) *Fetcher {
	return &Fetcher{
		Config: config,
	}
}

func (f *Fetcher) Fetch(url string) (Result, error) {
	// Load all the rules
	for _, pc := range f.Config.RuleProviders {
		provider, err := loadProvider(pc)
		if err != nil {
			continue
		}

		f.Config.Rules = append(f.Config.Rules, provider.Provide()...)
	}

	// Sort the rules
	sort.Sort(sort.Reverse(config.RuleSlice(f.Config.Rules)))

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
			doc := document.NewDocument(res.Request.URL.String(), res.Body)
			return f.parseDocument(doc), nil
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

func (f *Fetcher) parseDocument(doc *document.Document) Result {
	// Prepare document for parsing
	cleanDocument(doc)

	res := Result{
		Url:  doc.Url,
		Body: doc.Raw,
	}
	res.Content = make(map[string]interface{})

	// Iterate through all registered rules and find one that can be used
	for _, rule := range f.Config.Rules {
		for _, url := range rule.Urls {
			re := regexp.MustCompile(url)
			if re.MatchString(doc.Url) {
				res.Content = f.loadValues(rule.Values, doc)
				return res
			}
		}
	}

	return res
}

func (f *Fetcher) loadValues(values map[string]interface{}, doc *document.Document) interface{} {
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

			res := runExtractor(ec, doc)

			switch res := res.(type) {
			case map[string]interface{}:
				if err := mergo.Merge(&m, res); err != nil {
					panic(err.Error())
				}
			default:
				return res
			}
		} else {
			switch val := val.(type) {
			case map[string]interface{}:
				m[key] = f.loadValues(val, doc)
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
			return nil
		}

		return eres
	} else {
		panic(fmt.Sprintf("Extractor %s not found", id))
	}
}
