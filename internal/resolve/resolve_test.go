package resolve

import "testing"

func TestResolveExplicitArg_AutoPrefix(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		pattern string
		want    string
	}{
		{
			name:    "bare number gets prefixed",
			arg:     "123",
			pattern: `PROJ-\d+`,
			want:    "PROJ-123",
		},
		{
			name:    "already prefixed unchanged",
			arg:     "PROJ-123",
			pattern: `PROJ-\d+`,
			want:    "PROJ-123",
		},
		{
			name:    "non-numeric arg unchanged",
			arg:     "feature-branch",
			pattern: `PROJ-\d+`,
			want:    "feature-branch",
		},
		{
			name:    "no pattern returns arg as-is",
			arg:     "123",
			pattern: "",
			want:    "123",
		},
		{
			name:    "pattern without literal prefix",
			arg:     "123",
			pattern: `\d+`,
			want:    "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveExplicitArg(tt.arg, tt.pattern)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractLiteralPrefix(t *testing.T) {
	tests := []struct {
		pattern string
		want    string
	}{
		{`PROJ-\d+`, "PROJ-"},
		{`\d+`, ""},
		{`ABC`, "ABC"},
		{`FOO[0-9]+`, "FOO"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got := extractLiteralPrefix(tt.pattern)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripPlaceholderSegment(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "removes trailing placeholder",
			url:  "https://jira.com/browse/{ticket}",
			want: "https://jira.com/browse",
		},
		{
			name: "removes multiple placeholders",
			url:  "https://example.com/{repo}/browse/{ticket}",
			want: "https://example.com/browse",
		},
		{
			name: "no placeholders unchanged",
			url:  "https://example.com/browse/PROJ-123",
			want: "https://example.com/browse/PROJ-123",
		},
		{
			name: "all placeholders removed",
			url:  "https://example.com/{ticket}",
			want: "https://example.com",
		},
		{
			name: "preserves query string",
			url:  "https://example.com/browse/{ticket}?view=board",
			want: "https://example.com/browse?view=board",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripPlaceholderSegment(tt.url)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
