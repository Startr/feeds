# TODOs

## Theme toggle: support system preference as a third state

The startr.style theme toggle currently cycles between light and dark. It should
support a third option: "system" (follow `prefers-color-scheme`). When set to
system, remove the `data-theme` attribute entirely so the CSS `@media` queries
take over. The toggle would cycle: light → dark → system → light.

This also means the 24-hour localStorage expiry should reset to system default
when it expires, not stick on the last manual choice.

## Deferred from v0.1.0

- Startr.style npm publishing flow (ops team owns CI/CD)
- Startr.style Component Gallery scoping
- Mini mode layout variant for player widget
- Dynamic per-feed OG tags (requires server-side HTML templating, v0.2+)
