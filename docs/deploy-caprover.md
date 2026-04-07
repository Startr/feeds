---
title: "Deploying Startr/feeds to CapRover"
audience: operators
status: stable
applies_to: v0.0.x
last_updated: 2026-04-08
---

# Deploying Startr/feeds to CapRover

This walkthrough covers the **one-time** CapRover setup you need before `caprover deploy --default` (or `make it_deploy`) will produce a working feed at `https://feed.yourdomain.com/v1/your-show.xml`.

The deploy command itself is one line. Everything else on this page is dashboard clicks and CLI prereqs you only do once per app.

## Prerequisites

You need:

1. **A running CapRover server** — self-hosted on a VPS, DigitalOcean droplet, Hetzner box, or wherever. CapRover's [getting started guide](https://caprover.com/docs/get-started.html) walks through the initial server install. This document assumes you already have a CapRover dashboard reachable at `https://captain.yourdomain.com`.
2. **The `caprover` CLI installed locally** — needed for `caprover deploy --default`:
   ```bash
   npm install -g caprover
   caprover --version
   ```
3. **A local checkout of this repo** with the `captain-definition` file at the root. The deploy command builds *from your local working tree*, so the file you commit is the file CapRover uses.
4. **A domain you control** for the feed (e.g., `feed.yourdomain.com`), with DNS pointing at your CapRover server's IP. CapRover will provision a Let's Encrypt cert for it once the app is up.
5. **The upstream feed URL** for the show you're rewriting (e.g., `https://anchor.fm/s/9822e708/podcast/rss`).

## Step 1 — Authenticate the CapRover CLI against your server

```bash
caprover login
```

The CLI prompts for:
- **CapRover machine URL** — `https://captain.yourdomain.com`
- **Password** — the dashboard password you set during the CapRover install
- **Machine name** — a local alias (e.g., `startr-prod`). This is just for your `~/.caprover/config.json`.

Verify:

```bash
caprover list
```

You should see the machine you just added.

## Step 2 — Create the app in the CapRover dashboard

1. Open `https://captain.yourdomain.com` in a browser.
2. **Apps → Create A New App**.
3. App name: `feeds` (or `startr-feeds`, or whatever you want — this becomes part of the default `*.yourdomain.com` URL CapRover assigns).
4. Check **Has Persistent Data**. This is critical — without it, the cache state file and the rewritten XML output disappear on every container restart.
5. Click **Create New App**.

You now have an empty app. The next steps configure it.

## Step 3 — Add a persistent volume for `/app/pb_data`

The container writes two things to disk that need to survive restarts:

- `/app/pb_data/.feeds-state.json` — ETag + Last-Modified cache. Without this, every restart re-fetches the upstream feed in full instead of getting an HTTP 304.
- `/app/pb_data/public/your-show.xml` — the rewritten output your subscribers read.

In the CapRover dashboard:

1. Click into your `feeds` app → **App Configs** tab.
2. Scroll to **Persistent Directories**.
3. Add an entry:
   - **Path in App** — `/app/pb_data`
   - **Label** — `feeds-data` (or `feeds-pb-data`)
4. Click **Save & Update**.

CapRover will create a named Docker volume on the host and mount it at `/app/pb_data` inside the container.

## Step 4 — Set the `FEEDS_*` environment variables

Still in **App Configs**, scroll to **Environmental Variables**. Add one row per variable:

| Key | Value |
|---|---|
| `FEEDS_UPSTREAM` | `https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss` |
| `FEEDS_OUTPUT` | `/app/pb_data/public/your-show.xml` |
| `FEEDS_SELF_URL` | `https://feed.yourdomain.com/v1/your-show.xml` |
| `FEEDS_CHANNEL_TITLE` | `Your Show` |
| `FEEDS_CHANNEL_LINK` | `https://yourdomain.com/podcast` |
| `FEEDS_CHANNEL_IMAGE` | `https://yourdomain.com/podcast/cover.jpg` |
| `FEEDS_ITUNES_AUTHOR` | `Your Name` |
| `FEEDS_ITUNES_OWNER_EMAIL` | `you@yourdomain.com` |
| `FEEDS_INTERVAL` | `15m` |

Required: `FEEDS_UPSTREAM`, `FEEDS_OUTPUT`, `FEEDS_SELF_URL`, `FEEDS_CHANNEL_TITLE`, `FEEDS_CHANNEL_LINK`. The rest are optional but recommended — Apple Podcasts wants `<itunes:author>` and `<itunes:owner>` populated.

