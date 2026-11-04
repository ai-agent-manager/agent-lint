## MODIFIED Requirements

### Requirement: Homebrew Tap

A Homebrew formula for `agent-lint` SHALL be maintainable under `ai-agent-manager/homebrew-tap`, and SHALL be updated automatically after a tagged GitHub Release when a tap write token is configured.

#### Scenario: Formula install

- **WHEN** a user installs via `brew install ai-agent-manager/tap/agent-lint`
- **THEN** the formula points at a published release archive and verifies its checksum

#### Scenario: Automatic formula sync on release

- **WHEN** a `v*` tag Release completes successfully and tap credentials is available to the Release workflow
- **THEN** `Formula/agent-lint.rb` in `ai-agent-manager/homebrew-tap` is updated to that version’s darwin/linux archive URLs and SHA-256 checksums and pushed to the tap

#### Scenario: Missing tap token

- **WHEN** a `v*` tag Release completes but tap credentials is not configured
- **THEN** the GitHub Release still publishes successfully and the Homebrew sync step is skipped with an explicit warning
