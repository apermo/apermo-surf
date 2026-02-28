# Apermo Surf — CLI

A CLI tool for project-contextual link navigation. Open project URLs (environments, tools, docs) from the terminal with fuzzy matching, branch-aware placeholders, and interactive selection.

**Part of the [Apermo Surf](https://github.com/apermo) ecosystem:**
- **[apermo-surf](https://github.com/apermo/apermo-surf)** — CLI tool (this repo)
- **[apermo-surf-chrome](https://github.com/apermo/apermo-surf-chrome)** — Chrome extension

Both tools share the `.surf-links.yml` config format. Maintain your project links once, use them everywhere.

## Usage

```bash
# Open a link by fuzzy name
surf open prod          # opens production URL
surf open jira          # opens Jira (current ticket from branch)
surf open sentry        # opens Sentry

# Open a specific ticket
surf open jira 123      # opens PROJ-123 (auto-prefixes from pattern)
surf open jira PROJ-456 # opens PROJ-456 as-is

# Sub-links via compound names
surf open "jira board"  # opens Jira board sub-link

# Choose browser
surf open prod -b firefox

# List all links
surf links
surf links --env        # environments only
surf links --tools      # tools only

# Interactive picker (fzf integration)
surf open               # no args — interactive selection

# Create a new config
surf init               # interactive wizard
surf init --dist        # create .surf-links.yml.dist (shared template)

# Push config to Chrome extension
surf add-to-chrome
```

## Config

Create a `.surf-links.yml` in your project root (or run `surf init`):

```yaml
name: My Project

type: wordpress

environments:
  local: https://myproject.ddev.site
  staging: https://staging.example.com
  production: https://example.com

tools:
  jira:
    url: https://myorg.atlassian.net/browse/{ticket}
    pattern: "PROJ-\\d+"
    links:
      board: /boards/1
      backlog: /backlog
  github:
    url: https://github.com/myorg/myproject
    links:
      prs: /pulls
      issues: /issues
      actions: /actions
  sentry: https://sentry.io/organizations/myorg/projects/myproject

docs:
  confluence: https://myorg.atlassian.net/wiki/spaces/PROJ
  figma: https://figma.com/file/abc123
```

The config is discovered by walking up from cwd (like `.env` or `.git`).
A `.surf-links.yml.dist` file is used as fallback for team-shared templates.

- **`name`** — optional project display name
- **`type`** — standard CMS type (wordpress, typo3, laravel, drupal, shopware, magento, craft) auto-generates admin links per environment
- **`links`** — optional sub-links with paths relative to the parent URL

### Placeholders

| Placeholder | Source |
|-------------|--------|
| `{ticket}` | Extracted from git branch name using `pattern` |
| `{branch}` | Current git branch name |
| `{repo}` | Repository name from git remote |

### Ticket resolution order

1. **Explicit argument** — `surf open jira 123` → `PROJ-123`
2. **Branch name** — extract from current branch using `pattern`
3. **Fallback** — open the base URL

## Installation

```bash
brew tap apermo/tap
brew install apermo-surf
```

The binary is installed as `apermo-surf` with a `surf` symlink.

### Shell Completions

Enable tab completion for link names (one-time setup):

```bash
# Zsh (add to ~/.zshrc)
source <(surf completion zsh)

# Bash (add to ~/.bashrc)
source <(surf completion bash)

# Fish
surf completion fish | source
# To make persistent:
surf completion fish > ~/.config/fish/completions/surf.fish
```

After reloading your shell, `surf open <TAB>` will suggest link names from the current project's `.surf-links.yml`.

## Tech Stack

- Go
- [Cobra](https://github.com/spf13/cobra) for CLI
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) for config parsing
- [GoReleaser](https://goreleaser.com/) for builds and Homebrew distribution

## AI Disclaimer

This project is developed with major assistance from [Claude Code](https://docs.anthropic.com/en/docs/claude-code) (Anthropic).
Claude handles the bulk of the implementation — writing Go code, tests, CI workflows, and documentation — while
the maintainer reviews, steers, and makes final decisions. Projects with stricter rules regarding the use of AI-generated
code should refrain from forking or reusing code from this repository.

## License

MIT
