# Getting started

Shows a typical skills tree and what agent-lint **finds** when everything is valid.

## Layout

```text
getting-started/
  skills/
    greet-user/          ← Agent Skill (has SKILL.md)
      SKILL.md
      references/
        tone.md
    summarize-notes/     ← another Agent Skill
      SKILL.md
    backlog-triage/      ← Rovo-only (rovo-agent.yaml)
      rovo-agent.yaml
```

## Run

```bash
# From the agent-lint repository root
agent-lint check --root examples/getting-started/skills
```

Expected: exit code `0`.

## What gets discovered

| Directory | Files | Result |
|-----------|-------|--------|
| `greet-user` | `SKILL.md` | Validated as an Agent Skill |
| `summarize-notes` | `SKILL.md` | Validated as an Agent Skill |
| `backlog-triage` | `rovo-agent.yaml` only | Validated as a Rovo agent |

Only **immediate children** of `--root` are candidates. Nested folders inside a skill (for example `references/`) are not treated as separate skills.

## Try JSON output

```bash
agent-lint check --root examples/getting-started/skills --json
```

You should see empty `diagnostics` / zero error counts.
