# Checking skills

## Basic usage

```bash
# Default: scan ./skills
agent-lint check

# Custom skills root (each child directory is one skill)
agent-lint check --root path/to/skills

# With project config and ignore file
agent-lint check \
  --config .agent-lint.yaml \
  --ignore-file .agent-lint-ignore.yaml
```

If you omit `--config` / `--ignore-file`, agent-lint loads `.agent-lint.yaml` and `.agent-lint-ignore.yaml` from the current working directory when those files exist.

## What gets scanned

- Each **immediate child directory** of `--root` is a candidate skill.
- Directories with `SKILL.md` are linted as Agent Skills.
- Directories with only `rovo-agent.yaml` are skipped (Rovo checking is not implemented yet).
- Both files in one directory, or neither, is an error.

## Flags

| Flag | Description |
|------|-------------|
| `--root` | Directory containing skill folders (default: `skills`) |
| `--config` | Path to `.agent-lint.yaml` |
| `--ignore-file` | Path to `.agent-lint-ignore.yaml` |
| `--json` | Emit a machine-readable JSON report on stdout |
| `--quiet` / `-q` | Suppress the human summary line |
| `--strict` | Treat warnings as errors (exit `1`) |
| `--claude` | Enable Claude.ai-specific checks |
| `--allow-claude-code` | Allow Claude Code frontmatter extensions |

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success (no active errors; warnings allowed unless `--strict`) |
| `1` | Lint failures (active errors, or warnings with `--strict`) |
| `2` | Tool or usage error (missing root, bad config, unknown ignore rule ID, …) |

## JSON output

```bash
agent-lint check --json > report.json
```

The report includes active diagnostics, ignored diagnostics (for audit), and counts. Use this in CI scripts or coding agents. Each finding has a stable `ruleId` — see [Rules](rules.md).

## Example layouts

```text
my-repo/
  skills/
    pdf-processing/
      SKILL.md
      references/
        forms.md
    code-review/
      SKILL.md
  .agent-lint.yaml          # optional
  .agent-lint-ignore.yaml   # optional
```

```bash
cd my-repo
agent-lint check
```
