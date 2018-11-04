package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/business"
	"github.com/gin-gonic/gin"
)

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}

func transactionClear(context *gin.Context) {
	length, txs, err := business.ClearTransactions()
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
