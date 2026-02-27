# CLAUDE.md — apermo-surf (CLI)

## Overview

CLI tool for project-contextual link navigation. Written in Go. Opens project URLs from the terminal with fuzzy matching, branch-aware placeholders, and fzf integration.

## Ecosystem

Apermo Surf is a two-part project sharing the `.surf-links.yml` config format:

| Component | Repo | Purpose |
|-----------|------|---------|
| **CLI** | [apermo/apermo-surf](https://github.com/apermo/apermo-surf) | This repo — terminal-based link navigation |
| **Chrome** | [apermo/apermo-surf-chrome](https://github.com/apermo/apermo-surf-chrome) | Browser extension for URL detection and link navigation |

### Cross-repo work

When a change affects the shared config format (`.surf-links.yml`), create issues in **both** repos:
- CLI: `gh issue create --repo apermo/apermo-surf`
- Chrome: `gh issue create --repo apermo/apermo-surf-chrome`

The `surf add-to-chrome` command bridges both tools — the CLI pushes config to the Chrome extension.

## Tech Stack

- **Language**: Go
- **Config**: `.surf-links.yml` (YAML, discovered by walking up from cwd)
- **Build**: GoReleaser
- **Distribution**: Homebrew via `apermo/homebrew-tap`

## Project Structure

```
apermo-surf/
├── CLAUDE.md
├── README.md
├── go.mod
├── go.sum
├── main.go
├── cmd/                    # CLI commands (open, links, init, edit, add-to-chrome)
├── internal/
│   ├── config/             # .surf-links.yml discovery and parsing
│   ├── git/                # Branch name, ticket extraction
│   ├── browser/            # Cross-platform URL opening
│   ├── fuzzy/              # Fuzzy matching logic
│   └── chrome/             # Chrome extension bridge
├── .goreleaser.yml
└── .github/
    └── workflows/
        └── release.yml     # GoReleaser on tag push
```

## Shared Config Format

```yaml
# .surf-links.yml
environments:
  local: https://myproject.ddev.site
  staging: https://staging.example.com
  production: https://example.com

tools:
  jira:
    url: https://myorg.atlassian.net/browse/{ticket}
    pattern: "PROJ-\\d+"
  sentry: https://sentry.io/organizations/myorg/projects/myproject

docs:
  confluence: https://myorg.atlassian.net/wiki/spaces/PROJ
```

Links support two formats:
- **Simple**: `key: url`
- **Expanded**: `key: { url, pattern }` (for placeholder support)

## Commands

| Command | Description |
|---------|-------------|
| `surf open <name> [arg]` | Open a link by fuzzy name, optional ticket/ID argument |
| `surf links [--env\|--tools\|--docs]` | List all or filtered links |
| `surf init` | Interactive config setup |
| `surf edit` | Open config in `$EDITOR` |
| `surf add-to-chrome` | Push config to Chrome extension |
| `surf version` | Print version |

## Binary naming

- Homebrew formula: `apermo-surf`
- Binary: `apermo-surf` with `surf` symlink
- User types: `surf`

## Code Conventions

- Go standard project layout
- `internal/` for non-exported packages
- `cmd/` for CLI command definitions
- Error messages should be actionable ("no .surf-links.yml found — run `surf init` to create one")

## Commit Conventions

- Conventional Commits: `<type>(<scope>): <subject>`
- Subject: 50 chars max, body: 72 chars per line
- One topic per commit, atomic and cherry-pickable
- No Co-Authored-By lines
