# Startr/feeds

## v0.0.4

Self-hosted feed rewriter. Fetch an RSS/Atom feed from upstream, rewrite the branding and channel metadata to point at your own URL, republish through your own infrastructure. Your subscribers bind to your URL, forever.

Shipped as a single-shot Go binary. No daemon. No built-in scheduler. No cloud dependencies. Runs from cron, systemd, GitHub Actions, fly.io, k8s CronJob, Docker Compose, or a manual terminal invocation. Pick whichever scheduler fits your environment.

## Why this exists

Google's FeedBurner has been on a slow march to the graveyard since 2012. Major features were removed in 2021. No clean self-hostable replacement exists. Meanwhile, podcast publishers are locked into Spotify for Podcasters or Anchor's terms without an easy path to owning their feed URL.

Startr/feeds is the feed proxy / rewriter we wished existed:

- **Own the subscriber relationship.** Subscribers bind to `feed.yourdomain.com/v1/your-show.xml`, not a third-party URL. Migrate audio hosts without breaking a single subscriber.
- **Rent the audio host.** `<enclosure url>` elements are left untouched. Spotify or Anchor host the audio bytes for free. If you want to move to Archive.org, S3, or your own static host later, swap it with a config change and nobody notices.
- **No lock-in, anywhere.** AGPL-3.0. Single binary. Self-hostable on anything that runs a Go binary, including a Raspberry Pi. If you want to run this from cron on a box in your closet, nothing stops you.

## Current scope (v0.0.x)

