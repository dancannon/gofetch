package gofetch

import (
	neturl "net/url"
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/document"

	. "github.com/dancannon/gofetch/plugins"
	_ "github.com/dancannon/gofetch/plugins/base"
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
		// If the URLPattern was set then first match against the URL
		if rule.UrlPattern != "" {
			var re *regexp.Regexp

			// Check url against the url regular expression
			re = regexp.MustCompile(rule.UrlPattern)
			if !re.MatchString(url.String()) {
				continue
			}
		} else {

			// Otherwise match against the host and path separately
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
		}

		// Set the base page type
		res.PageType = rule.Type

		// Extract a single value and type if the rule contains a MultiExtractor
		if value, typ, err := f.extractValue(rule.Values, doc); err != nil {
			return Result{}, err
		} else {
			if typ != "" {
				res.PageType = typ
			}
			res.Content = value

			return res, nil
		}

		// Otherwise fallback to the normal extraction method.
		if value, err := f.extractValues(rule.Values, doc); err != nil {
			return Result{}, err
		} else {
			res.Content = value

			return res, nil
		}
	}

	return res, nil
}

// Check the first level of values for an extractor, if one is found then immediately
// return the result of the extractor.
func (f *Fetcher) extractValue(values []interface{}, doc *document.Document) (value interface{}, pageType string, err error) {
	for _, val := range values {
		var props map[string]interface{}

		if _, ok := val.(map[string]interface{}); !ok {
			return nil, "", fmt.Errorf("The value configuration is invalid")
		}
		props = val.(map[string]interface{})

		if typ, ok := props["type"]; ok {
			switch typ {
			case "extractor":
				if v, t, supported, err := f.runExtractor(true, props, doc); err != nil {
					return nil, "", err
				} else if !supported {
					continue
				} else {
					return v, t, nil
				}
			}
		}
	}

	value, err = f.extractValues(values, doc)
	return
}

func (f *Fetcher) extractValues(values []interface{}, doc *document.Document) (interface{}, error) {
	m := map[string]interface{}{}

	for _, val := range values {
		var name string
		var props map[string]interface{}

		if _, ok := val.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("The value configuration is invalid")
		}
		props = val.(map[string]interface{})

		if val, ok := props["name"].(string); !ok {
			continue
		} else {
			name = val
		}

		if typ, ok := props["type"]; ok {
			switch typ {
			case "extractor":
				if v, _, supported, err := f.runExtractor(false, props, doc); err != nil {
					return nil, err
				} else if !supported {
					continue
				} else {
					m[name] = v
				}
			case "value":
				if _, ok := props["value"]; !ok {
					continue
				}

				m[name] = props["value"]
			case "values":
				if _, ok := props["value"].([]interface{}); !ok {
					continue
				}

				v, err := f.extractValues(props["value"].([]interface{}), doc)
				if err != nil {
					return nil, err
				}
				m[name] = v
			}
		}
	}

	return m, nil
}

func (f *Fetcher) runExtractor(multi bool, props map[string]interface{}, doc *document.Document) (value interface{}, pageType string, supported bool, err error) {
	id, ok := props["id"].(string)
	if !ok {
		return nil, "", false, fmt.Errorf("The extractor configuration is invalid")
	}
	params, ok := props["params"].(map[string]interface{})
	if !ok {
		params = make(map[string]interface{})
	}

	if multi {
		if extractor := GetMultiExtractor(id); extractor != nil {
			// Check if the extractor is able to extract from the document
			if extractor, ok := extractor.(Supported); ok {
				if !extractor.Supports(*doc) {
					return nil, "", false, nil
				}
			}

			err := extractor.Setup(params)
			if err != nil {
				return nil, "", false, err
			}

			v, t, err := extractor.ExtractValues(*doc)
			if err != nil {
				return nil, "", false, err
			}

			return v, t, true, nil
		} else {
			return nil, "", false, nil
		}
	} else {
		if extractor := GetExtractor(id); extractor != nil {
			// Check if the extractor is able to extract from the document
			if extractor, ok := extractor.(Supported); ok {
				if !extractor.Supports(*doc) {
					return nil, "", false, nil
				}
			}

			err := extractor.Setup(params)
			if err != nil {
				return nil, "", false, err
			}

			v, err := extractor.Extract(*doc)
			if err != nil {
				return nil, "", false, err
			}

			return v, "", true, nil
		} else {
			return nil, "", false, fmt.Errorf("Extractor %s not found", id)
		}
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
