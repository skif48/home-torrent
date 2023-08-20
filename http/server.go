package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/skif48/home-torrent/config"
	"github.com/skif48/home-torrent/http/handlers"
)

type Server struct {
	router *gin.Engine
	conf   *config.Config

	s *http.Server
}

func NewServer(i *do.Injector) (*Server, error) {
	th := do.MustInvoke[*handlers.TorrentHandler](i)
	conf := do.MustInvoke[*config.Config](i)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	apiV1 := api.Group("/v1")
	{
		torrentGroup := apiV1.Group("/torrents")
		torrentGroup.POST("/preview", th.Preview)
		torrentGroup.POST("/preview-peers", th.RequestPeers)
	}

	return &Server{
		router: router,
		conf:   conf,
		s: &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HttpPort),
			Handler: router,
		},
	}, nil
}

func (s *Server) ListenAndServe() error {
	if err := s.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
