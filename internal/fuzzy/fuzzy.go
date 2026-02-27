package fuzzy

import (
	"strings"

	fuzzypkg "github.com/sahilm/fuzzy"
)

// BestMatch returns the best fuzzy match for pattern among names.
// Returns (match, nil) for a clear winner, ("", candidates) for ambiguous results,
// or ("", nil) for no match.
func BestMatch(pattern string, names []string) (string, []string) {
	// Exact match always wins
	lower := strings.ToLower(pattern)
	for _, name := range names {
		if strings.ToLower(name) == lower {
			return name, nil
		}
	}

	matches := fuzzypkg.Find(pattern, names)
	if len(matches) == 0 {
		return "", nil
	}

	if len(matches) == 1 {
		return matches[0].Str, nil
	}

	// Clear score gap → unambiguous best
	best := matches[0]
	second := matches[1]
	if best.Score > second.Score {
		return best.Str, nil
	}

	// Tied scores → ambiguous
	candidates := make([]string, len(matches))
	for i, m := range matches {
		candidates[i] = m.Str
	}
	return "", candidates
}
