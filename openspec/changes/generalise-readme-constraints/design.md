## Context

See proposal.md — Why. agentskills.io `metadata` is string→string; list-shaped constraints were a README-era leftover. This change unifies both files on string operators and dotted paths.

## Goals / Non-Goals

**Goals:**

- Parallel `readme_constraints` / `skill_constraints`
- String operators on both; list operators (`must_include` / `allowed`) on README only
- Dotted field paths (`metadata.phases`)
- Validate regexp syntax at config load
- File- and operator-specific rule IDs

**Non-Goals:**

- List operators on `skill_constraints` (agentskills.io metadata is string-valued)
- Compatibility shims for old `required_tags` / `allowed_phases` keys

## Decisions

### 1. Config shape

```yaml
skills:
  require_readme: true
  readme_constraints:
    - field: description
      matches: "(?i)fixture"
    - field: tags
      must_include: [agent-skill]
    - field: phases
      allowed: [visioning, discovery, elaboration, build_and_release, run_and_measure]
  skill_constraints:
    - field: metadata.phases
      matches: "^(visioning|discovery|elaboration|build_and_release|run_and_measure)$"
    - field: metadata.author
      not_equals: ""
```

Shared type `FrontmatterConstraint`. README may use string ops (exactly one) or list ops (`must_include` and/or `allowed`). Skill constraints accept string ops only. Load rejects mixes, empty operator sets, invalid regexps, and list ops on `skill_constraints`.

### 2. Path resolution

Split `field` on `.`. Walk maps; final value MUST be a YAML string (including empty). Missing path or non-string → `*-field-invalid`.

### 3. Rule IDs

| | README | SKILL.md |
|--|--------|----------|
| Missing / wrong type | `skill/readme-field-invalid` | `skill/skill-field-invalid` |
| `equals` failed | `skill/readme-field-equals` | `skill/skill-field-equals` |
| `not_equals` failed | `skill/readme-field-not-equals` | `skill/skill-field-not-equals` |
| `matches` failed | `skill/readme-field-matches` | `skill/skill-field-matches` |
| `not_matches` failed | `skill/readme-field-not-matches` | `skill/skill-field-not-matches` |
| `must_include` failed | `skill/readme-field-required-value` | _(n/a)_ |
| `allowed` failed | `skill/readme-field-allowed-value` | _(n/a)_ |

### 4. Allowlist carve-out

Add only the **root** segment of each `skill_constraints` field to the SKILL.md allowlist. `metadata.phases` → root `metadata` (already allowed). Top-level `phases` → carve-out `phases`.

### 5. Shared validation helper

One helper resolves the path, applies the operator, and emits the chosen rule ID trio/quintet for the file.

## Risks / Trade-offs

- [README fixtures still using YAML lists] → Update fixtures/docs to strings; list fields fail `*-field-invalid`.
- [Complex allowlists via regex] → Document alternation patterns; no first-class enum list.

## Migration Plan

Ship in the same pre-v1.0 release with no shim. Docs show `metadata.*` for SKILL.md and string operators for both lists.
