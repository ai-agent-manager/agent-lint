## Why

Rovo-only directories (`rovo-agent.yaml`, no `SKILL.md`) were skipped, so consumers kept a separate Node validator. agent-lint should lint them using the same JSON Schema agent-manager already publishes for Studio provisioning, so libraries share one source of truth.

## What Changes

- Validate Rovo-only directories instead of skipping them
- **Vendor** `rovo-agent.schema.json` from agent-manager and validate `rovo-agent.yaml` against it (JSON Schema 2020-12; **v2-beta only**)
- Emit `rovo/yaml-*` for parse failures and `rovo/schema-invalid` for schema violations
- Apply existing README config checks (`require_readme`, tags, phases) to Rovo dirs when enabled
- Update docs/examples; migrate example agents off archived v1
- **Non-goals:** CAF/`AGENT.md`, resolving/`$file` content loading beyond schema shape, Atlassian API calls, soft-passing schema errors

## Capabilities

### New Capabilities

- `rovo-agent`: Structural lint for `rovo-agent.yaml` via vendored JSON Schema + optional README via config

### Modified Capabilities

- `skill-discovery`: Rovo-only directories are linted by the Rovo checker rather than skipped
- `cli`: `check` reports Rovo diagnostics; `list-rules` includes `rovo/…` IDs

## Impact

- New package `internal/rovo` with embedded schema; Go JSON Schema dependency
- Discovery handoff from `internal/skill`; fixtures/docs/examples updates
- Consumers can drop Node Rovo validators after upgrading (note: schema is stricter than the old hand checks — e.g. name ≤30, ≤3 conversation starters, `knowledge` enum)
