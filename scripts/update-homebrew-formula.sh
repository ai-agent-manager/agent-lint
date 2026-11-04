#!/usr/bin/env bash
# Update Formula/agent-lint.rb in ai-agent-manager/homebrew-tap for a published release.
#
# Usage:
#   HOMEBREW_TAP_TOKEN=... ./scripts/update-homebrew-formula.sh v0.2.0
#   ./scripts/update-homebrew-formula.sh --dry-run v0.2.0
#
# Env:
#   HOMEBREW_TAP_TOKEN  Required to push (skip with warning if unset)
#   HOMEBREW_TAP_REPO   Override tap clone URL (default: ai-agent-manager/homebrew-tap)
#   DRY_RUN=1           Write formula to stdout / temp only; do not push
set -euo pipefail

DRY_RUN="${DRY_RUN:-0}"
VERSION_ARG=""

for arg in "$@"; do
  case "${arg}" in
    --dry-run) DRY_RUN=1 ;;
    -h|--help)
      sed -n '2,12p' "$0"
      exit 0
      ;;
    *)
      if [[ -n "${VERSION_ARG}" ]]; then
        echo "unexpected argument: ${arg}" >&2
        exit 2
      fi
      VERSION_ARG="${arg}"
      ;;
  esac
done

if [[ -z "${VERSION_ARG}" ]]; then
  echo "usage: $0 [--dry-run] <version>" >&2
  echo "  version examples: v0.2.0 or 0.2.0" >&2
  exit 2
fi

TAG="${VERSION_ARG}"
if [[ "${TAG}" != v* ]]; then
  TAG="v${TAG}"
fi
VERSION="${TAG#v}"

TAP_REPO="${HOMEBREW_TAP_REPO:-ai-agent-manager/homebrew-tap}"
RELEASE_REPO="${RELEASE_REPO:-ai-agent-manager/agent-lint}"
BASE_URL="https://github.com/${RELEASE_REPO}/releases/download/${TAG}"

if [[ -z "${HOMEBREW_TAP_TOKEN:-}" && "${DRY_RUN}" != "1" ]]; then
  echo "HOMEBREW_TAP_TOKEN not set; skipping Homebrew formula update" >&2
  exit 0
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

echo "Fetching checksums for ${TAG}..."
curl -fsSL "${BASE_URL}/checksums.txt" -o "${tmpdir}/checksums.txt"

checksum_for() {
  local asset="$1"
  local sum
  sum="$(awk -v f="${asset}" '$2 == f { print $1; exit }' "${tmpdir}/checksums.txt")"
  if [[ -z "${sum}" ]]; then
    echo "checksum entry for ${asset} not found in checksums.txt" >&2
    cat "${tmpdir}/checksums.txt" >&2
    exit 2
  fi
  printf '%s' "${sum}"
}

SHA_DARWIN_ARM64="$(checksum_for "agent-lint_${VERSION}_darwin_arm64.tar.gz")"
SHA_DARWIN_AMD64="$(checksum_for "agent-lint_${VERSION}_darwin_amd64.tar.gz")"
SHA_LINUX_ARM64="$(checksum_for "agent-lint_${VERSION}_linux_arm64.tar.gz")"
SHA_LINUX_AMD64="$(checksum_for "agent-lint_${VERSION}_linux_amd64.tar.gz")"

formula_path="${tmpdir}/agent-lint.rb"
cat > "${formula_path}" <<EOF
# typed: false
# frozen_string_literal: true

class AgentLint < Formula
  desc "Lint Agent Skills against the agentskills.io specification"
  homepage "https://github.com/${RELEASE_REPO}"
  version "${VERSION}"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "${BASE_URL}/agent-lint_${VERSION}_darwin_arm64.tar.gz"
      sha256 "${SHA_DARWIN_ARM64}"
    end
    on_intel do
      url "${BASE_URL}/agent-lint_${VERSION}_darwin_amd64.tar.gz"
      sha256 "${SHA_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_arm do
      url "${BASE_URL}/agent-lint_${VERSION}_linux_arm64.tar.gz"
      sha256 "${SHA_LINUX_ARM64}"
    end
    on_intel do
      url "${BASE_URL}/agent-lint_${VERSION}_linux_amd64.tar.gz"
      sha256 "${SHA_LINUX_AMD64}"
    end
  end

  def install
    bin.install "agent-lint"
  end

  test do
    system "#{bin}/agent-lint", "list-rules"
  end
end
EOF

if [[ "${DRY_RUN}" == "1" ]]; then
  echo "----- dry-run formula -----"
  cat "${formula_path}"
  exit 0
fi

clone_url="https://x-access-token:${HOMEBREW_TAP_TOKEN}@github.com/${TAP_REPO}.git"
tap_dir="${tmpdir}/tap"
git clone --depth 1 "${clone_url}" "${tap_dir}"
mkdir -p "${tap_dir}/Formula"
cp "${formula_path}" "${tap_dir}/Formula/agent-lint.rb"

git -C "${tap_dir}" config user.name "github-actions[bot]"
git -C "${tap_dir}" config user.email "41898282+github-actions[bot]@users.noreply.github.com"

if git -C "${tap_dir}" diff --quiet -- Formula/agent-lint.rb; then
  echo "Formula already at ${VERSION}; nothing to commit"
  exit 0
fi

git -C "${tap_dir}" add Formula/agent-lint.rb
git -C "${tap_dir}" commit -m "chore: bump agent-lint formula to ${VERSION}"
git -C "${tap_dir}" push origin HEAD
echo "Updated ${TAP_REPO} Formula/agent-lint.rb to ${VERSION}"
