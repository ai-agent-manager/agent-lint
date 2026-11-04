# CLI Specification

## Purpose

The `agent-lint` command-line interface is the sole user-facing entry point. It exposes non-interactive commands for validating skills and listing rule IDs, with stable exit codes suitable for CI and agent automation.

## Requirements

### Requirement: Check Command

The CLI SHALL provide a `check` command that validates skill directories under a configurable root.

#### Scenario: Default root

- **WHEN** the user runs `agent-lint check` with no `--root`
- **THEN** the tool scans the `skills` directory relative to the current working directory

#### Scenario: Custom root

- **WHEN** the user runs `agent-lint check --root <path>`
- **THEN** the tool scans the given skills root

#### Scenario: Missing root

- **WHEN** the configured skills root does not exist
- **THEN** the process exits with code `2`

### Requirement: List Rules Command

The CLI SHALL provide a `list-rules` command that prints every stable rule ID with its default severity and description.

#### Scenario: List rules

- **WHEN** the user runs `agent-lint list-rules`
- **THEN** each known rule ID is printed exactly once

### Requirement: Exit Codes

The CLI SHALL use stable exit codes: `0` success, `1` lint failures, `2` tool or usage errors.

#### Scenario: Clean check

- **WHEN** `check` completes with no active errors (and no active warnings if `--strict` is unset)
- **THEN** the process exits with code `0`

#### Scenario: Lint failures

- **WHEN** `check` finds one or more active error diagnostics
- **THEN** the process exits with code `1`

#### Scenario: Strict warnings

- **WHEN** `--strict` is set and one or more active warning diagnostics exist with no active errors
- **THEN** the process exits with code `1`

#### Scenario: Invalid ignore rule ID

- **WHEN** an ignore entry references an unknown rule ID
- **THEN** the process exits with code `2`

### Requirement: Output Modes

The CLI SHALL support human-readable output by default and machine-readable JSON via `--json`.

#### Scenario: JSON report

- **WHEN** `check` is run with `--json`
- **THEN** a single JSON document is written to stdout describing diagnostics, ignored findings, and counts

#### Scenario: Quiet human mode

- **WHEN** `check` is run with `--quiet` without `--json`
- **THEN** the summary line is suppressed and diagnostics are still shown

### Requirement: Feature Flags

The CLI SHALL expose opt-in flags for Claude.ai checks and Claude Code frontmatter extensions.

#### Scenario: Claude.ai profile

- **WHEN** `check` is run with `--claude`
- **THEN** Claude.ai-specific rules (reserved name words; angle brackets in description) are enforced

#### Scenario: Allow Claude Code fields

- **WHEN** `check` is run with `--allow-claude-code` or config sets `skills.allow_claude_code_fields`
- **THEN** documented Claude Code frontmatter extension keys are treated as allowed fields
