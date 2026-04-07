package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Startr/feeds/internal/cache"
	"github.com/Startr/feeds/internal/output/rss"
	"github.com/Startr/feeds/internal/pipeline"
	"github.com/Startr/feeds/internal/rewriter"
	"github.com/Startr/feeds/internal/source/spotify"
)

// newRewriteCmd returns the `feeds rewrite` Cobra command — a single-shot
// fetch, rewrite, and atomic write. Exits zero on success. On upstream
// failure, the existing output file is preserved (fail-loud behavior is
// implemented in the pipeline package).
func newRewriteCmd() *cobra.Command {
	var (
		sourceName       string
		upstream         string
		output           string
		selfURL          string
		title            string
		link             string
		image            string
		itunesAuthor     string
		itunesOwnerEmail string
		stateFile        string
		configPath       string
	)

	cmd := &cobra.Command{
		Use:   "rewrite",
		Short: "Single-shot fetch + rewrite + atomic write",
		Long: `feeds rewrite — single-shot fetch + rewrite + atomic write

Fetch an upstream feed, rewrite branding and channel metadata to point
at your own URL, and write the result to a static XML file. Exits zero
on success. Wire it up to any external scheduler (cron, systemd timer,
GitHub Actions, k8s CronJob).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath != "" {
				return fmt.Errorf("--config YAML support lands in v0.1.1; use flags or env vars")
			}

			if sourceName != "spotify" {
				return fmt.Errorf("unsupported --source %q (v0.0.x ships spotify only)", sourceName)
			}

			missing := []string{}
			if upstream == "" {
				missing = append(missing, "--upstream")
			}
			if output == "" {
				missing = append(missing, "--output")
			}
			if selfURL == "" {
				missing = append(missing, "--self-url")
			}
			if title == "" {
				missing = append(missing, "--channel-title")
			}
			if link == "" {
				missing = append(missing, "--channel-link")
			}
			if len(missing) > 0 {
				return fmt.Errorf("missing required flags: %v", missing)
			}

			// Ensure output directory exists.
			if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
				return fmt.Errorf("create output dir: %w", err)
			}

			src := spotify.New(upstream)
			rw := rewriter.New(rewriter.Branding{
				SelfURL:          selfURL,
				Title:            title,
				Link:             link,
				Image:            image,
				ITunesAuthor:     itunesAuthor,
				ITunesOwnerEmail: itunesOwnerEmail,
			})
			rw.Generator = fmt.Sprintf("Startr/feeds %s (https://github.com/Startr/feeds)", version)
			out := rss.NewRenderer(output)
			ch := cache.New(stateFile)

			return pipeline.Run(pipeline.Config{
				Source:   src,
				Rewriter: rw,
				Output:   out,
				Cache:    ch,
			})
		},
	}

	// Env vars are read as defaults so containerized deploys can configure
	// feeds without baking values into the CMD. CLI flags still override
	// env vars when explicitly passed.
	cmd.Flags().StringVar(&sourceName, "source", envString("FEEDS_SOURCE", "spotify"), "source adapter (env: FEEDS_SOURCE)")
	cmd.Flags().StringVar(&upstream, "upstream", envString("FEEDS_UPSTREAM", ""), "upstream feed URL (env: FEEDS_UPSTREAM, required)")
	cmd.Flags().StringVar(&output, "output", envString("FEEDS_OUTPUT", ""), "output XML path (env: FEEDS_OUTPUT, required)")
	cmd.Flags().StringVar(&selfURL, "self-url", envString("FEEDS_SELF_URL", ""), "public URL subscribers bind to (env: FEEDS_SELF_URL, required)")
	cmd.Flags().StringVar(&title, "channel-title", envString("FEEDS_CHANNEL_TITLE", ""), "channel title (env: FEEDS_CHANNEL_TITLE, required)")
	cmd.Flags().StringVar(&link, "channel-link", envString("FEEDS_CHANNEL_LINK", ""), "channel link (env: FEEDS_CHANNEL_LINK, required)")
	cmd.Flags().StringVar(&image, "channel-image", envString("FEEDS_CHANNEL_IMAGE", ""), "channel image URL (env: FEEDS_CHANNEL_IMAGE)")
	cmd.Flags().StringVar(&itunesAuthor, "itunes-author", envString("FEEDS_ITUNES_AUTHOR", ""), "iTunes author (env: FEEDS_ITUNES_AUTHOR)")
	cmd.Flags().StringVar(&itunesOwnerEmail, "itunes-owner-email", envString("FEEDS_ITUNES_OWNER_EMAIL", ""), "iTunes owner email (env: FEEDS_ITUNES_OWNER_EMAIL)")
	cmd.Flags().StringVar(&stateFile, "state", envString("FEEDS_STATE", ".feeds-state.json"), "cache state file for conditional GET (env: FEEDS_STATE)")
	cmd.Flags().StringVar(&configPath, "config", envString("FEEDS_CONFIG", ""), "YAML config file (env: FEEDS_CONFIG; v0.1.1+)")

	return cmd
}
