# Ignoring findings

Same kinds of problems as [findings](../findings/README.md), plus an ignore file that suppresses them by **path/glob** and **rule ID**.

## Layout

```text
ignoring-findings/
  .agent-lint-ignore.yaml
  skills/
    Bad_Name/
      SKILL.md
    empty-body/
      SKILL.md
    missing-link/
      SKILL.md
    legacy-demo/           ← suppressed entirely via empty rules: []
      SKILL.md
```

## Compare with and without ignores

```bash
# Without ignore file — expect failures (exit 1)
agent-lint check --root examples/ignoring-findings/skills
echo "without ignore: exit=$?"

# With ignore file — expect success (exit 0)
agent-lint check \
  --root examples/ignoring-findings/skills \
  --ignore-file examples/ignoring-findings/.agent-lint-ignore.yaml
echo "with ignore: exit=$?"
```

## What the ignore file does

See [`.agent-lint-ignore.yaml`](.agent-lint-ignore.yaml):

| Match | Rules | Effect |
|-------|-------|--------|
| `Bad_Name/SKILL.md` | `skill/name-format` | Suppresses the bad name format error only |
| `empty-body/SKILL.md` | `skill/body-empty` | Suppresses the empty-body warning |
| `missing-link/SKILL.md` | `skill/reference-missing` | Suppresses the broken markdown link |
| `legacy-demo/**` | `[]` (empty) | Suppresses **all** findings under that skill |

Paths in the ignore file are relative to your **working directory** when you run the command. The entries above assume you run from the **repository root** and pass `--root examples/ignoring-findings/skills`. Diagnostic paths look like `examples/ignoring-findings/skills/...`, so the ignore paths match that prefix.

If you `cd examples/ignoring-findings` first, use shorter paths (or a local ignore file whose paths match `skills/...`). This example’s ignore file is written for runs from the repo root.

## JSON audit trail

```bash
agent-lint check \
  --root examples/ignoring-findings/skills \
  --ignore-file examples/ignoring-findings/.agent-lint-ignore.yaml \
  --json
```

Ignored items still appear under `ignored` / `ignoredCount`. They do not fail the build.

## Prefer fixing

Ignores are for deliberate exceptions. For real skills, prefer renaming directories, splitting oversized bodies into `references/`, and fixing broken links — see [docs/ignore.md](../../docs/ignore.md).
