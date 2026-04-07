---
title: "Deploying Startr/feeds to CapRover"
audience: operators
status: stable
applies_to: v0.0.x+ (PocketBase standalone + JS hooks)
last_updated: 2026-04-09
---

# Deploying Startr/feeds to CapRover

This walkthrough covers the **one-time** CapRover setup you need before `make it_deploy` produces a working feed at `https://feed.yourdomain.com/v1/your-show.xml`.

## Prerequisites

1. **A running CapRover server** — [getting started guide](https://caprover.com/docs/get-started.html).
2. **The `caprover` CLI** — `npm install -g caprover`.
3. **A local checkout of this repo** with the `captain-definition` file at the root.
4. **A domain** for the feed (e.g., `feed.yourdomain.com`), DNS pointing at your CapRover server.
5. **The upstream feed URL** for the show you're rewriting.

## Step 1 — Authenticate the CLI

```bash
caprover login
```

Enter your CapRover machine URL, password, and a local alias.

## Step 2 — Create the app

In the CapRover dashboard (`https://captain.yourdomain.com`):

1. **Apps → Create A New App**.
2. App name: `feeds` (or `startr-feeds`).
3. Check **Has Persistent Data** — the SQLite database and cache state file need to survive restarts.
4. Click **Create New App**.

## Step 3 — Persistent volume

In **App Configs → Persistent Directories**:

- **Path in App:** `/app/pb_data`
- **Label:** `feeds-data`

Save & Update. This volume stores PocketBase's SQLite database and the feeds cache state file.

## Step 4 — Environment variables

In **App Configs → Environmental Variables**, add:

| Key | Value |
|---|---|
| `FEEDS_SOURCE_URL` | `https://anchor.fm/s/YOUR_SHOW_ID/podcast/rss` |
| `FEEDS_SLUG` | `your-show` |
| `FEEDS_DOMAIN` | `https://feed.yourdomain.com` |
| `FEEDS_TITLE` | `Your Show` |
| `FEEDS_WEBSITE` | `https://yourdomain.com/podcast` |
| `FEEDS_COVER_IMAGE` | `https://yourdomain.com/podcast/cover.jpg` |
| `FEEDS_ITUNES_AUTHOR` | `Your Name` |
| `FEEDS_ITUNES_OWNER_EMAIL` | `you@yourdomain.com` |
| `FEEDS_CRON` | `*/15 * * * *` |

Required: `FEEDS_SOURCE_URL`, `FEEDS_SLUG`, `FEEDS_TITLE`, `FEEDS_WEBSITE`.

The slug determines the feed's URL path — `your-show` → served at `/v1/your-show.xml`. `FEEDS_DOMAIN` sets the domain for the `<atom:link rel="self">` tag (required by Apple Podcasts).

Save & Update.

> **Fail-soft:** If required env vars are missing, the container starts anyway — PocketBase serves the admin UI at `/_/` and logs "no feed configured". You can add env vars after deploy.

## Step 5 — Map your domain

In **HTTP Settings**:

1. **Connect New Domain** — `feed.yourdomain.com`.
2. **Enable HTTPS** (Let's Encrypt).
3. Optionally **Force HTTPS**.

CapRover's reverse proxy routes `feed.yourdomain.com` → container port 8090. PocketBase handles everything from there — it serves the rewritten XML from `pb_public/` and the admin UI at `/_/`.

## Step 6 — Deploy

```bash
make it_deploy
```

(Shorthand for `caprover deploy --default`.)

Watch the deploy logs. When it prints `Deployed successfully`, check the app logs in the dashboard for:

```
Server started at http://0.0.0.0:8090
```

followed by the initial pipeline run.

## Step 7 — Create an admin account

The PocketBase admin UI at `https://feed.yourdomain.com/_/` requires an admin account. Create one from any machine that has the `feeds` binary:

```bash
# Inside the CapRover container (via dashboard terminal or SSH):
feeds superuser create you@yourdomain.com yourpassword
```

Or use the CapRover dashboard's **Terminal** tab to run the command inside the container.

## Step 8 — Verify

```bash
curl -I https://feed.yourdomain.com/v1/your-show.xml
```

Expect `200 OK` with `content-type: application/xml`.

```bash
curl -s https://feed.yourdomain.com/v1/your-show.xml | head -30
```

Verify your `<title>`, `<link>`, `<itunes:author>`, and `<generator>Startr/feeds ...` are present. The `<enclosure url>` elements still point at the upstream audio host — that's [intentional](../README.md#what-gets-rewritten-what-doesnt).

Validate with [Podbase Validator](https://podba.se/validate/) or [Apple Podcasts Connect](https://podcastsconnect.apple.com/).

## Subsequent deploys

```bash
make it_deploy
```

Env var changes don't require a code deploy — edit in the CapRover dashboard, Save & Update, and CapRover restarts the container. The persistent volume survives restarts.

## Troubleshooting

**404 for the feed URL** — the output file doesn't exist yet. Check that `FEEDS_OUTPUT` is set and points under `pb_public/`. Check app logs for "initial rewrite" errors.

**502 Bad Gateway** — container failing to start. Check app logs. Most common: missing required env vars (container logs "no feed configured" and serves admin UI only — but CapRover may timeout waiting for the health check).

**Logs show `upstream 304 Not Modified`** — correct steady-state behavior. The existing output file is still valid. To force a fresh fetch (e.g., after changing branding env vars), delete the cache state file and restart.

**Admin UI shows "No admin accounts found"** — run `feeds superuser create` per Step 7.
