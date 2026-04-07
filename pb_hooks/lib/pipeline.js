// =============================================================================
// Startr/feeds — rewrite pipeline module
//
// Loaded via require() inside PB 0.36 route/cron/bootstrap callbacks.
// PB 0.36 runs callbacks in isolated VMs where file-scope declarations from
// the parent hook file are not visible. This module is self-contained:
// it loads xml-js, defines all pipeline functions, and exports them via `this`.
// =============================================================================

// ---------------------------------------------------------------------------
// Load xml-js via the global object trick
// PB 0.36's goja VMs don't have `window`, but Function('return this')()
// returns the real global. Setting window on it lets the UMD bundle load.
// ---------------------------------------------------------------------------
var _global = Function('return this')();
_global.window = _global.window || {};
require(__hooks + "/lib/xml-js.js");

var xml2js = _global.window.xml2js;
var js2xml = _global.window.js2xml;
var VERSION = $os.getenv("FEEDS_VERSION") || "dev";

// ---------------------------------------------------------------------------
// XML helpers
// ---------------------------------------------------------------------------

function findElement(parent, name) {
  if (!parent.elements) return null;
  for (var i = 0; i < parent.elements.length; i++) {
    var el = parent.elements[i];
    if (el.type === "element" && el.name === name) return el;
  }
  return null;
}

function setChildText(parent, name, text) {
  if (text === null || text === undefined) return;
  if (!parent.elements) parent.elements = [];
  var el = findElement(parent, name);
  if (el) {
    el.elements = [{ type: "text", text: text }];
  } else {
    parent.elements.push({
      type: "element",
      name: name,
      elements: [{ type: "text", text: text }],
    });
  }
}

function setOrCreateAttrElement(parent, name, attr, value) {
  if (!parent.elements) parent.elements = [];
  var el = findElement(parent, name);
  if (el) {
    if (!el.attributes) el.attributes = {};
    el.attributes[attr] = value;
  } else {
    var attrs = {};
    attrs[attr] = value;
    parent.elements.push({
      type: "element",
      name: name,
      attributes: attrs,
    });
  }
}

// ---------------------------------------------------------------------------
// Rewrite pipeline for a single feed
// ---------------------------------------------------------------------------

function rewriteFeed(feed) {
  var headers = {};
  if (feed.etag) headers["If-None-Match"] = feed.etag;
  if (feed.lastModified) headers["If-Modified-Since"] = feed.lastModified;

  var res = $http.send({
    url:     feed.upstream,
    method:  "GET",
    headers: headers,
    timeout: 30,
  });

  if (res.statusCode === 304) {
    console.log("[feeds] " + feed.upstream + " -> 304 Not Modified, skipping");
    return { etag: feed.etag, lastModified: feed.lastModified, changed: false };
  }

  if (res.statusCode !== 200) {
    throw new Error("upstream returned HTTP " + res.statusCode);
  }

  var doc = xml2js(res.raw, { compact: false });

  var rss = findElement(doc, "rss");
  if (!rss) throw new Error("no <rss> root element found");

  var channel = findElement(rss, "channel");
  if (!channel) throw new Error("no <channel> element found");

  setChildText(channel, "title", feed.title);
  setChildText(channel, "link", feed.link);
  setChildText(channel, "generator", "Startr/feeds " + VERSION + " (https://github.com/Startr/feeds)");

  if (feed.image) {
    setChildText(channel, "itunes:image", null);
    setOrCreateAttrElement(channel, "itunes:image", "href", feed.image);
    var imageEl = findElement(channel, "image");
    if (imageEl) setChildText(imageEl, "url", feed.image);
  }

  if (feed.itunesAuthor) {
    setChildText(channel, "itunes:author", feed.itunesAuthor);
  }

  if (feed.itunesOwnerEmail) {
    var owner = findElement(channel, "itunes:owner");
    if (!owner) {
      owner = { type: "element", name: "itunes:owner", elements: [] };
      channel.elements.push(owner);
    }
    setChildText(owner, "itunes:email", feed.itunesOwnerEmail);
  }

  if (feed.selfUrl) {
    var atomLink = null;
    if (channel.elements) {
      for (var ai = 0; ai < channel.elements.length; ai++) {
        var el = channel.elements[ai];
        if (el.type === "element" &&
            el.name === "atom:link" &&
            el.attributes &&
            el.attributes.rel === "self") {
          atomLink = el;
          break;
        }
      }
    }
    if (atomLink) {
      atomLink.attributes.href = feed.selfUrl;
    } else {
      if (!rss.attributes) rss.attributes = {};
      if (!rss.attributes["xmlns:atom"]) {
        rss.attributes["xmlns:atom"] = "http://www.w3.org/2005/Atom";
      }
      channel.elements.splice(0, 0, {
        type: "element",
        name: "atom:link",
        attributes: {
          href: feed.selfUrl,
          rel: "self",
          type: "application/rss+xml",
        },
      });
    }
  }

  var xml = js2xml(doc, { compact: false, spaces: 2 });

  var parts = feed.outputPath.split("/");
  parts.pop();
  var dir = parts.join("/");
  if (dir) {
    try { $os.mkdirAll(dir, 0o755); } catch(e) { /* may already exist */ }
  }

  $os.writeFile(feed.outputPath, xml, 0o644);
  console.log("[feeds] wrote " + xml.length + " bytes -> " + feed.outputPath);

  var newEtag = res.headers && res.headers["Etag"] ? res.headers["Etag"][0] : (res.headers && res.headers["etag"] ? res.headers["etag"][0] : "");
  var newLastMod = res.headers && res.headers["Last-Modified"] ? res.headers["Last-Modified"][0] : (res.headers && res.headers["last-modified"] ? res.headers["last-modified"][0] : "");

  return { etag: newEtag, lastModified: newLastMod, changed: true };
}

