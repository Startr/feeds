/// <reference path="../pb_data/types.d.ts" />

// =============================================================================
// Startr/feeds — PocketBase JS hook
//
// PB 0.36 runs route/cron/bootstrap callbacks in isolated VMs where file-scope
// declarations are not visible. All pipeline logic lives in lib/pipeline.js,
// loaded via require() inside each callback. This file only registers hooks.
// =============================================================================

// ---------------------------------------------------------------------------
// CORS-enabled feed serving — lets <startr-player> fetch XML cross-origin.
// Self-contained (only uses $os.readFile), no pipeline dependency.
// ---------------------------------------------------------------------------

routerAdd("GET", "/v1/{slug}", function(e) {
  var slug = e.request.pathValue("slug");
  if (!slug || slug.indexOf(".xml") !== slug.length - 4) {
    return e.notFound();
  }
  slug = slug.slice(0, -4);
  if (!slug || slug.indexOf("..") !== -1 || slug.indexOf("/") !== -1) {
    return e.json(400, { message: "invalid slug" });
  }
  var xmlPath = "./pb_public/v1/" + slug + ".xml";
  var content;
  try {
    content = $os.readFile(xmlPath);
  } catch(err) {
    return e.json(404, { message: "feed not found" });
  }
  e.response.header().set("Access-Control-Allow-Origin", "*");
  e.response.header().set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS");
  e.response.header().set("Cache-Control", "public, max-age=300");
  return e.blob(200, "application/rss+xml; charset=utf-8", content);
});

routerAdd("OPTIONS", "/v1/{slug}", function(e) {
  var slug = e.request.pathValue("slug");
  if (!slug || slug.indexOf(".xml") !== slug.length - 4) {
    return e.notFound();
  }
  e.response.header().set("Access-Control-Allow-Origin", "*");
  e.response.header().set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS");
  e.response.header().set("Access-Control-Allow-Headers", "Content-Type");
  e.response.header().set("Access-Control-Max-Age", "86400");
  return e.noContent(204);
});

// ---------------------------------------------------------------------------
// Public feed list — no auth needed (feeds are already public static files)
// ---------------------------------------------------------------------------
routerAdd("GET", "/api/feeds/public", function(e) {
  var feeds = [];
  try {
    var records = $app.findAllRecords("feeds");
    for (var i = 0; i < records.length; i++) {
      feeds.push({ slug: records[i].getString("slug"), title: records[i].getString("title") });
    }
  } catch(err) {
    console.log("[feeds] collection query failed: " + err);
  }
  if (feeds.length === 0) {
    var slug = $os.getenv("FEEDS_SLUG");
    if (slug) feeds.push({ slug: slug, title: $os.getenv("FEEDS_TITLE") || slug });
  }
  return e.json(200, { feeds: feeds });
});

// ---------------------------------------------------------------------------
// Debug route — trigger pipeline manually (superuser auth required)
// ---------------------------------------------------------------------------
routerAdd("GET", "/debug/run", function(e) {
  if (!e.hasSuperuserAuth()) {
    return e.json(401, { message: "superuser auth required" });
  }
  try {
    var pipeline = require(__hooks + "/lib/pipeline.js");
    pipeline.runAllFeeds();
    return e.json(200, { status: "ok", version: pipeline.VERSION });
  } catch(err) {
    return e.json(500, { error: String(err) });
  }
});

// ---------------------------------------------------------------------------
// Cron — rewrite all feeds on schedule
// ---------------------------------------------------------------------------
var cronExpr = $os.getenv("FEEDS_CRON") || "*/15 * * * *";
cronAdd("feeds_rewrite", cronExpr, function() {
  try {
    var pipeline = require(__hooks + "/lib/pipeline.js");
    pipeline.runAllFeeds();
  } catch(err) {
    console.log("[feeds] cron rewrite failed: " + err);
  }
});

// ---------------------------------------------------------------------------
// Bootstrap — run once on server start so deploys get a fresh feed
// ---------------------------------------------------------------------------
onBootstrap(function(e) {
  e.next();
  try {
    var pipeline = require(__hooks + "/lib/pipeline.js");
    pipeline.runAllFeeds();
  } catch(err) {
    console.log("[feeds] initial run failed: " + err);
  }
});
