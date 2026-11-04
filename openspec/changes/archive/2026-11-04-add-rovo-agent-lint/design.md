## Context

Discovery already classifies Rovo-only dirs. agent-manager owns the canonical manifest schema at `schemas/rovo-agent.schema.json` (v2-beta only; v1 archived). See proposal.md for motivation.

## Goals / Non-Goals

**Goals:**

- Vendor that schema into agent-lint and validate parsed YAML against it
- Wire discovery so Rovo-only dirs invoke `internal/rovo`
- Reuse README helpers when `require_readme` is on

**Non-Goals:**

- CAF
- Runtime `$file` expansion (schema allows `{ $file }` shape; existence checks can follow later)
- Maintaining a parallel hand-rolled field matrix for v1

## Decisions

1. **Vendor path** — `internal/rovo/schemas/rovo-agent.schema.json` embedded with `go:embed`. Copy from agent-manager; document provenance + refresh steps in a short README next to the file.
2. **Library** — `github.com/santhosh-tekuri/jsonschema/v6` (Draft 2020-12). Alternative: hand-port Ajv — rejected.
7. **Regex engine** — `dlclark/regexp2` (ECMAScript) via `Compiler.UseRegexpEngine`, because the vendored `$file` pattern uses lookaheads RE2 cannot compile.
3. **Rule IDs** — `rovo/yaml-invalid`, `rovo/yaml-not-object` for parse; single `rovo/schema-invalid` for all schema errors (message includes JSON Pointer / schema path). Drop the earlier hand-rolled field-specific IDs introduced in this same change before release.
4. **v1** — Not supported; fails schema (`apiVersion` const). Examples use v2-beta.
5. **Handoff** — `skill.Check` calls `rovo.Check` for Rovo-only dirs; README checks stay in the walker.
6. **YAML → JSON** — Unmarshal with `yaml.v3` into `any`, then validate; rely on map/slice shapes JSON Schema expects.

## Risks / Trade-offs

- [Schema drift vs agent-manager] → Vendor explicitly; bump when agent-manager schema changes; note version in file header comment
- [Stricter than old Node library script] → Document in release notes; consumers may need small yaml fixes
- [Shared README IDs] → Same as before; accept for this slice

## Migration Plan

- Same PR as initial Rovo support (unreleased rule IDs may be renamed once)
- Refresh vendored schema when agent-manager changes it
