package torrent

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"vladusenko.io/home-torrent/apihelpers"
)

func PreviewHandler(ctx *gin.Context) {
	var err error
	var fileHeader *multipart.FileHeader
	var fileReader multipart.File
	var fileBytes []byte
	var torrent *Torrent

	if fileHeader, err = ctx.FormFile("torrent"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if fileReader, err = fileHeader.Open(); err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	if fileBytes, err = ioutil.ReadAll(fileReader); err != nil {
		apihelpers.InternalServerError(ctx)
		return
	}

	if torrent, err = Parse(fileBytes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Malformed torrent file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"files": torrent.Files})
}
