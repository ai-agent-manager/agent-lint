## 1. Config and rules

- [x] 1.1 Replace `RequiredTags` / `AllowedPhases` with shared `FrontmatterConstraint` lists (`readme_constraints`, `skill_constraints`) in `internal/config`, validate entries on load
- [x] 1.2 Replace tag/phase rule IDs with `skill/readme-field-*` and add `skill/skill-field-*` counterparts in `internal/rules`
- [x] 1.3 Add config unit tests for valid constraints and empty-entry rejection (both lists)

## 2. Validation

- [x] 2.1 Rewrite `validateReadme` to apply `readme_constraints` with shared string-list logic
- [x] 2.2 Apply `skill_constraints` in SKILL.md frontmatter validation and carve constrained fields out of the unknown-field allowlist
- [x] 2.3 Update skill tests and fixtures for README constraints; cover missing required value, allowlist miss, and invalid field type
- [x] 2.4 Add skill-constraint tests (valid carve-out, required/allowed/invalid, unconstrained unknown still fails)

## 3. Docs and examples

- [x] 3.1 Update `docs/configuration.md` and `docs/rules.md` for both constraint lists and rule IDs
- [x] 3.2 Confirm examples do not still show `required_tags` / `allowed_phases`
- [x] 3.3 Run `go test ./...` and `npx @fission-ai/openspec validate --specs --all`
