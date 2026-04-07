---
title: "DevOps Inbox"
priority: high
status: open
created: 2026-04-07
owner: devops
---

# DevOps Inbox

High-priority operational issues that need devops attention. Items here are blocking other work.

---

## [HIGH] dark-search returning 403 to agent requests

**Discovered:** 2026-04-07
**Status:** Open, blocking name verification work for the Startr/feeds podcast project
**Reporter:** Office hours session

### What's happening

The internal search engine at `https://dark-search.production.openco.ca/search?...` is returning HTTP 403 Forbidden when called from agent/automation contexts (Claude Code via WebFetch).

Example failing requests:
- `https://dark-search.production.openco.ca/search?q=%22Run+It+Local%22+podcast&format=json` → 403
- `https://dark-search.production.openco.ca/search?q=%22Off+Cloud%22+podcast&format=json` → 403

### What we need

1. Confirm whether dark-search is supposed to be reachable from agent contexts at all (auth model? IP allowlist? bearer token?)
2. The team has mentioned a **markdown endpoint** that should be used for agent search instead. Need:
   - The full URL pattern
   - Auth model (token? header? cookie?)
   - Expected response format (markdown body? JSON-wrapped markdown?)
   - Rate limits, if any
3. Whichever endpoint is canonical, document it somewhere agents can find it (CLAUDE.md, AGENTS.md, or a dedicated tools manifest)

### Workaround in the meantime

Agents are falling back to standard WebSearch. This is fine for public-web research but does not give us:
- Access to internal documents
- Privacy from external search providers (queries leak to the standard search vendor)
- Anything indexed only by dark-search

### Action items

- [ ] **devops:** Diagnose why dark-search returns 403 to agent requests
- [ ] **devops:** Provide the markdown endpoint URL and auth pattern
- [ ] **devops:** Document the canonical agent-search endpoint in `CLAUDE.md` so future sessions pick it up automatically
- [ ] **office-hours session:** Circle back in ~1 day to verify and re-run name research through the correct endpoint

### Related work blocked by this

- Name verification for `Run It Local`, `Off Cloud`, `Quiet Compute` (Startr/feeds podcast project) — currently using public WebSearch as a fallback, but the team prefers internal search for privacy
- Any future agent task that needs internal corpus search
