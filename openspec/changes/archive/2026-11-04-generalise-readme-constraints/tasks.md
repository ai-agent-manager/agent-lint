## 1. Config and rules

- [x] 1.1 Replace list operators with `FrontmatterConstraint` string operators (`equals` / `not_equals` / `matches` / `not_matches`), dotted `field`, validate exactly one operator + regexp on load
- [x] 1.2 Replace list-oriented rule IDs with `*-field-invalid` plus per-operator IDs for README and SKILL.md
- [x] 1.3 Add config tests for valid entries, multi/zero operator rejection, and bad regexp

## 2. Validation

- [x] 2.1 Resolve dotted paths to strings; apply operators via shared helper for README and SKILL.md
- [x] 2.2 Carve out only the root segment of each `skill_constraints` field from the allowlist
- [x] 2.3 Update fixtures/tests for string fields and operator failures; cover `metadata.*` paths

## 3. Docs

- [x] 3.1 Update `docs/configuration.md` and `docs/rules.md`
- [x] 3.2 Run `go test ./...` and OpenSpec validate
