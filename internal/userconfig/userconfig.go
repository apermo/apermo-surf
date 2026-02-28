package userconfig

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user-level settings from ~/.config/surf/config.yml.
type Config struct {
	Browser  string                   `yaml:"browser,omitempty"`
	Browsers map[string]BrowserConfig `yaml:"browsers,omitempty"`
}

// BrowserConfig defines a custom browser command.
type BrowserConfig struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}

// Load reads the user config from standard paths.
// Returns a zero-value Config if no file is found (not an error).
func Load() Config {
	for _, path := range configPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			continue
		}
		return cfg
	}
	return Config{}
}

func configPaths() []string {
	var paths []string

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "surf", "config.yml"))
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "surf", "config.yml"))
	}

	paths = append(paths, "/etc/surf/config.yml")
	return paths
}
