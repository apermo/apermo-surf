package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLink_UnmarshalYAML_Simple(t *testing.T) {
	cfg, err := parseYAML(t, `
environments:
  prod: https://example.com
`)
	if err != nil {
		t.Fatal(err)
	}
	link := cfg.Environments["prod"]
	if link.URL != "https://example.com" {
		t.Errorf("got URL %q, want https://example.com", link.URL)
	}
	if link.Pattern != "" {
		t.Errorf("got Pattern %q, want empty", link.Pattern)
	}
}

func TestLink_UnmarshalYAML_Expanded(t *testing.T) {
	cfg, err := parseYAML(t, `
tools:
  jira:
    url: https://jira.example.com/browse/{ticket}
    pattern: "PROJ-\\d+"
`)
	if err != nil {
		t.Fatal(err)
	}
	link := cfg.Tools["jira"]
	if link.URL != "https://jira.example.com/browse/{ticket}" {
		t.Errorf("got URL %q", link.URL)
	}
	if link.Pattern != `PROJ-\d+` {
		t.Errorf("got Pattern %q", link.Pattern)
	}
}

func TestConfig_Categories(t *testing.T) {
	cfg, err := parseYAML(t, `
environments:
  prod: https://example.com
docs:
  wiki: https://wiki.example.com
`)
	if err != nil {
		t.Fatal(err)
	}
	cats := cfg.Categories()
	if len(cats) != 2 {
		t.Fatalf("got %d categories, want 2", len(cats))
	}
	if cats[0].Name != "environments" {
		t.Errorf("first category = %q, want environments", cats[0].Name)
	}
	if cats[1].Name != "docs" {
		t.Errorf("second category = %q, want docs", cats[1].Name)
	}
}

func TestConfig_AllLinks(t *testing.T) {
	cfg, err := parseYAML(t, `
environments:
  prod: https://example.com
tools:
  ci: https://ci.example.com
docs:
  wiki: https://wiki.example.com
`)
	if err != nil {
		t.Fatal(err)
	}
	all := cfg.AllLinks()
	if len(all) != 3 {
		t.Fatalf("got %d links, want 3", len(all))
	}
}

func TestConfig_Validate_Empty(t *testing.T) {
	cfg := &Config{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty config")
	}
}

func TestConfig_Validate_MissingURL(t *testing.T) {
	cfg := &Config{
		Environments: map[string]Link{
			"prod": {URL: ""},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestFind_WalksUp(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dir, FileName)
	if err := os.WriteFile(configPath, []byte("environments:\n  x: http://x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	found, err := Find(nested)
	if err != nil {
		t.Fatal(err)
	}
	if found != configPath {
		t.Errorf("found %q, want %q", found, configPath)
	}
}

func TestFind_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := Find(dir)
	if err == nil {
		t.Error("expected error when no config exists")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	data := `
environments:
  prod: https://example.com
tools:
  jira:
    url: https://jira.example.com/browse/{ticket}
    pattern: "PROJ-\\d+"
`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environments["prod"].URL != "https://example.com" {
		t.Error("prod URL mismatch")
	}
	if cfg.Tools["jira"].Pattern != `PROJ-\d+` {
		t.Error("jira pattern mismatch")
	}
}

func parseYAML(t *testing.T, input string) (*Config, error) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		return nil, err
	}
	return Load(path)
}
