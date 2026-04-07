// Package main is the entrypoint for the feeds binary.
//
// feeds is built on PocketBase as a framework. PocketBase provides the
// HTTP server (static file serving from pb_public/), admin UI at /_/,
// built-in cron scheduler, SQLite, and graceful shutdown.
//
// The feed rewrite pipeline runs as a PocketBase JS hook
// (pb_hooks/feeds.pb.js). This Go binary is the server entrypoint only.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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

	// Serve static files from pb_public so the rewritten XML is
	// accessible over HTTP. PB framework mode doesn't auto-serve
	// pb_public — we register it explicitly.
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