- Two subcommands: `feeds rewrite` (single-shot for cron/CI) and `feeds serve` (long-running ticker for always-on hosts; fail-soft when no upstream is configured so empty starts don't crash deploys)
- One source adapter: `SpotifySource` (parses Spotify for Podcasters' auto-generated RSS)
- One output renderer: `RSS2PodcastRenderer` (writes RSS 2.0 + iTunes namespace XML)
- DOM-style rewriter via [`beevik/etree`](https://github.com/beevik/etree) — preserves iTunes namespace, podcast 2.0 namespace, and any unknown Spotify tags on round-trip
- `<channel><generator>` rewritten to identify Startr/feeds + version (replaces upstream "Anchor Podcasts" / "Spotify for Podcasters")
- HTTP conditional GET with `If-None-Match` / `If-Modified-Since` (98% of scheduled runs short-circuit on HTTP 304)
- Atomic write-rename output (readers never see a partial file)
- Fail-loud on upstream errors — last-good output is preserved

v0.2 wraps this in [PocketBase](https://pocketbase.io) for multi-feed orchestration, admin UI, transcripts, and guest metadata. v0.3+ adds Atom, JSON Feed, YouTube, text/video `media_type` support — at which point the tool is a self-hostable FeedBurner replacement.

## Install

### Container image (recommended)

Multi-arch (amd64 + arm64) is published to GHCR on every tag.

```bash
docker run --rm \
  -v $(pwd)/public:/out \
  ghcr.io/Startr/feeds:latest \
  feeds rewrite \
    --upstream      https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss \
    --output        /out/your-show.xml \
    --self-url      https://feed.yourdomain.com/v1/your-show.xml \
    --channel-title "Your Show" \
    --channel-link  https://yourdomain.com/podcast
```

Pin to a specific release with `ghcr.io/Startr/feeds:0.0.4` (or any tag from the [Releases](https://github.com/Startr/feeds/releases) page).

### From source with `go install`

```bash
go install github.com/Startr/feeds/cmd/feeds@latest
feeds --help
```

### From source with `make`

```bash
git clone https://github.com/Startr/feeds && cd feeds
make it_build       # produces startr/feeds:latest locally
```

### Binary tarballs (planned, v0.1.0+)

Pre-built binary tarballs from GitHub Releases, signed with [Sigstore cosign](https://www.sigstore.dev/) and a SLSA Level 3 provenance attestation, are planned for v0.1.0. v0.0.x ships only the GHCR container and `go install` paths.

## Quick start

Fetch the Spotify for Podcasters auto-feed for your show, rewrite branding to point at your own domain, output static XML to disk:

```bash
./feeds rewrite \
  --source spotify \
  --upstream https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss \
  --output ./public/v1/your-show.xml \
  --self-url https://feed.yourdomain.com/v1/your-show.xml \
  --channel-title "Your Show" \
  --channel-link https://yourdomain.com/podcast
```

The binary runs once and exits with status zero on success. Wire it up to any scheduler.

## Scheduler examples

Startr/feeds doesn't ship a built-in scheduler because schedulers are political. Cron is fine. Systemd timers are fine. GitHub Actions is fine if you trust GitHub. Pick yours.

See [`examples/`](./examples/) for ready-to-copy configs:

- `examples/github-actions.yml` — GitHub Actions on a 15-minute cadence
- `examples/systemd-timer/` — systemd service + timer unit files
- `examples/cron-on-raspberry-pi.txt` — plain crontab line
- `examples/flyio-scheduled-task.toml` — fly.io scheduled machine
- `examples/kubernetes-cronjob.yaml` — k8s CronJob resource
- `examples/docker-compose-cron.yml` — Docker Compose with a cron sidecar

## Configuration

v0.1.0 is flag-driven. YAML config file support lands in v0.1.1:

```yaml
# feeds.yaml (v0.1.1+)
source: spotify
upstream: https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss
output: ./public/v1/your-show.xml
self_url: https://feed.yourdomain.com/v1/your-show.xml
channel:
  title: Your Show
  link: https://yourdomain.com/podcast
  image: https://yourdomain.com/podcast/cover.jpg
itunes:
  author: Your Name
  owner_email: you@yourdomain.com
```

```bash
./feeds rewrite --config ./feeds.yaml    # v0.1.1+
```

CLI flags will always override YAML when both are set.

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
| Any unknown namespaced element | Left alone. Round-tripped via `beevik/etree`. |

Enclosure URLs staying untouched is how v0.1.0 runs at zero cost: Spotify hosts the audio bytes for free. If you want to rehost audio yourself, that's a v1.0+ feature via an opt-in rehosting config.

## Build, release, and deploy

This repo ships the `feeds` Go binary in a minimal Alpine container to GHCR, built with the canonical Sage CI/CD pattern (forked from `WEB-DB-sage-pb`). The Makefile is intentionally near-identical to its sibling — `make help` lists every target.

The binary has two subcommands:

- `feeds rewrite` — single-shot fetch + rewrite + atomic write. Exits zero on success. Wire it up to any scheduler.
- `feeds serve` — long-running ticker mode for always-on hosts (e.g., CapRover). v0.2 replaces the ticker with a PocketBase framework import for multi-feed orchestration and an admin UI.

### Prerequisites (one-time)

```bash
brew install git-flow-next        # Go rewrite of git-flow; the Makefile requires it
git flow init -d                  # default branch names: master/develop/feature/release/hotfix
gh auth login                     # for GHCR push
```

You also need a container runtime. The Makefile auto-detects `podman` (preferred) and falls back to `docker`.

### Local development

```bash
make it_build              # 2-stage Go build → startr/feeds:latest (runs go vet + tests)
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

Future versions of Startr/feeds will ship widgets (`<startrcast-player>`, `<startrcast-subscribe>`) alongside the rewriter. These widgets are AGPL-3.0.

**Important:** Embedding a widget via a remote `<script src>` from a hosted Startr/feeds instance does NOT infect your site with AGPL. You are referencing remote content, not redistributing or running modified code.

- **You embed `<script src="https://yourinstance/widgets/player.js">` on your site:** zero AGPL obligations. Your site code stays under whatever license you want.
- **You fork the widget code, modify it, and host it yourself:** AGPL applies. You must share your modifications under AGPL-3.0 if you run them as a network service.

This is exactly how [Plausible Analytics](https://plausible.io) operates. Tens of thousands of non-AGPL websites embed Plausible via a remote script with zero legal friction, because AGPL §13 applies to people *running modified versions* as a network service, not to people *referencing* the original.

### Contribution model

No CLA. Contributions accepted under the [Developer Certificate of Origin](./DCO.txt). Sign off your commits with `git commit -s`.

## Project name and trademark

The names `Startr/feeds`, `Startr`, and `startrcast` are project marks reserved separately from the code license. See [TRADEMARK.md](./TRADEMARK.md). Forks that make substantive changes should use a different name.

## Docs

Per-release notes live in [CHANGELOG.md](./CHANGELOG.md).

Design thinking and editorial context for the canonical `startr.media` instance live in [`docs/`](./docs/). These documents explain why the tool exists, what problem it solves, and how it fits into the broader Sage ecosystem. They are useful context for anyone evaluating whether to adopt Startr/feeds.
