package gofetch

import (
	"fmt"
	"github.com/dancannon/gofetch/config"
	"github.com/imdario/mergo"
	"io/ioutil"
	"net/http"
	"reflect"
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
	// Sort the rules
	sort.Sort(sort.Reverse(config.RuleSlice(f.Config.Rules)))

	// Make request
	response, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}

	doc := NewDocument(response.Request.URL.String(), response.Body)

	var result Result

	// Check the returned MIME type
	if isContentTypeParsable(response) {
		// If the page was HTML then parse the HTMl otherwise return the plain
		// text
		if isContentTypeHtml(response) {
			result = f.parseDocument(doc)
		} else {
			text, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return Result{}, err
			}

			result = Result{
				Url:      response.Request.URL.String(),
				PageType: "text",
				Content: map[string]interface{}{
					"text": text,
				},
			}
		}
	} else {
		// If the content cannot be parsed then guess the page type based on the
		// Content-Type header
		result = Result{
			Url:      response.Request.URL.String(),
			PageType: response.Header.Get("Content-Type"),
		}
	}

	// Validate the result
	err = f.validateResult(result)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (f *Fetcher) parseDocument(doc *Document) Result {
	// Prepare document for parsing
	cleanDocument(doc)

	res := Result{
		Url: doc.Url,
	}
	res.Content = make(map[string]interface{})

	// Iterate through all registered rules and find one that can be used
	for _, rule := range f.Config.Rules {
		for _, url := range rule.Urls {
			re := regexp.MustCompile(url)
			if re.MatchString(doc.Url) {
				res.PageType = rule.Type
				res.Content = f.loadValues(rule.Values, doc, &res)
				return res
			}
		}
	}

	return res
}

func (f *Fetcher) loadValues(values map[string]interface{}, doc *Document, r *Result) interface{} {
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

			res := runExtractor(ec, doc, r)

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
				m[key] = f.loadValues(val, doc, r)
			default:
				m[key] = val
			}
		}
	}

	return m
}

func runExtractor(config map[string]interface{}, doc *Document, res *Result) interface{} {
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

		eres, err := extractor.Extract(doc, res)
		if err != nil {
			return nil
		}

		return eres
	} else {
		panic(fmt.Sprintf("Extractor %s not found", id))
	}
}

func (f *Fetcher) validateResult(r Result) error {
	// Check that the result uses a known type
	for _, t := range f.Config.Types {
		if t.Id == r.PageType {
			if !t.Validate {
				return nil
			}

			return validateResultValues(t.Id, r.Content, t.Values)
		}
	}

	return fmt.Errorf("The page type %s does not exist", r.PageType)
}

func validateResultValues(pagetype string, values interface{}, typValues interface{}) error {
	// Check that both values have the same type
	if reflect.TypeOf(values) != reflect.TypeOf(typValues) {
		return fmt.Errorf("The result is not of the correct type")
	}

	switch typValues := typValues.(type) {
	// If the value is a map then validate each node
	case map[string]interface{}:
		seenNodes := []string{}

		valuesM := values.(map[string]interface{})

		for k, v := range typValues {
			// Ensure that the type value is of type map
			if v, ok := v.(map[string]interface{}); !ok {
				return fmt.Errorf("The result is not of the correct type")
			} else {
				// Check that the value has the node if it is required
				if required, ok := v["required"].(bool); ok && required {
					if _, ok := valuesM[k]; !ok {
						return fmt.Errorf("The type %s requires the field %s", pagetype, k)
					}
				}

				if _, ok := valuesM[k]; !ok {
					continue
				}

				// Validate any children nodes if they exist
				if childTypValues, ok := v["values"]; ok {
					err := validateResultValues(pagetype, valuesM[k], childTypValues)
					if err != nil {
						return err
					}
				}

				seenNodes = append(seenNodes, k)
			}
		}

		// Check that the result value doesnt have any extra nodes
		for k, _ := range valuesM {
			seen := false

			for _, sk := range seenNodes {
				if !seen && k == sk {
					seen = true
				}
			}

			if !seen {
				return fmt.Errorf("The type %s does not contain the field %s", pagetype, k)
			}
		}
	// If the value was not a map then the value does not validate
	default:
		return fmt.Errorf("The result is not of the correct type")
	}

	return nil
}
