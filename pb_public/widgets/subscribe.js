/**
 * <startr-subscribe> — embeddable podcast subscribe widget
 *
 * Usage:
 *   <script src="https://feed.startr.media/widgets/subscribe.js"></script>
 *   <startr-subscribe
 *     apple="https://podcasts.apple.com/podcast/id123"
 *     spotify="https://open.spotify.com/show/abc"
 *     rss="https://feed.startr.media/v1/show.xml">
 *   </startr-subscribe>
 *
 * Attributes:
 *   apple   — Apple Podcasts URL (optional)
 *   spotify — Spotify URL (optional)
 *   rss     — RSS feed URL (optional)
 *   accent  — Accent color hex (optional, also via --startr-accent CSS property)
 *
 * Only platforms with a URL attribute render. If all are omitted, nothing renders.
 *
 * AGPL-3.0 — https://github.com/Startr/feeds
 * Icons: canonical source at /widgets/icons.svg
 */
(function() {
  'use strict';

  // --- Icon SVG paths (source of truth: /widgets/icons.svg) ---
  var ICONS = {
    apple: '<path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>',
    spotify: '<path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z"/>',
    rss: '<circle cx="6.18" cy="17.82" r="2.18"/><path d="M4 4.44v2.83c7.03 0 12.73 5.7 12.73 12.73h2.83c0-8.59-6.97-15.56-15.56-15.56zm0 5.66v2.83c3.9 0 7.07 3.17 7.07 7.07h2.83c0-5.47-4.43-9.9-9.9-9.9z"/>'
  };

  var PLATFORMS = [
    { attr: 'apple',   label: 'Apple Podcasts', icon: 'apple',   viewBox: '0 0 24 24' },
    { attr: 'spotify', label: 'Spotify',        icon: 'spotify', viewBox: '0 0 24 24' },
    { attr: 'rss',     label: 'RSS Feed',       icon: 'rss',     viewBox: '0 0 24 24' }
  ];

  function svg(id, size) {
    return '<svg xmlns="http://www.w3.org/2000/svg" width="' + size + '" height="' + size +
      '" viewBox="' + (PLATFORMS.find(function(p) { return p.icon === id; }) || { viewBox: '0 0 24 24' }).viewBox +
      '" fill="currentColor" aria-hidden="true">' + (ICONS[id] || '') + '</svg>';
  }

  var STYLE = [
    ':host { display: inline-block; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }',
    '.badges { display: flex; flex-wrap: wrap; gap: 8px; }',
    'a { display: inline-flex; align-items: center; gap: 6px; padding: 8px 14px;',
    '  border-radius: 8px; text-decoration: none; font-size: 14px; font-weight: 600;',
    '  color: #fff; background: var(--startr-accent, #2563eb); transition: opacity 0.15s; }',
    'a:hover { opacity: 0.85; }',
    'a svg { flex-shrink: 0; }'
  ].join('\n');

  class StartrSubscribe extends HTMLElement {
    static get observedAttributes() { return ['apple', 'spotify', 'rss', 'accent']; }

    constructor() {
      super();
      this._shadow = this.attachShadow({ mode: 'closed' });
    }

    connectedCallback() { this._render(); }

    attributeChangedCallback() {
      if (this.isConnected) this._render();
    }

    _render() {
      var accent = this.getAttribute('accent');
      if (accent) {
        this.style.setProperty('--startr-accent', accent);
      }

      var html = '<style>' + STYLE + '</style><div class="badges">';
      var hasAny = false;

      for (var i = 0; i < PLATFORMS.length; i++) {
        var p = PLATFORMS[i];
        var url = this.getAttribute(p.attr);
        if (url) {
          hasAny = true;
          html += '<a href="' + _escAttr(url) + '" target="_blank" rel="noopener" aria-label="' +
            p.label + '">' + svg(p.icon, 16) + ' ' + p.label + '</a>';
        }
      }

      html += '</div>';
      this._shadow.innerHTML = hasAny ? html : '';
    }
  }

  function _escAttr(s) {
    return s.replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  customElements.define('startr-subscribe', StartrSubscribe);
})();
