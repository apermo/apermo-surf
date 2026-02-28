package config

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// standardTypes maps project type names to their admin paths.
var standardTypes = map[string]string{
	"craft":             "/admin",
	"drupal":            "/admin",
	"laravel":           "/admin",
	"magento":           "/admin",
	"shopware":          "/admin",
	"typo3":             "/typo3",
	"wordpress":         "/wp-admin",
	"wordpress-bedrock": "/wp/wp-admin",
}

// ProjectType represents a project type that generates admin links.
type ProjectType struct {
	Name      string
	AdminPath string
}

// UnmarshalYAML supports both a standard type name (string) and a custom
// mapping with an admin_path field.
func (pt *ProjectType) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		name := value.Value
		path, ok := standardTypes[name]
		if !ok {
			return fmt.Errorf("unknown project type %q", name)
		}
		pt.Name = name
		pt.AdminPath = path
		return nil
	}

	// Custom type: {name: ..., admin_path: ...}
	var raw struct {
		Name      string `yaml:"name"`
		AdminPath string `yaml:"admin_path"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	if raw.AdminPath == "" {
		return fmt.Errorf("custom project type requires admin_path")
	}
	pt.Name = raw.Name
	pt.AdminPath = raw.AdminPath
	return nil
}

// StandardTypeNames returns a sorted list of all standard type names.
func StandardTypeNames() []string {
	names := make([]string, 0, len(standardTypes))
	for name := range standardTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// NewStandardType creates a ProjectType from a known standard name.
func NewStandardType(name string) (*ProjectType, error) {
	path, ok := standardTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown project type %q", name)
	}
	return &ProjectType{Name: name, AdminPath: path}, nil
}

// MarshalYAML writes a standard type as a scalar string,
// or a custom type as a mapping with name and admin_path.
func (pt ProjectType) MarshalYAML() (interface{}, error) {
	if _, ok := standardTypes[pt.Name]; ok {
		return pt.Name, nil
	}
	return struct {
		Name      string `yaml:"name"`
		AdminPath string `yaml:"admin_path"`
	}{pt.Name, pt.AdminPath}, nil
}

// GenerateLinks creates admin links for each environment.
// Returns a map of link names to Links:
//   - "admin" → default environment (first alphabetically)
//   - "admin <env>" → per-environment admin links
func (pt *ProjectType) GenerateLinks(environments map[string]Link) map[string]Link {
	if pt == nil || len(environments) == 0 {
		return nil
	}

	links := make(map[string]Link)

	// Sort environment names for deterministic default
	envNames := make([]string, 0, len(environments))
	for name := range environments {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	for _, name := range envNames {
		env := environments[name]
		adminURL := env.URL + pt.AdminPath
		links["admin "+name] = Link{URL: adminURL}
	}

	// Default "admin" → first environment alphabetically
	if len(envNames) > 0 {
		defaultEnv := environments[envNames[0]]
		links["admin"] = Link{URL: defaultEnv.URL + pt.AdminPath}
	}

	return links
}
