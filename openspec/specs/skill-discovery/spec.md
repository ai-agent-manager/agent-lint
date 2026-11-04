# Skill Discovery Specification

## Purpose

agent-lint discovers lint targets by walking immediate child directories of the skills root and classifying each directory as an Agent Skill, a Rovo agent, or an invalid asset layout. This keeps skill validation focused while leaving room for a future Rovo checker.

## Requirements

### Requirement: Scan Immediate Children

The checker SHALL treat each immediate subdirectory of the skills root as a candidate asset directory.

#### Scenario: Nested directories are not scanned as top-level assets

- **WHEN** the skills root contains `foo/` and `foo/bar/`
- **THEN** only `foo` is considered a candidate asset directory for discovery

### Requirement: Exactly One Definition File

Each candidate directory SHALL contain exactly one of `SKILL.md` or `rovo-agent.yaml`.

#### Scenario: Skill directory

- **WHEN** a directory contains `SKILL.md` and not `rovo-agent.yaml`
- **THEN** it is validated as an Agent Skill

#### Scenario: Rovo-only directory skipped

- **WHEN** a directory contains `rovo-agent.yaml` and not `SKILL.md`
- **THEN** skill validation is skipped for that directory (no Rovo structural lint yet)

#### Scenario: Both definition files

- **WHEN** a directory contains both `SKILL.md` and `rovo-agent.yaml`
- **THEN** an error diagnostic with rule ID `skill/asset-type-conflict` is emitted

#### Scenario: Neither definition file

- **WHEN** a directory contains neither `SKILL.md` nor `rovo-agent.yaml`
- **THEN** an error diagnostic with rule ID `skill/asset-type-missing` is emitted

### Requirement: Exact SKILL.md Casing

A skill directory MUST present the definition file as exactly `SKILL.md`.

#### Scenario: Wrong casing warning

- **WHEN** a directory is missing `SKILL.md` but contains a differently cased `skill.md` (or similar)
- **THEN** a warning diagnostic with rule ID `skill/skill-md-casing` is emitted in addition to the missing-asset error

### Requirement: Duplicate Skill Names

Within a single scan root, declared skill `name` values SHALL be unique.

#### Scenario: Duplicate name

- **WHEN** two skill directories declare the same frontmatter `name`
- **THEN** an error diagnostic with rule ID `skill/name-duplicate` is emitted for the later skill
