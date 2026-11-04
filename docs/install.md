# Install

Pick one of the following. Prefer Homebrew or a pinned release binary in CI so you are not pulling an npm package.

## Homebrew

```bash
brew install ai-agent-manager/tap/agent-lint
agent-lint --version
```

The tap formula is updated automatically when a new `v*` GitHub Release is published.

## GitHub release binary

1. Open the latest [Release](https://github.com/ai-agent-manager/agent-lint/releases).
2. Download the archive for your OS and architecture (for example `agent-lint_0.1.1_darwin_arm64.tar.gz`).
3. Download `checksums.txt` from the same release and verify the archive SHA-256 matches.
4. Extract `agent-lint` and place it on your `PATH`.

Example (macOS arm64):

```bash
VERSION=0.1.1
curl -fsSL -O "https://github.com/ai-agent-manager/agent-lint/releases/download/v${VERSION}/agent-lint_${VERSION}_darwin_arm64.tar.gz"
curl -fsSL -O "https://github.com/ai-agent-manager/agent-lint/releases/download/v${VERSION}/checksums.txt"
# verify checksum, then:
tar -xzf "agent-lint_${VERSION}_darwin_arm64.tar.gz"
sudo mv agent-lint /usr/local/bin/
```

## From source

Requires Go 1.22+:

```bash
go install github.com/ai-agent-manager/agent-lint/cmd/agent-lint@latest
```

This installs whatever version `@latest` resolves to on the module proxy. For reproducible CI, prefer a tagged release binary instead.
