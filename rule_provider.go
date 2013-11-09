package gofetch

import (
	"errors"
)

type ruleProviderConfig struct {
	Id         string            `xml:"id"`
	Parameters []configParameter `xml:"parameter"`
}

type RuleProvider interface {
	Setup([]configParameter)
	Provide() []Rule
}

func loadProvider(key string, params []configParameter) (RuleProvider, error) {
	var provider RuleProvider

	switch key {
	case "directory":
		provider = &DirectoryRuleProvider{}
	default:
		return nil, errors.New("No provider was found for the given key")
	}

	provider.Setup(params)
	return provider, nil
}

type DirectoryRuleProvider struct {
	directory string
}

func (p *DirectoryRuleProvider) Setup(params []configParameter) {
	// Check that the provider has the correct parameters
	for _, param := range params {
		if param.Key == "directory" {
			p.directory = param.Value
			return
		}
	}

	panic("The rules directory must be passed to the Xml Rule Provider")
}

func (p *DirectoryRuleProvider) Provide() []Rule {
	return []Rule{}
}
