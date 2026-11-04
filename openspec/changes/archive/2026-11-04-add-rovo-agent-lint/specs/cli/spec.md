## MODIFIED Requirements

### Requirement: Check Command

The CLI SHALL provide a `check` command that validates skill and Rovo agent directories under a configurable root.

#### Scenario: Default root

- **WHEN** the user runs `agent-lint check` with no `--root`
- **THEN** the tool scans the `skills` directory relative to the current working directory

#### Scenario: Custom root

- **WHEN** the user runs `agent-lint check --root <path>`
- **THEN** the tool scans the given skills root

#### Scenario: Missing root

- **WHEN** the configured skills root does not exist
- **THEN** the process exits with code `2`

#### Scenario: Rovo diagnostics included

- **WHEN** the skills root contains a Rovo-only directory with an invalid `rovo-agent.yaml`
- **THEN** `check` reports the corresponding `rovo/…` diagnostics in human and JSON output

### Requirement: List Rules Command

The CLI SHALL provide a `list-rules` command that prints every stable rule ID with its default severity and description.

#### Scenario: List rules

- **WHEN** the user runs `agent-lint list-rules`
- **THEN** each known rule ID is printed exactly once

#### Scenario: Rovo rules listed

- **WHEN** the user runs `agent-lint list-rules`
- **THEN** every `rovo/…` rule ID is included in the output