// ---------------------------------------------------------------------------
// Feed config: collection records with env-var fallback
// ---------------------------------------------------------------------------

function resolveDomain(perFeedDomain) {
  if (perFeedDomain) return perFeedDomain.replace(/\/+$/, "");
  var envDomain = $os.getenv("FEEDS_DOMAIN");
  if (envDomain) return envDomain.replace(/\/+$/, "");
  try {
    var appUrl = $app.settings().meta.appURL;
    if (appUrl) return appUrl.replace(/\/+$/, "");
  } catch(e) {}
  return "";
}

function slugToConfig(slug, domain) {
  var outputPath = "./pb_public/v1/" + slug + ".xml";
  var selfUrl = domain ? domain + "/v1/" + slug + ".xml" : "";
  return { outputPath: outputPath, selfUrl: selfUrl };
}

function loadFeedConfigs() {
  var configs = [];

  try {
    var records = $app.findAllRecords("feeds");
    if (records && records.length > 0) {
      for (var i = 0; i < records.length; i++) {
        var r = records[i];
        var slug = r.getString("slug");
        var domain = resolveDomain(r.getString("domain"));
        var derived = slugToConfig(slug, domain);
        configs.push({
          id:               r.id,
          upstream:         r.getString("source_url"),
          outputPath:       derived.outputPath,
          selfUrl:          derived.selfUrl,
          title:            r.getString("title"),
          link:             r.getString("website"),
          image:            r.getString("cover_image"),
          itunesAuthor:     r.getString("author"),
          itunesOwnerEmail: r.getString("owner_email"),
          etag:             r.getString("etag"),
          lastModified:     r.getString("last_modified"),
        });
      }
      return configs;
    }
  } catch(e) {
    console.log("[feeds] collection lookup: " + e + " -- trying env vars");
  }

  var upstream = $os.getenv("FEEDS_SOURCE_URL") || $os.getenv("FEEDS_UPSTREAM");
  if (!upstream) return configs;

  var slug = $os.getenv("FEEDS_SLUG");
  if (slug) {
    var domain = resolveDomain("");
    var derived = slugToConfig(slug, domain);
    configs.push({
      id:               "__env__",
      upstream:         upstream,
      outputPath:       derived.outputPath,
      selfUrl:          derived.selfUrl,
      title:            $os.getenv("FEEDS_TITLE")       || $os.getenv("FEEDS_CHANNEL_TITLE") || "",
      link:             $os.getenv("FEEDS_WEBSITE")     || $os.getenv("FEEDS_CHANNEL_LINK") || "",
      image:            $os.getenv("FEEDS_COVER_IMAGE") || $os.getenv("FEEDS_CHANNEL_IMAGE") || "",
      itunesAuthor:     $os.getenv("FEEDS_ITUNES_AUTHOR") || "",
      itunesOwnerEmail: $os.getenv("FEEDS_ITUNES_OWNER_EMAIL") || "",
      etag:             "",
      lastModified:     "",
    });
  } else {
    configs.push({
      id:               "__env__",
      upstream:         upstream,
      outputPath:       $os.getenv("FEEDS_OUTPUT") || "./pb_public/feed.xml",
      selfUrl:          $os.getenv("FEEDS_SELF_URL") || "",
      title:            $os.getenv("FEEDS_TITLE")       || $os.getenv("FEEDS_CHANNEL_TITLE") || "",
      link:             $os.getenv("FEEDS_WEBSITE")     || $os.getenv("FEEDS_CHANNEL_LINK") || "",
      image:            $os.getenv("FEEDS_COVER_IMAGE") || $os.getenv("FEEDS_CHANNEL_IMAGE") || "",
      itunesAuthor:     $os.getenv("FEEDS_ITUNES_AUTHOR") || "",
      itunesOwnerEmail: $os.getenv("FEEDS_ITUNES_OWNER_EMAIL") || "",
      etag:             "",
      lastModified:     "",
    });
  }

  return configs;
}

function saveCacheState(feedConfig, state) {
  if (feedConfig.id === "__env__") return;
  try {
    var record = $app.findRecordById("feeds", feedConfig.id);
    record.set("etag", state.etag);
    record.set("last_modified", state.lastModified);
    $app.save(record);
  } catch(e) {
    console.log("[feeds] cache save failed: " + e);
  }
}

// ---------------------------------------------------------------------------
// Main runner
// ---------------------------------------------------------------------------

function runAllFeeds() {
  var configs = loadFeedConfigs();
  if (configs.length === 0) {
    console.log("[feeds] no feeds configured -- set FEEDS_SOURCE_URL or add a record in the feeds collection via /_/");
    return;
  }
  for (var i = 0; i < configs.length; i++) {
    var feed = configs[i];
    try {
      var state = rewriteFeed(feed);
      if (state.changed) saveCacheState(feed, state);
    } catch(e) {
      console.log("[feeds] rewrite failed for " + feed.upstream + ": " + e);
    }
  }
}

// ---------------------------------------------------------------------------
// Exports — PB's require() returns `this`, so assign to it
// ---------------------------------------------------------------------------
this.runAllFeeds = runAllFeeds;
this.rewriteFeed = rewriteFeed;
this.loadFeedConfigs = loadFeedConfigs;
this.xml2js = xml2js;
this.js2xml = js2xml;
this.VERSION = VERSION;
