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
3. [`tag-on-version.yaml`](../.github/workflows/tag-on-version.yaml) runs on the
   push to `main`. If no tag matches the new `VERSION`, it creates and pushes
   `v<VERSION>`.
4. That tag triggers [`release.yaml`](../.github/workflows/release.yaml), which
   runs [GoReleaser](../goreleaser.yaml) to build every platform and publish the
   `.deb` / `.rpm` / archives to GitHub Releases.

```
bump VERSION ─► merge to main ─► tag-on-version creates vX.Y.Z ─► release.yaml ─► GoReleaser publishes
```

So: **to cut a release, bump `VERSION` and merge.** Nothing else.

## Notes

- Re-running or merging without changing `VERSION` does nothing — the workflow
  skips when the tag already exists, so it's safe.
- The tag is pushed using the `TOKEN` repository secret (a PAT). This is
  required: a tag pushed with the default `GITHUB_TOKEN` would not trigger the
  release workflow.
- To release manually (e.g. a hotfix tag), you can still
  `git tag vX.Y.Z && git push origin vX.Y.Z` — just keep `VERSION` in sync.
- Follow semantic versioning: patch for fixes, minor for backward-compatible
  features, major for breaking changes.
