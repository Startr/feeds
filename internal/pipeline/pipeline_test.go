package pipeline

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// minimalFeed is the smallest valid Spotify-style RSS fixture that exercises
// the rewrite scope table.
const minimalFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
     xmlns:atom="http://www.w3.org/2005/Atom"
     xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <atom:link href="https://upstream.example/rss" rel="self" type="application/rss+xml"/>
    <title>Upstream Title</title>
    <link>https://upstream.example/show</link>
    <description>Upstream description</description>
    <itunes:author>Upstream Author</itunes:author>
    <item>
      <title>Episode 1</title>
      <guid>ep1</guid>
      <enclosure url="https://upstream.example/audio/1.mp3" length="100" type="audio/mpeg"/>
    </item>
  </channel>
</rss>`

// newTestConfig wires up a pipeline.Config pointing at a local httptest
// server and a temp directory, with realistic branding.
func newTestConfig(t *testing.T, handler http.HandlerFunc) (Config, *httptest.Server, string) {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	outputPath := filepath.Join(dir, "out.xml")
	statePath := filepath.Join(dir, "state.json")

	cfg := Config{
		Source: spotify.New(srv.URL),
		Rewriter: rewriter.New(rewriter.Branding{
			SelfURL:     "https://feed.example.com/v1/show.xml",
			Title:       "Rewritten Title",
			Link:        "https://example.com/show",
			ITunesAuthor: "Rewritten Author",
		}),
		Output: rss.NewRenderer(outputPath),
		Cache:  cache.New(statePath),
	}
	return cfg, srv, outputPath
}

// TestPipeline_FailLoudPreservesLastGood verifies that when the upstream
// fetch fails on a subsequent run, the previously-written "last good"
// output file is NOT overwritten or truncated. This is the contract that
// lets users wire feeds to their publishing infra without worrying that a
// transient upstream outage will zero out their subscribers' feed.
func TestPipeline_FailLoudPreservesLastGood(t *testing.T) {
	requestCount := 0
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			// First run: serve the minimal feed successfully.
			w.Header().Set("Content-Type", "application/rss+xml")
			_, _ = w.Write([]byte(minimalFeed))
			return
		}
		// Subsequent runs: simulate an upstream outage.
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}

	cfg, _, outputPath := newTestConfig(t, handler)

	// First run: must succeed and write the output file.
	if err := Run(cfg); err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	goodBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read last-good file: %v", err)
	}
	if len(goodBytes) == 0 {
		t.Fatalf("first run wrote an empty file")
	}

	// Second run: upstream returns 502. Pipeline MUST return an error.
	if err := Run(cfg); err == nil {
		t.Fatalf("second run returned nil error despite upstream 502")
	}

	// The last-good output file MUST still exist and MUST be byte-identical
	// to what the first run wrote.
	afterBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output after failure: %v", err)
	}
	if string(afterBytes) != string(goodBytes) {
		t.Errorf("last-good output was mutated on upstream failure\nwant %d bytes, got %d bytes", len(goodBytes), len(afterBytes))
	}
}

// TestPipeline_IdempotentNoOpWrite verifies that running the pipeline twice
// against an unchanged upstream does not cause a second atomic rename. This
// matters for CDN cache invalidation (fewer unnecessary touches = fewer
// cache purges) and for filesystem noise (inotify watchers, rsync diffs).
//
// The test server returns 200 on both requests with the same body and no
// caching headers, so the rewriter runs both times. The rss.Renderer must
// detect the byte-identical output and skip the rename.
func TestPipeline_IdempotentNoOpWrite(t *testing.T) {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// Explicitly no ETag / Last-Modified so conditional GET can't
		// short-circuit the pipeline. We want to exercise the rewriter's
		// no-op path, not the cache's 304 path.
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(minimalFeed))
	}

	cfg, _, outputPath := newTestConfig(t, handler)

	// First run: writes the file for real.
	if err := Run(cfg); err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	stat1, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("stat after first run: %v", err)
	}

	// Second run: same upstream, same bytes. The rss.Renderer should
	// detect the byte-identical output and skip the rename. os.SameFile
	// returns true iff both FileInfo values point at the same inode, so
	// if a rename happened the inode would differ and SameFile would
	// return false.
	if err := Run(cfg); err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	stat2, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("stat after second run: %v", err)
	}

	if !os.SameFile(stat1, stat2) {
		t.Errorf("second run replaced the output file (inode changed) when it should have been a byte-identical no-op")
	}
}
