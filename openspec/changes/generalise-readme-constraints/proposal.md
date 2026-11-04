## Why

`required_tags` and `allowed_phases` hard-code org-specific README frontmatter field names into agent-lint, even though neither belongs to agentskills.io. Libraries need the same two constraint operators (`must_include`, `allowed`) on arbitrary README and SKILL.md list fields without each taxonomy becoming a first-class product concept.

## What Changes

- Remove `skills.required_tags` and `skills.allowed_phases` from `.agent-lint.yaml` (pre-v1.0; no compatibility shim).
- Remove rule IDs `skill/readme-tags-invalid`, `skill/readme-required-tag`, `skill/readme-phases-invalid`, and `skill/readme-phase-value`.
- Add parallel config lists `skills.readme_constraints` and `skills.skill_constraints`: each entry names a frontmatter field and may declare `must_include` and/or `allowed` for string-list values.
- README constraints apply when `require_readme` is true; skill constraints always apply to `SKILL.md` when the list is non-empty.
- Fields named in `skill_constraints` are treated as allowed SKILL.md frontmatter keys (carve-out from the agentskills.io allowlist) so org taxonomy can live at the top level as string lists.
- Emit file-specific rule IDs for type/shape failures, missing required values, and values outside an allowlist (`skill/readme-field-*` and `skill/skill-field-*`).
- Update docs, fixtures, and configuration/examples to the new shape.

### Non-goals

- Single-string / scalar enum constraints, regex, or nested object schemas.
- Encoding list constraints inside `metadata` (string map only per agentskills.io).
- Skill installers, cross-agent flavours, or Rovo-specific metadata.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `configuration`: Replace tag/phase-specific README knobs with `readme_constraints` and add parallel `skill_constraints`; redefine rule IDs and validation semantics for string-list frontmatter constraints.
- `skill-frontmatter`: Fields listed in `skill_constraints` SHALL NOT emit `skill/frontmatter-unknown-field`.

## Impact

- `internal/config`, `internal/skill` frontmatter/README validation, `internal/rules`
- User docs (`docs/configuration.md`, `docs/rules.md`) and examples
- OpenSpec `configuration` and `skill-frontmatter` main specs (via archive later)
- Existing configs using `required_tags` / `allowed_phases` must migrate; ignore files referencing removed rule IDs must be updated
