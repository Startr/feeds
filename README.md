# Startr/feeds

## v0.0.4

Self-hosted feed rewriter. Fetch an RSS/Atom feed from upstream, rewrite the branding and channel metadata to point at your own URL, republish through your own infrastructure. Your subscribers bind to your URL, forever.

Built on [PocketBase](https://pocketbase.io) as a standalone binary with JS hooks. Ship a single container and get: HTTP-served rewritten feeds, admin dashboard at `/_/` for managing feed configs, automatic cron-scheduled rewrites, and multi-feed support from day one. No Go compiler needed — the container downloads a pre-built PocketBase binary and runs the rewrite pipeline as a JS hook.

## Why this exists

Google's FeedBurner has been on a slow march to the graveyard since 2012. Major features were removed in 2021. No clean self-hostable replacement exists. Meanwhile, podcast publishers are locked into Spotify for Podcasters or Anchor's terms without an easy path to owning their feed URL.

Startr/feeds is the feed proxy / rewriter we wished existed:

- **Own the subscriber relationship.** Subscribers bind to `feed.yourdomain.com/v1/your-show.xml`, not a third-party URL. Migrate audio hosts without breaking a single subscriber.
- **Rent the audio host.** `<enclosure url>` elements are left untouched. Spotify or Anchor host the audio bytes for free. If you want to move to Archive.org, S3, or your own static host later, swap it with a config change and nobody notices.
- **No lock-in, anywhere.** AGPL-3.0. Single container. Self-hostable on anything that runs Docker, including a Raspberry Pi. If you want to run this on a box in your closet, nothing stops you.

## Current scope (v0.0.x)

- **PocketBase standalone** — pre-built binary (renamed to `feeds` for branding). HTTP static file serving, admin UI at `/_/`, built-in cron scheduler, SQLite, graceful shutdown. One container does everything.
- **JS hook rewrite pipeline** (`pb_hooks/feeds.pb.js`) — the entire rewrite logic runs inside PocketBase's goja runtime. Zero Go compilation, build time in seconds.
- **Multi-feed from day one** — feed configs stored in PocketBase's "feeds" collection, managed via the admin UI at `/_/`. Single-feed deploys can use `FEEDS_*` env vars instead.
- XML parsing via [xml-js](https://github.com/nicknisi/xml-js) (vendored UMD bundle) — non-compact mode preserves iTunes namespace, podcast 2.0 namespace, and any unknown tags on round-trip
- `<channel><generator>` rewritten to identify Startr/feeds + version (replaces upstream "Anchor Podcasts" / "Spotify for Podcasters")
- HTTP conditional GET with `If-None-Match` / `If-Modified-Since` (98% of scheduled runs short-circuit on HTTP 304)
- Fail-loud per feed on upstream errors — last-good output is preserved, other feeds continue
- **Go pipeline** (`internal/`) kept as reference implementation and high-performance fallback for large-scale multi-feed deployments

v0.3+ adds Atom, JSON Feed, YouTube, text/video `media_type` support — at which point the tool is a self-hostable FeedBurner replacement.

## Install

### Container image (recommended)

Multi-arch (amd64 + arm64) is published to GHCR on every tag.

```bash
docker run --rm -p 8090:8090 \
  -e FEEDS_SOURCE_URL=https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss \
  -e FEEDS_SLUG=your-show \
  -e FEEDS_DOMAIN=https://feed.yourdomain.com \
  -e FEEDS_TITLE="Your Show" \
  -e FEEDS_WEBSITE=https://yourdomain.com/podcast \
  ghcr.io/Startr/feeds:latest
```

Feed is served at `http://localhost:8090/v1/your-show.xml`. Admin UI at `http://localhost:8090/_/`.

Pin to a specific release with `ghcr.io/Startr/feeds:0.0.4` (or any tag from the [Releases](https://github.com/Startr/feeds/releases) page).

### From source with `make`

```bash
git clone https://github.com/Startr/feeds && cd feeds
make it_build       # downloads PB binary + copies hooks → startr/feeds:latest
```

No Go compiler needed — the Dockerfile downloads a pre-built PocketBase binary and copies in the JS hooks. Build time: seconds.

## Quick start

Run the container with `FEEDS_*` env vars. PocketBase starts, runs the rewrite pipeline once immediately, then re-runs on the cron schedule (default: every 15 minutes):

```bash
docker run --rm -p 8090:8090 \
  -v feeds-data:/app/pb_data \
  -e FEEDS_SOURCE_URL=https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss \
  -e FEEDS_SLUG=your-show \
  -e FEEDS_DOMAIN=https://feed.yourdomain.com \
  -e FEEDS_TITLE="Your Show" \
  -e FEEDS_WEBSITE=https://yourdomain.com/podcast \
  ghcr.io/Startr/feeds:latest
```

The rewritten feed is served at `/v1/your-show.xml`. The PocketBase admin UI is at `/_/`.

For **multi-feed** setups, add feed configs in the admin UI instead of env vars — each record in the "feeds" collection is one feed to rewrite.

## Configuration

Feed config can come from two places:

1. **PocketBase "feeds" collection** (recommended for multi-feed) — manage via the admin UI at `/_/`. Each record is one feed to rewrite. The JS hook iterates over all records on each cron tick.
2. **`FEEDS_*` environment variables** (single-feed shortcut) — used as fallback when no collection records exist. Great for quick deploys.

### Environment variables

| Env var | Notes |
|---|---|
| `FEEDS_SOURCE_URL` | RSS feed URL to rewrite — **required** |
| `FEEDS_SLUG` | feed identifier, e.g. `my-show` → serves at `/v1/my-show.xml` — **required** |
| `FEEDS_DOMAIN` | your domain, e.g. `https://feed.example.com` — used for the `<atom:link rel="self">` URL |
| `FEEDS_TITLE` | your show title — **required** |
| `FEEDS_WEBSITE` | your show's homepage — **required** |
| `FEEDS_COVER_IMAGE` | cover art URL — optional |
| `FEEDS_ITUNES_AUTHOR` | iTunes author — optional |
| `FEEDS_ITUNES_OWNER_EMAIL` | iTunes owner email — optional |
| `FEEDS_CRON` | cron expression for rewrite schedule (default `*/15 * * * *`) |
| `FEEDS_VERSION` | version string for `<generator>` tag (default `dev`) |

PocketBase's own flags (`--http`, `--dir`, `--publicDir`, `--hooksDir`, `--migrationsDir`) are set in the Dockerfile CMD or passed on the CLI. Run `feeds serve --help` to see all available PB flags.

### CapRover deploy

> **First time deploying?** Read [`docs/deploy-caprover.md`](./docs/deploy-caprover.md) for the full one-time setup walkthrough.

The repo ships a [`captain-definition`](./captain-definition) at the root. CapRover builds from it and runs `feeds serve`. Set the `FEEDS_*` vars in the CapRover dashboard's **App Configs → Environment Variables** tab:

```
FEEDS_SOURCE_URL=https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss
FEEDS_SLUG=your-show
FEEDS_DOMAIN=https://feed.yourdomain.com
FEEDS_TITLE=Your Show
FEEDS_WEBSITE=https://yourdomain.com/podcast
FEEDS_COVER_IMAGE=https://yourdomain.com/podcast/cover.jpg
FEEDS_ITUNES_AUTHOR=Your Name
FEEDS_ITUNES_OWNER_EMAIL=you@yourdomain.com
FEEDS_CRON=*/15 * * * *
```

Save and restart. `feeds serve`:

1. Starts PocketBase (HTTP server, admin UI at `/_/`, SQLite).
2. **Runs the rewrite pipeline once immediately** — subscribers get a fresh feed at deploy time.
3. Re-runs on the `FEEDS_CRON` schedule. 98% of polls short-circuit on HTTP 304.
4. Serves the rewritten XML from `pb_public/` — no separate nginx or sidecar needed.
5. Shuts down cleanly on SIGTERM.


## What gets rewritten, what doesn't

| Element | Action |
|---|---|
| `<channel><title>` | Rewritten to your show title |
| `<channel><link>` | Rewritten to your show page |
| `<channel><image>` | Rewritten to your hosted cover art |
| `<channel><generator>` | Rewritten to identify Startr/feeds + version (replaces upstream "Anchor Podcasts" / "Spotify for Podcasters") |
| `<atom:link rel="self">` | Rewritten to your self URL (injected if missing — required by Apple Podcasts) |
| `<itunes:author>`, `<itunes:owner>`, `<itunes:image>` | Rewritten to your branding |
| `<item><enclosure url>` | **Left alone.** Points at the upstream audio host. This is intentional. |
| `<item><guid>`, `<item><title>`, `<item><description>`, `<item><pubDate>` | Left alone. Episode-level content. |
| Any unknown namespaced element | Left alone. Round-tripped via xml-js non-compact mode. |

Enclosure URLs staying untouched is how v0.1.0 runs at zero cost: Spotify hosts the audio bytes for free. If you want to rehost audio yourself, that's a v1.0+ feature via an opt-in rehosting config.

## Build, release, and deploy

This repo ships a pre-built PocketBase binary (renamed to `feeds`) with JS hooks in a minimal Alpine container to GHCR. No Go compiler needed — the Dockerfile downloads the PocketBase release, copies in `pb_hooks/`, `pb_migrations/`, and `pb_public/`, and that's it. Build time: seconds.

Key commands:

- `feeds serve` — PocketBase server + cron-driven JS hook rewrite pipeline. Serves rewritten feeds from `pb_public/`, admin UI at `/_/`, built-in scheduler. This is the primary deployment mode.
- `feeds superuser` — create or update admin accounts for the PocketBase dashboard.

### Prerequisites (one-time)

```bash
brew install git-flow-next        # Go rewrite of git-flow; the Makefile requires it
git flow init -d                  # default branch names: master/develop/feature/release/hotfix
gh auth login                     # for GHCR push
```

You also need a container runtime. The Makefile auto-detects `podman` (preferred) and falls back to `docker`.

### Local development

```bash
make it_build              # downloads PB binary + copies hooks → startr/feeds:latest
make it_run                # run on localhost:8090, data in volume startr-media-data
make it_run_dev            # run with bind-mounted pb_hooks/, pb_migrations/, pb_public/
make it_build_n_run        # build + run in one shot
```

### Release flow

The Makefile uses `git-flow-next` for branch management. Release version is auto-computed from the latest git tag.

<details>
<summary><b>First release ever</b> — no tags exist yet (kept for forks bootstrapping fresh)</summary>

```bash
git checkout develop
make first_release                 # creates release/0.0.1 branch
make bump_release_version          # rewrites the `## v0.0.0` heading in this README → `## v0.0.1`
git add -A && git commit -m "Bump to 0.0.1"
make it_build && make it_run       # smoke test
make ghcr_login                    # authenticate Docker against ghcr.io via gh CLI
make release_and_push_GHCR         # finish release + multi-arch GHCR push
```

</details>

**Subsequent releases** (pick the bump type):

```bash
git checkout develop
make minor_release                 # or patch_release / major_release
make bump_release_version
git add -A && git commit -m "Bump version"
make it_build && make it_run       # smoke test
make ghcr_login
make release_and_push_GHCR
```

`release_and_push_GHCR` chains `release_finish` (git-flow merge to `master`, tag, push `develop`/`master`/tags) → `it_build_multi_arch_push_GHCR` (multi-arch buildx push to `ghcr.io/Startr/feeds:<version>` and `:latest`). If the push fails after the release is finished, re-run `make it_build_multi_arch_push_GHCR` to retry just that step.

### Hotfix flow (emergency fix from `master`)

```bash
git checkout master
make hotfix                        # creates hotfix/x.y.z.1 branch (4th version component)
# ... fix the bug ...
make bump_release_version
git add -A && git commit -m "Fix + bump"
make it_build && make it_run       # smoke test
make ghcr_login
make hotfix_and_push_GHCR
```

## License

**AGPL-3.0.** See [LICENSE](./LICENSE).

### Embedding widgets: you are NOT bound by AGPL

Startr/feeds ships widgets (`<startr-player>`, `<startr-subscribe>`) alongside the rewriter. These widgets are AGPL-3.0.

**Important:** Embedding a widget via a remote `<script src>` from a hosted Startr/feeds instance does NOT infect your site with AGPL. You are referencing remote content, not redistributing or running modified code.

- **You embed `<script src="https://yourinstance/widgets/player.js">` on your site:** zero AGPL obligations. Your site code stays under whatever license you want.
- **You fork the widget code, modify it, and host it yourself:** AGPL applies. You must share your modifications under AGPL-3.0 if you run them as a network service.

This is exactly how [Plausible Analytics](https://plausible.io) operates. Tens of thousands of non-AGPL websites embed Plausible via a remote script with zero legal friction, because AGPL §13 applies to people *running modified versions* as a network service, not to people *referencing* the original.

### Contribution model

No CLA. Contributions accepted under the [Developer Certificate of Origin](./DCO.txt). Sign off your commits with `git commit -s`.

## Project name and trademark

The names `Startr/feeds`, `Startr`, `startrcast`, `startr-player`, and `startr-subscribe` are project marks reserved separately from the code license. See [TRADEMARK.md](./TRADEMARK.md). Forks that make substantive changes should use a different name.

## Docs

Per-release notes live in [CHANGELOG.md](./CHANGELOG.md).

Design thinking and editorial context for the canonical `startr.media` instance live in [`docs/`](./docs/). These documents explain why the tool exists, what problem it solves, and how it fits into the broader Sage ecosystem. They are useful context for anyone evaluating whether to adopt Startr/feeds.
