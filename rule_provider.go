package gofetch

import (
	"errors"
	"github.com/dancannon/gofetch/config"
)

type RuleProvider interface {
	Setup([]config.Parameter)
	Provide() []config.Rule
}

func loadProvider(key string, params []config.Parameter) (RuleProvider, error) {
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

func (p *DirectoryRuleProvider) Setup(params []config.Parameter) {
	// Check that the provider has the correct parameters
	for _, param := range params {
		if param.Key == "directory" {
			p.directory = param.Value
			return
		}
	}

	panic("The rules directory must be passed to the Xml Rule Provider")
}

func (p *DirectoryRuleProvider) Provide() []config.Rule {
	return []config.Rule{}
}
