---
title: "pb_hooks/lib — vendored JS libraries"
purpose: document how third-party JS libraries are managed for PocketBase hooks
---

# pb_hooks/lib — vendored JS libraries

PocketBase's JS runtime (goja) can't use npm at runtime — there's no `node_modules/` resolution only CommonJS `require()` from local files.

## How libraries get here

We vendor pre-built dist files directly into `pb_hooks/lib/` using `npm pack`:

```bash
cd /tmp
npm pack xml-js@1.6.11
tar xzf xml-js-1.6.11.tgz
cp package/dist/xml-js.js /path/to/repo/pb_hooks/lib/xml-js.js
rm -rf package xml-js-*.tgz
```

This pulls the published tarball from npm, extracts the UMD/CJS dist bundle, and commits it to the repo. This avoids `node_modules/` and package.json. The Dockerfile just COPYs `pb_hooks/` into the container and it works.

### Why not npm install?

The Docker image has no Node.js — it's just the PocketBase binary + Alpine. Adding Node.js to install one 250KB file would bloat the image and slow the build for no reason. Vendoring the dist file keeps the build fast (seconds instead of minutes) and the image minimal.

## Current libraries

| Library | Version | File | Purpose |
|---------|---------|------|---------|
| [xml-js](https://github.com/nicknisi/xml-js) | 1.6.11 | `xml-js.js` | XML ↔ JS object conversion for feed rewriting |

## Usage in goja

xml-js is a UMD bundle that attaches to `window`. In PocketBase's goja runtime:

```javascript
this.window = {};
require(`${__hooks}/lib/xml-js.js`);
const xml2js = this.window.xml2js;  // XML string → JS object
const js2xml = this.window.js2xml;  // JS object → XML string
```

## Adding a new library

1. Find the npm package and its dist/bundle file
2. `npm pack <package>@<version>` + extract the dist file
3. Copy to `pb_hooks/lib/`
4. Update the table above
5. Commit the file — it's vendored, not gitignored
