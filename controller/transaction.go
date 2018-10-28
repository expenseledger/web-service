package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

type transactionCreateForm struct {
	SrcWallet   string          `json:"src_wallet" binding:"required"`
	DstWallet   *string         `json:"dst_wallet"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Type        string          `json:"type" binding:"required"`
	Category    string          `json:"category" binding:"required"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

func transactionCreate(context *gin.Context) {
	var form transactionCreateForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var transaction model.Transaction
	copier.Copy(&transaction, &form)

	if err := transaction.Create(); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(transaction),
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
