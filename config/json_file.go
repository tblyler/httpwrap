package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// JSONFileSource gets the config from a JSON file
type JSONFileSource struct {
	filePath string
}

// NewJSONFileSource creates a new JSONFileSource instance for the given path
func NewJSONFileSource(filePath string) *JSONFileSource {
	return &JSONFileSource{
		filePath: filePath,
	}
}

// Config gets the config from this file source
func (jfs *JSONFileSource) Config(ctx context.Context) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(jfs.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config file at %s: %w", jfs.filePath, err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON config file at %s: %w", jfs.filePath, err)
	}

	return config, nil
}
