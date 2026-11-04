## Why

`required_tags` and `allowed_phases` hard-code org-specific README frontmatter field names into agent-lint, even though neither belongs to agentskills.io. Libraries need the same two constraint operators (`must_include`, `allowed`) on arbitrary README list fields without each taxonomy becoming a first-class product concept.

## What Changes

- Remove `skills.required_tags` and `skills.allowed_phases` from `.agent-lint.yaml` (pre-v1.0; no compatibility shim).
- Remove rule IDs `skill/readme-tags-invalid`, `skill/readme-required-tag`, `skill/readme-phases-invalid`, and `skill/readme-phase-value`.
- Add `skills.readme_constraints`: a list of per-field constraints applied to `README.md` frontmatter when `require_readme` is true.
- Each constraint names a frontmatter field and may declare `must_include` and/or `allowed` for string-list values.
- Emit generic rule IDs for type/shape failures, missing required values, and values outside an allowlist.
- Update docs, fixtures, and configuration/examples to the new shape.

### Non-goals

- Constraints on `SKILL.md` frontmatter (agentskills allowlist remains authoritative).
- Single-string / scalar enum constraints, regex, or nested object schemas.
- Skill installers, cross-agent flavours, or Rovo-specific metadata.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `configuration`: Replace tag/phase-specific README knobs with generic `readme_constraints`; redefine rule IDs and validation semantics for README list-field constraints.

## Impact

- `internal/config`, `internal/skill` README validation, `internal/rules`
- User docs (`docs/configuration.md`, `docs/rules.md`) and examples
- OpenSpec `configuration` main spec (via archive later)
- Existing configs using `required_tags` / `allowed_phases` must migrate; ignore files referencing removed rule IDs must be updated
