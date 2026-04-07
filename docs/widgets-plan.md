---
title: "Widgets Plan: v0.1.x Complete Package"
generated_by: /plan-ceo-review
date: 2026-04-10
branch: develop
status: ACTIVE
mode: SCOPE EXPANSION
---

# Widgets Plan: v0.1.x Complete Package

## Context

The feed rewriter infrastructure is done (PB standalone + JS hooks, multi-feed via
collections, admin UI, cron scheduling, HTTP conditional GET). The next phase ships
the complete "podcast in a box" package: embeddable player + subscribe widgets, then
the first episode.

This plan was produced by a `/plan-ceo-review` session on the `develop` branch.
The design doc at `docs/podcast-site-design.md` is the upstream source of truth for
the broader project vision.

## Vision

### 10x Check

Startr/feeds becomes the Plausible of podcast distribution. Not just a feed rewriter
with embed widgets, but a self-hosted podcast platform that gives any indie podcaster
everything Spotify for Podcasters gives them, minus the lock-in. Player widget as good
as Spotify's embed but open source. Subscribe widget with smart platform detection.
Analytics built into PocketBase. Transcript pages auto-generated from Whisper.
startr.media as the canonical instance AND the marketing site.

### Platonic Ideal

A podcaster discovers Startr/feeds. They run `docker run` with 4 env vars. In 30
seconds they have: their feed rewritten and self-hosted at a URL they own forever,
an admin dashboard with episode management, a player widget they can drop on any
website with one `<script>` tag, a subscribe widget with Apple/Spotify/RSS badges,
and a landing page that looks like someone who cares designed it. The whole thing
fits in one container. The user feels: "I can't believe this is free and I own it."

## Deliverables

### Core widgets (baseline scope)

1. **`<startr-player>`** -- vanilla web component, shadow DOM, reads RSS feed URL
2. **`<startr-subscribe>`** -- vanilla web component, configurable platform links
3. **Build/bundle step** producing `pb_public/widgets/player.js` and `pb_public/widgets/subscribe.js`
4. **Record first episode**, upload to Spotify for Podcasters
5. **Verify rewriter** picks it up, submit feed to Apple Podcasts and Spotify

### Accepted expansions (from CEO review ceremony)

| # | Feature | Effort | Why |
|---|---------|--------|-----|
| 3 | Episode picker in player | S-M | Turns player from single-play to mini podcast app. VCs can browse episodes. |
| 4 | Customizable accent color | S | CSS `--startr-accent` variable. Widget matches any site brand. |
| 5 | RSS auto-discovery meta tag | S | One `<link>` tag, unlocks feed discovery across web ecosystem. |
| 6 | Embed code generator in admin UI | M | The onboarding "aha moment" for non-technical podcasters. |
| 7 | Progress memory (localStorage) | S | Resume where you left off. Standard in every podcast app. |
| 8 | Keyboard shortcuts | S | Space=play/pause, arrows=skip, focus-scoped. Accessibility win. |
| 9 | Share button (Web Share API) | S | One-tap sharing. Podcast growth is word-of-mouth. |

### Deferred (TODOS)

| # | Feature | Effort | Why deferred |
|---|---------|--------|--------------|
| 1 | Dark mode auto-detection | S | Ship light-only first, add prefers-color-scheme later |
| 2 | Playback speed control | S | Ship basic player first, speed toggle later |
| 10 | Mini mode layout variant | S | Card mode sufficient for v0.1.x, add mode="mini" later |

## Scope Decisions

| # | Proposal | Effort | Decision | Reasoning |
|---|----------|--------|----------|-----------|
| 1 | Dark mode auto-detection | S | DEFERRED | User prefers to ship light-only first |
| 2 | Playback speed control | S | DEFERRED | Ship basic player first, speed later |
| 3 | Episode picker in player | S-M | ACCEPTED | Turns player from single-play to mini podcast app |
| 4 | Customizable accent color | S | ACCEPTED | CSS custom property, zero cost, any-brand match |
| 5 | RSS auto-discovery meta tag | S | ACCEPTED | One line, unlocks feed discovery |
| 6 | Embed code generator in admin UI | M | ACCEPTED | Onboarding "aha moment" |
| 7 | Progress memory (localStorage) | S | ACCEPTED | Standard podcast feature |
| 8 | Keyboard shortcuts | S | ACCEPTED | Accessibility + power-user |
| 9 | Share button (Web Share API) | S | ACCEPTED | Word-of-mouth growth |
| 10 | Mini mode layout variant | S | DEFERRED | Card mode sufficient for now |

## Implementation questions (resolve before building)

These surfaced during the CEO review's temporal interrogation:

1. **Where do widgets live in the repo?** Recommended: `widgets/player/` and
   `widgets/subscribe/` with a build step that outputs to `pb_public/widgets/`.
2. **Web component naming:** `<startr-player>` or `<feeds-player>`?
   Recommended: `<startr-player>` per the design doc.
3. **CORS for feed fetching:** The player widget fetches the RSS XML. If embedded
   on a different domain, CORS blocks the fetch. Options: (a) PB serves feeds with
   `Access-Control-Allow-Origin: *` header, (b) the widget takes a pre-parsed JSON
   prop instead of fetching XML. Recommended: option (a), add CORS header in PB hook.
4. **Shadow DOM styling:** Fully encapsulated with CSS custom properties for
   customization, or open shadow DOM? Recommended: closed shadow DOM + CSS custom
   properties (`--startr-accent`, etc.).
5. **Embed code generator implementation:** PB admin UI extension via JS hook, or
   a standalone page at `/embed`? Recommended: standalone page at `/embed` served
   from `pb_public/embed/index.html`, simpler than hooking into PB admin internals.
6. **Do all three Sage sites support arbitrary HTML/JS embeds?** Verify before
   building. If any site has a CMS approval workflow, account for that.

## Files to create/modify

### New files
- `widgets/player/startr-player.js` -- player web component source
- `widgets/subscribe/startr-subscribe.js` -- subscribe web component source
- `widgets/build.sh` (or similar) -- bundles to `pb_public/widgets/`
- `pb_public/widgets/player.js` -- built output
- `pb_public/widgets/subscribe.js` -- built output
- `pb_public/embed/index.html` -- embed code generator page

### Modified files
- `pb_public/index.html` -- add RSS auto-discovery meta tag
- `pb_hooks/feeds.pb.js` -- add CORS header for feed XML responses (if needed)
- `Dockerfile` -- include widgets in container
- `README.md` -- add widget documentation, embed instructions
- `CHANGELOG.md` -- add widget entries
- `docs/podcast-site-design.md` -- update status (v0.1.x scope expanded)

## Relationship to other docs

- **`docs/podcast-site-design.md`** -- the upstream design doc. Covers the full
  project vision from /office-hours + /plan-eng-review. This plan is a subset:
  specifically the v0.1.x widget work.
- **`README.md`** -- will need a "Widgets" section after implementation.
- **`CHANGELOG.md`** -- will need entries for each widget feature.
