// Package main is the entrypoint for the feeds binary.
//
// feeds has two subcommands:
//
//	feeds rewrite   single-shot fetch + rewrite + atomic write
//	feeds serve     long-running ticker for always-on hosts
//
// v0.1.0 ships only the SpotifySource + RSS2PodcastRenderer pair. The
// dispatcher is intentionally stdlib-only (no cobra, no urfave/cli) because
// two subcommands don't need a framework.
package main

import (
	"fmt"
	"os"
)

// version is stamped at build time via -ldflags. Defaults to "dev" for
// local/bare `go build` invocations.
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := os.Args[1]
	args := os.Args[2:]

	switch sub {
	case "rewrite":
		if err := cmdRewrite(args); err != nil {
			fmt.Fprintf(os.Stderr, "feeds rewrite: %v\n", err)
			os.Exit(1)
		}
	case "serve":
		if err := cmdServe(args); err != nil {
			fmt.Fprintf(os.Stderr, "feeds serve: %v\n", err)
			os.Exit(1)
		}
	case "-h", "--help", "help":
		usage()
	case "-v", "--version", "version":
		fmt.Println("feeds", version)
	default:
		fmt.Fprintf(os.Stderr, "feeds: unknown subcommand %q\n\n", sub)
		usage()
		os.Exit(2)
	}
}

// usage prints the top-level banner and subcommand list. The ASCII art is
// shown only on the top-level --help screen — subcommand --help (rewrite,
// serve) deliberately stays clean so its output is grep-friendly for
// scripting and CI logs. Box-drawing characters render in any UTF-8 capable
// terminal (iTerm2, GNOME Terminal, Windows Terminal, vscode, Docker logs).
func usage() {
	fmt.Fprintf(os.Stderr, `
   (((•)))   ┌─┐┌─┐┌─┐┌┬┐┌─┐
             ├┤ ├┤ ├┤  ││└─┐
             └  └─┘└─┘─┴┘└─┘

  Startr/feeds %s — self-hosted feed rewriter
  Own the subscriber URL. Rent the audio host.

Usage:
  feeds <subcommand> [flags]

Subcommands:
  rewrite    fetch an upstream feed, rewrite branding, write static XML
  serve      long-running ticker mode (v0.2 swaps this for PocketBase)

Run "feeds <subcommand> --help" for subcommand flags.
`, version)
}
