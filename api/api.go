package api

import (
	"github.com/gin-gonic/gin"
	"vladusenko.io/home-torrent/torrent"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	apiV1 := api.Group("/v1")
	{
		torrentGroup := apiV1.Group("/torrents")
		torrentGroup.POST("/preview", torrent.PreviewHandler)
	}

	return router
}
