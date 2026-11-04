## Why

Homebrew formula updates for `ai-agent-manager/homebrew-tap` are done by hand after each tagged release. That is easy to forget and drifts from published checksums. Automate formula sync as part of the Release workflow while keeping Formula install UX (not migrating to casks).

## What Changes

- After a successful GoReleaser publish on `v*` tags, update `Formula/agent-lint.rb` in `ai-agent-manager/homebrew-tap` from the release’s `checksums.txt`
- Add a maintainable sync script in this repo; wire it into `.github/workflows/release.yml`
- Require a tap write token secret (tap credentials); skip with a clear message when unset (forks)
- Update AGENTS.md / install docs so maintainers know the secret and that hand edits are no longer the primary path
- **Non-goals:** GoReleaser `homebrew_casks`, changing `brew install` UX, npm distribution

## Capabilities

### New Capabilities

- _(none)_

### Modified Capabilities

- `release-and-ci`: Homebrew tap formula SHALL be updated automatically on tagged releases when the tap token is configured

## Impact

- Release workflow + new script under `scripts/`
- Consumer install command unchanged
