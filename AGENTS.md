# AGENTS.md

Guidance for coding agents working on **agent-lint**.

## What this is

A small Go CLI that lints Agent Skills (`SKILL.md`) with tight adherence to the [agentskills.io specification](https://agentskills.io/specification). Frontmatter behaviour should stay compatible with [`skills-ref`](https://github.com/agentskills/agentskills/tree/main/skills-ref).

Body size limits (≤500 lines, &lt;5000 tokens) are **errors by default** even though the spec phrases them as recommendations. Only override via explicit config (`skills.max_body_lines` / `skills.max_body_tokens`).

Do not introduce an npm runtime or publish this tool to the npm registry. Distribution is GoReleaser binaries + Homebrew.

## Commands

```bash
go test ./...
go vet ./...
go build -o /tmp/agent-lint ./cmd/agent-lint

/tmp/agent-lint check --root testdata
/tmp/agent-lint list-rules
/tmp/agent-lint check --json --config .agent-lint.yaml --ignore-file .agent-lint-ignore.yaml
```

Exit codes: `0` success, `1` lint findings (or warnings under `--strict`), `2` tool/usage/config error.

Before finishing a change: `go test ./...` must be green. Add or update tests for behaviour you change.

## Architecture (where to edit)

| Concern | Package | Notes |
|--------|---------|--------|
| CLI flags / commands | `cmd/agent-lint` | Cobra only; keep thin |
| Skill walk + agentskills.io rules | `internal/skill` | Main validation logic |
| Frontmatter parse | `internal/frontmatter` | YAML between `---` fences |
| Stable rule IDs + default severities | `internal/rules` | Every diagnostic must use an ID from here |
| Config (`.agent-lint.yaml`) | `internal/config` | Defaults encode spec progressive-disclosure caps |
| Path/glob ignores | `internal/ignore` | Empty `rules: []` = ignore all findings on that path |
| Human / JSON reporters | `internal/report` | JSON must include ignored findings for audit |

Future asset types (e.g. Rovo) should be new packages implementing the same diagnostic model — do not fold unrelated formats into `internal/skill`.

## Rule design rules

- **Stable IDs** use `skill/<kebab-name>` (see `agent-lint list-rules`). Renaming an ID is a breaking change for ignore files.
- Prefer **errors** for agentskills.io normative constraints and for the enforced body/reference limits.
- Prefer **warnings** for extras beyond the spec (empty body, orphan `references/` files).
- Claude.ai / Claude Code extras are opt-in (`--claude`, `--allow-claude-code` / config). Do not make them default errors.
- Unknown ignore rule IDs must fail the tool (exit 2), not silently no-op.
- Token estimate is documented as `ceil(utf8_rune_count / 4)` — keep README and code aligned if you change it.

## Fixtures and tests

- Put skill fixtures under `testdata/<skill-dir>/`.
- Prefer table-driven tests that assert rule IDs, not only message substrings.
- When adding a rule, cover at least one failing and one passing case.

## Releases

- CI on `main` / PRs: vet, test, build, smoke CLI, `goreleaser check`.
- Tag `v*` pushes run GoReleaser and publish a GitHub Release with checksums.
- Homebrew formula lives in `ai-agent-manager/homebrew-tap`. GoReleaser brew upload is currently skipped unless configured; prefer updating the tap formula when cutting a release consumers need via Homebrew.
- Do not force-push release tags that already have published assets without an explicit human request.

## Commits and PRs

- Conventional commits: `feat:`, `fix:`, `docs:`, `ci:`, `chore:`, `test:`, `refactor:`.
- Keep PRs focused. Spec/rule behaviour changes should mention agentskills.io implications in the PR body.
- Do not commit secrets, local binaries, or `dist/`.

## User documentation

End-user how-to guides live in `docs/` (install, check, config, ignores, rules, CI). Keep that directory **user-facing** — not maintainer architecture, contribution process, or release engineering notes (those belong in `AGENTS.md`, OpenSpec, or workflow comments).

## Specs (OpenSpec)

Behaviour contracts live under `openspec/specs/`.

**All changes that are covered by an existing spec, or that change user-visible behaviour, MUST go through the OpenSpec flow** (`/opsx:propose` → apply → archive): propose a change, update or add the relevant spec deltas, implement against that change, then archive so main specs stay the source of truth. Do not land such work as a code-only PR that skips OpenSpec.

That includes (non-exhaustive): lint rules and severities, CLI flags/commands/exit codes, config or ignore-file semantics, reporters/JSON shape, and release/CI behaviour users rely on.

Validate with `npx @fission-ai/openspec validate --specs --all`.

## Out of scope (for now)

- Skill installers, cross-agent “flavors”, auto-fix, LLM/semantic quality reviews.
- Full Rovo `rovo-agent.yaml` validation (planned as a separate checker).
- Soft-passing body/token limits without an explicit config override.
