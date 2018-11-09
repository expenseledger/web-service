package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/business"
	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
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

func transactionGet(context *gin.Context) {
	var form transactionIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tx := model.Transaction{
		ID: form.ID,
	}

	if err := tx.Get(); err != nil {
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

func transactionDelete(context *gin.Context) {
	var form transactionIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tx := model.Transaction{
		ID: form.ID,
	}

	srcWallet, dstWallet, err := business.DeleteTransaction(&tx)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"transaction": tx,
	}
	if srcWallet != nil {
		data["src_wallet"] = srcWallet
	}
	if dstWallet != nil {
		data["dst_wallet"] = dstWallet
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(data),
	)
	return
}
