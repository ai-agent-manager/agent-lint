## 1. Config and rules

- [ ] 1.1 Replace `RequiredTags` / `AllowedPhases` with `ReadmeConstraints` in `internal/config`, validate entries on load
- [ ] 1.2 Replace tag/phase rule IDs with `skill/readme-field-invalid`, `skill/readme-field-required-value`, and `skill/readme-field-allowed-value` in `internal/rules`
- [ ] 1.3 Add config unit tests for valid constraints and empty-entry rejection

## 2. README validation

- [ ] 2.1 Rewrite `validateReadme` to apply `readme_constraints` with the shared string-list logic
- [ ] 2.2 Update skill tests and fixtures to use the new config shape
- [ ] 2.3 Cover missing required value, allowlist miss, and invalid field type via table-driven tests

## 3. Docs and examples

- [ ] 3.1 Update `docs/configuration.md` and `docs/rules.md` for the new keys and rule IDs
- [ ] 3.2 Update any example configs that still show `required_tags` / `allowed_phases`
- [ ] 3.3 Run `go test ./...` and `npx @fission-ai/openspec validate --specs --all`
