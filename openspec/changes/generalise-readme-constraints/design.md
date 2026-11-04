## Context

See proposal.md — Why. Today `SkillsConfig` has parallel `RequiredTags` / `AllowedPhases` fields and `validateReadme` hard-codes the `tags` / `phases` keys plus four rule IDs. SKILL.md frontmatter rejects unknown keys unless Claude Code extensions are enabled. No compatibility shim for the removed keys.

## Goals / Non-Goals

**Goals:**

- Parallel config lists (`readme_constraints`, `skill_constraints`) with the same operators
- File-specific stable rule IDs independent of field names
- Skill constraint field names carve out of the agentskills.io allowlist
- Fail fast on empty/invalid constraint entries at config load

**Non-Goals:**

- Aliases or migration for `required_tags` / `allowed_phases`
- Constraints beyond string lists
- Nesting taxonomy under `metadata`

## Decisions

### 1. Config shape (parallel keys)

```yaml
skills:
  require_readme: true
  readme_constraints:
    - field: tags
      must_include: [agent-skill]
    - field: phases
      allowed: [visioning, discovery, elaboration, build_and_release, run_and_measure]
  skill_constraints:
    - field: phases
      allowed: [visioning, discovery, elaboration, build_and_release, run_and_measure]
```

Shared Go type `FrontmatterConstraint` (`Field`, `MustInclude`, `Allowed`) used by both lists.

**Rejected:** Single list with a `target` discriminator (more YAML noise for little gain). Map-of-fields (Option A).

### 2. Presence semantics

If an entry has any non-empty `must_include` or `allowed`, the field must exist as a non-empty string array. No separate `required: true` flag in this change.

### 3. Rule IDs

| File | Invalid / missing list | Missing `must_include` | Outside `allowed` |
|------|------------------------|------------------------|-------------------|
| README.md | `skill/readme-field-invalid` | `skill/readme-field-required-value` | `skill/readme-field-allowed-value` |
| SKILL.md | `skill/skill-field-invalid` | `skill/skill-field-required-value` | `skill/skill-field-allowed-value` |

Delete the four tag/phase rule constants.

### 4. Allowlist carve-out

During SKILL.md frontmatter allowlist checks, every `field` named in `skill_constraints` is treated as allowed. Constraint validation still runs and may emit `skill/skill-field-*` diagnostics. Constrained fields that collide with agentskills keys (e.g. `name`) keep their existing type rules; string-list constraints on non-list agentskills fields will fail `skill-field-invalid` as expected.

### 5. Config validation

On `Load`, reject any constraint (in either list) with empty `field`, or with both `must_include` and `allowed` empty.

### 6. Validation loop

Shared helper applies string-list logic; callers pass the rule ID trio and path label (`README.md` / `SKILL.md`). README loop runs only when `require_readme` is true; skill loop runs whenever `skill_constraints` is non-empty.

## Risks / Trade-offs

- [Ignore files referencing old rule IDs] → Unknown IDs already fail load with exit 2; callers must update ignore files when upgrading.
- [Org taxonomy on SKILL.md diverges from agentskills.io] → Accepted: opt-in via config, documented as an extension; unconstrained unknown fields still error.
- [Old keys silently ignored] → Document removal in `docs/configuration.md`.

## Migration Plan

1. Ship the config/rule swap in the same release (pre-v1.0; no deprecation window).
2. Docs example shows tags + phases on README and optional skill constraints.
3. No runtime shim.
