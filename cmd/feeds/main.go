// Package main is the entrypoint for the feeds binary.
//
// feeds is built on PocketBase as a Go framework. PocketBase provides the
// HTTP server (static file serving from --publicDir), admin UI at /_/,
// built-in cron scheduler, Cobra CLI, SQLite, and graceful shutdown.
//
// Subcommands:
//
//	feeds serve     PocketBase server + cron-driven rewrite pipeline
//	feeds rewrite   single-shot fetch + rewrite + atomic write (no server)
//	feeds superuser create/update admin accounts for the dashboard
//
// Feed-specific config comes from FEEDS_* env vars (v0.0.x single-feed
// bridge). When multi-feed support lands, feed config moves into PocketBase
// collections managed via the admin UI — the pipeline internals stay
// identical, only the data source changes.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/pipeline"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// version is stamped at build time via -ldflags. Defaults to "dev" for
// local/bare `go build` invocations.
var version = "dev"

func main() {
	app := pocketbase.New()

	// Customize Cobra root command so `feeds --help` and `feeds --version`
	// show our branding instead of PocketBase's defaults.
	app.RootCmd.Use = "feeds"
	app.RootCmd.Version = version
	app.RootCmd.Long = fmt.Sprintf(`
   (((•)))   ┌─┐┌─┐┌─┐┌┬┐┌─┐
             ├┤ ├┤ ├┤  ││└─┐
             └  └─┘└─┘─┴┘└─┘

  Startr/feeds %s — self-hosted feed rewriter
  Own the subscriber URL. Rent the audio host.

  Built on PocketBase. Admin UI at /_/ when running "feeds serve".`, version)

	// Register the single-shot rewrite command (lightweight, no PB bootstrap).
	app.RootCmd.AddCommand(newRewriteCmd())

	// -----------------------------------------------------------------
	// OnServe hook: wired when `feeds serve` starts the HTTP server.
	// Registers static file serving + the cron-driven rewrite pipeline.
	// -----------------------------------------------------------------
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Serve static files from pb_public so the rewritten XML is
		// accessible over HTTP. PB framework mode doesn't auto-serve
		// pb_public — we register it explicitly. The {path...} catch-all
		// has lower priority than PB's internal routes (admin UI at /_/).
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// Build pipeline config from FEEDS_* env vars. If required vars
		// are missing, log and continue — PB still serves the admin UI,
		// which matches the v0.2 vision where the server starts empty and
		// waits for an admin to add a feed via the dashboard.
		cfg, err := buildPipelineConfig(app)
		if err != nil {
			log.Printf("feeds serve: no feed configured (%v) — admin UI only. Set the FEEDS_* env vars to enable the rewrite pipeline.", err)
			return se.Next()
		}

		// Run the pipeline once immediately so subscribers get a fresh
		// feed at deploy time — no waiting for the first cron tick.
		go func() {
			if err := pipeline.Run(cfg); err != nil {
				log.Printf("initial rewrite: %v", err)
			}
		}()

		// Schedule recurring rewrites via PB's built-in cron.
		cronExpr := envString("FEEDS_CRON", "*/15 * * * *")
		if err := app.Cron().Add("feeds_rewrite", cronExpr, func() {
			if err := pipeline.Run(cfg); err != nil {
				log.Printf("rewrite failed: %v", err)
			}
		}); err != nil {
			// Bad cron expression — fall back to every 15 minutes.
			log.Printf("invalid FEEDS_CRON %q: %v — falling back to */15 * * * *", cronExpr, err)
			app.Cron().MustAdd("feeds_rewrite", "*/15 * * * *", func() {
				if err := pipeline.Run(cfg); err != nil {
					log.Printf("rewrite failed: %v", err)
				}
			})
		}

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// buildPipelineConfig assembles a pipeline.Config from FEEDS_* env vars.
// Returns an error if any required var is missing (caller decides whether
// to fail hard or fail soft).
//
// v0.0.x: reads env vars (single-feed bridge).
// v0.2+:  will read from PB collections instead — same return type,
//
//	different data source.
func buildPipelineConfig(app *pocketbase.PocketBase) (pipeline.Config, error) {
	upstream := envString("FEEDS_UPSTREAM", "")
	output := envString("FEEDS_OUTPUT", "")
	selfURL := envString("FEEDS_SELF_URL", "")
	title := envString("FEEDS_CHANNEL_TITLE", "")
	link := envString("FEEDS_CHANNEL_LINK", "")

	missing := []string{}
	if upstream == "" {
		missing = append(missing, "FEEDS_UPSTREAM")
	}
	if output == "" {
		missing = append(missing, "FEEDS_OUTPUT")
	}
	if selfURL == "" {
		missing = append(missing, "FEEDS_SELF_URL")
	}
	if title == "" {
		missing = append(missing, "FEEDS_CHANNEL_TITLE")
	}
	if link == "" {
		missing = append(missing, "FEEDS_CHANNEL_LINK")
	}
	if len(missing) > 0 {
		return pipeline.Config{}, fmt.Errorf("missing required env vars: %v", missing)
	}

	// Ensure output directory exists (may be first boot).
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.Config{}, fmt.Errorf("create output dir: %w", err)
	}

	stateFile := envString("FEEDS_STATE", filepath.Join(app.DataDir(), ".feeds-state.json"))

	src := spotify.New(upstream)
	rw := rewriter.New(rewriter.Branding{
		SelfURL:          selfURL,
		Title:            title,
		Link:             link,
		Image:            envString("FEEDS_CHANNEL_IMAGE", ""),
		ITunesAuthor:     envString("FEEDS_ITUNES_AUTHOR", ""),
		ITunesOwnerEmail: envString("FEEDS_ITUNES_OWNER_EMAIL", ""),
	})
	rw.Generator = fmt.Sprintf("Startr/feeds %s (https://github.com/Startr/feeds)", version)

	return pipeline.Config{
		Source:   src,
		Rewriter: rw,
		Output:   rss.NewRenderer(output),
		Cache:    cache.New(stateFile),
	}, nil
}

// envString returns the value of env var key if non-empty, otherwise
// fallback. Used to give FEEDS_* env vars higher precedence than literal
// defaults while still letting CLI flags (when present) override both.
func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
