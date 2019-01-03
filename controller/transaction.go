package controller

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/pkg/type/date"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type transactionIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}

type txCreateForm struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Category    string          `json:"category" binding:"required"`
	Description string          `json:"description"`
	Date        date.Date       `json:"date"`
}

type txExpenseForm struct {
	From string `json:"from" binding:"required"`
	txCreateForm
}

type txIncomeForm struct {
	To string `json:"to" binding:"required"`
	txCreateForm
}

type txTransferForm struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	txCreateForm
}

// func transactionClear(context *gin.Context) {
// 	var txs model.Transactions
// 	length, err := txs.Clear()
// 	if err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	items := itemList{
// 		Length: length,
// 		Items:  txs,
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(items),
// 	)
// 	return
// }

func createExpense(context *gin.Context) {
	var form txExpenseForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	tx, err := model.CreateTransction(
		form.Amount,
		constant.TransactionTypes().Expense,
		form.From,
		"",
		form.Category,
		form.Description,
		form.Date,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, tx)
	return
}

func createIncome(context *gin.Context) {
	var form txIncomeForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	tx, err := model.CreateTransction(
		form.Amount,
		constant.TransactionTypes().Income,
		"",
		form.To,
		form.Category,
		form.Description,
		form.Date,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, tx)
	return
}

func createTransfer(context *gin.Context) {
	var form txTransferForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	tx, err := model.CreateTransction(
		form.Amount,
		constant.TransactionTypes().Transfer,
		form.From,
		form.To,
		form.Category,
		form.Description,
		form.Date,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, tx)
	return
}

func clearTransactions(context *gin.Context) {
	txs, err := model.ClearTransactions()
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(txs),
		Items:  txs,
	}

	buildSuccessContext(context, items)
	return
}

// func transactionGet(context *gin.Context) {
// 	var form transactionIdentifyForm
// 	if err := context.ShouldBindJSON(&form); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	tx := model.Transaction{
// 		ID: form.ID,
// 	}

// 	if err := tx.Get(); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(tx),
// 	)
// 	return
// }

// func transactionDelete(context *gin.Context) {
// 	var form transactionIdentifyForm
// 	if err := context.ShouldBindJSON(&form); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	tx := model.Transaction{
// 		ID: form.ID,
// 	}

// 	srcWallet, dstWallet, err := business.DeleteTransaction(&tx)
// 	if err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"transaction": tx,
// 	}
// 	if srcWallet != nil {
// 		data["src_wallet"] = srcWallet
// 	}
// 	if dstWallet != nil {
// 		data["dst_wallet"] = dstWallet
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(data),
// 	)
// 	return
// }
