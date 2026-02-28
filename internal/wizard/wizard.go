package wizard

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/apermo/apermo-surf/internal/config"
)

// Wizard guides the user through creating a .surf-links.yml config.
type Wizard struct {
	scanner *bufio.Scanner
	out     io.Writer
}

// New creates a wizard that reads from in and writes prompts to out.
func New(in io.Reader, out io.Writer) *Wizard {
	return &Wizard{
		scanner: bufio.NewScanner(in),
		out:     out,
	}
}

// AskDist asks whether to create a .dist file.
// Returns true for .dist, false for regular config.
func (w *Wizard) AskDist() (bool, error) {
	fmt.Fprint(w.out, "Create .surf-links.yml.dist (shared template) or .surf-links.yml (local)? [dist/local] (local): ")
	line, err := w.readLine()
	if err != nil {
		return false, err
	}
	return strings.EqualFold(line, "dist"), nil
}

// Run walks the user through building a Config.
func (w *Wizard) Run() (*config.Config, error) {
	cfg := &config.Config{}

	// Project type (optional)
	pt, err := w.askProjectType()
	if err != nil {
		return nil, err
	}
	cfg.Type = pt

	// Environments (at least one required)
	envs, err := w.askLinks("environment", true)
	if err != nil {
		return nil, err
	}
	cfg.Environments = envs

	// Tools (optional, with pattern support)
	fmt.Fprintln(w.out, "\nTools (optional, press Enter to skip):")
	tools, err := w.askLinksWithPattern("tool")
	if err != nil {
		return nil, err
	}
	cfg.Tools = tools

	// Docs (optional)
	fmt.Fprintln(w.out, "\nDocs (optional, press Enter to skip):")
	docs, err := w.askLinks("doc", false)
	if err != nil {
		return nil, err
	}
	cfg.Docs = docs

	return cfg, nil
}

func (w *Wizard) askProjectType() (*config.ProjectType, error) {
	names := config.StandardTypeNames()
	fmt.Fprintf(w.out, "Project type? [%s] (skip): ", strings.Join(names, ", "))
	line, err := w.readLine()
	if err != nil {
		return nil, err
	}
	if line == "" {
		return nil, nil
	}

	pt, err := config.NewStandardType(line)
	if err != nil {
		fmt.Fprintf(w.out, "  unknown type %q, skipping\n", line)
		return nil, nil
	}
	return pt, nil
}

func (w *Wizard) askLinks(label string, required bool) (map[string]config.Link, error) {
	links := make(map[string]config.Link)
	if required {
		fmt.Fprintf(w.out, "\nAdd %ss (at least one required):\n", label)
	}

	for {
		fmt.Fprintf(w.out, "  %s name (Enter to finish): ", label)
		name, err := w.readLine()
		if err != nil {
			return nil, err
		}
		if name == "" {
			if required && len(links) == 0 {
				fmt.Fprintln(w.out, "  at least one environment is required")
				continue
			}
			break
		}

		u, err := w.askURL()
		if err != nil {
			return nil, err
		}
		links[name] = config.Link{URL: u}
	}

	if len(links) == 0 {
		return nil, nil
	}
	return links, nil
}

func (w *Wizard) askLinksWithPattern(label string) (map[string]config.Link, error) {
	links := make(map[string]config.Link)

	for {
		fmt.Fprintf(w.out, "  %s name (Enter to finish): ", label)
		name, err := w.readLine()
		if err != nil {
			return nil, err
		}
		if name == "" {
			break
		}

		u, err := w.askURL()
		if err != nil {
			return nil, err
		}

		fmt.Fprint(w.out, "    pattern (Enter to skip): ")
		pattern, err := w.readLine()
		if err != nil {
			return nil, err
		}

		if pattern != "" {
			for !validatePattern(pattern) {
				fmt.Fprintln(w.out, "    invalid regex, try again")
				fmt.Fprint(w.out, "    pattern (Enter to skip): ")
				pattern, err = w.readLine()
				if err != nil {
					return nil, err
				}
				if pattern == "" {
					break
				}
			}
		}

		link := config.Link{URL: u}
		if pattern != "" {
			link.Pattern = pattern
		}
		links[name] = link
	}

	if len(links) == 0 {
		return nil, nil
	}
	return links, nil
}

func (w *Wizard) askURL() (string, error) {
	for {
		fmt.Fprint(w.out, "    url: ")
		u, err := w.readLine()
		if err != nil {
			return "", err
		}
		if validateURL(u) {
			return u, nil
		}
		fmt.Fprintln(w.out, "    invalid URL (must start with http:// or https://)")
	}
}

func (w *Wizard) readLine() (string, error) {
	if !w.scanner.Scan() {
		if err := w.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(w.scanner.Text()), nil
}

func validateURL(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func validatePattern(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}
