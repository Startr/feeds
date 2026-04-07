# Changelog

All notable changes to Startr/feeds will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `CHANGELOG.md` (this file).
- Environment variable support for every flag on both `feeds rewrite` and `feeds serve`. Each flag now reads its default from a `FEEDS_*` env var (e.g. `FEEDS_UPSTREAM`, `FEEDS_OUTPUT`, `FEEDS_SELF_URL`, `FEEDS_CHANNEL_TITLE`, `FEEDS_INTERVAL`, …) so containerized deploys (CapRover, fly.io, k8s) can configure feeds without baking values into the CMD line. CLI flags still override env vars when explicitly passed — precedence is CLI flag > env var > literal default. The `envString` and `envDuration` helpers in `cmd/feeds/main.go` keep the per-flag wiring to one line. Bad `FEEDS_INTERVAL` values fall back silently to the default rather than crashing the container at startup. The README has a new "Environment variables" table and a CapRover deploy snippet.
- `docs/deploy-caprover.md` — operator walkthrough for the one-time CapRover dashboard setup that has to happen before `caprover deploy --default` produces a working feed. Covers caprover CLI install, server login, app creation, persistent volume for `/app/pb_data`, env var configuration, domain + HTTPS via Let's Encrypt, two options for serving the static XML off the persistent volume (sidecar nginx vs CapRover's built-in nginx custom config), first deploy, verification with `curl` + Apple Podcasts validators, subsequent deploys, rollback, and a troubleshooting section for the common 404 / 502 / stale 304 cache failures. The README now links to it from the CapRover deploy section.

### Changed
- README updated to match v0.0.x reality: dropped the unshipped Sigstore cosign + GitHub Releases tarball install instructions, replaced the broken `--config /config.yaml` container example with the working flag form, added `<channel><generator>` to the rewrite scope table, renamed "Scope for v0.1.0" to "Current scope (v0.0.x)", and collapsed the "First release ever" block into a `<details>` since it has been done four times already (kept for forks bootstrapping fresh).
- `captain-definition` simplified. The CMD line previously hardcoded `--http=0.0.0.0:8090 --dir=/app/pb_data --hooks-dir=/app/pb_hooks --migrations-dir=/app/pb_migrations --public-dir=/app/pb_public` — every one of those values is now either the binary's built-in default or settable via a `FEEDS_*` env var, so the CMD reduces to `["feeds", "serve"]`. CapRover env var changes no longer require rebuilding the captain image. The README's CapRover deploy section spells out the four-step startup behavior (`serve` runs the rewrite pipeline once immediately, then on the configured `FEEDS_INTERVAL` tick) so it's clear that no separate `feeds rewrite` invocation is needed at deploy time.

### Fixed
- `feeds rewrite --help` and `feeds serve --help` now exit with code 0 instead of code 1. The `flag.ErrHelp` sentinel was previously propagated as a regular error, causing the help screens to print `feeds serve: flag: help requested` to stderr and exit with a failure code despite working correctly. Both subcommand files now check `errors.Is(err, flag.ErrHelp)` after `fs.Parse` and return nil.

## [0.0.4] - 2026-04-07

CI/CD shakedown release. No functional changes from v0.0.3 — multi-arch push to `ghcr.io/Startr/feeds:0.0.4` was the goal of the release.

## [0.0.3] - 2026-04-07

CI/CD shakedown release. No functional changes from v0.0.2.

## [0.0.2] - 2026-04-07

First feature drop after the bootstrap.

### Added
- ASCII banner art on `feeds --help` (top-level only — subcommand `--help` screens stay clean for scripting and CI grep).
- `<channel><generator>` is now rewritten to identify Startr/feeds + version, replacing upstream's "Anchor Podcasts" / "Spotify for Podcasters" attribution. The generator string is built from the binary's stamped version at runtime.
- `Generator` field on `internal/rewriter.Rewriter`, with corresponding test (`TestRewrite_GeneratorRewritten`) and a `<generator>Anchor Podcasts</generator>` element added to the test fixture so the test verifies replace-existing behavior matches what real anchor.fm feeds look like.

