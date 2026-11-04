## ADDED Requirements

### Requirement: Constrained Root Keys Are Allowed

The root path segment of each `skills.skill_constraints` `field` SHALL be accepted by the SKILL.md allowlist check and SHALL NOT emit `skill/frontmatter-unknown-field` solely for being outside the agentskills.io field table. Paths under already-allowed roots (for example `metadata.phases`) require no additional carve-out.

#### Scenario: Top-level constrained org field present

- **WHEN** `skill_constraints` includes `field: phases` with a valid string operator and SKILL.md frontmatter contains string `phases` satisfying that operator
- **THEN** no `skill/frontmatter-unknown-field` diagnostic is emitted for `phases`

#### Scenario: Metadata path does not widen allowlist unnecessarily

- **WHEN** `skill_constraints` only references `metadata.phases` and SKILL.md also contains an unconstrained unknown top-level key `owner`
- **THEN** an error with rule ID `skill/frontmatter-unknown-field` is emitted for `owner`

#### Scenario: Unconstrained unknown field still rejected

- **WHEN** SKILL.md frontmatter contains a key that is neither an agentskills.io field, a Claude Code extension (when enabled), nor the root of a `skill_constraints` field path
- **THEN** an error with rule ID `skill/frontmatter-unknown-field` is emitted
