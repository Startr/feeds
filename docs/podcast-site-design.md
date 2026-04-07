---
title: "Podcast Site for Sage Ecosystem"
generated_by: /office-hours
date: 2026-04-06
branch: develop
repo: workspace
status: APPROVED
mode: Startup (Intrapreneurship)
---

# Design: Podcast Site for Sage Ecosystem

## Problem Statement

People discovering sage.is, sage.education, and startr.style are asking for more — specifically videos, podcasts, and supplementary content to go deeper on what these projects do, how, and why. The team already has organic roundtable conversations happening when visitors are present in person. The gap is capture and distribution, not content creation.

The podcast site needs to serve two roles simultaneously: depth content for warm ecosystem visitors, AND a discovery surface for new audiences (including VCs and parents evaluating the products).

**Scope evolution (two stages):**

*Stage 1 — Open source positioning:* What started as "build a podcast site" became "build an open-source podcast library tool." `startr.media` is both the team's own podcast hub AND the canonical instance of a tool other teams can self-host. The team is its own first user. Other teams (indie podcasters, devrel orgs, edu groups) become downstream users.

*Stage 2 — FeedBurner replacement positioning:* The deeper realization is that the underlying engine is a feed proxy/rewriter, not a podcast-specific tool. Google's FeedBurner has been on a slow march to the graveyard since 2012 (major features removed in 2021). No clean self-hostable replacement exists. The architecture being built here — fetch any feed, normalize it, rewrite branding, republish through your own URL — is the FeedBurner shape. By making the data model and adapters format-agnostic from day one, the same engine can handle podcasts (v0.1), text blogs (v0.3+), vlogs/video channels (v0.3+), and arbitrary RSS/Atom feeds (v0.3+).

**Critical scope discipline:** v0.1 still ships a podcast. v0.1 still serves the three Sage sites. v0.1 still answers "do you have a podcast?" this week. The FeedBurner positioning shapes the *architecture*, not the v0.1 surface area. The failure mode to avoid is letting the bigger vision eat the shipping deadline.

**Roundtable format:** The existing format is an informal, multi-person conversation (2-4 people) that happens organically when visitors are present. Freeform, no script, typically 20-45 minutes. There is a natural host/facilitator. The podcast captures this format, not replaces it.

**First episode format exception:** Episode 1 is a two-host pilot — Izzy Plante (Co-founder & CEO) and Alex Somma (Co-founder & CTO) — covering content from `sage.is/resources` and `sage.education`. Two-host format is simpler to coordinate for the pilot and works well for an introductory episode. Roundtable format with guests resumes from episode 2.

## Demand Evidence

- Multiple users and community members have directly asked "do you have a podcast?" or "do you have videos?"
- When people meet the team in person, casual roundtables happen organically — demand exists, format exists, distribution does not
- VCs evaluating the ecosystem have reached out wanting to understand the team's industry perspective
- Technical parents of students without product access are asking for supplementary content to evaluate whether the products are right for their kids

## Status Quo

People who want this content currently:
- Get pointed to existing written docs and blog posts (not the same)
- Find competitor podcasts or YouTube channels in the EdTech / dev tools space
- Get ad-hoc 1:1 DMs or calls with team members (unscalable)
- Meet the team in person for casual roundtable conversations (the actual product, unrecorded)

## Target User & Narrowest Wedge

**Primary:** Existing visitors across all three Sage sites — builders, developers, educators — who want to understand the philosophy and reasoning behind the products.

**Secondary (high-value):** Two groups that don't yet have product access:
- **Technical parents, homeschool families, and "tech mums."** This group has surged recently — homeschool enrollment hit 5.4% of K-12 in 2024-25 and is climbing, and the technically literate subset of those families are actively looking for AI tools they control rather than vendors they trust. Sage is currently working with this group directly, and the inbound interest has spiked. Existing parent-AI podcasts (AI for Kids, AI Parenting Podcast, Happy Homeschooler) cover basic AI literacy but NONE address the privacy/sovereignty/operator angle. There is no podcast home for tech-literate homeschool parents who want to understand the gap between vendor AI marketing and what's actually happening with their kids' data. That gap is a real opportunity.
- VCs evaluating the team's thinking and industry thesis

**Tertiary (open-source downstream users):** Other teams who want to own their feed without building the infrastructure themselves.

In Stage 1 (podcast-only): Indie podcasters fed up with Spotify lock-in. Devrel teams who want a beautiful self-hosted podcast library. Education groups publishing audio content.

In Stage 2 (FeedBurner-shaped): Anyone publishing any kind of feed who wants to own their subscriber relationship — bloggers leaving Substack, vloggers wanting independent distribution, newsletter operators, podcasters, RSS-driven aggregators. The audience is much larger because there is no good self-hostable FeedBurner replacement and millions of publishers still rely on a Google service that hasn't shipped a feature in years.

**Narrowest wedge:** One recorded roundtable, audio file uploaded somewhere (Spotify for Podcasters free tier is fastest), RSS feed self-hosted at `feed.startr.media`, and an embedded HTML5 audio player on a `/podcast` subpage of the matching site. That's it. That's the MVP. The "do you have a podcast?" answer goes from no to yes.

**Prerequisites status:**
- Recording setup: DONE (microphones, software, acoustics all in place)
- Domain: DONE (`startr.media` is owned)
- DNS: IN PROGRESS — moving `startr.media` into the `dawnofthegeeks@openco.ca` Cloudflare account
- Hosts identified: Izzy Plante (CEO) + Alex Somma (CTO), both Sage co-founders
- First episode source content: Identified — recent `sage.is/resources` articles (April 2026)
- Local repo: `WEB-Startr.media` (intranet name)
- Public repo: `Startr/feeds` on GitHub (going live today or tomorrow)
- Remaining: Confirm embed capability on all three sites — identify who has deploy access and whether each site supports arbitrary HTML embeds (iframe). If any site has a CMS approval workflow, account for that in the Phase 1 timeline.
- Placement rule: Phase 1 embeds should go on a dedicated `/podcast` subpage (not in the main content flow) to minimize risk to existing site traffic.

## Education Stakeholder Lens

A first-class lens that runs through the entire project. Edu stakeholders are not a niche audience to bolt on at the end. They are central to the show's positioning and the broader Sage philosophy.

