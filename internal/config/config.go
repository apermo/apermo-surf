package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Link represents a project URL, either as a simple string or with a pattern.
type Link struct {
	URL     string `yaml:"url"`
	Pattern string `yaml:"pattern,omitempty"`
}

func (l *Link) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		l.URL = value.Value
		return nil
	}

	// Expanded format: {url: ..., pattern: ...}
	type plain Link
	return value.Decode((*plain)(l))
}

// MarshalYAML writes a Link as a scalar string when there is no pattern,
// or as a mapping {url, pattern} when a pattern is set.
func (l Link) MarshalYAML() (interface{}, error) {
	if l.Pattern == "" {
		return l.URL, nil
	}
	return struct {
		URL     string `yaml:"url"`
		Pattern string `yaml:"pattern"`
	}{l.URL, l.Pattern}, nil
}

// Category groups links under a name (environments, tools, docs).
type Category struct {
	Name  string
	Links map[string]Link
}

// Config is the top-level .surf-links.yml structure.
type Config struct {
	Type         *ProjectType    `yaml:"type,omitempty"`
	Environments map[string]Link `yaml:"environments,omitempty"`
	Tools        map[string]Link `yaml:"tools,omitempty"`
	Docs         map[string]Link `yaml:"docs,omitempty"`
}

// Categories returns the non-empty categories in display order.
func (c *Config) Categories() []Category {
	var cats []Category
	if len(c.Environments) > 0 {
		cats = append(cats, Category{Name: "environments", Links: c.Environments})
	}
	if len(c.Tools) > 0 {
		cats = append(cats, Category{Name: "tools", Links: c.Tools})
	}
	if len(c.Docs) > 0 {
		cats = append(cats, Category{Name: "docs", Links: c.Docs})
	}
	return cats
}

// AllLinks returns a flat map of all link names to their Link values.
// Generated type links are added first; explicit links override them.
func (c *Config) AllLinks() map[string]Link {
	all := make(map[string]Link)

	// Generated links first (so explicit links can override)
	if c.Type != nil {
		for k, v := range c.Type.GenerateLinks(c.Environments) {
			all[k] = v
		}
	}

	for k, v := range c.Environments {
		all[k] = v
	}
	for k, v := range c.Tools {
		all[k] = v
	}
	for k, v := range c.Docs {
		all[k] = v
	}
	return all
}

// Validate checks that the config has at least one link and all links have URLs.
func (c *Config) Validate() error {
	all := c.AllLinks()
	if len(all) == 0 {
		return fmt.Errorf("config has no links defined")
	}
	for name, link := range all {
		if link.URL == "" {
			return fmt.Errorf("link %q has no url", name)
		}
	}
	return nil
}
