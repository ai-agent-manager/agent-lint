## Why

`required_tags` and `allowed_phases` hard-code org-specific README frontmatter field names into agent-lint, even though neither belongs to agentskills.io. Libraries need configurable checks on arbitrary frontmatter string fields (including dotted paths such as `metadata.phases`) without baking taxonomy into the product. SKILL.md stays aligned with agentskills.io’s string-valued `metadata` map rather than inventing YAML arrays.

## What Changes

- Remove `skills.required_tags` and `skills.allowed_phases` (pre-v1.0; no compatibility shim).
- Remove the old tag/phase rule IDs.
- Add parallel config lists `skills.readme_constraints` and `skills.skill_constraints`.
- Each entry names a dotted `field` path and exactly one string operator: `equals`, `not_equals`, `matches`, or `not_matches` (`matches` / `not_matches` use Go/RE2 regular expressions).
- README constraints apply when `require_readme` is true; skill constraints apply whenever configured.
- The root key of each `skill_constraints` field path is carved out of the SKILL.md allowlist when needed (so top-level org keys still work; `metadata.*` needs no carve-out).
- Emit file- and operator-specific rule IDs for type failures and operator failures.
- Update docs, fixtures, and examples.

### Non-goals

- List/array operators (`must_include`, `allowed` over YAML sequences).
- Nested object schemas beyond dotted string paths.
- Skill installers, cross-agent flavours, or Rovo-specific metadata.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `configuration`: Replace tag/phase knobs with string-operator constraint lists for README and SKILL.md.
- `skill-frontmatter`: Root keys of `skill_constraints` field paths SHALL NOT alone trigger `skill/frontmatter-unknown-field`.

## Impact

- `internal/config`, `internal/skill`, `internal/rules`
- User docs and fixtures
- OpenSpec `configuration` and `skill-frontmatter` (via archive later)
