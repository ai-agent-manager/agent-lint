## ADDED Requirements

### Requirement: Constrained Fields Are Allowed

Frontmatter keys named as `field` in `skills.skill_constraints` SHALL be accepted by the SKILL.md allowlist check and SHALL NOT emit `skill/frontmatter-unknown-field` solely for being outside the agentskills.io field table.

#### Scenario: Constrained org field present

- **WHEN** `skill_constraints` includes `field: phases` and SKILL.md frontmatter contains `phases` as a valid string list within the configured allowlist
- **THEN** no `skill/frontmatter-unknown-field` diagnostic is emitted for `phases`

#### Scenario: Unconstrained unknown field still rejected

- **WHEN** SKILL.md frontmatter contains a key that is neither an agentskills.io field, a Claude Code extension (when enabled), nor named in `skill_constraints`
- **THEN** an error with rule ID `skill/frontmatter-unknown-field` is emitted
