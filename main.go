package main

import (
	"strconv"

	"vladusenko.io/home-torrent/api"
	"vladusenko.io/home-torrent/config"
	"vladusenko.io/home-torrent/defaults"
	"vladusenko.io/home-torrent/log"
)

func main() {
	var err error
	var conf *config.Config

	if conf, err = config.GetConfig(defaults.DEFAULT_CONFIG_PATH); err != nil {
		panic(err)
	}

	log.Configure(conf.Logging)
	logger := log.GetLogger()
	router := api.SetupRouter()
	logger.Info().Msg("GIN Router has been set up")
	logger.Info().Msgf("Setting up http listener on port %d", conf.HttpPort)

	router.Run(":" + strconv.Itoa(conf.HttpPort))
}
