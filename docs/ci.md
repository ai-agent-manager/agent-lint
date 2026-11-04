# Using agent-lint in CI

Run a **pinned** release binary and verify its checksum. Avoid `go install @latest` in CI if you need reproducible builds.

## GitHub Actions (example)

```yaml
name: Lint skills

on:
  pull_request:
  push:
    branches: [main]

jobs:
  agent-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install agent-lint
        env:
          AGENT_LINT_VERSION: "0.1.1"
        run: |
          set -euo pipefail
          os=linux
          arch=amd64
          asset="agent-lint_${AGENT_LINT_VERSION}_${os}_${arch}.tar.gz"
          base="https://github.com/ai-agent-manager/agent-lint/releases/download/v${AGENT_LINT_VERSION}"
          curl -fsSL "$base/checksums.txt" -o checksums.txt
          curl -fsSL "$base/$asset" -o "$asset"
          expected="$(awk -v f="$asset" '$2 == f { print $1 }' checksums.txt)"
          actual="$(sha256sum "$asset" | awk '{ print $1 }')"
          test "$expected" = "$actual"
          tar -xzf "$asset"
          sudo mv agent-lint /usr/local/bin/

      - name: Check skills
        run: agent-lint check --root skills --config .agent-lint.yaml --ignore-file .agent-lint-ignore.yaml
```

Adjust `AGENT_LINT_VERSION`, OS/arch, and paths for your repo. Fail the job on exit code `1` (lint) or `2` (tool/config).

## Tips

- Commit `.agent-lint.yaml` (and an ignore file only if needed) so local and CI behaviour match.
- Use `--json` if a later step needs to parse findings.
- Use `--strict` if you want warnings to fail the build.
- Bump the pinned version deliberately when you adopt new rules.
