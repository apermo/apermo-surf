package userconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())

	cfg := Load()
	if cfg.Browser != "" {
		t.Errorf("expected empty browser, got %q", cfg.Browser)
	}
	if cfg.Browsers != nil {
		t.Error("expected nil browsers map")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "surf")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `
browser: firefox
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()
	if cfg.Browser != "firefox" {
		t.Errorf("expected browser=firefox, got %q", cfg.Browser)
	}
}

func TestLoad_CustomBrowsers(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "surf")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := `
browser: my-browser
browsers:
  my-browser:
    command: /usr/bin/custom-browser
    args: ["--new-tab"]
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()
	if cfg.Browser != "my-browser" {
		t.Errorf("expected browser=my-browser, got %q", cfg.Browser)
	}
	bc, ok := cfg.Browsers["my-browser"]
	if !ok {
		t.Fatal("expected my-browser in browsers map")
	}
	if bc.Command != "/usr/bin/custom-browser" {
		t.Errorf("expected command=/usr/bin/custom-browser, got %q", bc.Command)
	}
	if len(bc.Args) != 1 || bc.Args[0] != "--new-tab" {
		t.Errorf("unexpected args: %v", bc.Args)
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "surf")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(configDir, "config.yml"), []byte("{{invalid"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := Load()
	if cfg.Browser != "" {
		t.Errorf("expected empty config for malformed YAML, got browser=%q", cfg.Browser)
	}
}
