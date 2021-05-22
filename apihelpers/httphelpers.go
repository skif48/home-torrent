package apihelpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
}
