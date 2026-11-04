# Findings

Intentionally broken skills so you can see **how agent-lint reports problems**.

## Layout

```text
findings/
  skills/
    Bad_Name/              ← name format + directory mismatch
      SKILL.md
    empty-body/            ← empty instructions after frontmatter
      SKILL.md
    both-definitions/      ← SKILL.md and rovo-agent.yaml together
      SKILL.md
      rovo-agent.yaml
    missing-link/          ← markdown link to a file that does not exist
      SKILL.md
```

## Run

```bash
agent-lint check --root examples/findings/skills
echo "exit=$?"
```

Expected: exit code `1`, with errors (and at least one warning) similar to:

| Path | Rule ID | Why |
|------|---------|-----|
| `Bad_Name/SKILL.md` | `skill/name-format` | Uppercase / underscore not allowed |
| `both-definitions` | `skill/asset-type-conflict` | Both definition files present |
| `empty-body/SKILL.md` | `skill/body-empty` | Warning: no body after frontmatter |
| `missing-link/SKILL.md` | `skill/reference-missing` | Link target not on disk |

## JSON view

```bash
agent-lint check --root examples/findings/skills --json | head
```

Each diagnostic includes `ruleId`, `severity`, `path`, and `message`. Use those `ruleId` values when writing an ignore file (see [ignoring-findings](../ignoring-findings/README.md)).

## List all rule IDs

```bash
agent-lint list-rules
```
