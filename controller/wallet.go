package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func walletCreate(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"user": "Kohpai", "value": 12})
}
