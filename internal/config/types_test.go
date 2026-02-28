package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestProjectType_UnmarshalYAML_Standard(t *testing.T) {
	var pt ProjectType
	if err := yaml.Unmarshal([]byte(`wordpress`), &pt); err != nil {
		t.Fatal(err)
	}
	if pt.Name != "wordpress" {
		t.Errorf("got Name %q, want wordpress", pt.Name)
	}
	if pt.AdminPath != "/wp-admin" {
		t.Errorf("got AdminPath %q, want /wp-admin", pt.AdminPath)
	}
}

func TestProjectType_UnmarshalYAML_Bedrock(t *testing.T) {
	var pt ProjectType
	if err := yaml.Unmarshal([]byte(`wordpress-bedrock`), &pt); err != nil {
		t.Fatal(err)
	}
	if pt.AdminPath != "/wp/wp-admin" {
		t.Errorf("got AdminPath %q, want /wp/wp-admin", pt.AdminPath)
	}
}

func TestProjectType_UnmarshalYAML_Unknown(t *testing.T) {
	var pt ProjectType
	if err := yaml.Unmarshal([]byte(`unknown-type`), &pt); err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestProjectType_UnmarshalYAML_Custom(t *testing.T) {
	data := `
name: custom-cms
admin_path: /backend
`
	var pt ProjectType
	if err := yaml.Unmarshal([]byte(data), &pt); err != nil {
		t.Fatal(err)
	}
	if pt.Name != "custom-cms" {
		t.Errorf("got Name %q, want custom-cms", pt.Name)
	}
	if pt.AdminPath != "/backend" {
		t.Errorf("got AdminPath %q, want /backend", pt.AdminPath)
	}
}

func TestProjectType_UnmarshalYAML_CustomMissingPath(t *testing.T) {
	data := `name: custom-cms`
	var pt ProjectType
	if err := yaml.Unmarshal([]byte(data), &pt); err == nil {
		t.Error("expected error for missing admin_path")
	}
}

func TestProjectType_MarshalYAML_Standard(t *testing.T) {
	pt := ProjectType{Name: "wordpress", AdminPath: "/wp-admin"}
	val, err := pt.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	s, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}
	if s != "wordpress" {
		t.Errorf("got %q, want wordpress", s)
	}
}

func TestProjectType_MarshalYAML_Custom(t *testing.T) {
	pt := ProjectType{Name: "custom-cms", AdminPath: "/backend"}
	val, err := pt.MarshalYAML()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := val.(string); ok {
		t.Fatal("expected mapping for custom type, got scalar")
	}
}

func TestStandardTypeNames(t *testing.T) {
	names := StandardTypeNames()
	if len(names) == 0 {
		t.Fatal("expected non-empty standard type names")
	}
	// Verify sorted
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("not sorted: %q comes after %q", names[i], names[i-1])
		}
	}
}

func TestNewStandardType(t *testing.T) {
	pt, err := NewStandardType("wordpress")
	if err != nil {
		t.Fatal(err)
	}
	if pt.AdminPath != "/wp-admin" {
		t.Errorf("got %q, want /wp-admin", pt.AdminPath)
	}

	_, err = NewStandardType("nonexistent")
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestProjectType_GenerateLinks_MultiEnv(t *testing.T) {
	pt := &ProjectType{Name: "wordpress", AdminPath: "/wp-admin"}
	envs := map[string]Link{
		"production": {URL: "https://example.com"},
		"staging":    {URL: "https://staging.example.com"},
		"local":      {URL: "https://local.example.com"},
	}

	links := pt.GenerateLinks(envs)

	// Should have default "admin" + per-env links
	if len(links) != 4 {
		t.Fatalf("got %d links, want 4", len(links))
	}

	// Default admin → first alphabetically (local)
	if links["admin"].URL != "https://local.example.com/wp-admin" {
		t.Errorf("admin URL = %q, want local env", links["admin"].URL)
	}

	if links["admin production"].URL != "https://example.com/wp-admin" {
		t.Errorf("admin production URL = %q", links["admin production"].URL)
	}

	if links["admin staging"].URL != "https://staging.example.com/wp-admin" {
		t.Errorf("admin staging URL = %q", links["admin staging"].URL)
	}
}

func TestProjectType_GenerateLinks_SingleEnv(t *testing.T) {
	pt := &ProjectType{Name: "wordpress", AdminPath: "/wp-admin"}
	envs := map[string]Link{
		"production": {URL: "https://example.com"},
	}

	links := pt.GenerateLinks(envs)

	if len(links) != 2 {
		t.Fatalf("got %d links, want 2", len(links))
	}

	if links["admin"].URL != "https://example.com/wp-admin" {
		t.Errorf("admin URL = %q", links["admin"].URL)
	}
}

func TestProjectType_GenerateLinks_Nil(t *testing.T) {
	var pt *ProjectType
	links := pt.GenerateLinks(map[string]Link{"x": {URL: "http://x"}})
	if links != nil {
		t.Error("expected nil for nil ProjectType")
	}
}
