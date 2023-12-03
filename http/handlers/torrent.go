package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/skif48/home-torrent/apihelpers"
	"github.com/skif48/home-torrent/config"
	"github.com/skif48/home-torrent/torrent"
)

type TorrentHandler struct {
	config *config.Config
}

func NewTorrentHandler(i *do.Injector) (*TorrentHandler, error) {
	config := do.MustInvoke[*config.Config](i)

	return &TorrentHandler{
		config: config,
	}, nil
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

	peers, err := torrent.RequestPeers([]byte(t.config.Torrent.PeerId), t.config.Torrent.PeerPort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"peers": peers})
}
