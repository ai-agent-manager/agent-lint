## MODIFIED Requirements

### Requirement: Config-Gated README Validation

When `skills.require_readme` is true, each skill directory SHALL include a `README.md` with valid frontmatter `name` and `description`, and SHALL satisfy every entry in `skills.readme_constraints` when that list is non-empty.

Each `readme_constraints` entry SHALL name a frontmatter `field` and MAY include `must_include` and/or `allowed` string lists. When either list is non-empty for an entry, the named field SHALL be present as a non-empty array of strings. When `must_include` is non-empty, every listed value SHALL appear in the field. When `allowed` is non-empty, every value in the field SHALL be a member of that allowlist. Entries with both lists empty SHALL be rejected at config load (exit code `2`).

#### Scenario: Missing README

- **WHEN** `require_readme` is true and `README.md` is absent
- **THEN** an error with rule ID `skill/readme-missing` is emitted

#### Scenario: Required list value missing

- **WHEN** a constraint has `field: tags` and `must_include: [agent-skill]` and README `tags` omit `agent-skill`
- **THEN** an error with rule ID `skill/readme-field-required-value` is emitted

#### Scenario: Value outside allowlist

- **WHEN** a constraint has `field: phases` and a non-empty `allowed` list, and README `phases` contains a value outside that set
- **THEN** an error with rule ID `skill/readme-field-allowed-value` is emitted

#### Scenario: Field not a string list

- **WHEN** a constraint applies to a field that is missing, empty, or not an array of strings
- **THEN** an error with rule ID `skill/readme-field-invalid` is emitted

#### Scenario: Empty constraint entry rejected

- **WHEN** a `readme_constraints` entry has neither `must_include` nor `allowed` values
- **THEN** configuration loading fails with exit code `2`
