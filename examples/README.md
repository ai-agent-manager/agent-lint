# Examples

Self-contained mini-projects that show how **agent-lint** discovers skills, reports findings, and applies ignores.

From the repository root (with `agent-lint` on your `PATH`):

```bash
# Clean pass — discovery only
agent-lint check --root examples/getting-started/skills

# Intentional failures — see how findings look
agent-lint check --root examples/findings/skills

# Same failures, suppressed via ignore file
agent-lint check \
  --root examples/ignoring-findings/skills \
  --ignore-file examples/ignoring-findings/.agent-lint-ignore.yaml
```

| Example | What it demonstrates |
|---------|----------------------|
| [getting-started](getting-started/README.md) | Layout, discovery, a clean `check` |
| [findings](findings/README.md) | Common errors the linter reports |
| [ignoring-findings](ignoring-findings/README.md) | Path/glob ignores and rule IDs |

User docs: [`docs/`](../docs/README.md).
