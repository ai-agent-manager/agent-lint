# Ignoring findings

Use an ignore file when you need a deliberate, temporary, or false-positive suppression. Prefer fixing the skill when you can.

Default path: `.agent-lint-ignore.yaml` (or `--ignore-file`). You can also embed the same entries under `ignore.ignores` in `.agent-lint.yaml`.

## Format

```yaml
ignores:
  # Empty rules list = ignore ALL findings under this path
  - path: skills/legacy-migration/**
    rules: []

  # Only suppress named rule IDs for this file
  - path: skills/big-skill/SKILL.md
    rules:
      - skill/body-line-limit
      - skill/body-token-limit

  # Glob across many skills
  - path: "skills/*-evals/**"
    rules:
      - skill/name-directory-mismatch
```

## Semantics

- **`path`**: exact relative path or glob (`*`, `**`, `?`), resolved from the working directory / scan context. Matches skill dirs, files, or trees.
- **`rules`**: list of stable rule IDs from `agent-lint list-rules`.
  - **`[]` (empty)** → suppress every finding on matching paths.
  - **One or more IDs** → suppress only those codes; other findings still fail the build.
- Unknown rule IDs are a **tool error** (exit `2`) so typos do not silently pass.
- Ignored findings do not count toward exit `1`. With `--json`, they still appear under `ignored` for audit.

## Tips

- Keep ignores narrow (specific skill paths) when possible.
- Comment why an ignore exists; treat broad `skills/**` ignores as migration debt.
- Re-run `agent-lint list-rules` after upgrading agent-lint so IDs stay accurate.
