# Rules

agent-lint emits diagnostics with stable **rule IDs**. Use those IDs in ignore files and when triaging CI failures.

```bash
agent-lint list-rules
```

## What is checked (summary)

### Layout

| Rule ID | Severity | Meaning |
|---------|----------|---------|
| `skill/asset-type-conflict` | error | Both `SKILL.md` and `rovo-agent.yaml` present |
| `skill/asset-type-missing` | error | Neither definition file present |
| `skill/skill-md-casing` | warning | Wrong casing for `SKILL.md` |
| `skill/name-duplicate` | error | Same `name` declared twice in one scan |

### Frontmatter (agentskills.io)

| Rule ID | Severity | Meaning |
|---------|----------|---------|
| `skill/frontmatter-missing` | error | Missing `---` frontmatter |
| `skill/frontmatter-invalid-yaml` | error | Frontmatter is not valid YAML |
| `skill/frontmatter-not-object` | error | Frontmatter is not a mapping |
| `skill/name-missing` | error | Missing `name` |
| `skill/name-format` | error | `name` fails format rules (length, charset, hyphens) |
| `skill/name-directory-mismatch` | error | `name` ≠ parent directory |
| `skill/description-missing` | error | Missing `description` |
| `skill/description-length` | error | `description` &gt; 1024 characters |
| `skill/compatibility-length` | error | `compatibility` &gt; 500 characters |
| `skill/metadata-format` | error | `metadata` is not string→string |
| `skill/allowed-tools-format` | error | `allowed-tools` is not a space-separated string |
| `skill/frontmatter-unknown-field` | error | Field outside the agentskills.io allowlist |

Allowed frontmatter fields by default: `name`, `description`, `license`, `compatibility`, `metadata`, `allowed-tools`.

### Body and references

| Rule ID | Severity | Meaning |
|---------|----------|---------|
| `skill/body-empty` | warning | No instructions after frontmatter |
| `skill/body-line-limit` | error | Body exceeds line cap (default 500) |
| `skill/body-token-limit` | error | Body exceeds token cap (default 5000) |
| `skill/reference-depth` | error | Relative ref deeper than one directory level |
| `skill/reference-missing` | error | Markdown link target missing on disk |
| `skill/reference-orphan` | warning | File under `references/` never mentioned |

**Token estimate:** `ceil(utf8_rune_count / 4)`.

Markdown links are checked for existence. Plain-text path mentions in prose are not required to exist (so examples do not fail the lint).

### Config-gated README rules

When `require_readme` / tags / phases are configured: `skill/readme-missing`, `skill/readme-frontmatter-invalid`, `skill/readme-name-missing`, `skill/readme-description-missing`, `skill/readme-tags-invalid`, `skill/readme-required-tag`, `skill/readme-phases-invalid`, `skill/readme-phase-value`.

### Opt-in Claude.ai rules

With `--claude`: `skill/claude-reserved-name`, `skill/claude-angle-brackets`.

## Fixing common failures

- **Body too long** — move detail into `references/` and link it from `SKILL.md` with a clear “when to read” cue ([progressive disclosure](https://agentskills.io/specification#progressive-disclosure)).
- **Unknown field** — move non-spec keys under `metadata` (string values), or enable `--allow-claude-code` / `allow_claude_code_fields` if you intentionally use Claude Code extensions.
- **Name mismatch** — rename the directory or the `name` field so they match.
