# Release And CI Specification

## Purpose

agent-lint is verified on every push/PR and published as checksummed multi-platform binaries on version tags, with optional Homebrew formula updates in `ai-agent-manager/homebrew-tap`. Dependency updates are grouped weekly via Dependabot.

## Requirements

### Requirement: Continuous Integration Gate

Pushes to `main` and pull requests SHALL run vet, unit tests, a binary build, a CLI smoke check, and a GoReleaser configuration check.

#### Scenario: PR validation

- **WHEN** a pull request is opened or updated
- **THEN** the CI workflow runs `go vet`, `go test`, `go build`, a smoke invocation of the built binary, and `goreleaser check`

### Requirement: Tagged Releases

Pushing a tag matching `v*` SHALL run GoReleaser to publish a GitHub Release with archives and `checksums.txt`.

#### Scenario: Tag push

- **WHEN** a tag such as `v0.1.1` is pushed
- **THEN** the Release workflow builds platform archives and publishes them to the GitHub Release for that tag

#### Scenario: Checksums published

- **WHEN** a release is published
- **THEN** a `checksums.txt` asset is available listing SHA-256 hashes for each archive

### Requirement: Supported Platforms

Release builds SHALL include linux and darwin (amd64 and arm64) and windows amd64, with CGO disabled.

#### Scenario: Archive naming

- **WHEN** archives are produced
- **THEN** names follow `agent-lint_<version>_<os>_<arch>` with `.tar.gz` (or `.zip` on Windows)

### Requirement: Homebrew Tap

A Homebrew formula for `agent-lint` SHALL be maintainable under `ai-agent-manager/homebrew-tap`, and SHALL be updated automatically after a tagged GitHub Release when tap credentials are configured for the Release workflow.

#### Scenario: Formula install

- **WHEN** a user installs via `brew install ai-agent-manager/tap/agent-lint`
- **THEN** the formula points at a published release archive and verifies its checksum

#### Scenario: Automatic formula sync on release

- **WHEN** a `v*` tag Release completes successfully and tap credentials are available to the Release workflow
- **THEN** `Formula/agent-lint.rb` in `ai-agent-manager/homebrew-tap` is updated to that version’s darwin/linux archive URLs and SHA-256 checksums and pushed to the tap

#### Scenario: Missing tap credentials

- **WHEN** a `v*` tag Release completes but tap credentials are not configured
- **THEN** the GitHub Release still publishes successfully and the Homebrew sync step is skipped with an explicit warning

### Requirement: Dependabot Grouping

Dependabot SHALL open weekly grouped updates for Go modules and GitHub Actions.

#### Scenario: Grouped Go updates

- **WHEN** Dependabot runs for the `gomod` ecosystem
- **THEN** dependency updates are grouped into a single `go-modules` pull request pattern group

#### Scenario: Grouped Actions updates

- **WHEN** Dependabot runs for the `github-actions` ecosystem
- **THEN** Action version updates are grouped into a single `github-actions` pull request pattern group

### Requirement: Apache-2.0 License

The project SHALL be licensed under Apache License 2.0.

#### Scenario: License file present

- **WHEN** a consumer inspects the repository root
- **THEN** a `LICENSE` file containing the Apache License 2.0 text is present
