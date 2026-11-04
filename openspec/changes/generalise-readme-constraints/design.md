## Context

See proposal.md — Why. Today `SkillsConfig` has parallel `RequiredTags` / `AllowedPhases` fields and `validateReadme` hard-codes the `tags` / `phases` keys plus four rule IDs. No compatibility shim: treat the old keys as never having existed.

## Goals / Non-Goals

**Goals:**

- One config list (`readme_constraints`) driving README string-list checks
- Generic, stable rule IDs independent of field names
- Fail fast on empty/invalid constraint entries at config load

**Non-Goals:**

- Aliases or migration for `required_tags` / `allowed_phases`
- Constraints beyond string lists on README frontmatter
- Applying constraints to `SKILL.md`

## Decisions

### 1. Config shape (Option B)

```yaml
skills:
  require_readme: true
  readme_constraints:
    - field: tags
      must_include: [agent-skill]
    - field: phases
      allowed: [visioning, discovery, elaboration, build_and_release, run_and_measure]
```

Go model: `[]ReadmeConstraint` with `Field`, `MustInclude`, `Allowed`. Prefer a list over a field→constraint map so order is stable in diagnostics and future operators attach cleanly.

**Rejected:** Keep `required_tags` as a special case (leaves taxonomy baked in). Map-of-fields (Option A) — slightly terser YAML but worse for multi-operator extension and diagnostic ordering.

### 2. Presence semantics

If an entry has any non-empty `must_include` or `allowed`, the field must exist as a non-empty string array (same as today's implied behaviour). Do not add a separate `required: true` flag in this change.

### 3. Rule IDs

| ID | When |
|----|------|
| `skill/readme-field-invalid` | Missing / wrong type / empty list when constrained |
| `skill/readme-field-required-value` | `must_include` value absent |
| `skill/readme-field-allowed-value` | Value not in `allowed` |

Messages name the field and value so humans do not need field-specific IDs. Delete the four tag/phase rule constants.

### 4. Config validation

On `Load`, reject any constraint with empty `field`, or with both `must_include` and `allowed` empty. Unknown YAML keys under old names are ignored by the YAML decoder (no special error) — documented as removed, not detected.

### 5. Validation loop

In `validateReadme`, after name/description checks, iterate `readme_constraints` and apply the shared string-list logic. No field-name switches.

## Risks / Trade-offs

- [Ignore files referencing old rule IDs] → Unknown IDs already fail load with exit 2; callers must update ignore files when upgrading.
- [Org configs still using old keys] → Silent no-op of old keys until they migrate; accepted because no backwards compatibility was requested — document clearly in `docs/configuration.md`.

## Migration Plan

1. Ship the config/rule swap in the same release (pre-v1.0; no deprecation window).
2. Docs example shows tags + phases expressed as constraints.
3. No runtime shim.
