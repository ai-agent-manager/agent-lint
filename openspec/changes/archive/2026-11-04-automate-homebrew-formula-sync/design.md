## Context

GoReleaser already publishes archives + `checksums.txt`. The tap formula is a hand-maintained Ruby Formula. `brews` is deprecated; we deliberately keep Formulas (option 2) rather than casks. See proposal.md.

## Goals / Non-Goals

**Goals:**

- Idempotent script that rewrites `Formula/agent-lint.rb` for a given version using release checksums
- Invoke from Release workflow after GoReleaser succeeds
- Preserve `brew install ai-agent-manager/tap/agent-lint`

**Non-Goals:**

- `homebrew_casks` migration
- Opening PRs instead of committing (direct push to tap `main` is fine for this org tap)

## Decisions

1. **Script in agent-lint** — `scripts/update-homebrew-formula.sh` generates the full formula from a here-doc template (version + four platform sha256s). Avoid fragile in-place sed on an evolving file.
2. **Checksums source** — Download `checksums.txt` from the just-published GitHub Release for the tag (same assets GoReleaser uploaded).
3. **Auth** — Clone/push tap using configured credentials from the Release workflow environment. Skip (exit 0 + warn) if unset.
4. **Platforms in formula** — darwin/linux × amd64/arm64 only (match current formula; windows stays release-only).
5. **Workflow** — New step after GoReleaser; pass tag version via `github.ref_name`.

## Risks / Trade-offs

- [Missing tap credentials] → Skip with warning so release assets still publish
- [Tap push conflicts] → Rare; fail the step loudly so maintainers notice
- [Formula template drift] → Keep template next to script comments; smoke-test locally against a dry-run flag

## Migration Plan

- Merge, next `v*` tag updates the tap when credentials are configured
- No consumer-facing change
