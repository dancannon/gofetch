package gofetch

import (
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	neturl "net/url"

	_ "github.com/dancannon/gofetch/plugins/javascript"
	_ "github.com/dancannon/gofetch/plugins/oembed"
	_ "github.com/dancannon/gofetch/plugins/opengraph"
	_ "github.com/dancannon/gofetch/plugins/selector"
	_ "github.com/dancannon/gofetch/plugins/selector_text"
	_ "github.com/dancannon/gofetch/plugins/text"
	_ "github.com/dancannon/gofetch/plugins/title"
	_ "github.com/dancannon/gofetch/plugins/url_mapper"

	"fmt"
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

	doc := document.NewDocument(response.Request.URL.String(), response.Body)

	var result Result

	// Check the returned MIME type
	if isContentTypeParsable(response) {
		// If the page was HTML then parse the HTMl otherwise return the plain
		// text
		if isContentTypeHtml(response) {
			result, err = f.parseDocument(doc)
			if err != nil {
				return Result{}, err
			}
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

func (f *Fetcher) parseDocument(doc *document.Document) (Result, error) {
	// Prepare document for parsing
	cleanDocument(doc)

	res := Result{
		Url:      doc.Url,
		PageType: "unknown",
	}
	res.Content = make(map[string]interface{})

	// Parse the request URL
	url, err := neturl.Parse(doc.Url)
	if err != nil {
		return Result{}, err
	}

	// Iterate through all registered rules and find one that can be used
	for _, rule := range f.Config.Rules {
		var re *regexp.Regexp
		// Clean host
		re = regexp.MustCompile(".*?://")
		host := re.ReplaceAllString(strings.TrimLeft(url.Host, "www."), "")
		ruleHost := re.ReplaceAllString(strings.TrimLeft(rule.Host, "www."), "")

		// Check host
		if host != ruleHost {
			continue
		}

		// Check path against the path regular expression
		re = regexp.MustCompile(rule.PathPattern)
		if !re.MatchString(url.RequestURI()) {
			continue
		}

		// Set the base page type
		res.PageType = rule.Type

		value, typ, err := f.extractTopLevelValues(rule.Values, doc)
		if err != nil {
			return Result{}, err
		}
		if typ != "" {
			res.PageType = typ
		}
		res.Content = value

		return res, nil
	}

	return res, nil
}

// Check the first level of values for an extractor, if one is found then immediately
// return the result of the extractor.
func (f *Fetcher) extractTopLevelValues(values map[string]interface{}, doc *document.Document) (interface{}, string, error) {
	for key, val := range values {
		// If value is an extractor reference then run the extractor and merge
		// the result. Extractors start with a '@'
		if strings.Index(key, "@") == 0 || strings.Index(key, "_") == 0 {
			// Ensure that the extractor config is in the right format
			ec, ok := val.(map[string]interface{})
			if !ok {
				return nil, "", fmt.Errorf("The extractor configuration is invalid")
			}

			id, ok := ec["id"].(string)
			if !ok {
				return nil, "", fmt.Errorf("The extractor configuration is invalid")
			}
			params, ok := ec["params"].(map[string]interface{})
			if !ok {
				params = make(map[string]interface{})
			}

			if extractor := GetMultiExtractor(id); extractor != nil {
				err := extractor.Setup(params)
				if err != nil {
					return nil, "", err
				}

				return extractor.ExtractValues(*doc)
			} else {
				break
			}
		}
	}

	value, err := f.extractValues(values, doc)
	return value, "", err
}

func (f *Fetcher) extractValues(values map[string]interface{}, doc *document.Document) (interface{}, error) {
	m := map[string]interface{}{}

	for key, val := range values {
		// If value is an extractor reference then run the extractor and merge
		// the result. Extractors start with a '@'
		if strings.Index(key, "@") == 0 || strings.Index(key, "_") == 0 {
			// Ensure that the extractor config is in the right format
			ec, ok := val.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("The extractor configuration is invalid")
			}

			id, ok := ec["id"].(string)
			if !ok {
				return nil, fmt.Errorf("The extractor configuration is invalid")
			}
			params, ok := ec["params"].(map[string]interface{})
			if !ok {
				params = make(map[string]interface{})
			}

			if extractor := GetExtractor(id); extractor != nil {
				err := extractor.Setup(params)
				if err != nil {
					return nil, err
				}

				return extractor.Extract(*doc)
			} else {
				return nil, fmt.Errorf("Extractor %s not found", id)
			}
		} else {
			switch val := val.(type) {
			case map[string]interface{}:
				v, err := f.extractValues(val, doc)
				if err != nil {
					return nil, err
				}
				m[key] = v
			default:
				m[key] = val
			}
		}
	}

	return m, nil
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
