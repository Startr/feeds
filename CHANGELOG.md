# Changelog

All notable changes to Startr/feeds will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **PocketBase standalone + JS hooks.** The entire rewrite pipeline now runs as a PocketBase JS hook (`pb_hooks/feeds.pb.js`). The container downloads a pre-built PocketBase binary (renamed to `feeds` for branding) — no Go compiler, no compilation, build time in seconds. PocketBase provides: HTTP server, static file serving (`pb_public/`), admin UI at `/_/`, cron scheduler, SQLite, and graceful shutdown.
- **Multi-feed from day one.** Feed configs are stored in a PocketBase "feeds" collection, managed via the admin UI at `/_/`. Each record is one feed to rewrite. Single-feed deploys can still use `FEEDS_*` env vars as a fallback.
- `pb_hooks/feeds.pb.js` — the rewrite pipeline as a PocketBase JS hook. Fetches upstream RSS with conditional GET, rewrites channel-level branding via xml-js DOM manipulation, writes rewritten XML to `pb_public/`, stores cache state (ETag, Last-Modified) back to the collection record.
- `pb_hooks/lib/xml-js.js` — vendored UMD bundle of [xml-js](https://github.com/nicknisi/xml-js) v1.6.11 for XML ↔ JS object conversion. Vendored via `npm pack` (see `pb_hooks/lib/LIBRARIES.md`).
- `pb_migrations/1744156800_create_feeds.js` — creates the "feeds" collection schema on first boot.
- `CHANGELOG.md` (this file).
- `FEEDS_CRON` env var for configuring the rewrite schedule using standard cron expressions (default `*/15 * * * *`). Replaces `FEEDS_INTERVAL`.
- `FEEDS_VERSION` env var for stamping the `<generator>` tag (default `dev`).
- `feeds superuser` command (from PocketBase) for creating/updating admin accounts for the dashboard.
- `docs/deploy-caprover.md` — operator walkthrough for CapRover deployment.

### Changed
- **Architecture: Go binary → PB standalone + JS hooks.** The Dockerfile no longer compiles Go code. It downloads the pre-built PocketBase binary and copies in `pb_hooks/`, `pb_migrations/`, and `pb_public/`. The Go pipeline code in `internal/` is kept as a reference implementation and high-performance fallback for large-scale multi-feed deployments.
- **Scheduler:** `time.Ticker` → PocketBase's built-in cron via `cronAdd()` in the JS hook.
- **Static file serving:** PocketBase serves `pb_public/` over HTTP. The JS hook writes output XML there. No external web server needed.
- **Env var scope:** Feed-specific config uses `FEEDS_*` env vars (UPSTREAM, OUTPUT, SELF_URL, CHANNEL_*, ITUNES_*, CRON). PocketBase infrastructure flags (`--http`, `--dir`, `--publicDir`) come from the CLI / Dockerfile CMD — no env var duplication.
- `FEEDS_INTERVAL` removed (replaced by `FEEDS_CRON`). `FEEDS_HTTP`, `FEEDS_DIR`, `FEEDS_PUBLIC_DIR`, `FEEDS_HOOKS_DIR`, `FEEDS_MIGRATIONS_DIR` removed (PocketBase handles these via its own flags).
- Dockerfile rewritten: single-stage Alpine, downloads PocketBase binary, copies hooks/migrations/public. No Go build stage.
- Makefile updated: removed `--build-arg VERSION=` from build targets (no Go compilation), updated `it_run_dev` for PB standalone.
- README updated to reflect PB standalone + JS hooks architecture, multi-feed support, simplified install.

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
