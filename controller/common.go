package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func bindJSON(context *gin.Context, form interface{}) (err error) {
	if err = context.ShouldBindJSON(form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}
	return
}

func buildFailedContext(context *gin.Context, err error) {
	context.JSON(
		http.StatusBadRequest,
		buildNonsuccessResponse(err, nil),
	)
}

func buildSuccessContext(context *gin.Context, data interface{}) {
	context.JSON(
		http.StatusOK,
		buildSuccessResponse(data),
	)
}
