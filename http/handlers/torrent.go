package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/home-torrent/apihelpers"
	"github.com/home-torrent/config"
	"github.com/home-torrent/torrent"
)

type TorrentHandler struct {
	config *config.Config
}

func NewTorrentHandler(config *config.Config) *TorrentHandler {
	return &TorrentHandler{
		config: config,
	}
}

func (t *TorrentHandler) Preview(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("torrent")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	torrent, err := torrent.Parse(fileBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Malformed torrent file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"files": torrent.Files})
}

func (t *TorrentHandler) RequestPeers(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("torrent")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	torrent, err := torrent.Parse(fileBytes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Malformed torrent file"})
		return
	}

	peers, err := torrent.RequestPeers(t.config.Torrent.PeerId, t.config.Torrent.PeerPort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"peers": peers})
}
