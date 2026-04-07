package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/pipeline"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// cmdRewrite handles `feeds rewrite ...` — a single-shot fetch, rewrite, and
// atomic write. Exits zero on success. On upstream failure, the existing
// output file is preserved (fail-loud behavior is implemented in the
// pipeline package).
func cmdRewrite(args []string) error {
	fs := flag.NewFlagSet("rewrite", flag.ContinueOnError)

	// Env vars are read as defaults so containerized deploys (CapRover,
	// fly.io, k8s) can configure feeds without baking values into the CMD.
	// CLI flags still override env vars when explicitly passed. See the
	// README's "Environment variables" section for the full table.
	var (
		sourceName  = fs.String("source", envString("FEEDS_SOURCE", "spotify"), "source adapter (env: FEEDS_SOURCE; v0.1.0: spotify only)")
		upstream    = fs.String("upstream", envString("FEEDS_UPSTREAM", ""), "upstream feed URL (env: FEEDS_UPSTREAM, required)")
		output      = fs.String("output", envString("FEEDS_OUTPUT", ""), "output XML path (env: FEEDS_OUTPUT, required)")
		selfURL     = fs.String("self-url", envString("FEEDS_SELF_URL", ""), "public URL subscribers bind to (env: FEEDS_SELF_URL, required)")
		title       = fs.String("channel-title", envString("FEEDS_CHANNEL_TITLE", ""), "channel title to inject (env: FEEDS_CHANNEL_TITLE, required)")
		link        = fs.String("channel-link", envString("FEEDS_CHANNEL_LINK", ""), "channel link to inject (env: FEEDS_CHANNEL_LINK, required)")
		image       = fs.String("channel-image", envString("FEEDS_CHANNEL_IMAGE", ""), "channel image URL (env: FEEDS_CHANNEL_IMAGE, optional)")
		itunesAuthr = fs.String("itunes-author", envString("FEEDS_ITUNES_AUTHOR", ""), "iTunes author (env: FEEDS_ITUNES_AUTHOR, optional)")
		itunesOwner = fs.String("itunes-owner-email", envString("FEEDS_ITUNES_OWNER_EMAIL", ""), "iTunes owner email (env: FEEDS_ITUNES_OWNER_EMAIL, optional)")
		stateFile   = fs.String("state", envString("FEEDS_STATE", ".feeds-state.json"), "path to cache state file for conditional GET (env: FEEDS_STATE)")
		configPath  = fs.String("config", envString("FEEDS_CONFIG", ""), "YAML config file (env: FEEDS_CONFIG; v0.1.1+, not yet supported)")
	)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `feeds rewrite — single-shot fetch + rewrite + atomic write

Usage:
  feeds rewrite --source spotify \
                --upstream   https://anchor.fm/s/YOUR_ID/podcast/rss \
                --output     ./public/v1/your-show.xml \
                --self-url   https://feed.yourdomain.com/v1/your-show.xml \
                --channel-title "Your Show" \
                --channel-link  https://yourdomain.com/podcast

Flags:
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		// `feeds rewrite --help` triggers fs.Usage() (which already
		// printed) and returns flag.ErrHelp. Treat that as success — the
		// user got what they asked for, exit 0 instead of dumping the
		// error to stderr and exiting 1.
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	if *configPath != "" {
		return fmt.Errorf("--config YAML support lands in v0.1.1; use flags in v0.1.0")
	}

	if *sourceName != "spotify" {
		return fmt.Errorf("unsupported --source %q (v0.1.0 ships spotify only)", *sourceName)
	}

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

	src := spotify.New(*upstream)
	rw := rewriter.New(rewriter.Branding{
		SelfURL:          *selfURL,
		Title:            *title,
		Link:             *link,
		Image:            *image,
		ITunesAuthor:     *itunesAuthr,
		ITunesOwnerEmail: *itunesOwner,
	})
	// Stamp the rewriter binary into the output XML's <generator> element so
	// the rewritten feed identifies the tool that produced it (replacing
	// upstream's "Anchor Podcasts" / "Spotify for Podcasters"). The version
	// var is stamped at build time via -ldflags from the Makefile.
	rw.Generator = fmt.Sprintf("Startr/feeds %s (https://github.com/Startr/feeds)", version)
	out := rss.NewRenderer(*output)
	ch := cache.New(*stateFile)

	return pipeline.Run(pipeline.Config{
		Source:   src,
		Rewriter: rw,
		Output:   out,
		Cache:    ch,
	})
}
