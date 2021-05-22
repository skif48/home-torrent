package main

import (
	"vladusenko.io/home-torrent/api"
)

func main() {
	router := api.SetupRouter()

	router.Run(":8080")
}
