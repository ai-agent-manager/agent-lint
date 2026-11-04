# agent-lint

Lint [Agent Skills](https://agentskills.io/specification) (`SKILL.md`) with tight adherence to the agentskills.io specification — including hard caps on `SKILL.md` body size — and structurally check Rovo agents (`rovo-agent.yaml`).

Designed as a small Go CLI with checksummed GitHub release binaries and Homebrew install (no npm runtime dependency).

**User documentation:** [`docs/`](docs/README.md) · **Examples:** [`examples/`](examples/README.md)

## Install

```bash
brew install ai-agent-manager/tap/agent-lint
```

Other options (release binary, `go install`): [docs/install.md](docs/install.md).

## Quick start

```bash
agent-lint check --root skills
agent-lint check --config .agent-lint.yaml --ignore-file .agent-lint-ignore.yaml --json
agent-lint list-rules
```

Exit codes: `0` success, `1` lint failures, `2` tool/usage error.

Full flag reference and examples: [docs/check.md](docs/check.md).

## What it checks

Normative authority: [agentskills.io/specification](https://agentskills.io/specification). Frontmatter behaviour aims for parity with [`skills-ref`](https://github.com/agentskills/agentskills/tree/main/skills-ref).

Highlights:

- Required `SKILL.md` frontmatter (`name`, `description`) with agentskills.io naming rules
- Spec allowlist for frontmatter fields (Claude Code extensions opt-in)
- **Body ≤ 500 lines** and **&lt; 5000 tokens** as errors by default (`ceil(rune_count / 4)` estimator)
- Relative markdown references: one level deep; linked paths must exist
- Rovo `rovo-agent.yaml` against the vendored agent-manager JSON Schema (Studio `v2-beta`)

Rule catalogue: [docs/rules.md](docs/rules.md). Config and ignores: [docs/configuration.md](docs/configuration.md), [docs/ignore.md](docs/ignore.md).

## Roadmap

- Common Agent Format (`AGENT.md`) checker
- `--fix` for common issues
- Optional GitHub Action wrapper

## License

Apache License 2.0
