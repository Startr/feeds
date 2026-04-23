# TODOs

## Theme toggle: support system preference as a third state

The startr.style theme toggle currently cycles between light and dark. It should
support a third option: "system" (follow `prefers-color-scheme`). When set to
system, remove the `data-theme` attribute entirely so the CSS `@media` queries
take over. The toggle would cycle: light → dark → system → light.

This also means the 24-hour localStorage expiry should reset to system default
when it expires, not stick on the last manual choice.

## Deferred from v0.2 CEO Review (2026-04-12)

### Feed health indicator on portal cards

Add a `last_status` field to the feeds collection. Pipeline updates it on every
cron run: "ok" (successful rewrite), "not_modified" (304), or "error:{message}"
(fetch failed). Portal cards show a green/yellow/red dot. Operator sees at a
glance if an upstream feed went down.

**Why:** Currently, if a feed fetch fails, the only evidence is a `console.log`
line nobody reads. This is the observability gap for self-hosted operators.

**Effort:** S (human: ~2 hrs / CC: ~10 min) | **Priority:** P2
**Depends on:** v0.2 portal cards (needs the card UI to display the dot on)

## Deferred from v0.1.0

- Startr.style npm publishing flow (ops team owns CI/CD)
- Startr.style Component Gallery scoping
- Mini mode layout variant for player widget
- ~~Dynamic per-feed OG tags (requires server-side HTML templating, v0.2+)~~ Replaced by static OG pages in v0.2
