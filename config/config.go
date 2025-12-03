package config

import (
	"encoding/json"
	"os"
)

// Route defines a mapping from a path to a backend URL.
type Route struct {
	Path       string `json:"path"`
	BackendURL string `json:"backend_url"`
}

// Config holds the application configuration.
type Config struct {
	Port   string  `json:"port"`
	Routes []Route `json:"routes"`
}

// LoadConfig loads configuration from a JSON file.
// It supports environment variable expansion in the JSON file (e.g. ${VAR})
// and overrides the Port if the PORT environment variable is set.
func LoadConfig(filename string) (*Config, error) {
	// 1. Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 2. Expand environment variables in the content
	content := os.ExpandEnv(string(data))

	// 3. Decode JSON
	config := &Config{}
	err = json.Unmarshal([]byte(content), config)
	if err != nil {
		return nil, err
	}

	// 4. Override Port from environment variable if present
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	return config, nil
}
