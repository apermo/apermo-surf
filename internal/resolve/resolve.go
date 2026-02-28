package resolve

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/apermo/apermo-surf/internal/config"
	"github.com/apermo/apermo-surf/internal/git"
)

// Result holds a resolved URL and any warnings generated during resolution.
type Result struct {
	URL      string
	Warnings []string
}

// Resolve replaces placeholders in a link's URL with git-derived values.
// configDir is the directory containing .surf-links.yml (used as git context).
// explicitArg overrides {ticket} when non-empty (resolution: explicit → branch → fallback).
func Resolve(link config.Link, configDir string, explicitArg string) Result {
	rawURL := link.URL

	if !strings.Contains(rawURL, "{") {
		return Result{URL: rawURL}
	}

	var warnings []string

	branch, _ := git.Branch(configDir)
	repo, _ := git.Repo(configDir)

	// Ticket resolution: explicit arg → branch extraction → empty
	var ticket string
	if explicitArg != "" {
		ticket = resolveExplicitArg(explicitArg, link.Pattern)
	} else {
		ticket, _ = git.Ticket(branch, link.Pattern)
	}

	replacements := map[string]string{
		"{branch}": branch,
		"{repo}":   repo,
		"{ticket}": ticket,
	}

	for placeholder, value := range replacements {
		if !strings.Contains(rawURL, placeholder) {
			continue
		}
		if value == "" {
			warnings = append(warnings, fmt.Sprintf("could not resolve %s", placeholder))
		} else {
			rawURL = strings.ReplaceAll(rawURL, placeholder, value)
		}
	}

	rawURL = stripPlaceholderSegment(rawURL)

	return Result{URL: rawURL, Warnings: warnings}
}

// resolveExplicitArg applies auto-prefix logic to the explicit ticket argument.
// If arg is a bare number and the pattern has a literal prefix, it prepends the prefix.
func resolveExplicitArg(arg, pattern string) string {
	if pattern == "" {
		return arg
	}

	prefix := extractLiteralPrefix(pattern)
	if prefix == "" {
		return arg
	}

	// Already has the prefix — use as-is
	if strings.HasPrefix(arg, prefix) {
		return arg
	}

	// Bare number → auto-prefix
	for _, r := range arg {
		if r < '0' || r > '9' {
			return arg
		}
	}
	return prefix + arg
}

// extractLiteralPrefix returns the literal prefix of a regex pattern,
// stopping at the first metacharacter.
func extractLiteralPrefix(pattern string) string {
	meta := `\[(.+?{^$|`
	var prefix strings.Builder
	for _, r := range pattern {
		if strings.ContainsRune(meta, r) {
			break
		}
		prefix.WriteRune(r)
	}
	return prefix.String()
}

// stripPlaceholderSegment removes path segments containing unresolved {…} placeholders.
func stripPlaceholderSegment(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	parts := strings.Split(u.Path, "/")
	var cleaned []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.Contains(p, "{") && strings.Contains(p, "}") {
			continue
		}
		cleaned = append(cleaned, p)
	}

	if len(cleaned) == 0 {
		u.Path = ""
	} else {
		u.Path = "/" + strings.Join(cleaned, "/")
	}

	return u.String()
}