### Changed
- `feeds serve` now fails soft when required flags (`--upstream`, `--output`, `--self-url`, `--channel-title`, `--channel-link`) are missing. Instead of erroring out, it logs a warning and idles on the signal context until SIGTERM. This keeps `make it_run` and CapRover smoke deploys working without per-feed config, matching the v0.2 vision where PocketBase collections will be the source of feed config and starting up empty is normal.
- Binary version is now stamped from the Makefile's existing `RELEASE_VERSION` cascade (release branch name → git tag) via `-ldflags`, removing a parallel "version from README heading" code path. One source of truth for "what version is this", rooted in git, not in file contents.
- Dockerfile switched from `go mod download` to `go mod tidy` so the build self-bootstraps `go.sum` from a cold start (no committed `go.sum` in v0.0.x).
- Dockerfile rewritten from a 3-stage PocketBase build to a 2-stage Go build (`golang:1.25-alpine` → `alpine:3.21`). The build runs `go vet` and `go test ./...` before the final `go build`, so test failures fail the image.
- `make it_run_dev` updated to invoke `feeds serve` with kebab-case flags instead of `pocketbase serve` with camelCase flags.

### Fixed
- Hardcoded `"v0.1.0"` literal in the `cmd/feeds/serve.go` log line replaced with the dynamic `version` package var. The log line now correctly reflects the build's stamped version.

## [0.0.1] - 2026-04-07

Initial release. Bootstrap scope for v0.0.x.

### Added
- `feeds rewrite` subcommand: single-shot fetch + rewrite + atomic write. Exits zero on success. Wire to any scheduler.
- `feeds serve` subcommand: long-running ticker mode for always-on hosts. v0.2 will swap the ticker for a PocketBase framework import for multi-feed orchestration and an admin UI. The flag surface (`--http`, `--dir`, `--hooks-dir`, `--migrations-dir`, `--public-dir`) is reserved now for forward-compat with v0.2's PocketBase serve flags so CapRover deployments survive the swap.
- `internal/rewriter`: DOM-style XML rewriter via [`beevik/etree`](https://github.com/beevik/etree). Preserves iTunes namespace, podcast 2.0 namespace, and any unknown namespaced elements on round-trip.
- `internal/source/spotify`: HTTP fetcher with conditional GET (`If-None-Match` / `If-Modified-Since`). 304 responses short-circuit the rest of the pipeline so steady-state polls cost ~one HTTP request and zero disk I/O.
- `internal/output/rss`: atomic write-rename output via `os.CreateTemp` + `os.Rename`. Readers never see a partial file. Idempotent no-op on byte-equal output (no rename, no inode churn).
- `internal/cache`: JSON state file for ETag and Last-Modified persistence between runs. Self-healing on corrupt JSON (treated as first run).
- `internal/pipeline`: orchestrator. Fail-loud on errors with last-good output preserved (the existing output file is not touched on failure).
- Five gating tests baked into the Docker build's `go test ./...` step:
  1. iTunes / podcast namespace preservation on round-trip.
  2. Enclosure URL preservation — the most important test in the codebase. Rewriting it would force us to rehost audio bytes, which is explicitly punted to v1.0+.
  3. `atom:link rel="self"` rewritten when present, injected when missing (Apple Podcasts requires this element).
  4. Fail-loud on upstream error preserves last-good output.
  5. Idempotent no-op write — second run of identical input doesn't rotate the file inode.
- `examples/`: six scheduler snippets that the README references — GitHub Actions cron (`github-actions.yml`), systemd timer (`systemd-timer/feeds-rewrite.{service,timer}`), raspberry-pi crontab (`cron-on-raspberry-pi.txt`), fly.io scheduled task (`flyio-scheduled-task.toml`), Kubernetes CronJob (`kubernetes-cronjob.yaml`), and Docker Compose with cron sidecar (`docker-compose-cron.yml`).
- Multi-arch (amd64 + arm64) Docker image published to `ghcr.io/Startr/feeds` on every release tag.

[Unreleased]: https://github.com/Startr/feeds/compare/v0.0.4...HEAD
[0.0.4]: https://github.com/Startr/feeds/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/Startr/feeds/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/Startr/feeds/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/Startr/feeds/releases/tag/v0.0.1
