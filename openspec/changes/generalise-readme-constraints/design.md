## Context

See proposal.md — Why. agentskills.io `metadata` is string→string; list-shaped constraints were a README-era leftover. This change unifies both files on string operators and dotted paths.

## Goals / Non-Goals

**Goals:**

- Parallel `readme_constraints` / `skill_constraints`
- Exactly one operator per entry: `equals` | `not_equals` | `matches` | `not_matches`
- Dotted field paths (`metadata.phases`)
- Validate regexp syntax at config load
- File- and operator-specific rule IDs

**Non-Goals:**

- YAML sequence operators
- Compatibility shims for `must_include` / `allowed` / old keys

## Decisions

### 1. Config shape

```yaml
skills:
  require_readme: true
  readme_constraints:
    - field: name
      equals: Valid Skill
    - field: description
      matches: "(?i)fixture"
  skill_constraints:
    - field: metadata.phases
      matches: "^(visioning|discovery|elaboration|build_and_release|run_and_measure)$"
    - field: metadata.author
      not_equals: ""
```

Shared type `FrontmatterConstraint` with `Field` and optional pointer/presence for the four operators; load rejects 0 or >1 operators and invalid regexps.

### 2. Path resolution

Split `field` on `.`. Walk maps; final value MUST be a YAML string (including empty). Missing path or non-string → `*-field-invalid`.

### 3. Rule IDs

| | README | SKILL.md |
|--|--------|----------|
| Missing / not string | `skill/readme-field-invalid` | `skill/skill-field-invalid` |
| `equals` failed | `skill/readme-field-equals` | `skill/skill-field-equals` |
| `not_equals` failed | `skill/readme-field-not-equals` | `skill/skill-field-not-equals` |
| `matches` failed | `skill/readme-field-matches` | `skill/skill-field-matches` |
| `not_matches` failed | `skill/readme-field-not-matches` | `skill/skill-field-not-matches` |

### 4. Allowlist carve-out

Add only the **root** segment of each `skill_constraints` field to the SKILL.md allowlist. `metadata.phases` → root `metadata` (already allowed). Top-level `phases` → carve-out `phases`.

### 5. Shared validation helper

One helper resolves the path, applies the operator, and emits the chosen rule ID trio/quintet for the file.

## Risks / Trade-offs

- [README fixtures still using YAML lists] → Update fixtures/docs to strings; list fields fail `*-field-invalid`.
- [Complex allowlists via regex] → Document alternation patterns; no first-class enum list.

## Migration Plan

Ship in the same pre-v1.0 release with no shim. Docs show `metadata.*` for SKILL.md and string operators for both lists.
