// Package spotify implements a Source adapter for Spotify for Podcasters'
// auto-generated RSS feeds.
//
// It's also a perfectly ordinary RSS fetcher — Spotify's feed is a regular
// RSS 2.0 document with iTunes namespace. The "Spotify" name just marks the
// source of the upstream URL in logs and future metadata.
package spotify

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Source holds configuration for fetching a Spotify for Podcasters RSS feed.
type Source struct {
	URL    string
	Client *http.Client
}

// New returns a Source for the given upstream URL with a sensible default
// HTTP client (30s timeout, follows redirects).
func New(url string) *Source {
	return &Source{
		URL: url,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchResult is what a conditional GET returns.
type FetchResult struct {
	Body         []byte // nil if NotModified == true
	ETag         string
	LastModified string
	NotModified  bool // true on HTTP 304
}

// Fetch performs an HTTP GET against the upstream feed URL, including
// conditional-GET headers if prior ETag / Last-Modified values are provided.
// Returns NotModified=true on HTTP 304 so callers can short-circuit the
// rewrite pipeline.
func (s *Source) Fetch(prevETag, prevLastModified string) (*FetchResult, error) {
	req, err := http.NewRequest(http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("User-Agent", "feeds/0.1.0 (+https://github.com/Startr/feeds)")
	req.Header.Set("Accept", "application/rss+xml, application/xml;q=0.9, */*;q=0.8")
	if prevETag != "" {
		req.Header.Set("If-None-Match", prevETag)
	}
	if prevLastModified != "" {
		req.Header.Set("If-Modified-Since", prevLastModified)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstream fetch: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotModified:
		return &FetchResult{NotModified: true}, nil
	case http.StatusOK:
		// fall through
	default:
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return &FetchResult{
		Body:         body,
		ETag:         resp.Header.Get("ETag"),
		LastModified: resp.Header.Get("Last-Modified"),
	}, nil
}
