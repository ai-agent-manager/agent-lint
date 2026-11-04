# Skill Frontmatter Specification

## Purpose

Agent Skill `SKILL.md` files must open with YAML frontmatter that conforms to the agentskills.io field table and naming rules. This capability is the normative frontmatter contract agent-lint enforces.

## Requirements

### Requirement: Frontmatter Delimiters

`SKILL.md` SHALL begin with YAML frontmatter delimited by opening and closing `---` lines and parse to a YAML mapping.

#### Scenario: Missing opening fence

- **WHEN** `SKILL.md` does not start with `---`
- **THEN** an error with rule ID `skill/frontmatter-missing` is emitted

#### Scenario: Invalid YAML

- **WHEN** the frontmatter block is not valid YAML
- **THEN** an error with rule ID `skill/frontmatter-invalid-yaml` is emitted

### Requirement: Required Name Field

Frontmatter SHALL include a non-empty `name` that satisfies agentskills.io naming constraints and matches the parent directory name.

#### Scenario: Valid name

- **WHEN** `name` is 1–64 characters, uses lowercase alphanumeric characters and single hyphens only (no leading, trailing, or consecutive hyphens), and equals the parent directory name
- **THEN** no name-format or directory-mismatch errors are emitted for that field

#### Scenario: Invalid format

- **WHEN** `name` violates the agentskills.io format rules
- **THEN** an error with rule ID `skill/name-format` is emitted

#### Scenario: Directory mismatch

- **WHEN** `name` does not equal the parent directory name
- **THEN** an error with rule ID `skill/name-directory-mismatch` is emitted

#### Scenario: Missing name

- **WHEN** `name` is absent or empty
- **THEN** an error with rule ID `skill/name-missing` is emitted

### Requirement: Required Description Field

Frontmatter SHALL include a non-empty `description` of at most 1024 characters.

#### Scenario: Description too long

- **WHEN** `description` exceeds 1024 characters
- **THEN** an error with rule ID `skill/description-length` is emitted

#### Scenario: Missing description

- **WHEN** `description` is absent or empty
- **THEN** an error with rule ID `skill/description-missing` is emitted

### Requirement: Optional Compatibility Field

When `compatibility` is present, it SHALL be a non-empty string of at most 500 characters.

#### Scenario: Compatibility too long

- **WHEN** `compatibility` exceeds 500 characters
- **THEN** an error with rule ID `skill/compatibility-length` is emitted

### Requirement: Optional Metadata Map

When `metadata` is present, it SHALL be a mapping of string keys to string values.

#### Scenario: Non-string metadata value

- **WHEN** a metadata entry is not a string
- **THEN** an error with rule ID `skill/metadata-format` is emitted

### Requirement: Optional Allowed-Tools Field

When `allowed-tools` is present, it SHALL be a space-separated string (not a comma-separated list or non-string).

#### Scenario: Comma-separated tools

- **WHEN** `allowed-tools` contains a comma
- **THEN** an error with rule ID `skill/allowed-tools-format` is emitted

### Requirement: Frontmatter Allowlist

Only agentskills.io fields (`name`, `description`, `license`, `compatibility`, `metadata`, `allowed-tools`) SHALL be accepted by default.

#### Scenario: Unknown field

- **WHEN** frontmatter contains a key outside the allowlist and Claude Code extensions are not enabled
- **THEN** an error with rule ID `skill/frontmatter-unknown-field` is emitted

#### Scenario: Claude Code extensions allowed

- **WHEN** Claude Code extensions are enabled and frontmatter contains a documented Claude Code extension key
- **THEN** that key does not produce `skill/frontmatter-unknown-field`
