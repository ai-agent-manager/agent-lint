# agent-lint

Lint [Agent Skills](https://agentskills.io/specification) (`SKILL.md`) with tight adherence to the agentskills.io specification — including hard caps on `SKILL.md` body size.

Designed as a small Go CLI with checksummed GitHub release binaries and Homebrew install (no npm runtime dependency).

## Install

### Homebrew

```bash
brew install ai-agent-manager/tap/agent-lint
```

### GitHub release binary

Download the archive for your OS/arch from [Releases](https://github.com/ai-agent-manager/agent-lint/releases), verify against `checksums.txt`, and place `agent-lint` on your `PATH`.

### From source

```bash
go install github.com/ai-agent-manager/agent-lint/cmd/agent-lint@latest
```

## Usage

```bash
agent-lint check --root skills
agent-lint check --config .agent-lint.yaml --ignore-file .agent-lint-ignore.yaml --json
agent-lint list-rules
```

Exit codes: `0` success, `1` lint failures, `2` tool/usage error.

### Flags

| Flag | Description |
|------|-------------|
| `--root` | Directory containing skill folders (default: `skills`) |
| `--config` | Path to `.agent-lint.yaml` |
| `--ignore-file` | Path to `.agent-lint-ignore.yaml` |
| `--json` | Machine-readable JSON report |
| `--quiet` / `-q` | Suppress summary line |
| `--strict` | Treat warnings as errors |
| `--claude` | Enable Claude.ai-specific checks |
| `--allow-claude-code` | Allow Claude Code frontmatter extensions |

## What it checks

Normative authority: [agentskills.io/specification](https://agentskills.io/specification). Frontmatter behaviour aims for parity with [`skills-ref`](https://github.com/agentskills/agentskills/tree/main/skills-ref).

### Errors (agentskills.io)

- Required `SKILL.md` with YAML frontmatter
- `name`: 1–64 chars, lowercase alphanumeric + hyphens, no leading/trailing/consecutive hyphens, must match directory name
- `description`: required, max 1024 characters
- `compatibility`: max 500 characters when present
- `metadata`: string→string map when present
- `allowed-tools`: space-separated string when present
- Unknown frontmatter fields (spec allowlist only)
- Duplicate names within the scan root
- **Body ≤ 500 lines** and **&lt; 5000 tokens** (enforced as errors; estimator = `ceil(rune_count / 4)`)
- Relative file references at most one directory level deep; referenced paths must exist

### Warnings

- Empty body
- Orphan files under `references/`

### Library config (optional)

`.agent-lint.yaml` can require README frontmatter tags/phases and tune body caps:

```yaml
skills:
  max_body_lines: 500
  max_body_tokens: 5000
  require_readme: true
  required_tags: ["agent-skill"]
  allowed_phases: [visioning, discovery, elaboration, build_and_release, run_and_measure]
  allow_claude_code_fields: false
```

### Ignore file

`.agent-lint-ignore.yaml`:

```yaml
ignores:
  - path: skills/legacy/**
    rules: []   # empty = ignore all findings on matching paths
  - path: skills/big-skill/SKILL.md
    rules:
      - skill/body-line-limit
      - skill/body-token-limit
```

Unknown rule IDs in the ignore file are a tool error (exit 2). Run `agent-lint list-rules` for IDs.

Directories with only `rovo-agent.yaml` are skipped (Rovo checker planned later). Both `SKILL.md` and `rovo-agent.yaml` in one directory is an error.

## Token estimator

Body token count uses a documented rough estimator: `ceil(utf8_rune_count / 4)`. This matches common heuristic estimators used by other skill linters and fails closed when over the configured max (default 5000).

## Roadmap

- Rovo agent checker
- `--fix` for common issues
- Optional GitHub Action wrapper

## License

Apache License 2.0
