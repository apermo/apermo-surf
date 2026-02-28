package wizard

import (
	"bytes"
	"strings"
	"testing"
)

func TestWizard_FullFlow(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",               // skip name
		"",               // skip project type
		"prod",           // env name
		"https://example.com", // env url
		"",               // finish envs
		"",               // finish tools
		"",               // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environments["prod"].URL != "https://example.com" {
		t.Errorf("prod URL = %q", cfg.Environments["prod"].URL)
	}
	if cfg.Type != nil {
		t.Error("expected nil project type")
	}
	if cfg.Tools != nil {
		t.Error("expected nil tools")
	}
	if cfg.Name != "" {
		t.Errorf("expected empty name, got %q", cfg.Name)
	}
}

func TestWizard_WithName(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"My Project",     // name
		"",               // skip type
		"prod",           // env
		"https://example.com",
		"",               // finish envs
		"",               // finish tools
		"",               // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "My Project" {
		t.Errorf("name = %q, want My Project", cfg.Name)
	}
}

func TestWizard_WithProjectType(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",               // skip name
		"wordpress",
		"prod",
		"https://example.com",
		"",
		"",
		"",
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type == nil || cfg.Type.Name != "wordpress" {
		t.Error("expected wordpress project type")
	}
}

func TestWizard_WithToolAndPattern(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",                    // skip name
		"",                    // skip type
		"prod",                // env
		"https://example.com", // env url
		"",                    // finish envs
		"jira",                // tool name
		"https://jira.example.com/browse/{ticket}", // tool url
		`PROJ-\d+`,           // pattern
		"",                    // finish tools
		"",                    // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tools["jira"].Pattern != `PROJ-\d+` {
		t.Errorf("pattern = %q", cfg.Tools["jira"].Pattern)
	}
}

func TestWizard_InvalidURLRecovery(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",                    // skip name
		"",                    // skip type
		"prod",                // env
		"not-a-url",           // invalid
		"https://example.com", // valid retry
		"",                    // finish envs
		"",                    // finish tools
		"",                    // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environments["prod"].URL != "https://example.com" {
		t.Errorf("prod URL = %q", cfg.Environments["prod"].URL)
	}
	if !strings.Contains(out.String(), "invalid URL") {
		t.Error("expected invalid URL message in output")
	}
}

func TestWizard_InvalidPatternRecovery(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",                    // skip name
		"",                    // skip type
		"prod",                // env
		"https://example.com", // env url
		"",                    // finish envs
		"jira",                // tool name
		"https://jira.example.com", // tool url
		"[invalid",            // bad pattern
		`\d+`,                 // valid retry
		"",                    // finish tools
		"",                    // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Tools["jira"].Pattern != `\d+` {
		t.Errorf("pattern = %q", cfg.Tools["jira"].Pattern)
	}
	if !strings.Contains(out.String(), "invalid regex") {
		t.Error("expected invalid regex message")
	}
}

func TestWizard_EmptyEnvsError(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",                    // skip name
		"",                    // skip type
		"",                    // try to finish envs (rejected)
		"prod",                // env name
		"https://example.com", // env url
		"",                    // finish envs
		"",                    // finish tools
		"",                    // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environments["prod"].URL != "https://example.com" {
		t.Errorf("prod URL = %q", cfg.Environments["prod"].URL)
	}
	if !strings.Contains(out.String(), "at least one environment") {
		t.Error("expected 'at least one environment' message")
	}
}

func TestWizard_AskDist_Default(t *testing.T) {
	input := strings.NewReader("\n")
	var out bytes.Buffer
	w := New(input, &out)

	dist, err := w.AskDist()
	if err != nil {
		t.Fatal(err)
	}
	if dist {
		t.Error("expected local (false) as default")
	}
}

func TestWizard_AskDist_Explicit(t *testing.T) {
	input := strings.NewReader("dist\n")
	var out bytes.Buffer
	w := New(input, &out)

	dist, err := w.AskDist()
	if err != nil {
		t.Fatal(err)
	}
	if !dist {
		t.Error("expected dist (true)")
	}
}

func TestWizard_WithDocs(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		"",                       // skip name
		"",                       // skip type
		"prod",                   // env
		"https://example.com",    // env url
		"",                       // finish envs
		"",                       // finish tools
		"wiki",                   // doc name
		"https://wiki.example.com", // doc url
		"",                       // finish docs
	}, "\n") + "\n")

	var out bytes.Buffer
	w := New(input, &out)

	cfg, err := w.Run()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Docs["wiki"].URL != "https://wiki.example.com" {
		t.Errorf("wiki URL = %q", cfg.Docs["wiki"].URL)
	}
}
