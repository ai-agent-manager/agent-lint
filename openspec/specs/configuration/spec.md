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

When `skills.require_readme` is true, each skill directory SHALL include a `README.md` with valid frontmatter `name` and `description`, and SHALL satisfy every entry in `skills.readme_constraints` when that list is non-empty.

Each `readme_constraints` entry SHALL name a dotted frontmatter `field` path and SHALL use either string operators or list operators, not both.

String operators: exactly one of `equals`, `not_equals`, `matches`, or `not_matches`. The resolved value SHALL be a YAML string. `matches` / `not_matches` use a Go regular expression.

List operators: `must_include` and/or `allowed`. The resolved value SHALL be a non-empty YAML array of strings. When `must_include` is set, every listed value SHALL appear in the array. When `allowed` is set, every array value SHALL be a member of that allowlist.

Entries with an empty `field`, no operators, mixed string and list operators, more than one string operator, or an invalid regular expression SHALL be rejected at config load (exit code `2`).

#### Scenario: Missing README

- **WHEN** `require_readme` is true and `README.md` is absent
- **THEN** an error with rule ID `skill/readme-missing` is emitted

#### Scenario: Equals failure

- **WHEN** a constraint has `field: name` and `equals: Expected` and README `name` is not exactly `Expected`
- **THEN** an error with rule ID `skill/readme-field-equals` is emitted

#### Scenario: Matches failure

- **WHEN** a constraint has `matches` set and the README field string does not match the pattern
- **THEN** an error with rule ID `skill/readme-field-matches` is emitted

#### Scenario: Field missing or wrong type

- **WHEN** a string-operator path is missing or not a string, or a list-operator path is missing or not a non-empty string array
- **THEN** an error with rule ID `skill/readme-field-invalid` is emitted

#### Scenario: Must include failure

- **WHEN** a constraint has `must_include: [agent-skill]` and the README field array omits `agent-skill`
- **THEN** an error with rule ID `skill/readme-field-required-value` is emitted

#### Scenario: Allowed failure

- **WHEN** a constraint has a non-empty `allowed` list and the README field array contains a value outside that set
- **THEN** an error with rule ID `skill/readme-field-allowed-value` is emitted

#### Scenario: Invalid constraint entry rejected

- **WHEN** a `readme_constraints` entry sets zero operators, mixes string and list operators, sets multiple string operators, or an invalid `matches` / `not_matches` pattern
- **THEN** configuration loading fails with exit code `2`

### Requirement: Skill Frontmatter Constraints

When `skills.skill_constraints` is non-empty, each skill's `SKILL.md` frontmatter SHALL satisfy every entry. Entries SHALL use dotted `field` paths and exactly one string operator (`equals`, `not_equals`, `matches`, `not_matches`). List operators (`must_include`, `allowed`) SHALL be rejected at config load (exit code `2`). Failures SHALL use `skill/skill-field-invalid`, `skill/skill-field-equals`, `skill/skill-field-not-equals`, `skill/skill-field-matches`, or `skill/skill-field-not-matches` as appropriate.

#### Scenario: Metadata equals failure

- **WHEN** a skill constraint has `field: metadata.author` and `equals: agent-lint` and the metadata value differs
- **THEN** an error with rule ID `skill/skill-field-equals` is emitted

#### Scenario: Metadata matches allowlist pattern

- **WHEN** a skill constraint has `field: metadata.phases` and `matches: "^(build_and_release)$"` and the value is `build_and_release`
- **THEN** no `skill/skill-field-matches` diagnostic is emitted

#### Scenario: Field missing or not a string on SKILL.md

- **WHEN** a skill constraint path is missing or does not resolve to a string
- **THEN** an error with rule ID `skill/skill-field-invalid` is emitted

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