**Who counts as an edu stakeholder for this project:**
- **Homeschool parents** (and the rapidly growing "tech mum" demographic — technically literate parents who care about both the substance of their kid's learning and the infrastructure underneath it)
- **Microschool operators** running small private learning environments
- **Public school teachers** who care about what's actually happening to their students' data, even when their district has approved a tool
- **University faculty** dealing with vendor AI policy
- **Education researchers** studying the impact of AI on learning

**The shared concern:** every one of these stakeholders is being told that "approved" and "compliant" tools are safe. Microsoft 365 for Education is FERPA-compliant. Google Classroom is signed off by district legal. The Canadian-hosted versions are "sovereign." Every layer of approval is a paper trail, not a control. The actual control still sits with US hyperscalers operating under US extraterritorial law.

**The autonomy distinction matters more in education than anywhere else:**

1. **Children cannot consent.** A homeschooler's seven-year-old cannot meaningfully agree to have their writing assignments sent to Meta as training data. The compliance frameworks pretend they can.
2. **Learning data is uniquely sensitive.** Kids ask AI tools things they would never write in a school journal. Emotional questions, identity questions, questions about their family. Once that's in a vendor's logs, it's there forever and may already be in the next training run.
3. **Subpoena risk is real for kids too.** A US prosecutor investigating *anything* can compel a US cloud provider to hand over Canadian children's data without notification. The CLOUD Act does not stop at adults.
4. **Vendor lock-in is generational.** A school that adopts a vendor AI tool today is shaping how an entire cohort of students thinks about whether software is something they can shape or just something they consume. That's a worldview decision, not a procurement decision.

**How this lens shapes the project:**

- **Show editorial:** Every episode should be readable to a tech-literate homeschool parent. Not "dumbed down" — assumed-intelligent. The voice does not change between engineers and parents.
- **Sage.education positioning:** The product story is not "an open-source AI tool for schools." It's "the only AI tool a school can adopt where the answer to 'whose law governs this data' is the school's own."
- **Distribution:** The podcast should be promoted in homeschool communities (Reddit r/homeschool, homeschool Slack groups, microschool networks) as well as the standard tech communities. Different distribution channels, same content, no rewrite needed.
- **Episode topic selection:** Future episodes should regularly cover stories that hit both audiences. Edu-AI policy news. Lawsuits involving children's data. Open-source models that run on hardware a homeschool parent could buy. Honest cost comparisons for school districts considering self-hosting.
- **The open-source tool (Startr/feeds):** Eventual v0.3+ support for text feeds means homeschool curriculum publishers can self-host their own content distribution without depending on Substack or Medium. That's a real edu use case.

## Constraints

- Intrapreneurship context — needs to ship fast, justify itself quickly
- Content already exists in oral form (roundtables) — production lift is capture and edit, not creation
- Three existing sites that should link to/embed podcast content
- VCs and parents are a distinct secondary audience with higher stakes — production quality matters more for them
- Every episode must have a transcript from day one (auto-generated, lightly reviewed for proper nouns). This is both accessibility and SEO — not a Phase 3 luxury
- Budget: Phase 1 can be $0/month using self-hosted RSS feed + Spotify for Podcasters free tier + Whisper for transcripts. Paid tooling (Transistor, paid analytics) is an upgrade path, not a Phase 1 requirement.
- Domain ownership: `startr.media` is already owned by the team and will be used as the podcast network hub
- Cadence: Weekly publishing from the start

## Premises

1. The podcast serves both cold audience discovery AND depth engagement for warm visitors. Not one or the other.
2. The roundtable conversations already happening are the content. The work is capturing and distributing them, not creating new shows.
3. VCs and technical parents are a distinct secondary audience — production quality and clarity of thesis matter more for them than for the builder/educator audience.

## Technical Architecture (decided)

The team has settled on building a self-hosted podcast library tool. `startr.media` is the canonical instance; the tool itself will be open-sourced so other teams can self-host their own libraries.

**Tech stack:**
- **Language:** Go. Chosen over Python for: single-binary distribution (critical for open-source dev tool adoption), shared ecosystem with PocketBase (the planned v0.2 backend), excellent stdlib XML support for RSS/Atom work.
- **v0.1 backend:** Standalone Go binary. Fetches the source feed (initially Spotify for Podcasters' auto-generated feed), rewrites all branding/links/metadata to point at `startr.media`, outputs static RSS XML. Single-shot, scheduler-agnostic — runs once and exits. Any external scheduler can drive it: cron, systemd timer, GitHub Actions, fly.io scheduled task, k8s CronJob, or a manual terminal run. The canonical `startr.media` instance happens to use GitHub Actions; downstream self-hosters pick whatever fits their environment. The binary ships ~6 example scheduler configs in `examples/`.
- **v0.2 backend:** PocketBase. The Go rewriter becomes a hook inside PocketBase. PocketBase handles feed records, item records, transcripts, guest metadata, admin UI, REST/Realtime API. Single binary deployment to fly.io or Hetzner. SQLite under the hood.
- **Frontend:** Static site (Astro or Next.js export) consuming PocketBase API in v0.2. In v0.1, the static RSS XML is the only "frontend" — embedded HTML5 audio players on the three existing sites consume it.
- **Audio file hosting:** Spotify for Podcasters CDN initially (free, fast). v1.0 adds Archive.org rehosting as a config option for permanent independence.

**Format-agnostic data model (load-bearing for the FeedBurner long-game):**
- Use generic terms throughout: `Feed` (collection) and `Item` (entry). Not `Show` and `Episode`. The terms map directly to the RSS/Atom specs and don't bake in podcast assumptions.
- v0.1.0 deliberately does NOT add a `media_type` discriminator on `Item` (YAGNI — v0.1.0 only ships `SpotifySource` and only handles podcast audio). v0.3+ adds `media_type: audio | video | text | image` when the second source type lands. Adding it later is a one-line struct field plus a switch in the renderer; it costs nothing to defer and would be dead code in v0.1.0.
- Format-specific fields (audio duration, video resolution, text body, podcast chapters, iTunes metadata) currently round-trip through `etree` as raw XML elements and never need to be modeled in Go structs. v0.2+ (PocketBase) introduces typed columns when the admin UI needs to read/write them.

**Pluggable source adapters:**
- The rewriter's pipeline is `Source → normalize → Output`, not "fetch from Spotify."
- v0.1 ships one source: `SpotifySource` (parses Spotify for Podcasters' auto-generated RSS).
- v0.3+ adds: `RSSSource` (any generic RSS feed), `AtomSource`, `JSONFeedSource`, `YouTubeSource`, `SubstackSource`, etc.
- Each source implements a small Go interface that converts external data into the generic `Item` shape. Adding a new source is a single file with no changes to the rest of the pipeline.

**Pluggable output renderers:**
- v0.1 outputs RSS 2.0 with the iTunes podcast namespace (what Apple Podcasts and Spotify expect for audio).
- v0.3+ adds Atom output, JSON Feed output, podcast 2.0 namespace, plain RSS for blogs.
- The same `Feed`/`Item` data renders into multiple formats from the same template engine.

These three architectural decisions — generic data model, pluggable sources, pluggable outputs — cost roughly 20 lines of Go in v0.1 and prevent a category of future rewrites. Without them, expanding to text/video later is a rewrite. With them, it is additive.
- **License:** **AGPL-3.0 for the entire repo, including widgets.** Closes the SaaS loophole — anyone who modifies the code and runs it for others must share their changes. Standard for self-hostable infrastructure (Mastodon, Plausible, Nextcloud, Matomo, Bitwarden all use AGPL). Aligns with the Sage open-knowledge philosophy.
  - **Why widgets can also be AGPL:** Embedding a remote script (`<script src="https://startr.media/widgets/player.js">`) does not "infect" the host site with copyleft. The host site references remote content; it does not redistribute or run modified AGPL code. AGPL §13 applies to people *running modified versions* as a network service, not to people *referencing* the original. This is exactly how Plausible Analytics operates: AGPL-licensed, embedded by tens of thousands of non-AGPL websites via a remote `<script>` tag, no legal friction in years of operation.
  - **The clean line:** If you embed via `<script src>` from a startrcast instance, you have zero AGPL obligations. If you fork the widget code, modify it, and serve it from your own infrastructure, AGPL applies and you must share your modifications. The README must explain this distinction explicitly so downstream users are not confused.
  - **Trademark:** The project name (`startr/cast` or whatever it lands on) is reserved separately from the code license. Forks must use a different name. Add a `TRADEMARK.md` to the repo. Standard practice — Mozilla/Firefox is the canonical precedent.
  - **CLA decision (deferrable):** No CLA for v0.1; accept contributions under DCO (Developer Certificate of Origin) like Plausible does. Reconsider if a future business model needs dual-licensing flexibility.

**Key architecture principle:** Own the feed XML, treat audio host as interchangeable. Subscribers bind to `feed.startr.media/v1/[show-slug].xml`, never to a third-party URL. Audio files can move between Spotify, Archive.org, S3, or self-hosted without breaking any subscriber.

**Feed URL scheme:** `feed.startr.media/v1/[show-slug].xml` — versioned from day one. The `/v1/` path lets the team evolve the feed format (custom namespaces, additional metadata) in `/v2/` later without breaking existing subscribers. This is a one-character decision that prevents a category of future pain.

**Versioning roadmap:**

- **v0.1.0 (this week):** ~200-line Go binary. Fetches Spotify feed, rewrites it, outputs static XML at `feed.startr.media/v1/[show-slug].xml`. Single-shot scheduler-agnostic binary; the canonical `startr.media` instance is driven by GitHub Actions, but `examples/` ships configs for cron, systemd, fly.io, k8s, and Docker Compose so self-hosters never get locked into one runtime. No widgets yet. Ships in a day. AGPL-3.0.
- **v0.1.1: Play widget + test backfill.** Embeddable HTML5 audio player as a web component. `<startrcast-player feed="https://feed.startr.media/v1/sage-resources.xml" episode="latest"></startrcast-player>` drops in any HTML page. Pure web component, no React/Vue dependency, encapsulated styles via shadow DOM. AGPL-3.0. Distributed via remote `<script src="https://startr.media/widgets/player.js">` so embedders are not bound by AGPL (Plausible model). **Also:** backfill the ~29 non-critical tests deferred from v0.1.0 (source fetch happy path, cache module, CLI flag parsing, disk-full / permission-denied / concurrent-write edge cases). Target: feature complete with full test coverage in the same week as v0.1.0 ship.
- **v0.1.2: Subscribe widget.** Embeddable subscribe button group. `<startrcast-subscribe feed="..." apple="..." spotify="..." rss="..."></startrcast-subscribe>` renders the standard "Listen on Apple Podcasts / Spotify / RSS" badges. Configurable per-show. AGPL-3.0. Same remote-script distribution.
- **v0.1.3+:** Polish patches as needed. Once the three sites have the feed + play + subscribe widgets and the team is publishing regularly, v0.1 is functionally complete.
- **v0.2.0 (after first 3-5 episodes work):** PocketBase backend wrapping the rewriter. Editable feed and item records. Admin UI for show notes, guests, transcripts. Astro/Next.js frontend at `startr.media` consumes the PocketBase API. This is "beautiful and open-sourceable." Data model is already generic (Feed/Item, not Show/Episode) so the future expansion is invisible to existing users.
- **v1.0 (open source release — podcast-shaped):** README with install instructions, deploy-button to fly.io, example configs, AGPL documented with embed-vs-fork explainer, Archive.org rehosting option for permanent audio independence. "Show HN: We built our own podcast library because Spotify owns your feed."
- **v0.3+ (post-v1.0 — FeedBurner expansion):** Add `RSSSource`, `AtomSource`, `JSONFeedSource`, `YouTubeSource`. Add output renderers for Atom and JSON Feed. Add `media_type=text` and `media_type=video` support to the admin UI and frontend. This unlocks the broader positioning: not just a podcast tool, but a self-hostable feed proxy / FeedBurner replacement. Second "Show HN: We built the FeedBurner replacement Google won't ship."

## Approaches Considered

### Approach A: Embedded Podcast Hub (Minimal Viable)
  Summary: Self-host the RSS feed XML at `feed.startr.media`. Use a free audio host (Spotify for Podcasters or Archive.org) for the audio files. Embed a player on the three existing sites. No new site to build for Phase 1.
  Effort: S
  Risk: Low
  Pros:
    - Ships in days, not weeks
    - Answers "do you have a podcast?" immediately
    - Free audio hosting; own the feed from day one
    - No new site to build for Phase 1
  Cons:
    - No standalone brand presence VCs/parents can bookmark yet
    - No unified discovery surface for new audiences
  Reuses: Existing three sites as distribution channels

### Approach B: Standalone Podcast Site
  Summary: Use the already-owned `startr.media` domain as a media hub for the whole ecosystem.
  Effort: M
  Risk: Medium
  Pros:
    - Own presence, bookmarkable by VCs and parents
    - Can cross-reference all three projects in one place
    - Scales as content library grows
  Cons:
    - More to build and maintain
    - Empty site is worse than no site — needs content first
  Reuses: RSS feed from hosting platform; design language from existing sites

### Approach C: Multi-format Content Engine
  Summary: Each roundtable becomes podcast + YouTube video + transcript blog post + newsletter excerpt. One recording, four distribution assets.
  Effort: L
  Risk: Medium
  Pros:
    - Maximum distribution surface
    - SEO via transcripts
    - Meets audiences where they are (audio, video, text)
  Cons:
    - High production overhead if done poorly
    - Can become a bottleneck that stops content from shipping
  Reuses: Existing newsletter(s), blog infrastructure on each site

## Recommended Approach

**Phase the ladder: A → B → C**

Phase 1 / v0.1 (this week): Build a small Go binary (~200 lines) that fetches the Spotify for Podcasters feed and rewrites it to point at `feed.startr.media/v1/[show-slug].xml`. Output static RSS XML. Deploy XML to Cloudflare Pages or Backblaze B2. Single-shot binary, scheduler-agnostic — the canonical `startr.media` instance uses GitHub Actions on a 15-minute cadence, but the binary runs equally well from cron, systemd, fly.io, or a manual terminal. Submit feed URL to Apple Podcasts and Spotify. Ship v0.1.0.

Then ship v0.1.1 (play widget — web component) and v0.1.2 (subscribe widget — web component) as fast follows. Embed both on a `/podcast` subpage of the readiest of the three sites. The "do you have a podcast?" answer becomes yes — and you own the feed from day one.

Phase 2 / v0.2 (after first 3-5 episodes work, ~45-60 days from first episode): Wrap the rewriter in PocketBase. PocketBase handles episode records, transcripts, guests, admin UI, REST API. Astro or Next.js frontend at `startr.media` consumes the PocketBase API. This is the beautiful, public-facing hub. VCs and parents can bookmark it.

Phase 3 / v1.0 (open source release, ongoing): Polish for downstream users. README, install docs, fly.io deploy button, example configs, MIT license. Archive.org rehosting option for permanent audio independence. Submit to Hacker News. Build a documented content workflow (tool decision by episode 3 — Notion, Linear, or a markdown checklist in the repo): record → edit → upload audio → PocketBase ingests → static site rebuilds → YouTube → transcript post → newsletter excerpt.

The trap to avoid: do not wait for the multi-format engine before shipping the first episode.

## Questions - Answers

- **Which of the three sites is the primary home for the podcast brand? Or is a new domain the right call?**
	- We have a domain: `startr.media`. This becomes the network hub.
- **Is there existing recording equipment / workflow, or does that need to be set up?**
	- Yes, there is existing recording equipment.
	- Simple workflows have been done for other productions.
	- One should be specified for this production (see Phase 3 checklist workflow).
- **What's the editorial angle — is it one show, or three shows (one per project)?**
	- Three interwoven shows based on articles from `sage.is/resources`, `sage.education/posts`, and `startr.style` (third needs article publishing in place).
	- *Still open:* Are these three separate feeds (three Apple Podcasts listings) or one umbrella show with three series? See Open Questions below — this decision affects feed architecture and should be made before the first episode ships.
- **Who owns publishing cadence? Weekly, bi-weekly, as-recorded?**
	- Initially weekly. Across three shows this means roughly one episode per show every three weeks, or three separate weekly episodes if capacity allows.
- **What recording setup exists? Microphones, recording software, room acoustics. This is a Phase 1 blocker — cannot record without it.**
	- Microphones, recording software, and room acoustics are all dealt with. No longer a blocker.
- **YouTube channel strategy: one unified Sage channel, or one per project?**
	- Still to be decided. Linked to the show structure decision above — if three separate feeds, three YouTube channels likely makes sense. If one umbrella show, one channel.

## Open Questions (decisions needed before shipping)

1. **Show structure:** Three distinct feeds (three Apple Podcasts subscribe buttons, three brands) OR one umbrella "Startr Media" feed with three labeled series? This is the most important decision — it sets the feed architecture, the subscribe flow, the YouTube strategy, and how content is cross-promoted. Recommendation: start with ONE umbrella feed with three clearly labeled series. Easier to launch, simpler to grow, and a single "Subscribe" CTA is always better than three. Can split later if any series gets enough pull to stand on its own. The Go rewriter can output one feed or three with a config flag — the architecture supports both.
2. **First episode source:** Which of the three sites has the article content and roundtable participants ready to record this week? Ship from the readiest.
3. **YouTube channel strategy:** Unified "Startr Media" channel (one home for all video) or one per project. Recommendation: one unified channel — same logic as show structure.
4. **Open source project name:** DECIDED — `Startr/feeds`. Direct, descriptive, on-brand with the Sage/startr ecosystem, doesn't bake in a specific positioning. Local intranet name during development: `WEB-Startr.media`. Public GitHub repo (`Startr/feeds`) goes live today or tomorrow.
5. **License confirmed:** AGPL-3.0 for the entire repo (server + widgets). Document in the README: (a) the SaaS-loophole reasoning, (b) the embed-vs-fork distinction, (c) the Plausible precedent, (d) trademark reservation for the project name.
6. **Audio rehosting in v1.0:** Archive.org integration is recommended for permanent audio independence, but adds workflow complexity. Confirm before v1.0 ships.
7. **Phase 3 workflow tool:** Notion, Linear, markdown checklist. Decision deadline: episode 3.

## Success Criteria

- **v0.1.0:** Feed rewriter shipped, first episode live, RSS XML self-hosted at `feed.startr.media/v1/[show-slug].xml`, accepted by Apple Podcasts and Spotify, within 2 weeks. Go rewriter binary committed to a public AGPL-3.0 repo.
- **v0.1.1:** Play widget shipped as AGPL-3.0 web component (distributed via remote `<script src>` so embedders are not bound by AGPL), embedded on at least one of the three sites.
- **v0.1.2:** Subscribe widget shipped as AGPL-3.0 web component, embedded alongside the play widget.
- **v0.2:** PocketBase backend running, beautiful static frontend at `startr.media`, episode admin workflow used by the team, within 45-60 days of first episode.
- **v1.0:** Open source release with README, install docs, deploy button, MIT license. At least one external team has installed and run an instance. Submitted to Hacker News.
- Each new episode generates at least 3 content assets (audio + video + text).
- Qualitative: VCs and parents can be pointed to a URL that gives them what they need without a 1:1 call.
- **Metrics per phase:** v0.1: plays per episode, referral source (which of the three sites drives listens). v0.2: unique visitors to `startr.media`, returning visitor rate. v1.0: GitHub stars on the open source repo, number of external installs, newsletter signups from transcript pages.
- **Kill switch:** If after 5 episodes total plays are below 50 per episode, pause production and reassess whether the format, audience, or distribution is the problem. (The open source tool can ship regardless — the podcast and the tool have separate viability.)

## Distribution Plan

**Architecture principle:** Own the feed, rent the audio host.

- **RSS feed XML:** Self-hosted at `feed.startr.media/[show-slug].xml`. Either a static file or a tiny generator in the repo. The feed is the asset you must own forever. Three feeds (one per show) since this is a 3-show network.
- **Audio file hosting:** Interchangeable. Start with Spotify for Podcasters (free, fast to set up) or Archive.org (free, permanent archive). Enclosure URLs in the self-hosted feed point at whatever host you're using. Swap hosts anytime without breaking subscriber feeds.
- **Why this works:** Subscribers bind to the feed URL (`feed.startr.media/sage-resources.xml`), not the audio URL. You can migrate audio files between Anchor, Archive.org, S3, or a paid host like Transistor without anyone noticing.
- **When to upgrade audio hosts:** Move to Transistor (~$19/mo) only if you hit one of these signals: (a) you want real analytics, (b) Anchor's free tier adds friction, (c) you cross 1000+ listens per episode and want infrastructure you can lean on.
- **Apple Podcasts / Spotify distribution:** Submit the self-hosted feed URL to Apple Podcasts Connect and Spotify for Podcasters directly. They ingest your feed; you never lose control.
- **Video:** YouTube channel (structure TBD — see Questions - Answers)
- **Web:** Embedded player via iframe on existing sites (Phase 1), `startr.media` static site hub (Phase 2)
- **Newsletter:** Excerpt per episode, link to full episode on `startr.media`
- **Transcripts:** Auto-generated via Whisper (free, runs locally) or the audio host's built-in transcription, lightly reviewed for proper nouns, published alongside each episode from day one

CI/CD: Existing deployment pipelines for the three sites cover Phase 1 embeds. Feed XML and Phase 2 hub deploy via Vercel/Netlify with auto-deploy on push.

## The Assignment

Recording gear is ready. Domain is ready. This week:

1. Decide the show structure (one umbrella feed with three series, or three distinct feeds — see Open Questions). Recommended: one umbrella feed.
2. Pick the project name (NOT `startr/cast` — too podcast-specific given the Stage 2 FeedBurner positioning). See Open Questions for candidates.
3. Create a public GitHub repo under the chosen name. AGPL-3.0 license at repo root. Add `TRADEMARK.md` reserving the project name. Add a README stub with the embed-vs-fork explainer. Use DCO for contributions (no CLA).
4. Write the v0.1.0 Go rewriter (~200 lines) with format-agnostic foundations: `Feed` and `Item` types (NOT `Show`/`Episode`), a `Source` interface with `SpotifySource` as the first implementation, an `Output` interface with `RSS2PodcastRenderer` as the first implementation. Parse with **`beevik/etree`** (DOM-style, preserves all unknown elements and namespaces on round-trip — `encoding/xml` is wrong for a rewriter because it silently drops anything not declared in a struct). Rewrite branding/links/images/channel metadata to point at `feed.startr.media/v1/[show-slug].xml`. **Leave `<enclosure url>` alone** — Spotify hosts the audio bytes for free, and rewriting the URL would force us to proxy or re-host. Output transformed XML via atomic write-rename. The architecture is generic; only the v0.1 surface is podcast-shaped.
5. Set up DNS: `feed.startr.media` → static file host (Vercel/Netlify/S3+CloudFront). Make sure the `/v1/` path prefix is in place from day one.
6. Set up the scheduler of choice for the canonical `startr.media` instance. Default plan is GitHub Actions on a 15-minute cadence, but the binary is single-shot and scheduler-agnostic — pick whatever fits the environment (`examples/` ships configs for cron, systemd, fly.io, k8s, Docker Compose).
7. Record the first roundtable drawing from the readiest of sage.is/resources, sage.education/posts, or startr.style.
8. Upload audio to Spotify for Podcasters. Wait for Spotify's auto-feed to update. Confirm the rewriter picks up the new episode and publishes the rewritten feed at `feed.startr.media/v1/[show-slug].xml`.
9. Submit the rewritten feed URL to Apple Podcasts Connect and Spotify for Podcasters directly.
10. Ship v0.1.1: build the play widget as a vanilla web component (`<startrcast-player>`, or whatever the renamed prefix becomes) under AGPL-3.0, in `widgets/player/`. No framework dependencies. Reads from the feed URL. Ship the bundled JS at `https://startr.media/widgets/player.js` for remote `<script src>` embedding.
11. Ship v0.1.2: build the subscribe widget as a vanilla web component under AGPL-3.0, in `widgets/subscribe/`. Configurable per-platform. Same remote-script distribution.
12. Embed both widgets on a `/podcast` subpage of the matching site.

Everything else is v0.2 (PocketBase backend) or later.

## Engineering Review Decisions (2026-04-07)

These are the implementation specifics locked during the `/plan-eng-review` pass. They sit downstream of the high-level architecture and are binding for the v0.1.0 build. Anything not on this list is fair game for the implementer.

### 12 locked decisions

| # | Decision | Locked-in choice | Rationale |
|---|---|---|---|
| 1 | XML library | `beevik/etree` for round-trip preservation. Hybrid fallback to `encoding/xml` only if etree proves problematic. | A rewriter must preserve every element from upstream (iTunes namespace, podcast 2.0 namespace, Spotify custom tags). `encoding/xml` requires declaring a struct field for every element it should keep — anything not declared is dropped silently. `beevik/etree` is DOM-style: parse → mutate → serialize, with full round-trip preservation. This is the single most important technical decision in the build. |
| 2 | Interface shape | Two-method byte-oriented. `Source.Fetch(ctx) ([]byte, error)` returns raw upstream XML; `Output.Render(ctx, []byte) error` takes rewritten XML and writes it. | Byte-oriented keeps the contract narrow and lets etree handle the structure. No premature struct modeling. |
| 3 | Deployment target | Static host: Cloudflare Pages or Backblaze B2. Either works; both are pure static-file. | Both are scheduler-agnostic, both are cheap, both have first-class CDN. Final pick during DNS setup. |
| 4 | **Cron mechanics (revised)** | Binary is single-shot, scheduler-agnostic, no built-in daemon. Canonical `startr.media` happens to use GitHub Actions; downstream users pick anything. README ships ~6 example configs in `examples/`. | Locking the cron into GitHub Actions would force Raspberry Pi self-hosters to install a GitHub Actions runner. The architecture must walk the show's "real autonomy" thesis. Examples directory is the contract: cron, systemd, fly.io, k8s CronJob, Docker Compose, GitHub Actions all supported equally. |
| 5 | Logging | Go `log/slog` stdlib structured logging. JSON output by default in production, text output in dev. | Stdlib, no dependency. Structured fields make scheduler log-parsing trivial. |
| 6 | Config | CLI flags primary, optional `--config /path/to/feeds.yaml` for repeated invocations. CLI flags always win over YAML. | Smallest possible surface area for v0.1.0. YAML is just sugar for the same flags. |
| 7 | Error handling | `errors.Join` and `fmt.Errorf("%w: ...", err)` wrap. Single error return per pipeline stage. Non-zero exit on any error. | Idiomatic Go. Scheduler-friendly: any non-zero exit is a failure the scheduler must surface. |
| 8 | Upstream failure | Fail loud, keep last good output. Atomic write-rename via `os.Rename` on a temp file in the same directory as the target. | Subscribers continue to see the last successful feed even during outages. Operator sees the failure via the scheduler. No silent corruption. |
| 9 | PocketBase v0.2 transition | PocketBase regenerates static XML to the same `/v1/[show-slug].xml` path. The URL contract never changes. | Subscribers bound to v0.1 URLs continue working forever. The backend swap is invisible. |
| 10 | Distribution pipeline | Full pipeline: GoReleaser builds linux/darwin/windows × amd64/arm64 binaries, Sigstore cosign keyless signing, SLSA Level 3 provenance attestation, ghcr.io container image, GitHub Releases page. ~30 minutes of CC time vs. ~5 minutes for go-install only. | The delta is ~25 minutes of CC for an open-source dev tool to look professionally distributed from day one. Self-hosters can verify checksum + cosign signature without installing the Go toolchain. |
| 11 | HTTP caching | Conditional GET. Store `ETag` and `Last-Modified` from each successful fetch in `.feeds-state.json`. Send `If-None-Match` / `If-Modified-Since` on the next fetch. On HTTP 304, short-circuit: no parse, no rewrite, no write, exit zero with "upstream not modified" log line. | Spotify's auto-feed updates infrequently (a few times per week). At a 15-minute cadence, ~98% of fetches will be 304s. This single feature cuts bandwidth and write churn by ~50x. |
| 12 | `media_type` discriminator | YAGNI for v0.1.0. Defer until v0.3 when the second source type (text or video) lands. | v0.1.0 only ships `SpotifySource` and only handles podcast audio. Adding the field now is dead code. Adding it in v0.3 is a single struct field plus a switch in the renderer. |

### Project layout

```
Startr/feeds/
├── cmd/feeds/main.go             # CLI entry point: flag parsing, top-level orchestration
├── internal/
│   ├── source/                   # Source adapters
│   │   └── spotify.go            # SpotifySource: conditional GET, etree parse
│   ├── output/                   # Output renderers (named to avoid stdlib conflict)
│   │   └── rss.go                # RSS2PodcastRenderer: atomic write-rename
│   ├── rewriter/                 # Core etree-based rewrite logic
│   │   └── rewrite.go            # All XML element rewrites live here
│   ├── pipeline/                 # Source → rewrite → output orchestration
│   │   └── pipeline.go           # Failure handling, last-good preservation
│   └── cache/                    # State file + conditional GET helpers
│       └── cache.go              # .feeds-state.json read/write
├── examples/                     # Scheduler examples (the contract for portability)
│   ├── github-actions.yml
│   ├── systemd-timer/feeds.service
│   ├── systemd-timer/feeds.timer
│   ├── cron-on-raspberry-pi.txt
│   ├── flyio-scheduled-task.toml
│   ├── kubernetes-cronjob.yaml
│   └── docker-compose-cron.yml
├── feeds.example.yaml            # Example config file
├── README.md                     # Install, run, embed-vs-fork explainer
├── LICENSE                       # AGPL-3.0
├── TRADEMARK.md                  # Project name reservation
├── DCO.txt                       # Developer Certificate of Origin
├── go.mod
├── .goreleaser.yaml              # Release pipeline
└── .github/workflows/release.yml # Triggers GoReleaser on tag push
```

The canonical `startr.media` instance also ships a `.feeds-state.json` (gitignored) for HTTP cache state.

### Rewrite scope (what gets touched, what stays)

| XML element | Action | Why |
|---|---|---|
| `<channel><title>` | REWRITE | Show branding |
| `<channel><link>` | REWRITE | Point at the startr.media show page |
| `<channel><description>` | REWRITE | Optionally append "Rewritten by Startr/feeds" attribution |
| `<channel><image><url>` | REWRITE | Point at startr.media-hosted cover art |
| `<atom:link rel="self">` | REWRITE | **REQUIRED by Apple Podcasts spec.** Must point at `feed.startr.media/v1/[show-slug].xml`. If missing from upstream, INJECT it. |
| `<itunes:author>` | REWRITE | Show branding |
| `<itunes:owner>` | REWRITE | Show branding |
| `<itunes:image>` | REWRITE | Same as channel image |
| `<item><enclosure url>` | **LEAVE ALONE** | **Load-bearing.** Spotify hosts the audio bytes for free. Rewriting this URL would force us to proxy or rehost audio, which kills the v0.1.0 cost model. Tested explicitly. |
| `<item><guid>` | LEAVE ALONE | Subscriber-stable identifier — rewriting it triggers duplicate-episode detection in podcast clients |
| `<item><pubDate>` | LEAVE ALONE | Source of truth, no transform needed |
| `<item><title>` | LEAVE ALONE | Episode-level content, not branding |
| `<item><description>` | LEAVE ALONE | Episode-level content |
| `<item><itunes:*>` (episode-level) | LEAVE ALONE | Episode metadata, not channel branding |
| Any unknown namespaced element | LEAVE ALONE | etree round-trip preserves automatically; this is the whole reason etree was chosen over `encoding/xml` |

### Test requirements

**v0.1.0 ships with ONLY the 5 critical tests.** Decision locked 2026-04-07. The other ~29 tests backfill in v0.1.1. This is a deliberate speed-over-coverage call: get the release out, the 5 critical tests gate the worst failure modes, and the backfill sprint happens immediately after.

1. **`TestRewrite_PreservesITunesNamespace`** — Parse a Spotify feed with iTunes namespace declaration, run the rewriter, serialize, re-parse. Assert iTunes namespace declaration survives. (If this fails, Apple Podcasts rejects the feed.)
2. **`TestRewrite_DoesNotRewriteEnclosureURL`** — Parse a feed, rewrite, assert every `<enclosure url>` is byte-identical to upstream. (If this fails, audio doesn't play.)
3. **`TestRewrite_AtomSelfLinkRewritten`** — Parse a feed, rewrite, assert `<atom:link rel="self">` points at `feed.startr.media/v1/[show-slug].xml`. Test the inject-when-missing path too. (Required by Apple Podcasts spec.)
4. **`TestPipeline_FailLoudPreservesLastGood`** — Set up an output file with known content. Stub the source to return an error. Run the pipeline. Assert non-zero exit AND output file is unchanged byte-for-byte.
5. **`TestPipeline_IdempotentNoOpWrite`** — Run the pipeline twice with unchanged upstream. Assert second run does not modify the output file's mtime (idempotency check downstream of the 304 short-circuit).

**v0.1.1 test backfill scope** (ships immediately after v0.1.0, same week if possible):

| File | Tests needed | Priority |
|---|---|---|
| `internal/source/spotify_test.go` | ~6 (fetch happy path, 304, 5xx, malformed XML, conditional GET headers, context timeout) | High |
| `internal/output/rss_test.go` | ~6 (write happy path, atomic rename, disk-full simulation, permission denied, idempotent no-op, concurrent write) | High |
| `internal/rewriter/rewrite_test.go` | ~10 (the 3 critical tests above + channel title rewrite, channel link rewrite, channel image rewrite, iTunes author rewrite, missing channel title, empty feed, unknown namespace preservation) | **Critical** |
| `internal/pipeline/pipeline_test.go` | ~6 (the 2 critical tests above + happy path, source error, output error, state file corrupted) | **Critical** |
| `internal/cache/cache_test.go` | ~4 (load happy path, load missing file, load corrupted file, save happy path) | Medium |
| `cmd/feeds/main_test.go` | ~2 (CLI flag parsing, --config YAML load) | Low |

**~34 tests total across both releases.** v0.1.0 ships with the 5 critical rewriter/pipeline tests only. The remaining ~29 backfill in v0.1.1. The full table plus edge cases lives in the test plan artifact at `agent-develop-eng-review-test-plan-20260407-120000.md` and is the input for the `/qa` skill.

**Risks accepted by deferring to v0.1.1:**

- Source fetch path has no unit tests for a week. A Spotify schema change between v0.1.0 ship and v0.1.1 backfill would fail loud in production (good: fail-loud guarantees last-good preservation) but without a test to catch drift in CI.
- Cache module has no unit tests for a week. A state-file regression would self-heal on next run but might silently disable conditional GETs, causing full fetches every 15 minutes instead of mostly 304s. Bandwidth cost is marginal; detection is via log volume.
- CLI flag parsing has no unit tests. A flag regression would crash the binary at startup in a way the scheduler surfaces immediately.
- No edge-case coverage on disk-full, permission denied, concurrent writes. All three rely on atomic write-rename correctness, which is kernel-guaranteed but not test-verified.

The 5 critical tests are the ones that map to production-breaking, invisible-failure scenarios. The backfill tests map to either fail-loud scenarios (caught by the scheduler) or low-impact regressions (caught by observation). The tradeoff is defensible for a v0.1.0 ship.

### Performance characterization

| Metric | v0.1.0 expected | Concern level |
|---|---|---|
| Compiled binary size | ~10 MB (Go static binary, etree adds ~200KB) | None |
| Cold-start latency | <50ms (no JIT, no runtime) | None |
| Fetch latency (Spotify 200) | 200-500ms typical | None |
| Fetch latency (Spotify 304) | 50-150ms typical | None — and 304s are ~98% of runs |
| Parse latency (etree, 50-200KB feed) | <5ms | None |
| Rewrite latency (dozens of element edits) | <1ms | None |
| Write latency (atomic write-rename, SSD) | 10-50ms | None |
| Total run wall-clock (200 path) | <2 seconds | None |
| Total run wall-clock (304 path) | <500ms | None |
| Memory (peak) | 20-40 MB, dominated by etree DOM | None |
| Bandwidth per run (200) | 50-200 KB | None |
| Bandwidth per run (304) | ~1 KB headers | None |

**At 15-minute cadence × 1 feed:** 96 runs/day, ~92 of those are 304s, ~4 are full fetches. Annualized bandwidth: a few MB. This is not a perf surface for v0.1.0.

**One scaling consideration to flag for v0.3+:** etree holds the entire feed in memory as a DOM. For pathological feeds (10,000+ items, ~10MB XML — possible for some text feeds, never for podcasts), peak memory could hit 50-100MB. v0.1.0 only handles podcast feeds (max ~500 items, ~5MB), so this is irrelevant. Revisit when text feeds land in v0.3.

**No optimization work in v0.1.0.** Premature.

### Failure modes table

| Failure | Likelihood | Impact | Detection | Mitigation |
|---|---|---|---|---|
| Spotify upstream 5xx | Medium (monthly) | Low | Scheduler reports failure | Last-good preserved, fail loud, no manual action needed (Path 4 in test plan) |
| Spotify upstream malformed XML | Low | Low | etree parse error, fail loud | Last-good preserved |
| Spotify silently changes RSS schema | Low | Medium | Apple/Spotify reject feed downstream | The 5 critical tests run against fixture data; add a CI smoke test that fetches the real Spotify feed weekly to catch drift |
| `beevik/etree` dependency abandoned | Low | Medium | Go mod tooling notice | Hybrid fallback to `encoding/xml` is in the design (Decision #1). Migration path is real but multi-day. |
| State file `.feeds-state.json` corrupted | Low | Low | Cache.Load returns empty state, heals on next run | Self-healing — next run does a full fetch with no conditional headers, then writes a fresh state file |
| Disk full on output write | Low | Low | os.Rename fails, last-good preserved | Atomic write-rename pattern |
| Permission denied on output path | Low | Low | os.Rename fails, last-good preserved | Same |
| Concurrent runs (scheduler + manual) | Medium | Low | Atomic rename guarantees one wins | No corruption possible. Both runs may succeed identically; or one is a 304 no-op. |
| GitHub Actions outage (canonical instance) | Low | Medium | Scheduler doesn't fire | Self-hoster docs include cron/systemd alternatives — same binary, different scheduler |
| Cosign key compromise | N/A | N/A | N/A | Sigstore keyless signing removes the need to hold a key |
| AGPL compliance challenge | Low | High | Lawyer letter | Licensed clearly at repo root, embed-vs-fork doc in README, DCO on contributions, trademark reserved separately |
| iTunes namespace declaration drops on round-trip | Medium → 0 (with test) | High | TestRewrite_PreservesITunesNamespace | The test prevents this from ever shipping |
| `<atom:link rel="self">` not rewritten or missing | Medium → 0 (with test) | High | TestRewrite_AtomSelfLinkRewritten | The test prevents this from ever shipping |
| `<enclosure url>` accidentally rewritten | Low → 0 (with test) | **Critical** | TestRewrite_DoesNotRewriteEnclosureURL | The test prevents this from ever shipping |

The three "→ 0 with test" rows are why the 5 critical tests are gating: they convert the most expensive failure modes from "possible in production" to "impossible to ship."

### Worktree parallelization strategy

Six worktrees can build in parallel because the interface contracts (Decision #2) are locked from this review. Each worktree builds against fakes of the other layers and unit-tests in isolation. Final integration happens in worktree A.

| Worktree | Scope | Files | Depends on |
|---|---|---|---|
| A | CLI + state + pipeline skeleton | `cmd/feeds/main.go`, `internal/cache/`, `internal/pipeline/` | Nothing — defines interfaces other worktrees implement |
| B | Spotify source | `internal/source/spotify.go`, fixtures | Worktree A's interfaces |
| C | Rewriter (the 3 critical rewriter tests) | `internal/rewriter/rewrite.go`, `internal/rewriter/rewrite_test.go`, fixtures | Nothing — pure function on `[]byte` |
| D | RSS output renderer | `internal/output/rss.go`, atomic write-rename helper | Nothing — pure function on `[]byte` |
| E | Examples directory | `examples/*` (6 scheduler configs) | Nothing — docs only |
| F | Distribution pipeline | `.goreleaser.yaml`, `.github/workflows/release.yml`, cosign config, ghcr.io publishing | Worktree A merged (needs main.go to build) |

Worktrees C and D are highest leverage because they contain the critical tests and the load-bearing logic. Run them first in parallel.

## First Episode Plan (Pilot)

**Show working title:** `Run It Local` (proposed — see name research below). The previous suggestion `Self-Hosted` is unavailable (active competing podcast). `Run It Local` is double-coded — it reads as an operator phrase for engineers and as a parenting phrase for homeschool families, both audiences finding their own true reading. Open to alternatives.

**Hosts:** Izzy Plante (CEO) and Alex Somma (CTO). CEO/CTO pairing is a proven format — Izzy leads on context, vision, and "why this matters." Alex leads on technical detail and "what's actually happening under the hood." They riff naturally; no formal interview structure.

**Show thesis: Real autonomy vs. compliance theater.**

There is a difference between *performative* sovereignty and *real* sovereignty.

Performative sovereignty is when Microsoft announces a $19 billion Canadian AI data center investment and brands it as "protecting Canada's digital sovereignty" — while remaining a US company subject to the CLOUD Act, which lets US authorities subpoena Canadian data without notice and without Canadian judicial review. It's when Google launches "Google Sovereign Cloud" with a "Google Data Boundary." It's when a school district tells parents their kids' data is safe because Microsoft 365 for Education is "FERPA-compliant" and "Canadian-hosted." Policy critics already have a name for this: **sovereignty-washing**.

Real sovereignty is when you run the model on your own hardware, hold the keys yourself, and no third party can be subpoenaed for your data because no third party has it. It is not about where the server is. It is about who controls access, who holds the keys, and who can say no.

The show is about that gap. From people who actually run their own GPUs and ship infrastructure that respects this distinction. There are no other podcasts doing this from an operator perspective.

**Why this thesis works for both audiences:**

- **Engineers and infrastructure builders** read this as operator content: how do you actually self-host, what does the threat model look like, what are the real cost tradeoffs, who's lying about what.
- **Tech-literate parents and homeschool families** read this as parenting content: my kid's homework prompts, my kid's emotional health questions, my kid's writing — none of this should be in a vendor's training data or subject to a foreign government's subpoena. "Compliance" is not the same as "safety."
- Both audiences are concerned with the same underlying problem: vendors and governments are calling something "sovereign" or "compliant" or "safe" when it isn't, and the people who notice are the people who actually understand the stack.

The show speaks to both without code-switching. The voice does not change between audiences. The same story lands for both, for related but distinct reasons.

**Target runtime:** 30-35 minutes. Long enough to establish substance, short enough for a single sitting. First episodes should be slightly more produced/structured than later episodes because new listeners arrive cold.

**Source content (from `sage.is/resources`, all April 2026):**
- Anchor story: **The Prompt You Thought Was Private** (Apr 2) — Perplexity AI lawsuit, prompts shared with Meta and Google via hidden tracking scripts. The cleanest illustration of the thesis: "approved" privacy tools that are doing the opposite of what they advertise.
- Supporting story: **The $70,000 Illusion** (Apr 4) — vendor pricing vs. actual self-hosted infrastructure costs. The economic case for real sovereignty: it's not just safer, it's cheaper.
- Supporting story: **The Thinking Layer** (Apr 2) — open-source AI reasoning is here, who owns the hardware running your thinking. The technical foundation that makes real sovereignty possible right now (it wasn't a year ago).

**External news peg to weave in:** Microsoft's $19B / $7.5B Canadian AI data center investment, marketed as protecting Canadian "digital sovereignty" while still being a US company subject to the CLOUD Act. Google's "Sovereign Cloud" with "Google Data Boundary." Both are the textbook example of sovereignty-washing the show is calling out. Episode 1 should reference at least one of these by name to ground the abstract thesis in current news.

**Episode structure:**

| Section | Time | What happens |
|---|---|---|
| Cold open | 0:00-1:00 | A surprising line that hooks the listener. Something like: "Microsoft just spent $19 billion telling Canada they're protecting your sovereignty. They're not. We're going to talk about who actually is, and why your kid's homework prompts matter to this story." |
| Intro | 1:00-3:30 | Who Izzy is. Who Alex is. What Sage is in one sentence. What this podcast is — operator-flavored AI commentary about the gap between what companies and governments call "sovereign" and what actually is. |
| Why this exists | 3:30-7:00 | "We write about this stuff at sage.is/resources. People kept asking if we had a podcast. So here we are." Brief, no fluff. Mention the parent/homeschool audience explicitly — this is for them too. |
| Story 1: Perplexity (anchor) | 7:00-17:00 | Walk through what happened. The hidden tracking. Why this is particularly bad for AI prompts (you don't search "what's the weather," you search things you wouldn't put in writing anywhere else). What it means for healthcare, education, legal — AND for the homeschool parent whose kid asked an AI for help with a writing assignment and just had that prompt sent to Meta. The "compliance theater" beat: Perplexity has a privacy policy. It's been reviewed by lawyers. It's "compliant." It's also doing exactly the thing the policy says it doesn't. Compliance is not safety. |
| Story 2: $70K Illusion + sovereignty-washing | 17:00-25:00 | Cost analysis from the article: what you actually pay vendors vs. what self-hosting costs. Tie it directly to the Microsoft Canada announcement: Microsoft's $19B "sovereign" investment is structured so that Canadian taxpayers and Canadian customers pay US vendor margins for infrastructure that the US can subpoena. Real sovereignty is cheaper AND safer. Numbers, specifics, named villains. |
| Story 3: Thinking Layer | 25:00-30:00 | Brief but important — open-source reasoning models are now actually shipping. The hardware to run them is now affordable. The thing that makes real sovereignty *technically possible* right now is brand new. If the answer used to be "you have to use a vendor because nothing else is good enough," that excuse expired in 2026. Sage is positioned for this moment. |
| Close | 30:00-33:00 | Where to find Sage (sage.is, sage.education). Where to find these articles (links in show notes). The show going forward will cover one or two stories per episode through this same lens. How to get in touch. One specific call to action: "If you're a parent or a teacher, the next time someone tells you a tool is 'compliant' or 'safe,' ask them whose law actually governs your kid's data. If the answer involves any US company, the answer is the CLOUD Act." Subscribe call. |

**Editorial principles for episode 1:**
- Each story has a "what happened" / "why it matters" / "what we'd do differently" arc
- No defensive corporate tone. The Sage voice is direct and a little prickly. Honor it.
- Quote the actual article titles in show notes; link to each one
- Speak to both audiences without code-switching: every story should land for both engineers and parents because the underlying concern is the same
- Avoid jargon when a plain word will do. "The site sends your prompts to Meta" is better than "the application exfiltrates user input via third-party tracking pixels" — the technical audience reads both and understands; the parent audience reads only the first
- End with one specific thing the listener should do (read the article, try the platform, share with a colleague who's still using a vendor AI tool, or share with another parent in their homeschool group)

**Show notes deliverable:**
- Episode title (working: "Episode 1: The Prompt You Thought Was Private")
- 2-3 sentence description
- Bullet list of stories covered with timestamps
- Direct links to all three source articles
- Bios for Izzy and Alex (2-3 sentences each)
- Subscribe links (Apple Podcasts, Spotify, RSS via `feed.startr.media/v1/[show-slug].xml`)

**Pre-record checklist:**
- Both Izzy and Alex have read all three articles fresh (not just remember writing them)
- Outline lives in a shared doc — bullet points only, not a script
- Test recording (5 min) to confirm levels and vibe before the real take
- Record in one continuous take if possible; edit only for breath/dead air

**Post-record:**
- Light edit only: remove dead air, trim umms, no aggressive pacing changes
- Auto-transcribe via Whisper; review for proper nouns (especially "Perplexity," "Meta," "Sage")
- Generate cover art (placeholder if needed for episode 1; iterate later)
- Upload to Spotify for Podcasters; wait for auto-feed
- Run the v0.1.0 Go rewriter to publish to `feed.startr.media/v1/self-hosted.xml`
- Submit feed URL to Apple Podcasts Connect and Spotify for Podcasters
- Embed player on `sage.is/podcast` (the readiest of the three sites since the content draws from sage.is)

**Show name candidates (decide before recording — researched 2026-04-07):**

This research was done with standard WebSearch. To be redone with the team's internal `dark-search.production.openco.ca` engine on the next pass for verification.

Researched candidates, with availability and rationale:

| Name | Status | Notes |
|---|---|---|
| `Self-Hosted` | ❌ **Hard NO** | Active podcast on selfhosted.show (Jupiter Broadcasting), 144+ episodes, hosts named Chris and Alex Kretzschmar. The Alex overlap with Alex Somma would create direct confusion. |
| `Sovereign Stack` | ❌ Taken | sovereignstack.tech is an active project. |
| `Sovereign Computing` | ❌ Taken AND directly competitive | The Sovereign Computing Show (Jordan Bravo, Atlanta freedom-tech hackerspace) covers self-hosting, AI privacy. Worth listening to as an adjacent voice. |
| `Stay Local` | ❌ Taken | German university careers podcast. Different niche but the name is in active use. |
| `Quiet Stack` | ⚠️ Twitter handle taken | @Quiet_Stack on X is an onchain product foundry. No podcast yet but brand confusion likely. |
| `Run It Local` | ✅ **Available — recommended** | No podcast found. Double-coded: operator framing for engineers ("run it on your own hardware"), parent framing for homeschoolers ("we run it locally instead of trusting vendors"). Same words, two true readings. |
| `Off Cloud` | ✅ Available | Bold, direct, clear positioning. Slightly more tech-sided than parent-sided. |
| `Quiet Compute` | ✅ Available | More poetic. "Quiet" carries parenting resonance, "compute" grounds it technically. |

**Recommendation:** `Run It Local`. The double-coding is genuine (not strained), it works for the operator audience without alienating the parent/homeschool audience, and it captures the action — not just the philosophy but what the listener actually does. Available across podcast directories at the time of research. Verify trademark and check the .show / .fm / .com domains before committing.

**Adjacent shows worth knowing about as competitive landscape:**
- **The Sovereign Computing Show** (Jordan Bravo) — closest direct competitor in the operator-AI-privacy niche
- **Self-Hosted** (Jupiter Broadcasting) — broader self-hosting, less AI-focused
- **AI Parenting Podcast** — addresses the parent audience but at a basic literacy level, no operator angle

## Minimum Production Standards

- Audio: 128kbps minimum, mono is fine. Use a USB condenser mic or better (not laptop mic).
- Episode length: 15-45 minutes. Shorter is fine. Longer needs a strong reason.
- Metadata per episode: title, 2-3 sentence description, guest names and roles, show notes with links mentioned.
- Transcript: auto-generated (Whisper is free and runs locally; most hosted services also offer this), lightly reviewed for proper nouns. Published with episode.
- Intro/outro: keep it under 15 seconds each. Name the show, name the project, done.

## What I noticed about how you think

- You said "people are asking us more about what we do and how we do it and why we do it" — that's not vague audience interest, that's pull. You already have demand; you're just not capturing it yet.
- You mentioned the casual roundtables unprompted. That's the product trying to surface. Most founders would have said "we need to create content." You named the format that already works.
- When I presented three approaches, you pushed back and asked for a combination. That's good instinct — you weren't satisfied with a false choice.
- You named VCs and technical parents as distinct audiences without prompting. That's market segmentation thinking, not feature-list thinking.
- When I claimed Spotify locks you into their feed, you pushed back with your CTO's actual mechanism. You weren't deferential. You were right, and I was wrong. That kind of correction is rare and valuable — it kept the design from being built on a false constraint.
- You expanded the scope from "podcast site" to "open source podcast library tool" without losing the simplest-and-fastest instinct. That's the hardest balance in product work — ambition without scope creep. You named PocketBase specifically, which means you've already done the architectural thinking, you just wanted a sounding board.
- During the eng review, when I had locked the cron mechanics into GitHub Actions, you caught the lock-in: "locking in GitHub actions would mean our needing to walk in a GitHub actions runner for if somebody wanted to say keep this on a pie." That's the same instinct the show's thesis is about — don't tell people they're autonomous while subtly tethering them to a hyperscaler. The architecture now walks the thesis: single-shot scheduler-agnostic binary with six example configs in `examples/`. You eat your own dogfood.

## GSTACK Engineering Review Report

Generated by `/plan-eng-review` on 2026-04-07. Branch: `develop`. Artifact reviewed: `agent-develop-design-20260406-174520.md`.

### Review readiness dashboard

| Dimension | Status | Notes |
|---|---|---|
| Architecture | LOCKED | 12 decisions recorded in "Engineering Review Decisions" section. Project layout defined. Interface contract narrow (byte-oriented). |
| Code quality | LOCKED | Rewrite scope table defines what gets touched vs. left alone. `<enclosure url>` explicitly protected (load-bearing for free audio hosting). |
| Tests | LOCKED | v0.1.0 ships with 5 critical tests only (decision 2026-04-07). ~29 remaining backfill in v0.1.1 alongside the play widget. Risks documented in the test requirements section. |
| Performance | NO CONCERN | Sub-2-second wall clock, <40MB memory, ~98% of runs short-circuit on HTTP 304. Not a perf surface for v0.1.0. Flagged: etree DOM scaling for text feeds in v0.3+. |
| Failure modes | MAPPED | 13 failure modes tabled. 3 highest-impact modes (iTunes namespace drop, atom:self missing, enclosure rewrite) converted from "possible" to "impossible to ship" by the 5 critical tests. |
| Parallelization | DESIGNED | 6 worktrees can build in parallel. Contracts locked from this review; each worktree fakes adjacent layers and unit-tests in isolation. |
| Distribution | LOCKED | Full pipeline: GoReleaser multi-arch builds, Sigstore cosign keyless signing, SLSA L3 provenance, ghcr.io container, GitHub Releases. ~30 min CC delta over a bare `go install`. |

**Completion status: DONE.** All decisions locked, all open TODOs resolved, ready to implement.

### NOT in scope for v0.1.0 (explicit rejections)

- **Built-in scheduler / daemon mode.** Explicitly rejected this review. Would tie users to one process model and break the "real autonomy" thesis. Revisit only if a use case emerges that truly needs it, and even then as a separate `feeds serve` subcommand that ships alongside the single-shot `feeds rewrite`.
- **Multi-feed orchestration.** v0.1.0 handles one feed per invocation. v0.2 (PocketBase) handles many.
- **Audio rehosting.** v1.0+. Spotify hosts the audio bytes for free in v0.1.0; that's the whole cost model.
- **Web UI / admin.** v0.2 with PocketBase.
- **Metrics / observability dashboard.** v0.2+. v0.1.0 emits structured slog; the scheduler's log is the observability surface.
- **Custom per-item branding overrides.** v0.2+.
- **Multi-source merging / federation.** v0.3+.
- **Non-RSS sources (Atom, JSON Feed, YouTube).** v0.3+. The `Source` interface is designed to make this additive.
- **Non-audio media types (text, video).** v0.3+. The `media_type` field is intentionally NOT added in v0.1.0 (YAGNI).
- **Webhook notifications on update.** v0.2+.
- **API for programmatic control.** v0.2+ via PocketBase.
- **Authentication / multi-tenant.** v0.2+ via PocketBase.
- **Performance optimization (streaming parser, worker pool, caching layer).** v0.3+ if needed. v0.1.0 runs in under 2 seconds with <40MB memory — premature.
- **`media_type` discriminator field.** YAGNI. Add when the second source type lands in v0.3.

### What already exists in `/workspace`

- This design document and its canonical gstack copy (`agent-develop-design-20260406-174520.md`).
- The test plan artifact for `/qa` consumption.
- `DEVOPS-INBOX.md` (unrelated, tracks a dark-search 403 blocker).
- `CLAUDE.md` with skill routing and frontmatter conventions.
- `resources.md` from the previous office-hours session.
- **Zero Go code.** The `Startr/feeds` GitHub repo does not exist yet. The domain `startr.media` is owned; DNS is being moved into the `dawnofthegeeks@openco.ca` Cloudflare account. The Spotify for Podcasters auto-feed does not yet exist because no episodes have been recorded. This is a greenfield build against an approved design.

### Open TODOs (surfaced during review)

1. **Test coverage commitment scope.** ~~Ship all 34 tests with v0.1.0, or ship the 5 critical ones and backfill the rest in v0.1.1?~~ **RESOLVED 2026-04-07.** Decision: ship v0.1.0 with the 5 critical tests only. Backfill the remaining ~29 in v0.1.1 alongside the play widget. Risks documented in the test requirements section; tradeoff accepted in favor of ship velocity.

No other open TODOs. All 12 architecture decisions, rewrite scope, tests, performance characterization, failure modes, and parallelization strategy are locked.

### Review log entry

Appended to `~/.gstack/projects/workspace/develop-reviews.jsonl`:

```json
{"ts":"2026-04-07T15:42:00Z","skill":"plan-eng-review","branch":"develop","artifact":"agent-develop-design-20260406-174520.md","outcome":"DONE_WITH_CONCERNS","decisions_locked":12,"critical_tests":5,"total_tests":34,"failure_modes":13,"worktrees":6,"open_todos":1}
```

### What to do next

1. Create the `Startr/feeds` GitHub repo with AGPL-3.0, `TRADEMARK.md`, `DCO.txt`, and a README stub with the embed-vs-fork explainer.
2. Run the 6 worktrees in parallel. Start with C (rewriter + 3 critical tests) and D (output renderer) because they contain the highest-leverage logic and load-bearing tests.
3. When all 6 worktrees merge, tag v0.1.0 and let GoReleaser ship the full distribution pipeline.
4. Route the output XML to Cloudflare Pages or Backblaze B2 at `feed.startr.media/v1/[show-slug].xml`.
5. Submit to Apple Podcasts Connect and Spotify for Podcasters.
6. Record Episode 1 (`Run It Local`, pilot).
7. Ship v0.1.1 (play widget + ~29 backfill tests) and v0.1.2 (subscribe widget) as fast follows.

The `/qa` skill will consume the test plan artifact once there's code to test. Note that v0.1.0 will only run the 5 critical tests against the codebase; the rest land in the v0.1.1 backfill sprint.
