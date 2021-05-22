package main

import (
	"strconv"

	"vladusenko.io/home-torrent/api"
	"vladusenko.io/home-torrent/config"
	"vladusenko.io/home-torrent/defaults"
)

func main() {
	var err error
	var conf *config.Config

	if conf, err = config.GetConfig(defaults.DEFAULT_CONFIG_PATH); err != nil {
		panic(err)
	}

	router := api.SetupRouter()

	router.Run(":" + strconv.Itoa(conf.HttpPort))
}
