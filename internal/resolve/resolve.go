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
func Resolve(link config.Link, configDir string) Result {
	rawURL := link.URL

	if !strings.Contains(rawURL, "{") {
		return Result{URL: rawURL}
	}

	var warnings []string

	branch, _ := git.Branch(configDir)
	repo, _ := git.Repo(configDir)
	ticket, _ := git.Ticket(branch, link.Pattern)

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
