package http

import (
	"github.com/gin-gonic/gin"
	"github.com/home-torrent/http/handlers"
)

func SetupRouter(th *handlers.TorrentHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	apiV1 := api.Group("/v1")
	{
		torrentGroup := apiV1.Group("/torrents")
		torrentGroup.POST("/preview", th.Preview)
		torrentGroup.POST("/preview-peers", th.RequestPeers)
	}

	return router
}
