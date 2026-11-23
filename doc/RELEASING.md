# Releasing ProxyX

Releases are **automated**. You bump one file; CI does the rest.

## The `VERSION` file is the source of truth

The repository root holds a [`VERSION`](../VERSION) file with the current
version, e.g. `0.1.5`. The running binary reports it via `proxyx version`
(injected at build time with ldflags — never hardcoded).

## How a release happens

1. Open a PR that bumps `VERSION` (e.g. `0.1.5` → `0.1.6`), alongside your
   changes.
2. Merge the PR to `main`.
3. [`auto-release.yaml`](../.github/workflows/auto-release.yaml) runs on the push
   to `main`. If no tag matches the new `VERSION`, it creates and pushes
   `v<VERSION>` and then runs [GoReleaser](../goreleaser.yaml) to build every
   platform and publish the `.deb` / `.rpm` / archives to GitHub Releases — all
   in one job.

```
bump VERSION ─► merge to main ─► auto-release tags vX.Y.Z + publishes via GoReleaser
```

So: **to cut a release, bump `VERSION` and merge.** Nothing else.

## No PAT required

Both workflows use the built-in `GITHUB_TOKEN` (with `contents: write`), not a
personal access token. Tag creation and the release run in the **same** job in
`auto-release.yaml` precisely so that no cross-workflow trigger — and therefore
no PAT — is needed (a tag pushed with `GITHUB_TOKEN` does not trigger other
workflows).

## Notes

- Re-running, or merging without changing `VERSION`, does nothing: the workflow
  skips when the tag already exists, so it is safe.
- You can also run `auto-release` manually from the **Actions** tab
  (workflow_dispatch) to release whatever `VERSION` currently names.
- Manual escape hatch: `git tag vX.Y.Z && git push origin vX.Y.Z` triggers
  [`release.yaml`](../.github/workflows/release.yaml). Keep `VERSION` in sync.
- Follow semantic versioning: patch for fixes, minor for backward-compatible
  features, major for breaking changes.
