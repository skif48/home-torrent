package main

import (
	"github.com/samber/do"
	"github.com/skif48/home-torrent/config"
	"github.com/skif48/home-torrent/http"
	"github.com/skif48/home-torrent/http/handlers"
)

func main() {
	injector := do.New()
	do.Provide(injector, config.NewConfig)
	do.Provide(injector, handlers.NewTorrentHandler)
	do.Provide(injector, http.NewServer)

	server := do.MustInvoke[*http.Server](injector)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
