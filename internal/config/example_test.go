package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ExampleLoad() {
	dir, _ := os.MkdirTemp("", "surf-example-*")
	defer os.RemoveAll(dir)

	data := `environments:
  prod: https://example.com
`
	path := filepath.Join(dir, ".surf-links.yml")
	os.WriteFile(path, []byte(data), 0o644)

	cfg, _ := Load(path)
	fmt.Println(cfg.Environments["prod"].URL)
	// Output: https://example.com
}

func ExampleConfig_AllLinks() {
	cfg := &Config{
		Type: &ProjectType{Name: "wordpress", AdminPath: "/wp-admin"},
		Environments: map[string]Link{
			"prod": {URL: "https://example.com"},
		},
	}
	all := cfg.AllLinks()

	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("%s: %s\n", name, all[name].URL)
	}
	// Output:
	// admin: https://example.com/wp-admin
	// admin prod: https://example.com/wp-admin
	// prod: https://example.com
}
