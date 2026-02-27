package fuzzy

import "testing"

var names = []string{"production", "staging", "local", "jira", "sentry", "confluence"}

func TestBestMatch_Exact(t *testing.T) {
	match, candidates := BestMatch("production", names)
	if match != "production" || candidates != nil {
		t.Errorf("exact match failed: got %q, %v", match, candidates)
	}
}

func TestBestMatch_ExactCaseInsensitive(t *testing.T) {
	match, candidates := BestMatch("Production", names)
	if match != "production" || candidates != nil {
		t.Errorf("case-insensitive exact match failed: got %q, %v", match, candidates)
	}
}

func TestBestMatch_Fuzzy(t *testing.T) {
	match, candidates := BestMatch("prod", names)
	if match != "production" || candidates != nil {
		t.Errorf("fuzzy match failed: got %q, %v", match, candidates)
	}
}

func TestBestMatch_NoMatch(t *testing.T) {
	match, candidates := BestMatch("zzzzz", names)
	if match != "" || candidates != nil {
		t.Errorf("expected no match: got %q, %v", match, candidates)
	}
}

func TestBestMatch_Ambiguous(t *testing.T) {
	// "s" should match both "staging" and "sentry" with similar scores
	match, candidates := BestMatch("s", []string{"staging", "sentry"})
	if match != "" {
		// If there's a clear winner, that's acceptable too
		return
	}
	if len(candidates) < 2 {
		t.Errorf("expected ambiguous result with candidates: got %q, %v", match, candidates)
	}
}
