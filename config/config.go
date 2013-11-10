package config

import (
	"encoding/xml"
	"github.com/davecgh/go-spew/spew"
	"os"
	"path/filepath"
)

type Config struct {
	RuleProviders []ruleProviderConfig `xml:"rules>provider"`
	Rules         []Rule               `xml:"rules>rule"`
}

type Parameter struct {
	Key   string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type ParameterSlice []Parameter

func (p ParameterSlice) ToMap() map[string]string {
	m := map[string]string{}

	for _, e := range p {
		m[e.Key] = e.Value
	}

	return m
}

func LoadConfig(path string) Config {
	var config Config

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(absPath)
	if err != nil {
		panic("Error opening file")
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		spew.Dump(err)
		panic("Error decoding config file")
	}

	return config
}
