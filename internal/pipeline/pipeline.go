// Package pipeline orchestrates the source → rewriter → output flow with
// conditional-GET caching and fail-loud error semantics.
//
// The contract:
//
//  1. If the upstream returns 304 Not Modified, do nothing and return nil.
//  2. If the upstream returns 200, rewrite and atomically write the result.
//     Save the new ETag / Last-Modified to the cache.
//  3. If the upstream fetch, parse, or rewrite fails, return the error.
//     The existing output file is left untouched (last-good preservation
//     is inherent to the atomic write-rename pattern — we only swap the
//     file on success).
package pipeline

import (
	"fmt"
	"log"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// Config bundles the components needed for one rewrite run.
type Config struct {
	Source   *spotify.Source
	Rewriter *rewriter.Rewriter
	Output   *rss.Renderer
	Cache    *cache.Cache
}

// Run executes one fetch + rewrite + write cycle.
func Run(cfg Config) error {
	state, err := cfg.Cache.Load()
	if err != nil {
		// Non-fatal: we'll just do a full GET.
		log.Printf("cache load: %v (treating as first run)", err)
		state = cache.State{}
	}

	fetch, err := cfg.Source.Fetch(state.ETag, state.LastModified)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	if fetch.NotModified {
		log.Printf("upstream 304 Not Modified — no rewrite needed")
		return nil
	}

	rewritten, err := cfg.Rewriter.Rewrite(fetch.Body)
	if err != nil {
		// Don't touch the output file — last-good preservation.
		return fmt.Errorf("rewrite: %w", err)
	}

	changed, err := cfg.Output.Write(rewritten)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if changed {
		log.Printf("wrote %d bytes to %s", len(rewritten), cfg.Output.Path)
	} else {
		log.Printf("output unchanged at %s (idempotent no-op)", cfg.Output.Path)
	}

	// Only save the new cache state on successful write.
	newState := cache.State{
		ETag:         fetch.ETag,
		LastModified: fetch.LastModified,
	}
	if err := cfg.Cache.Save(newState); err != nil {
		// Non-fatal: worst case is one wasted full-GET next run.
		log.Printf("cache save: %v", err)
	}

	return nil
}
