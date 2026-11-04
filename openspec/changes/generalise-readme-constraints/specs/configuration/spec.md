## MODIFIED Requirements

### Requirement: Config-Gated README Validation

When `skills.require_readme` is true, each skill directory SHALL include a `README.md` with valid frontmatter `name` and `description`, and SHALL satisfy every entry in `skills.readme_constraints` when that list is non-empty.

Each `readme_constraints` entry SHALL name a dotted frontmatter `field` path and SHALL use either string operators or list operators, not both.

String operators: exactly one of `equals`, `not_equals`, `matches`, or `not_matches`. The resolved value SHALL be a YAML string. `matches` / `not_matches` use a Go regular expression.

List operators: `must_include` and/or `allowed`. The resolved value SHALL be a non-empty YAML array of strings. When `must_include` is set, every listed value SHALL appear in the array. When `allowed` is set, every array value SHALL be a member of that allowlist.

Entries with an empty `field`, no operators, mixed string and list operators, more than one string operator, or an invalid regular expression SHALL be rejected at config load (exit code `2`).

#### Scenario: Missing README

- **WHEN** `require_readme` is true and `README.md` is absent
- **THEN** an error with rule ID `skill/readme-missing` is emitted

#### Scenario: Equals failure

- **WHEN** a constraint has `field: name` and `equals: Expected` and README `name` is not exactly `Expected`
- **THEN** an error with rule ID `skill/readme-field-equals` is emitted

#### Scenario: Matches failure

- **WHEN** a constraint has `matches` set and the README field string does not match the pattern
- **THEN** an error with rule ID `skill/readme-field-matches` is emitted

#### Scenario: Field missing or wrong type

- **WHEN** a string-operator path is missing or not a string, or a list-operator path is missing or not a non-empty string array
- **THEN** an error with rule ID `skill/readme-field-invalid` is emitted

#### Scenario: Must include failure

- **WHEN** a constraint has `must_include: [agent-skill]` and the README field array omits `agent-skill`
- **THEN** an error with rule ID `skill/readme-field-required-value` is emitted

#### Scenario: Allowed failure

- **WHEN** a constraint has a non-empty `allowed` list and the README field array contains a value outside that set
- **THEN** an error with rule ID `skill/readme-field-allowed-value` is emitted

#### Scenario: Invalid constraint entry rejected

- **WHEN** a `readme_constraints` entry sets zero operators, mixes string and list operators, sets multiple string operators, or an invalid `matches` / `not_matches` pattern
- **THEN** configuration loading fails with exit code `2`

## ADDED Requirements

### Requirement: Skill Frontmatter Constraints

When `skills.skill_constraints` is non-empty, each skill's `SKILL.md` frontmatter SHALL satisfy every entry. Entries SHALL use dotted `field` paths and exactly one string operator (`equals`, `not_equals`, `matches`, `not_matches`). List operators (`must_include`, `allowed`) SHALL be rejected at config load (exit code `2`). Failures SHALL use `skill/skill-field-invalid`, `skill/skill-field-equals`, `skill/skill-field-not-equals`, `skill/skill-field-matches`, or `skill/skill-field-not-matches` as appropriate.

#### Scenario: Metadata equals failure

- **WHEN** a skill constraint has `field: metadata.author` and `equals: agent-lint` and the metadata value differs
- **THEN** an error with rule ID `skill/skill-field-equals` is emitted

#### Scenario: Metadata matches allowlist pattern

- **WHEN** a skill constraint has `field: metadata.phases` and `matches: "^(build_and_release)$"` and the value is `build_and_release`
- **THEN** no `skill/skill-field-matches` diagnostic is emitted

#### Scenario: Field missing or not a string on SKILL.md

- **WHEN** a skill constraint path is missing or does not resolve to a string
- **THEN** an error with rule ID `skill/skill-field-invalid` is emitted
