package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Link represents a project URL, either as a simple string or with a pattern.
// Links is an optional map of sub-link names to relative paths.
type Link struct {
	URL     string            `yaml:"url"`
	Pattern string            `yaml:"pattern,omitempty"`
	Links   map[string]string `yaml:"links,omitempty"`
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

// MarshalYAML writes a Link as a scalar string when there is no pattern
// or sub-links, or as a mapping otherwise.
func (l Link) MarshalYAML() (interface{}, error) {
	if l.Pattern == "" && len(l.Links) == 0 {
		return l.URL, nil
	}
	return struct {
		URL     string            `yaml:"url"`
		Pattern string            `yaml:"pattern,omitempty"`
		Links   map[string]string `yaml:"links,omitempty"`
	}{l.URL, l.Pattern, l.Links}, nil
}

// Category groups links under a name (environments, tools, docs).
type Category struct {
	Name  string
	Links map[string]Link
}

// Config is the top-level .surf-links.yml structure.
type Config struct {
	Name         string          `yaml:"name,omitempty"`
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
// Sub-links are expanded into compound names (e.g. "jira board").
func (c *Config) AllLinks() map[string]Link {
	all := make(map[string]Link)

	// Generated links first (so explicit links can override)
	if c.Type != nil {
		for k, v := range c.Type.GenerateLinks(c.Environments) {
			all[k] = v
		}
	}

	for _, links := range []map[string]Link{c.Environments, c.Tools, c.Docs} {
		for k, v := range links {
			all[k] = v
			for sub, path := range v.Links {
				subURL := strings.TrimRight(v.URL, "/") + path
				all[k+" "+sub] = Link{URL: subURL}
			}
		}
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
