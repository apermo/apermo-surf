# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.3.1] - Unreleased

### Changed

- Updated README with full config format documentation

### Fixed

- Removed stale Kong reference from tech stack
- Fixed release workflow writing notes to working tree

## [0.3.0] - 2026-02-28

### Added

- `.surf-links.yml.dist` fallback for team-shared template configs (#19)
- Configurable browser selection via `--browser` flag and
  `~/.config/surf/config.yml` (#14)
- Interactive config wizard `surf init` with `--dist` flag (#20)
- Optional `name` field in config for project display name (#18)
- Sub-links on link entries for nested navigation shortcuts (#17)
- Craft CMS added to standard project types (#16)
- Testable examples for `go doc` output (#21)
- Fuzz tests for regex-heavy functions (#22)
- GoReleaser config and GitHub Actions release workflow (#9)
- CI workflows, PR validation, and branch protection (#23)

## [0.2.0] - 2026-02-28

### Added

- Placeholder substitution for `{branch}`, `{ticket}`, and
  `{repo}` in link URLs (#5)
- Explicit ticket argument: `surf open jira 123` (#6)
- Interactive picker with fzf integration, numbered list
  fallback (#7)
- Dynamic shell completions for bash, zsh, and fish (#8)
- Configurable project type with auto-generated admin links
  for WordPress, TYPO3, Laravel, Drupal, Shopware, and
  Magento (#15)

## [0.1.0] - 2026-02-27

### Added

- Go CLI scaffold with Cobra framework (#1)
- `.surf-links.yml` config parsing with simple and expanded
  link formats (#2)
- `surf links` command with category filtering (#3)
- `surf open` command with fuzzy name matching (#4)
