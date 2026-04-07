/// <reference path="../pb_data/types.d.ts" />

// Creates the "feeds" collection for storing feed configurations.
// Each record is one feed to rewrite. The cron hook iterates over all records.
//
// Field names are chosen for clarity in the PocketBase admin UI — a podcaster
// who has never seen a terminal should be able to fill these in.

migrate(
  // up
  function(app) {
    const collection = new Collection({
      name: "feeds",
      type: "base",
      fields: [
        // ── What to rewrite ──────────────────────────────────────
        {
          name: "source_url",
          type: "url",
          required: true,
          presentable: true,
        },
        // ── Where to publish ─────────────────────────────────────
        // The slug is the feed's unique identifier. It determines both:
        //   • output file:  pb_public/v1/{slug}.xml
        //   • public URL:   https://{domain}/v1/{slug}.xml
        {
          name: "slug",
          type: "text",
          required: true,
          presentable: true,
          min: 1,
          max: 128,
          pattern: "^[a-zA-Z0-9][a-zA-Z0-9_-]*$",
        },
        // Optional domain override. If blank, uses FEEDS_DOMAIN env var
        // or the PocketBase application URL from admin settings.
        // e.g. https://feed.example.com
        {
          name: "domain",
          type: "url",
          required: false,
        },
        // ── Branding ─────────────────────────────────────────────
        {
          name: "title",
          type: "text",
          required: true,
          presentable: true,
        },
        {
          name: "website",
          type: "url",
          required: true,
        },
        {
          name: "cover_image",
          type: "url",
          required: false,
        },
        {
          name: "author",
          type: "text",
          required: false,
        },
        {
          name: "owner_email",
          type: "email",
          required: false,
        },
        // ── Schedule ─────────────────────────────────────────────
        {
          name: "schedule",
          type: "text",
          required: false,
        },
        // ── Internal cache state (hidden from API, visible in admin) ──
        {
          name: "etag",
          type: "text",
          required: false,
          hidden: true,
        },
        {
          name: "last_modified",
          type: "text",
          required: false,
          hidden: true,
        },
      ],
    });

    app.save(collection);
  },

  // down
  function(app) {
    const collection = app.findCollectionByNameOrId("feeds");
    app.delete(collection);
  }
);
