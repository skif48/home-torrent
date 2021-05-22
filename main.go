package main

import (
	"strconv"

	"vladusenko.io/home-torrent/api"
	"vladusenko.io/home-torrent/config"
)

const DEFAULT_CONFIG_PATH string = "./config.yml"

func main() {
	var err error
	var conf *config.Config

	if conf, err = config.GetConfig(DEFAULT_CONFIG_PATH); err != nil {
		panic(err)
	}

	router := api.SetupRouter()

	router.Run(":" + strconv.Itoa(conf.HttpPort))
}
