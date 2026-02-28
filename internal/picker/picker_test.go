package picker

import (
	"bytes"
	"strings"
	"testing"
)

func TestPickWithList_DefaultChoice(t *testing.T) {
	names := []string{"production", "staging", "local"}
	urls := []string{"https://example.com", "https://staging.example.com", "https://local.example.com"}

	idx, err := pickWithList(names, urls, strings.NewReader("\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if idx != 0 {
		t.Errorf("got index %d, want 0", idx)
	}
}

func TestPickWithList_ExplicitChoice(t *testing.T) {
	names := []string{"production", "staging"}
	urls := []string{"https://example.com", "https://staging.example.com"}

	idx, err := pickWithList(names, urls, strings.NewReader("2\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 {
		t.Errorf("got index %d, want 1", idx)
	}
}

func TestPickWithList_OutOfRange(t *testing.T) {
	names := []string{"production"}
	urls := []string{"https://example.com"}

	_, err := pickWithList(names, urls, strings.NewReader("5\n"), &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for out-of-range choice")
	}
}

func TestPickWithList_InvalidInput(t *testing.T) {
	names := []string{"production"}
	urls := []string{"https://example.com"}

	_, err := pickWithList(names, urls, strings.NewReader("abc\n"), &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for non-numeric input")
	}
}

func TestPickWithList_Output(t *testing.T) {
	names := []string{"prod", "staging"}
	urls := []string{"https://example.com", "https://staging.example.com"}

	var buf bytes.Buffer
	_, err := pickWithList(names, urls, strings.NewReader("1\n"), &buf)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "prod") {
		t.Error("output should contain 'prod'")
	}
	if !strings.Contains(output, "staging") {
		t.Error("output should contain 'staging'")
	}
	if !strings.Contains(output, "Pick a link") {
		t.Error("output should contain prompt")
	}
}
