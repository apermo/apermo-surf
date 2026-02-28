package resolve

import "testing"

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
