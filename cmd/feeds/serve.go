package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/pipeline"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// cmdServe handles `feeds serve ...` — long-running ticker mode for always-on
// hosts (CapRover, fly.io machines, k8s Deployment). v0.1.0 is stdlib-only:
// no PocketBase, no HTTP admin UI, no multi-feed orchestration. It just
// re-runs the rewrite pipeline on a fixed interval.
//
// v0.2.0 replaces this entire file with a PocketBase framework import. The
// flag surface is chosen to be forward-compatible with PB's serve flags
// (--http, --dir, --hooks-dir, --migrations-dir, --public-dir) so CapRover
// deployments survive the swap.
func cmdServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)

	var (
		httpAddr   = fs.String("http", "0.0.0.0:8090", "HTTP bind address (reserved for v0.2 PocketBase admin; unused in v0.1.0)")
		dataDir    = fs.String("dir", "/app/pb_data", "data directory — state file lives here")
		interval   = fs.Duration("interval", 15*time.Minute, "how often to re-run the rewrite pipeline")
		upstream   = fs.String("upstream", "", "upstream feed URL (required in v0.1.0 until PB collections land)")
		output     = fs.String("output", "", "output XML path (required)")
		selfURL    = fs.String("self-url", "", "public URL subscribers bind to (required)")
		title      = fs.String("channel-title", "", "channel title to inject")
		link       = fs.String("channel-link", "", "channel link to inject")
		image      = fs.String("channel-image", "", "channel image URL (optional)")
		itunesAuth = fs.String("itunes-author", "", "iTunes author (optional)")
		itunesOwnr = fs.String("itunes-owner-email", "", "iTunes owner email (optional)")
	)

	// Register flags we parse for forward-compat with v0.2 PocketBase but
	// don't read in v0.1.0. Statement form (discarding return value) keeps
	// the unused-var checker happy.
	fs.String("hooks-dir", "/app/pb_hooks", "reserved for v0.2 PocketBase hooks")
	fs.String("migrations-dir", "/app/pb_migrations", "reserved for v0.2 PocketBase migrations")
	fs.String("public-dir", "/app/pb_public", "reserved for v0.2 PocketBase public assets")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `feeds serve — long-running ticker (v0.2 swaps this for PocketBase)

Usage:
  feeds serve --upstream https://... --output /app/pb_data/public/show.xml \
              --self-url https://feed.example.com/v1/show.xml \
              --channel-title "Show" --channel-link https://example.com \
              --interval 15m

Flags:
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// The --http flag is parsed for forward-compat with v0.2 PocketBase but
	// not used in v0.1.0. Log it so operators know it's ignored.
	log.Printf("feeds serve v0.1.0: ticker mode, interval=%s, http=%s (ignored until v0.2)", *interval, *httpAddr)

	missing := []string{}
	if *upstream == "" {
		missing = append(missing, "--upstream")
	}
	if *output == "" {
		missing = append(missing, "--output")
	}
	if *selfURL == "" {
		missing = append(missing, "--self-url")
	}
	if *title == "" {
		missing = append(missing, "--channel-title")
	}
	if *link == "" {
		missing = append(missing, "--channel-link")
	}
	if len(missing) > 0 {
		fs.Usage()
		return fmt.Errorf("missing required flags: %v", missing)
	}

	// Ensure the data dir exists (state file lives here by default).
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	stateFile := *dataDir + "/.feeds-state.json"

	src := spotify.New(*upstream)
	rw := rewriter.New(rewriter.Branding{
		SelfURL:          *selfURL,
		Title:            *title,
		Link:             *link,
		Image:            *image,
		ITunesAuthor:     *itunesAuth,
		ITunesOwnerEmail: *itunesOwnr,
	})
	out := rss.NewRenderer(*output)
	ch := cache.New(stateFile)

	cfg := pipeline.Config{
		Source:   src,
		Rewriter: rw,
		Output:   out,
		Cache:    ch,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run once immediately, then on the ticker.
	if err := pipeline.Run(cfg); err != nil {
		log.Printf("initial run failed: %v", err)
	}

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("feeds serve: shutdown signal received")
			return nil
		case <-ticker.C:
			if err := pipeline.Run(cfg); err != nil {
				// Fail-loud but keep running — the last-good output is
				// preserved by the pipeline itself.
				log.Printf("rewrite failed: %v", err)
			}
		}
	}
}
