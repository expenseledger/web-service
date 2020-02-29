package controller

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/pkg"
	"github.com/expenseledger/web-service/pkg/type/date"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type txIdentifyForm struct {
	ID string `json:"id" binding:"required"`
}

type txListForm struct {
	Wallet string `json:"wallet" binding:"required"`
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

func createExpense(context *gin.Context) {
	var form txExpenseForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	wallet, err := model.GetWallet(form.From, userId)
	if err != nil {
		buildFailedContext(context, err)
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
		userId,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	err = wallet.Expend(tx)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	data := map[string]interface{}{
		"transaction": tx,
		"src_wallet":  wallet,
	}

	buildSuccessContext(context, data)
}

func createIncome(context *gin.Context) {
	var form txIncomeForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	wallet, err := model.GetWallet(form.To, userId)
	if err != nil {
		buildFailedContext(context, err)
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
		userId,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	err = wallet.Receive(tx)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	data := map[string]interface{}{
		"transaction": tx,
		"dst_wallet":  wallet,
	}

	buildSuccessContext(context, data)
}

func createTransfer(context *gin.Context) {
	var form txTransferForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	srcWallet, err := model.GetWallet(form.From, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	dstWallet, err := model.GetWallet(form.To, userId)
	if err != nil {
		buildFailedContext(context, err)
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
		userId,
	)

	if err != nil {
		buildFailedContext(context, err)
		return
	}

	err = srcWallet.Expend(tx)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	err = dstWallet.Receive(tx)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	data := map[string]interface{}{
		"transaction": tx,
		"src_wallet":  srcWallet,
		"dst_wallet":  dstWallet,
	}

	buildSuccessContext(context, data)
}

func getTransaction(context *gin.Context) {
	var form txIdentifyForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	tx, err := model.GetTransaction(form.ID, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, tx)
}

func clearTransactions(context *gin.Context) {
	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	txs, err := model.ClearTransactions(userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(txs),
		Items:  txs,
	}

	buildSuccessContext(context, items)
}

func listTransactions(context *gin.Context) {
	var form txListForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	txs, err := model.ListTransactions(form.Wallet, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(txs),
		Items:  txs,
	}

	buildSuccessContext(context, items)
}

func listTransactionTypes(context *gin.Context) {
	types := constant.ListTransactionTypes()
	items := itemList{
		Length: len(types),
		Items:  types,
	}

	buildSuccessContext(context, items)
}

func deleteTransaction(context *gin.Context) {
	var form txIdentifyForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	tx, err := model.DeleteTransaction(form.ID, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	var srcWallet, dstWallet *model.Wallet
	if from := tx.From; from != "" {
		srcWallet, err = model.GetWallet(from, userId)
		if err != nil {
			buildFailedContext(context, err)
			return
		}
		if err = srcWallet.Receive(tx); err != nil {
			buildFailedContext(context, err)
			return
		}
	}
	if to := tx.To; to != "" {
		dstWallet, err = model.GetWallet(to, userId)
		if err != nil {
			buildFailedContext(context, err)
			return
		}
		if err = dstWallet.Expend(tx); err != nil {
			buildFailedContext(context, err)
			return
		}
	}

	data := map[string]interface{}{
		"transaction": tx,
		"src_wallet":  srcWallet,
		"dst_wallet":  dstWallet,
	}

	buildSuccessContext(context, data)
}
