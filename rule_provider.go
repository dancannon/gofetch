package gofetch

import (
	"errors"
	"github.com/dancannon/gofetch/config"
)

type RuleProvider interface {
	Setup(map[string]string)
	Provide() []config.Rule
}

func loadProvider(pc config.ProviderConfig) (RuleProvider, error) {
	var provider RuleProvider

	switch pc.Id {
	case "directory":
		provider = &DirectoryRuleProvider{}
	default:
		return nil, errors.New("No provider was found for the given key")
	}

	provider.Setup(pc.Parameters)
	return provider, nil
}

type DirectoryRuleProvider struct {
	directory string
}

func (p *DirectoryRuleProvider) Setup(params map[string]string) {
	// Check that the provider has the correct parameters
	if directory, ok := params["directory"]; ok {
		p.directory = directory
		return
	}

	panic("The rules directory must be passed to the Xml Rule Provider")
}

func (p *DirectoryRuleProvider) Provide() []config.Rule {
	return []config.Rule{}
}
