package config

import (
	"encoding/json"
	"github.com/dancannon/gofetch/sandbox"
	"github.com/davecgh/go-spew/spew"
	"os"
	"path/filepath"
)

type Config struct {
	Plugins []sandbox.SandboxConfig `json:"plugins"`
	Rules   []Rule                  `json:"rules"`
	Types   []Type                  `json:"types"`
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

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		spew.Dump(err)
		panic("Error decoding config file")
	}

	return config
}
