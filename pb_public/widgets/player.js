/**
 * <startr-player> — embeddable podcast player widget
 *
 * Usage:
 *   <script src="https://feed.startr.media/widgets/player.js"></script>
 *   <startr-player feed="https://feed.startr.media/v1/show.xml"></startr-player>
 *
 * Attributes:
 *   feed    — RSS feed URL to fetch and parse (required)
 *   episode — 0-based episode index (optional, default: 0 = latest)
 *   accent  — Accent color hex (optional, also via --startr-accent CSS property)
 *
 * Features: episode picker, progress memory (localStorage), keyboard shortcuts
 * (space=play/pause, arrows=skip +-15s), share button (Web Share API + clipboard).
 *
 * AGPL-3.0 — https://github.com/Startr/feeds
 * Icons: canonical source at /widgets/icons.svg
 */
(function() {
  'use strict';

  // --- Icon SVG paths (source of truth: /widgets/icons.svg) ---
  var ICONS = {
    play:         '<path d="M8 5.14v14l11-7z"/>',
    pause:        '<path d="M6 5h4v14H6zm8 0h4v14h-4z"/>',
    'skip-back':  '<path d="M11 18V6l-8.5 6zm.5-6l8.5 6V6z"/>',
    'skip-fwd':   '<path d="M4 18l8.5-6L4 6zm9-12v12l8.5-6z"/>',
    share:        '<path d="M4 12v8a2 2 0 002 2h12a2 2 0 002-2v-8" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><polyline points="16 6 12 2 8 6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><line x1="12" y1="2" x2="12" y2="15" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>'
  };

  function icon(id, size) {
    return '<svg xmlns="http://www.w3.org/2000/svg" width="' + size + '" height="' + size +
      '" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">' + (ICONS[id] || '') + '</svg>';
  }

  // --- Styles ---
  var STYLE = [
    ':host { display: block; max-width: 400px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }',
    '* { box-sizing: border-box; margin: 0; padding: 0; }',
    '.card { border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden; background: #fff; }',
    '',
    '/* Header: artwork + metadata */',
    '.header { display: flex; gap: 12px; padding: 16px; }',
    '.artwork { width: 80px; height: 80px; border-radius: 8px; object-fit: cover; background: #f3f4f6; flex-shrink: 0; }',
    '.meta { display: flex; flex-direction: column; justify-content: center; min-width: 0; }',
    '.show-title { font-size: 12px; color: #666; text-transform: uppercase; letter-spacing: 0.03em; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }',
    '.ep-title { font-size: 15px; font-weight: 600; color: #1a1a1a; margin-top: 2px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }',
    '.ep-date { font-size: 12px; color: #999; margin-top: 4px; }',
    '',
    '/* Controls */',
    '.controls { display: flex; align-items: center; gap: 8px; padding: 0 16px 8px; }',
    '.btn { background: none; border: none; cursor: pointer; color: #1a1a1a; padding: 6px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }',
    '.btn:hover { background: #f3f4f6; }',
    '.btn-play { width: 44px; height: 44px; background: var(--startr-accent, #2563eb); color: #fff; border-radius: 50%; }',
    '.btn-play:hover { opacity: 0.85; background: var(--startr-accent, #2563eb); }',
    '.time { font-size: 12px; color: #666; margin-left: auto; white-space: nowrap; }',
    '',
    '/* Scrubber */',
    '.scrubber-wrap { padding: 0 16px 12px; }',
    '.scrubber { -webkit-appearance: none; appearance: none; width: 100%; height: 4px; border-radius: 2px; background: #e5e7eb; outline: none; cursor: pointer; }',
    '.scrubber::-webkit-slider-thumb { -webkit-appearance: none; width: 14px; height: 14px; border-radius: 50%; background: var(--startr-accent, #2563eb); cursor: pointer; }',
    '.scrubber::-moz-range-thumb { width: 14px; height: 14px; border-radius: 50%; background: var(--startr-accent, #2563eb); border: none; cursor: pointer; }',
    '',
    '/* Footer: episode picker + share */',
    '.footer { display: flex; align-items: center; gap: 8px; padding: 8px 16px 12px; border-top: 1px solid #e5e7eb; }',
    '.picker { flex: 1; min-width: 0; font-size: 13px; padding: 6px 8px; border: 1px solid #e5e7eb; border-radius: 6px; background: #fff; color: #1a1a1a; font-family: inherit; cursor: pointer; }',
    '.picker:focus { outline: 2px solid var(--startr-accent, #2563eb); outline-offset: -1px; }',
    '.btn-share { font-size: 12px; display: flex; align-items: center; gap: 4px; padding: 6px 10px; border: 1px solid #e5e7eb; border-radius: 6px; background: #fff; color: #1a1a1a; cursor: pointer; font-family: inherit; }',
    '.btn-share:hover { background: #f3f4f6; }',
    '.tooltip { font-size: 11px; color: var(--startr-accent, #2563eb); font-weight: 600; }',
    '',
    '/* Loading + error states */',
    '.loading, .error { padding: 24px 16px; text-align: center; font-size: 14px; color: #666; }',
    '.error { color: #dc2626; }'
  ].join('\n');

  // --- Component ---
  class StartrPlayer extends HTMLElement {
    static get observedAttributes() { return ['feed', 'episode', 'accent']; }

    constructor() {
      super();
      this._shadow = this.attachShadow({ mode: 'closed' });
      this._audio = null;
      this._episodes = [];
      this._showTitle = '';
      this._showImage = '';
      this._currentIndex = 0;
      this._progressTimer = 0;
      this._boundKeyHandler = this._handleKeyboard.bind(this);
    }

    connectedCallback() {
      var accent = this.getAttribute('accent');
      if (accent) this.style.setProperty('--startr-accent', accent);
      this._fetchFeed();
    }

    disconnectedCallback() {
      this._saveProgress();
      if (this._audio) {
        this._audio.pause();
        this._audio.src = '';
      }
    }

    attributeChangedCallback(name) {
      if (!this.isConnected) return;
      if (name === 'feed') this._fetchFeed();
      if (name === 'accent') {
        var a = this.getAttribute('accent');
        if (a) this.style.setProperty('--startr-accent', a);
      }
    }

    // --- Feed ---

    _fetchFeed() {
      var feedUrl = this.getAttribute('feed');
      if (!feedUrl) {
        this._renderEmpty();
        return;
      }
      this._feedUrl = feedUrl;

      this._shadow.innerHTML = '<style>' + STYLE + '</style><div class="card"><div class="loading">Loading feed\u2026</div></div>';

      var self = this;
      fetch(feedUrl)
        .then(function(r) {
          if (!r.ok) throw new Error('HTTP ' + r.status);
          return r.text();
        })
        .then(function(xml) {
          self._parseFeed(xml);
          var startIndex = parseInt(self.getAttribute('episode'), 10) || 0;
          self._currentIndex = Math.min(startIndex, self._episodes.length - 1);
          self._render();
          self._loadEpisode(self._currentIndex);
        })
        .catch(function(err) {
          self._shadow.innerHTML = '<style>' + STYLE + '</style><div class="card"><div class="error">Could not load feed: ' + _esc(err.message) + '</div></div>';
        });
    }

    _parseFeed(xmlText) {
      var doc = new DOMParser().parseFromString(xmlText, 'application/xml');
      var channel = doc.querySelector('channel');
      if (!channel) throw new Error('No <channel> in feed');

      this._showTitle = _text(channel, 'title');
      this._showImage = _attr(channel, 'itunes\\:image', 'href') || _text(channel, 'image > url') || '';

      var items = channel.querySelectorAll('item');
      var limit = Math.min(items.length, 10);
      this._episodes = [];
      for (var i = 0; i < limit; i++) {
        var it = items[i];
        this._episodes.push({
          title:   _text(it, 'title'),
          guid:    _text(it, 'guid') || _text(it, 'title'),
          audio:   _attr(it, 'enclosure', 'url') || '',
          image:   _attr(it, 'itunes\\:image', 'href') || this._showImage,
          pubDate: _text(it, 'pubDate')
        });
      }
    }

    // --- Rendering ---

    _render() {
      var ep = this._episodes[this._currentIndex];
      if (!ep) return;

      var date = ep.pubDate ? _formatDate(ep.pubDate) : '';

      var html = '<style>' + STYLE + '</style>';
      html += '<div class="card" tabindex="0">';

      // Header
      html += '<div class="header">';
      html += '<img class="artwork" src="' + _escAttr(ep.image) + '" alt="" loading="lazy">';
      html += '<div class="meta">';
      html += '<div class="show-title">' + _esc(this._showTitle) + '</div>';
      html += '<div class="ep-title">' + _esc(ep.title) + '</div>';
      if (date) html += '<div class="ep-date">' + _esc(date) + '</div>';
      html += '</div></div>';

      // Controls
      html += '<div class="controls">';
      html += '<button class="btn btn-skip" data-action="back" aria-label="Back 15 seconds">' + icon('skip-back', 20) + '</button>';
      html += '<button class="btn btn-play" data-action="play" aria-label="Play">' + icon('play', 22) + '</button>';
      html += '<button class="btn btn-skip" data-action="fwd" aria-label="Forward 15 seconds">' + icon('skip-fwd', 20) + '</button>';
      html += '<span class="time"><span class="time-current">0:00</span> / <span class="time-total">0:00</span></span>';
      html += '</div>';

      // Scrubber
      html += '<div class="scrubber-wrap">';
      html += '<input type="range" class="scrubber" min="0" max="100" value="0" step="0.1" aria-label="Seek">';
      html += '</div>';

      // Footer: episode picker + share
      html += '<div class="footer">';
      html += '<select class="picker" aria-label="Choose episode">';
      for (var i = 0; i < this._episodes.length; i++) {
        var e = this._episodes[i];
        var sel = i === this._currentIndex ? ' selected' : '';
        html += '<option value="' + i + '"' + sel + '>' + _esc(e.title) + '</option>';
      }
      html += '</select>';
      html += '<button class="btn-share" data-action="share">' + icon('share', 14) + ' Share</button>';
      html += '</div>';

      html += '</div>';
      this._shadow.innerHTML = html;
      this._bindEvents();
    }

    _renderEmpty() {
      var html = '<style>' + STYLE + '</style>';
      html += '<div class="card">';
      html += '<div class="header">';
      html += '<div class="artwork"></div>';
      html += '<div class="meta">';
      html += '<div class="show-title">Your Show</div>';
      html += '<div class="ep-title">Add a feed to start playing</div>';
      html += '<div class="ep-date">startr-player</div>';
      html += '</div></div>';
      html += '<div class="controls">';
      html += '<button class="btn btn-skip" disabled>' + icon('skip-back', 20) + '</button>';
      html += '<button class="btn btn-play" disabled>' + icon('play', 22) + '</button>';
      html += '<button class="btn btn-skip" disabled>' + icon('skip-fwd', 20) + '</button>';
      html += '<span class="time">0:00 / 0:00</span>';
      html += '</div>';
      html += '<div class="scrubber-wrap">';
      html += '<input type="range" class="scrubber" min="0" max="100" value="0" disabled>';
      html += '</div>';
      html += '<div class="footer">';
      html += '<select class="picker" disabled><option>No episodes yet</option></select>';
      html += '<button class="btn-share" disabled>' + icon('share', 14) + ' Share</button>';
      html += '</div>';
      html += '</div>';
      this._shadow.innerHTML = html;
    }

    _bindEvents() {
      var self = this;
      var card = this._shadow.querySelector('.card');

      // Button clicks
      card.addEventListener('click', function(e) {
        var btn = e.target.closest('[data-action]');
        if (!btn) return;
        var action = btn.getAttribute('data-action');
        if (action === 'play')  self._togglePlay();
        if (action === 'back')  self._seek(-15);
        if (action === 'fwd')   self._seek(15);
        if (action === 'share') self._share();
      });

      // Episode picker
      var picker = this._shadow.querySelector('.picker');
      picker.addEventListener('change', function() {
        self._selectEpisode(parseInt(picker.value, 10));
      });

      // Scrubber
      var scrubber = this._shadow.querySelector('.scrubber');
      var scrubbing = false;
      scrubber.addEventListener('input', function() {
        scrubbing = true;
        if (self._audio && self._audio.duration) {
          self._audio.currentTime = (scrubber.value / 100) * self._audio.duration;
        }
      });
      scrubber.addEventListener('change', function() { scrubbing = false; });

      // Audio events
      if (this._audio) {
        this._audio.addEventListener('timeupdate', function() {
          if (scrubbing) return;
          self._updateTime();
          // Throttled progress save (every 5s)
          var now = Date.now();
          if (now - self._progressTimer > 5000) {
            self._progressTimer = now;
            self._saveProgress();
          }
        });
        this._audio.addEventListener('pause', function() { self._saveProgress(); self._updatePlayBtn(); });
        this._audio.addEventListener('play', function() { self._updatePlayBtn(); });
        this._audio.addEventListener('ended', function() {
          self._clearProgress();
          self._updatePlayBtn();
        });
        this._audio.addEventListener('loadedmetadata', function() { self._updateTime(); });
      }

      // Keyboard (focus-scoped)
      card.addEventListener('keydown', this._boundKeyHandler);
    }

    // --- Playback ---

    _loadEpisode(index) {
      var ep = this._episodes[index];
      if (!ep || !ep.audio) return;

      if (!this._audio) {
        this._audio = new Audio();
        this._audio.preload = 'metadata';
      } else {
        this._audio.pause();
      }

      this._audio.src = ep.audio;
      this._currentIndex = index;

      // Re-bind audio events (new episode)
      this._render();
      this._restoreProgress();
    }

    _togglePlay() {
      if (!this._audio) return;
      if (this._audio.paused) {
        this._audio.play().catch(function() {});
      } else {
        this._audio.pause();
      }
    }

    _seek(seconds) {
      if (!this._audio) return;
      this._audio.currentTime = Math.max(0, Math.min(this._audio.duration || 0, this._audio.currentTime + seconds));
    }

    _selectEpisode(index) {
      if (index === this._currentIndex) return;
      this._saveProgress();
      this._currentIndex = index;
      this._loadEpisode(index);
    }

    _updatePlayBtn() {
      var btn = this._shadow.querySelector('.btn-play');
      if (!btn || !this._audio) return;
      btn.innerHTML = this._audio.paused ? icon('play', 22) : icon('pause', 22);
      btn.setAttribute('aria-label', this._audio.paused ? 'Play' : 'Pause');
    }

    _updateTime() {
      if (!this._audio) return;
      var cur = this._shadow.querySelector('.time-current');
      var tot = this._shadow.querySelector('.time-total');
      var scrubber = this._shadow.querySelector('.scrubber');
      if (cur) cur.textContent = _formatTime(this._audio.currentTime);
      if (tot) tot.textContent = _formatTime(this._audio.duration || 0);
      if (scrubber && this._audio.duration) {
        scrubber.value = (this._audio.currentTime / this._audio.duration) * 100;
      }
    }

    // --- Progress memory (localStorage) ---

    _storageKey() {
      var ep = this._episodes[this._currentIndex];
      if (!ep || !this._feedUrl) return null;
      return 'startr:progress:' + this._feedUrl + ':' + ep.guid;
    }

    _saveProgress() {
      if (!this._audio || !this._audio.currentTime) return;
      var key = this._storageKey();
      if (!key) return;
      try {
        localStorage.setItem(key, JSON.stringify({
          time: this._audio.currentTime,
          ts: Date.now()
        }));
      } catch(e) { /* quota exceeded or private mode */ }
    }

    _restoreProgress() {
      var key = this._storageKey();
      if (!key || !this._audio) return;
      try {
        var saved = localStorage.getItem(key);
        if (saved) {
          var data = JSON.parse(saved);
          if (data.time > 0) {
            this._audio.currentTime = data.time;
          }
        }
      } catch(e) { /* corrupt or unavailable */ }
    }

    _clearProgress() {
      var key = this._storageKey();
      if (!key) return;
      try { localStorage.removeItem(key); } catch(e) {}
    }

    // --- Keyboard ---

    _handleKeyboard(e) {
      switch(e.code) {
        case 'Space':
          e.preventDefault();
          this._togglePlay();
          break;
        case 'ArrowLeft':
          e.preventDefault();
          this._seek(-15);
          break;
        case 'ArrowRight':
          e.preventDefault();
          this._seek(15);
          break;
      }
    }

    // --- Share ---

    _share() {
      var ep = this._episodes[this._currentIndex];
      if (!ep) return;
      var data = {
        title: ep.title,
        text: ep.title + ' \u2014 ' + this._showTitle,
        url: location.href
      };

      var self = this;
      if (navigator.share) {
        navigator.share(data).catch(function() {});
      } else if (navigator.clipboard) {
        navigator.clipboard.writeText(data.url).then(function() {
          self._showTooltip('Copied!');
        }).catch(function() {});
      }
    }

    _showTooltip(text) {
      var btn = this._shadow.querySelector('.btn-share');
      if (!btn) return;
      var tip = document.createElement('span');
      tip.className = 'tooltip';
      tip.textContent = text;
      btn.parentNode.insertBefore(tip, btn.nextSibling);
      setTimeout(function() { tip.remove(); }, 2000);
    }
  }

  // --- Utilities ---

  // querySelector can't find namespaced elements (itunes:image) in XML docs.
  // Fall back to getElementsByTagNameNS for ns:local selectors.
  function _qsel(parent, selector) {
    // Handle "ns:local" selectors via namespace-aware lookup.
    // Can't use querySelector for namespaced XML elements — use getElementsByTagNameNS.
    // Detect the actual namespace URI from the document (feeds vary between
    // itunes.com and itunes.apple.com).
    var m = selector.match(/^itunes\\?:(\w+)$/);
    if (m) {
      var root = parent.ownerDocument ? parent.ownerDocument.documentElement : parent;
      var ns = root.getAttribute('xmlns:itunes') || 'http://www.itunes.com/dtds/podcast-1.0.dtd';
      var els = parent.getElementsByTagNameNS(ns, m[1]);
      return els.length > 0 ? els[0] : null;
    }
    try { return parent.querySelector(selector); } catch(e) { return null; }
  }

  function _text(parent, selector) {
    var el = _qsel(parent, selector);
    return el ? (el.textContent || '').trim() : '';
  }

  function _attr(parent, selector, attr) {
    var el = _qsel(parent, selector);
    return el ? (el.getAttribute(attr) || '') : '';
  }

  function _formatTime(s) {
    if (!s || !isFinite(s)) return '0:00';
    s = Math.floor(s);
    var h = Math.floor(s / 3600);
    var m = Math.floor((s % 3600) / 60);
    var sec = s % 60;
    var pad = sec < 10 ? '0' : '';
    if (h > 0) return h + ':' + (m < 10 ? '0' : '') + m + ':' + pad + sec;
    return m + ':' + pad + sec;
  }

  function _formatDate(dateStr) {
    try {
      var d = new Date(dateStr);
      if (isNaN(d.getTime())) return dateStr;
      return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
    } catch(e) { return dateStr; }
  }

  function _esc(s) {
    if (!s) return '';
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  function _escAttr(s) {
    return _esc(s).replace(/"/g, '&quot;');
  }

  customElements.define('startr-player', StartrPlayer);
})();
