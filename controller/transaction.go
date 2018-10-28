package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}

func transactionGet(context *gin.Context) {
	var form transactionIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var tx model.Transaction
	if err := tx.Get(form.ID); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(tx),
	)
	return
}

func transactionClear(context *gin.Context) {
	var txs model.Transactions
	length, err := txs.Clear()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  txs,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}
