# agent-lint documentation

User guide for [agent-lint](https://github.com/ai-agent-manager/agent-lint) — a CLI that lints [Agent Skills](https://agentskills.io/specification) (`SKILL.md`) against the agentskills.io format (including hard caps on skill body size) and structurally validates Rovo agents (`rovo-agent.yaml`).

## Contents

| Guide | What it covers |
|-------|----------------|
| [Install](install.md) | Homebrew, release binaries, `go install` |
| [Checking skills](check.md) | `agent-lint check`, flags, exit codes, JSON output |
| [Configuration](configuration.md) | `.agent-lint.yaml` options |
| [Ignoring findings](ignore.md) | `.agent-lint-ignore.yaml` paths, globs, and rule IDs |
| [Rules](rules.md) | What is checked and how to list rule IDs |
| [CI](ci.md) | Running agent-lint in GitHub Actions and other pipelines |

Worked examples (runnable mini-projects): [`examples/`](../examples/README.md).

## Quick start

```bash
brew install ai-agent-manager/tap/agent-lint

# From a repo whose skills live under ./skills
agent-lint check

# Or point at a directory of skill folders
agent-lint check --root path/to/skills
```

Exit codes: `0` success, `1` lint failures, `2` tool or usage error.
