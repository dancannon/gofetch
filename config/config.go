package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Rules []Rule `json:"rules"`
	Types []Type `json:"types"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	absPath, err := filepath.Abs(path)
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error opening file: %s", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding config file: %s", err)
	}

	return config, nil
}