`FEEDS_OUTPUT` must live under the persistent volume path (`/app/pb_data/...`). If you write it anywhere else, the file disappears on restart.

Click **Save & Update**.

> **Note on safety:** v0.0.x `feeds serve` is fail-soft on missing flags. If you forget to set the required env vars, the container comes up, logs `feeds serve: no feed configured (missing [...]) — idling`, and waits on SIGTERM instead of crash-looping. So you can spin up the app first and add env vars after — the container won't fight you.

## Step 5 — Map your subscriber domain

In the same app → **HTTP Settings** tab:

1. **Connect New Domain** — enter `feed.yourdomain.com` and click connect. CapRover verifies the DNS A record points at your CapRover server.
2. Once it's connected, click **Enable HTTPS** next to that domain. CapRover provisions a Let's Encrypt certificate. (You may need to set a contact email at the CapRover server level first if you've never done HTTPS before — the dashboard prompts you.)
3. Optionally check **Force HTTPS by redirecting all HTTP traffic to HTTPS**.
4. **Container HTTP Port** — leave at `80`. The Startr/feeds container exposes `8090`, but in v0.0.x there's no HTTP server in the binary itself. CapRover's built-in nginx is what serves the static XML file off the persistent volume. (See Step 6.)

## Step 6 — Configure nginx to serve the static XML

This is the step most people miss the first time. The container *writes* the rewritten feed to `/app/pb_data/public/your-show.xml`, but in v0.0.x there's no HTTP server inside the container itself. You need CapRover's nginx layer to serve that file at the URL your subscribers will hit.

Two options:

### Option A — sidecar static server (simplest)

Add a second app, e.g., `feeds-web`, using a tiny static-file image:

```
FROM nginx:alpine
```

In `feeds-web`'s **App Configs**, add the **same** persistent volume (`feeds-data` → `/usr/share/nginx/html`, set to `/app/pb_data/public` if your CapRover version supports subpath mounts). Then map `feed.yourdomain.com` to `feeds-web` instead of `feeds`. The `feeds` app stays headless and just writes XML files; `feeds-web` serves them.

### Option B — CapRover's built-in nginx custom config (single app)

In the `feeds` app → **HTTP Settings** → expand **Edit Default NGINX Configurations**. Add a `location` block that serves files from the persistent volume:

```nginx
location /v1/ {
    alias /captain/data/feeds-data/public/;
    add_header Cache-Control "public, max-age=300";
    types {
        application/xml xml;
        application/rss+xml rss;
    }
    default_type application/xml;
}
```

The exact path to the persistent volume on the CapRover host is `/captain/data/<volume-label>/` — adjust `feeds-data` if you used a different label in Step 3. CapRover's nginx runs on the host, so it can read the volume directly without the file going through the `feeds` container at all.

Save & Update. Restart nginx if CapRover doesn't auto-reload.

> **Why two options?** Sidecar (A) is portable and doesn't require touching CapRover's nginx config. Built-in nginx (B) is one less app to manage and one less container running. Both work. Pick whichever your operational model prefers.

## Step 7 — Deploy

From your local checkout of this repo:

```bash
make it_deploy
```

which is shorthand for:

```bash
caprover deploy --default
```

The CLI uses the `--default` machine + app from `~/.caprover/config.json`. If you have multiple apps registered, run `caprover deploy` without `--default` and pick interactively, or pass `-a feeds` to target a specific app.

CapRover will:

1. Tar your local checkout (excluding files in `.gitignore` and `.captain-ignore`).
2. Upload it to the CapRover server.
3. Build the image using the Dockerfile referenced in `captain-definition` — which in our case is `FROM ghcr.io/startr/feeds:latest` plus the `pb_hooks/`, `pb_migrations/`, `pb_public/` COPY layers and the `feeds serve` CMD.
4. Roll out the new container, attach the persistent volume + env vars, and start it.

Watch the deploy logs in the terminal. When CapRover prints `Deployed successfully`, switch to the dashboard's **App Logs** tab and confirm you see:

```
feeds serve <version>: ticker mode, interval=15m, http=0.0.0.0:8090 (ignored until v0.2)
```

followed by the pipeline running once on startup.

## Step 8 — Verify

```bash
curl -I https://feed.yourdomain.com/v1/your-show.xml
```

Expect `200 OK` and `content-type: application/xml`.

```bash
curl -s https://feed.yourdomain.com/v1/your-show.xml | head -40
```

Expect to see your `<title>` (from `FEEDS_CHANNEL_TITLE`), your `<link>` (from `FEEDS_CHANNEL_LINK`), your `<itunes:author>`, and `<generator>Startr/feeds <version> (https://github.com/Startr/feeds)</generator>`. The `<enclosure url>` elements will still point at the upstream audio host — that's intentional ([README → "What gets rewritten, what doesn't"](../README.md#what-gets-rewritten-what-doesnt)).

Validate against Apple Podcasts' parser using [Podbase Validator](https://podba.se/validate/) or Apple's own [Podcasts Connect Validator](https://podcastsconnect.apple.com/) (login required). Both flag missing iTunes namespace tags, malformed `pubDate`, missing `atom:link rel="self"`, etc. Startr/feeds injects `atom:link` if upstream forgets it, but other field issues come from upstream and need fixing there.

Once validation passes, **submit `https://feed.yourdomain.com/v1/your-show.xml` to Apple Podcasts Connect, Spotify for Podcasters, Pocket Casts, Overcast, etc.** From now on, your subscribers bind to *your* URL — you can swap the upstream host whenever you want and nobody re-subscribes.

## Subsequent deploys

Once Steps 1–6 are done, future deploys are just:

```bash
make it_deploy
```

Env var changes don't require a code deploy at all — edit them in the CapRover dashboard and click **Save & Update**, and CapRover restarts the container with the new values. Because the persistent volume survives restarts, the cache state file and last-good output XML are preserved across restarts, so subscribers never see a partial or stale feed during the rollover.

If you bump the binary version (via `make minor_release` / `make patch_release` → `make release_and_push_GHCR`), the new image lands on `ghcr.io/Startr/feeds:latest`, and the next `make it_deploy` will pull it in. CapRover does not auto-pull `:latest` on its own — it pulls during build, which happens when you deploy.

## Troubleshooting

**`curl` returns 404 for `/v1/your-show.xml`** — nginx isn't serving the file. Check Step 6. Also confirm the file actually exists on the host: SSH to the CapRover server and `ls /captain/data/feeds-data/public/` (substitute your volume label).

**`curl` returns 502** — the `feeds` container is failing to start. Check **App Logs** in the dashboard. Most common cause: the binary is fail-soft idling because env vars are missing. Look for `no feed configured (missing [...])` in the logs.

**The feed shows the upstream's title/branding instead of yours** — the rewriter isn't running. Check `FEEDS_OUTPUT` matches the path nginx is serving from. Also check `FEEDS_UPSTREAM` is set — without it, serve idles and never writes.

**Logs show `upstream 304 Not Modified — no rewrite needed`** — this is **correct** steady-state behavior. Conditional GET is doing its job. The existing output file on disk is still valid. To force a fresh fetch + rewrite (e.g., after changing env vars that affect the output), delete the cache state file: `rm /captain/data/feeds-data/.feeds-state.json` and restart the app.

**Apple Podcasts shows the old feed for a long time after you update env vars** — Apple caches feeds aggressively (typically 24 hours). This is on Apple's end, not yours. Your `<channel><lastBuildDate>` and HTTP `Last-Modified` headers are correct; Apple just won't notice immediately.

**You need to roll back** — CapRover keeps the previous N images. **Apps → feeds → Deployment → Versions** lets you click a previous version and redeploy it. Persistent volume data is unchanged, so subscribers see continuous service.

## What this doesn't cover

- **Multiple feeds in one app.** v0.0.x is one feed per container. To run two feeds, create two CapRover apps (`feeds-show1`, `feeds-show2`), each with its own `FEEDS_*` env vars and its own persistent volume. v0.2 will add multi-feed orchestration via PocketBase collections so a single app can serve many feeds.
- **Custom rehosting of audio enclosures.** v0.0.x leaves `<enclosure url>` untouched on purpose — Spotify hosts the audio bytes for free. Audio rehosting is a v1.0+ opt-in feature.
- **Authenticated upstream feeds.** v0.0.x assumes the upstream is publicly fetchable. Auth headers / cookies are not yet plumbed through.
- **Webhook-driven updates.** v0.0.x is purely interval-based. There's no `POST /refresh` endpoint to trigger an immediate re-fetch from CI. v0.2's PocketBase serve mode will expose admin endpoints for this.

For anything beyond the scope of this walkthrough, see the [main README](../README.md) and [CHANGELOG](../CHANGELOG.md).
