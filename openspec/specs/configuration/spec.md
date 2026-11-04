# Configuration Specification

## Purpose

Libraries and monorepos customize agent-lint through `.agent-lint.yaml` (limits, README conventions, Claude Code field policy) and suppress known findings through `.agent-lint-ignore.yaml` (or embedded `ignore` in the main config) without disabling rules globally.

## Requirements

### Requirement: Default Configuration

When no config file is present, agent-lint SHALL use agentskills.io progressive-disclosure defaults of 500 body lines and 5000 body tokens, with README requirements disabled.

#### Scenario: No config file

- **WHEN** `check` runs without `--config` and no `.agent-lint.yaml` exists in the working directory
- **THEN** body caps default to 500 / 5000 and README checks are not required

### Requirement: Auto-Discover Config

When `--config` is omitted, agent-lint SHALL load `.agent-lint.yaml` from the working directory if present.

#### Scenario: Config discovered

- **WHEN** `.agent-lint.yaml` exists in the working directory
- **THEN** its `skills` settings are applied

### Requirement: Config-Gated README Validation

When `skills.require_readme` is true, each skill directory SHALL include a `README.md` with valid frontmatter name and description, and SHALL satisfy configured tag and phase constraints when those lists are non-empty.

#### Scenario: Missing README

- **WHEN** `require_readme` is true and `README.md` is absent
- **THEN** an error with rule ID `skill/readme-missing` is emitted

#### Scenario: Required tag

- **WHEN** `required_tags` includes `agent-skill` and README tags omit it
- **THEN** an error with rule ID `skill/readme-required-tag` is emitted

#### Scenario: Invalid phase

- **WHEN** `allowed_phases` is configured and README lists a phase outside that set
- **THEN** an error with rule ID `skill/readme-phase-value` is emitted

### Requirement: Ignore File Discovery

When `--ignore-file` is omitted, agent-lint SHALL load `.agent-lint-ignore.yaml` from the working directory if present, and SHALL also honour `ignore.ignores` entries embedded in the main config.

#### Scenario: Combined ignore sources

- **WHEN** both config-embedded ignores and an ignore file define entries
- **THEN** all entries are applied

### Requirement: Path And Glob Matching

Each ignore entry SHALL match diagnostics whose path equals the entry path, is under that directory prefix, or matches the entry as a glob (`*`, `**`, `?`).

#### Scenario: Exact file ignore

- **WHEN** an ignore entry path is `skills/foo/SKILL.md` and a diagnostic is emitted for that path
- **THEN** the entry is eligible to suppress that diagnostic

#### Scenario: Glob ignore

- **WHEN** an ignore entry path is `skills/*-evals/**`
- **THEN** diagnostics under matching eval skill directories are eligible for suppression

### Requirement: Rule List Semantics

An ignore entry's `rules` list SHALL select which rule IDs to suppress on matching paths; an empty list SHALL suppress all findings on matching paths.

#### Scenario: Empty rules list

- **WHEN** `rules` is `[]` for a matching path
- **THEN** every diagnostic on that path is marked ignored

#### Scenario: Specific rules

- **WHEN** `rules` lists `skill/body-line-limit`
- **THEN** only diagnostics with that rule ID on matching paths are ignored

#### Scenario: Unknown rule ID

- **WHEN** `rules` contains an ID not present in the rule registry
- **THEN** configuration loading fails with exit code `2`
