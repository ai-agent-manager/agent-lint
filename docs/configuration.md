# Configuration

Create `.agent-lint.yaml` in your project root (or pass `--config`).

## Example

```yaml
skills:
  max_body_lines: 500
  max_body_tokens: 5000
  require_readme: true
  readme_constraints:
    - field: tags
      must_include:
        - agent-skill
    - field: phases
      allowed:
        - visioning
        - discovery
        - elaboration
        - build_and_release
        - run_and_measure
  allow_claude_code_fields: false
```

## Options

| Key | Default | Meaning |
|-----|---------|---------|
| `skills.max_body_lines` | `500` | Max lines in the `SKILL.md` body (after frontmatter). Enforced as an **error**. |
| `skills.max_body_tokens` | `5000` | Max estimated tokens in the body. Estimator: `ceil(utf8_rune_count / 4)`. Enforced as an **error**. |
| `skills.require_readme` | `false` | Require each skill directory to include `README.md` with frontmatter. |
| `skills.readme_constraints` | _(empty)_ | When non-empty (and README is required), each entry constrains a README frontmatter string-list `field` via `must_include` and/or `allowed`. Entries with neither list set are a config error (exit `2`). |
| `skills.allow_claude_code_fields` | `false` | Allow Claude Code–only frontmatter keys (same effect as `--allow-claude-code`). |
| `skills.claude_code_fields` | _(built-in list)_ | Optional override of which Claude Code keys are allowed when the flag/config above is enabled. |
| `ignore.ignores` | _(empty)_ | Optional ignore entries embedded in this file (same shape as [Ignoring findings](ignore.md)). |

### README constraints

Each `readme_constraints` entry:

- **`field`** — frontmatter key to check (required)
- **`must_include`** — every listed value must appear in the field
- **`allowed`** — every value in the field must be a member of this list

When either list is non-empty, the field must be present as a non-empty array of strings.

## Defaults without a config file

- Body limits stay at **500 lines** / **5000 tokens**.
- README / constraint checks are **off**.
- Only agentskills.io frontmatter fields are allowed unless you opt into Claude Code fields.

## Spec authority

Frontmatter and naming rules follow [agentskills.io/specification](https://agentskills.io/specification). Body size limits are recommendations in the spec; agent-lint treats them as errors by default so oversized skills fail CI unless you raise the caps in config.
