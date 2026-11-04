# Skill Progressive Disclosure Specification

## Purpose

agentskills.io progressive disclosure expects a concise `SKILL.md` body with deeper material loaded on demand from referenced files. agent-lint enforces the published size guidance as hard errors by default and validates relative references from the skill root.

## Requirements

### Requirement: Body Line Limit

The `SKILL.md` body (content after the closing frontmatter fence) SHALL contain at most the configured maximum number of lines (default **500**).

#### Scenario: Body exceeds line cap

- **WHEN** the body has more than `skills.max_body_lines` lines
- **THEN** an error with rule ID `skill/body-line-limit` is emitted citing the agentskills.io guidance

#### Scenario: Config override

- **WHEN** `skills.max_body_lines` is set higher than 500 in config
- **THEN** the configured limit is used for the error threshold

### Requirement: Body Token Limit

The `SKILL.md` body SHALL be under the configured maximum token estimate (default **5000**), using `ceil(utf8_rune_count / 4)`.

#### Scenario: Body at or above token cap

- **WHEN** the estimated token count is greater than or equal to `skills.max_body_tokens`
- **THEN** an error with rule ID `skill/body-token-limit` is emitted

### Requirement: Non-Empty Body Warning

An empty body after frontmatter SHALL be reported as a warning.

#### Scenario: Empty body

- **WHEN** the body is empty or only whitespace
- **THEN** a warning with rule ID `skill/body-empty` is emitted

### Requirement: Markdown Link Targets Exist

Relative markdown link targets in `SKILL.md` SHALL resolve to an existing path under the skill directory.

#### Scenario: Missing linked file

- **WHEN** a markdown link points to a relative path that does not exist under the skill root
- **THEN** an error with rule ID `skill/reference-missing` is emitted

#### Scenario: External links skipped

- **WHEN** a link target is `http:`, `https:`, `mailto:`, or an in-page `#` anchor
- **THEN** no existence check is performed for that target

#### Scenario: Placeholder targets skipped

- **WHEN** a link target is a placeholder such as `url` or contains a glob `*`
- **THEN** no existence check is performed for that target

### Requirement: One-Level Reference Depth

Relative references that are real skill-local paths (markdown links, or plain paths that exist) SHALL be at most one directory level deep from `SKILL.md`.

#### Scenario: Nested reference path

- **WHEN** a validated reference path has more than one directory segment before the filename (e.g. `references/sub/file.md`)
- **THEN** an error with rule ID `skill/reference-depth` is emitted

#### Scenario: Trailing slash normalized

- **WHEN** a reference path ends with `/`
- **THEN** the trailing slash does not by itself inflate the measured depth

### Requirement: Orphan References Warning

Files under `references/` that are never mentioned from `SKILL.md` SHALL be reported as warnings.

#### Scenario: Unmentioned references file

- **WHEN** a file exists under `references/` and is not mentioned from `SKILL.md`
- **THEN** a warning with rule ID `skill/reference-orphan` is emitted

### Requirement: Plain Path Mentions Are Non-Authoritative

Plain-text path mentions (outside markdown links) SHALL NOT by themselves require on-disk existence, so documentation examples do not fail the lint.

#### Scenario: Example path in prose

- **WHEN** the body mentions `scripts/example.py` as plain text and that file does not exist
- **THEN** no `skill/reference-missing` diagnostic is emitted solely for that plain mention
