package resolve

import "testing"

func FuzzExtractLiteralPrefix(f *testing.F) {
	f.Add(`PROJ-\d+`)
	f.Add(`[A-Z]+-\d+`)
	f.Add(``)
	f.Add(`hello`)
	f.Add(`^start`)
	f.Add(`(group)`)
	f.Add(`a.b.c`)
	f.Add(`prefix\suffix`)
	f.Fuzz(func(t *testing.T, pattern string) {
		extractLiteralPrefix(pattern)
	})
}

func FuzzStripPlaceholderSegment(f *testing.F) {
	f.Add("https://example.com/browse/{ticket}")
	f.Add("https://example.com")
	f.Add("https://example.com/{a}/{b}")
	f.Add("")
	f.Add("not-a-url")
	f.Add("https://example.com/path/{ticket}/sub")
	f.Add("://broken")
	f.Fuzz(func(t *testing.T, rawURL string) {
		stripPlaceholderSegment(rawURL)
	})
}

func FuzzResolveExplicitArg(f *testing.F) {
	f.Add("PROJ-123", `PROJ-\d+`)
	f.Add("123", `PROJ-\d+`)
	f.Add("", "")
	f.Add("123", "")
	f.Add("abc", `PROJ-\d+`)
	f.Add("42", `[A-Z]+-\d+`)
	f.Add("PROJ-123", "")
	f.Fuzz(func(t *testing.T, arg, pattern string) {
		resolveExplicitArg(arg, pattern)
	})
}
