# Diagnostics Reporting Specification

## Purpose

Every lint finding is a structured diagnostic with a stable rule ID, severity, path, and message. Reporters present active findings to humans or machines and keep ignored findings auditable without failing the build.

## Requirements

### Requirement: Diagnostic Shape

Each diagnostic SHALL include a stable `ruleId`, a severity of `error` or `warning`, a repo-relative `path`, and a human-readable `message`.

#### Scenario: Rule ID present

- **WHEN** any finding is emitted
- **THEN** its `ruleId` is one of the IDs published by `agent-lint list-rules`

### Requirement: Default Severities

Rule severities SHALL follow the registry defaults unless promoted by `--strict` for exit-code purposes.

#### Scenario: Warning does not fail by default

- **WHEN** only warning diagnostics are active and `--strict` is unset
- **THEN** the check exits with code `0`

### Requirement: Ignore Filtering For Exit Status

Ignored diagnostics SHALL not count toward error/warning totals used for exit code `1`.

#### Scenario: All findings ignored

- **WHEN** every diagnostic is marked ignored via the ignore file
- **THEN** the check exits with code `0`

### Requirement: JSON Audit Trail

JSON output SHALL include both active diagnostics and ignored diagnostics with counts.

#### Scenario: Ignored array populated

- **WHEN** at least one diagnostic is ignored and `--json` is set
- **THEN** the JSON document contains those findings under `ignored` and sets `ignoredCount` accordingly

### Requirement: Human Reporter

Human output SHALL label each active diagnostic with its severity and rule ID.

#### Scenario: Colour when TTY

- **WHEN** stdout is a TTY
- **THEN** severity labels may use ANSI colour

#### Scenario: No colour when not a TTY

- **WHEN** stdout is not a TTY
- **THEN** output does not rely on colour for meaning
