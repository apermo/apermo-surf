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

# List all links
surf links
surf links --env        # environments only
surf links --tools      # tools only

# Interactive picker (fzf integration)
surf open               # no args — interactive selection

# Push config to Chrome extension
surf add-to-chrome
```

## Config

Create a `.surf-links.yml` in your project root:

```yaml
environments:
  local: https://myproject.ddev.site
  staging: https://staging.example.com
  production: https://example.com

tools:
  jira:
    url: https://myorg.atlassian.net/browse/{ticket}
    pattern: "PROJ-\\d+"
  sentry: https://sentry.io/organizations/myorg/projects/myproject
  ci: https://ci.example.com/myorg/myproject

docs:
  confluence: https://myorg.atlassian.net/wiki/spaces/PROJ
  figma: https://figma.com/file/abc123
```

The config is discovered by walking up from cwd (like `.env` or `.git`).

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
- [Cobra](https://github.com/spf13/cobra) or [Kong](https://github.com/alecthomas/kong) for CLI
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) for config parsing
- [GoReleaser](https://goreleaser.com/) for builds and Homebrew distribution

## License

MIT
