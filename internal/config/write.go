package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Write marshals the config to YAML and writes it to path.
func Write(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
