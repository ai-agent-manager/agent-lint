## 1. Sync script

- [x] 1.1 Add `scripts/update-homebrew-formula.sh` (version arg, fetch checksums, rewrite Formula, commit/push)
- [x] 1.2 Support dry-run / skip when `HOMEBREW_TAP_TOKEN` unset
- [x] 1.3 Make script executable; document usage in header comments

## 2. Release workflow

- [x] 2.1 Add post-GoReleaser step invoking the script with tag version and tap token
- [x] 2.2 Ensure step does not block asset publish when token missing (skip path)

## 3. Docs

- [x] 3.1 Update AGENTS.md release notes for automated tap sync + required secret
- [x] 3.2 Note automation in `docs/install.md` or release-and-ci maintainer guidance as appropriate
- [x] 3.3 Validate OpenSpec change
