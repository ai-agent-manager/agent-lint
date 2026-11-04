## Purpose

Structurally validates Atlassian Rovo / Studio agent definitions (`rovo-agent.yaml`) against the vendored agent-manager JSON Schema so Rovo-only skill-root directories produce actionable diagnostics under the same check pipeline as Agent Skills.

## ADDED Requirements

### Requirement: Parse rovo-agent.yaml

The checker SHALL parse `rovo-agent.yaml` as a YAML object and emit error diagnostics when the file is unreadable, invalid YAML, or not a mapping.

#### Scenario: Invalid YAML

- **WHEN** `rovo-agent.yaml` cannot be parsed as YAML
- **THEN** an error diagnostic with rule ID `rovo/yaml-invalid` is emitted

#### Scenario: Non-object root

- **WHEN** `rovo-agent.yaml` parses to a non-mapping value
- **THEN** an error diagnostic with rule ID `rovo/yaml-not-object` is emitted

### Requirement: JSON Schema validation

The checker SHALL validate the parsed document against the vendored Rovo agent JSON Schema (agent-manager `rovo-agent.schema.json`, Draft 2020-12, v2-beta only). Each schema violation SHALL produce an error diagnostic with rule ID `rovo/schema-invalid` whose message identifies the failing location and constraint.

#### Scenario: Missing required field

- **WHEN** a document omits a schema-required property such as `instructions`
- **THEN** an error diagnostic with rule ID `rovo/schema-invalid` is emitted

#### Scenario: Constraint violation

- **WHEN** a document violates a schema constraint (for example `name` longer than 30 characters, or more than 3 `conversationStarters`)
- **THEN** an error diagnostic with rule ID `rovo/schema-invalid` is emitted

#### Scenario: Archived apiVersion

- **WHEN** `apiVersion` is not exactly `rovo.atlassian.com/v2-beta`
- **THEN** an error diagnostic with rule ID `rovo/schema-invalid` is emitted

### Requirement: Vendored schema asset

The Rovo checker SHALL ship with an embedded copy of the agent-manager Rovo agent schema used for validation (no network fetch at lint time).

#### Scenario: Offline validation

- **WHEN** `check` runs without network access
- **THEN** Rovo schema validation still executes using the embedded schema

### Requirement: README config for Rovo directories

When configuration enables README requirements, Rovo-only directories SHALL be subject to the same README presence and frontmatter checks as skill directories (including required tags and allowed phases when configured).

#### Scenario: Missing README under require_readme

- **WHEN** `skills.require_readme` is true and a Rovo-only directory has no `README.md`
- **THEN** an error diagnostic with rule ID `skill/readme-missing` is emitted for that directory
