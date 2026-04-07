package main

import (
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

	var (
		sourceName  = fs.String("source", "spotify", "source adapter (v0.1.0: spotify only)")
		upstream    = fs.String("upstream", "", "upstream feed URL (required)")
		output      = fs.String("output", "", "output XML path (required)")
		selfURL     = fs.String("self-url", "", "public URL subscribers bind to (required)")
		title       = fs.String("channel-title", "", "channel title to inject")
		link        = fs.String("channel-link", "", "channel link to inject")
		image       = fs.String("channel-image", "", "channel image URL (optional)")
		itunesAuthr = fs.String("itunes-author", "", "iTunes author (optional)")
		itunesOwner = fs.String("itunes-owner-email", "", "iTunes owner email (optional)")
		stateFile   = fs.String("state", ".feeds-state.json", "path to cache state file for conditional GET")
		configPath  = fs.String("config", "", "YAML config file (v0.1.1+, not yet supported)")
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
	out := rss.NewRenderer(*output)
	ch := cache.New(*stateFile)

	return pipeline.Run(pipeline.Config{
		Source:   src,
		Rewriter: rw,
		Output:   out,
		Cache:    ch,
	})
}
