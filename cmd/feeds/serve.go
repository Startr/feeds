package main

import (
	"context"
	"errors"
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
		// `feeds serve --help` triggers fs.Usage() (which already printed)
		// and returns flag.ErrHelp. Treat that as success — the user got
		// what they asked for, exit 0 instead of dumping the error to
		// stderr and exiting 1.
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	// The --http flag is parsed for forward-compat with v0.2 PocketBase but
	// not used in v0.1.0. Log it so operators know it's ignored. The version
	// string comes from the package-level `version` var, stamped at build
	// time via -ldflags from the Makefile's RELEASE_VERSION cascade.
	log.Printf("feeds serve %s: ticker mode, interval=%s, http=%s (ignored until v0.2)", version, *interval, *httpAddr)

	// Set up signal handling first — both the configured (ticker) and idle
	// (no-config) paths need it.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Validate the flags v0.1.0 needs to actually do work. If any are
	// missing, fail-soft: log the situation and idle on the signal context
	// instead of erroring out. This matches the v0.2 vision where PocketBase
	// collections will be the source of feed config and starting up with no
	// flags is normal — the server comes up empty and waits for an admin to
	// add a feed. For v0.1.0 it keeps `make it_run` (no flags passed) and
	// CapRover deploys-without-config working as smoke tests of the binary
	// itself, instead of crashing the container loop on first boot.
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
		log.Printf("feeds serve: no feed configured (missing %v) — idling. Set the missing flags to start the rewrite pipeline. Send SIGTERM or Ctrl+C to stop.", missing)
		<-ctx.Done()
		log.Printf("feeds serve: shutdown signal received")
		return nil
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
	// Stamp the rewriter binary into the output XML's <generator> element so
	// the rewritten feed identifies the tool that produced it (replacing
	// upstream's "Anchor Podcasts" / "Spotify for Podcasters"). The version
	// var is stamped at build time via -ldflags from the Makefile.
	rw.Generator = fmt.Sprintf("Startr/feeds %s (https://github.com/Startr/feeds)", version)
	out := rss.NewRenderer(*output)
	ch := cache.New(stateFile)

	cfg := pipeline.Config{
		Source:   src,
		Rewriter: rw,
		Output:   out,
		Cache:    ch,
	}

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
