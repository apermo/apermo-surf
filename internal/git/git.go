package git

import (
	"os/exec"
	"regexp"
	"strings"
)

// Branch returns the current git branch name.
// Returns ("", nil) when not in a git repository.
func Branch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}

// Ticket extracts a ticket identifier from a branch name using the given regex pattern.
// Returns ("", nil) when the pattern doesn't match or inputs are empty.
func Ticket(branch, pattern string) (string, error) {
	if branch == "" || pattern == "" {
		return "", nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.FindString(branch), nil
}

// Repo returns the repository name derived from the git remote URL.
// Returns ("", nil) when not in a git repository or no remote is configured.
func Repo(dir string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return repoNameFromURL(strings.TrimSpace(string(out))), nil
}

// repoNameFromURL extracts the repository name from a git remote URL.
func repoNameFromURL(raw string) string {
	raw = strings.TrimSuffix(raw, ".git")

	// SSH format: git@github.com:user/repo
	if i := strings.LastIndex(raw, ":"); i != -1 && !strings.Contains(raw, "://") {
		raw = raw[i+1:]
	}

	// Take last path component
	if i := strings.LastIndex(raw, "/"); i != -1 {
		raw = raw[i+1:]
	}

	return raw
}
