package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Version  string        `json:"version"`
	Paths    PathConfig    `json:"paths"`
	Commands CommandConfig `json:"commands"`
}

type PathConfig struct {
	Main  string `json:"main"`
	Sync  string `json:"sync"`
	Build string `json:"build"`
}

type CommandConfig struct {
	Install Command            `json:"install"`
	Build   Command            `json:"build"`
	Release Command            `json:"release"`
	Checks  map[string]Command `json:"checks"`
}

type Command struct {
	Cmd         []string          `json:"cmd"`
	Description string            `json:"description"`
	Env         map[string]string `json:"env,omitempty"`
}

// Load reads a configuration file from the specified path and unmarshals it into a Config struct.
// It returns a pointer to the Config struct if successful, or an error if any reading or unmarshalling errors occur.
// Parameters:
//   - configPath: the path to the configuration file
//
// Returns:
//   - *Config: a pointer to the Config struct containing the configuration data
//   - error: an error if any reading or unmarshalling issues occur
func Load(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("failed to validate config: %v", err)
	}

	return &config, nil
}

func validate(config Config) error {
	if config.Paths.Main == "" {
		return fmt.Errorf("MainDir is required in the config")
	}
	if config.Paths.Sync == "" {
		return fmt.Errorf("SyncDir is required in the config")
	}
	if config.Paths.Build == "" {
		return fmt.Errorf("BuildDir is required in the config")
	}
	if len(config.Commands.Install.Cmd) < 2 {
		return fmt.Errorf("InstallCommand must have at least two elements")
	}

	return nil
}
