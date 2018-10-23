package controller

import "github.com/gin-gonic/gin"

type itemList struct {
	Length int         `json:"length"`
	Items  interface{} `json:"items"`
}

func buildSuccessResponse(data interface{}) map[string]interface{} {
	return gin.H{
		"success": true,
		"data":    data,
	}
}

func buildNonsuccessResponse(
	err error,
	data interface{},
) map[string]interface{} {
	return gin.H{
		"success": false,
		"message": err.Error(),
		"context": data,
	}
}
