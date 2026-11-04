## 1. Vendor schema and rules

- [x] 1.1 Copy agent-manager `rovo-agent.schema.json` into `internal/rovo/schemas/` with provenance note
- [x] 1.2 Replace hand-rolled `rovo/…` field rules with `rovo/yaml-invalid`, `rovo/yaml-not-object`, `rovo/schema-invalid`

## 2. Schema validation

- [x] 2.1 Add `santhosh-tekuri/jsonschema/v6` dependency and compile embedded schema once
- [x] 2.2 Rewrite `internal/rovo.Check` to parse YAML then validate against the schema
- [x] 2.3 Map each schema error to a `rovo/schema-invalid` diagnostic with location + message

## 3. Discovery wiring

- [x] 3.1 Call `rovo.Check` for Rovo-only directories instead of skipping
- [x] 3.2 Keep conflict / missing-asset behaviour unchanged
- [x] 3.3 Apply README config checks for Rovo dirs

## 4. Fixtures and tests

- [x] 4.1 Update fixtures for v2-beta schema (drop/replace v1-valid fixture; cover schema violations)
- [x] 4.2 Update unit/integration tests for `rovo/schema-invalid`
- [x] 4.3 `go test ./...` green; smoke against a consumer skills root with Rovo agents

## 5. Docs and examples

- [x] 5.1 Document schema-based Rovo rules and vendored schema source
- [x] 5.2 Migrate getting-started Rovo example to v2-beta
- [x] 5.3 Validate OpenSpec change
