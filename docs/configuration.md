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
      equals: agent-skill
    - field: phases
      matches: "^(visioning|discovery|elaboration|build_and_release|run_and_measure)$"
  skill_constraints:
    - field: metadata.phases
      matches: "^(visioning|discovery|elaboration|build_and_release|run_and_measure)$"
    - field: metadata.author
      not_equals: ""
  allow_claude_code_fields: false
```

## Options

| Key | Default | Meaning |
|-----|---------|---------|
| `skills.max_body_lines` | `500` | Max lines in the `SKILL.md` body (after frontmatter). Enforced as an **error**. |
| `skills.max_body_tokens` | `5000` | Max estimated tokens in the body. Estimator: `ceil(utf8_rune_count / 4)`. Enforced as an **error**. |
| `skills.require_readme` | `false` | Require each skill directory to include `README.md` with frontmatter. |
| `skills.readme_constraints` | _(empty)_ | When non-empty (and README is required), each entry constrains a dotted frontmatter string path with exactly one operator. |
| `skills.skill_constraints` | _(empty)_ | Same operators for `SKILL.md`. Prefer `metadata.*` paths to stay on-spec; the root key of each path is carved out of the allowlist when needed. |
| `skills.allow_claude_code_fields` | `false` | Allow Claude Code–only frontmatter keys (same effect as `--allow-claude-code`). |
| `skills.claude_code_fields` | _(built-in list)_ | Optional override of which Claude Code keys are allowed when the flag/config above is enabled. |
| `ignore.ignores` | _(empty)_ | Optional ignore entries embedded in this file (same shape as [Ignoring findings](ignore.md)). |

### Frontmatter constraints

Each `readme_constraints` / `skill_constraints` entry:

- **`field`** — dotted path to a string value (e.g. `tags`, `metadata.phases`)
- **exactly one** of:
  - **`equals`** — value must equal this string
  - **`not_equals`** — value must not equal this string
  - **`matches`** — value must match this Go/RE2 regular expression
  - **`not_matches`** — value must not match this regular expression

Missing paths or non-string values (including YAML arrays) fail as invalid. Invalid regexps or multiple operators fail config load (exit `2`).

## Defaults without a config file

- Body limits stay at **500 lines** / **5000 tokens**.
- README / constraint checks are **off**.
- Only agentskills.io frontmatter fields are allowed unless you opt into Claude Code fields or a `skill_constraints` root-key carve-out.

## Spec authority

Frontmatter and naming rules follow [agentskills.io/specification](https://agentskills.io/specification). Body size limits are recommendations in the spec; agent-lint treats them as errors by default so oversized skills fail CI unless you raise the caps in config.
