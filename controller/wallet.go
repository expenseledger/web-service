package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/business"
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

type walletCreateForm struct {
	Name    string          `json:"name" binding:"required"`
	Type    string          `json:"type" binding:"required"`
	Balance decimal.Decimal `json:"balance"`
}

type walletIdentifyForm struct {
	Name string `json:"name" binding:"required"`
}

type txCreateForm struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Category    string          `json:"category" binding:"required"`
	Description string          `json:"description"`
	Date        date.Date       `json:"date"`
}

type walletTxRxForm struct {
	Name string `json:"name" binding:"required"`
	txCreateForm
}

type walletTransferForm struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	txCreateForm
}

func walletCreate(context *gin.Context) {
	var form walletCreateForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var wallet model.Wallet
	copier.Copy(&wallet, &form)

	if err := wallet.Create(); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(wallet),
	)
	return
}

func walletGet(context *gin.Context) {
	var form walletIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	wallet := model.Wallet{
		Name: form.Name,
	}

	if err := wallet.Get(); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(wallet),
	)
	return
}

func walletDelete(context *gin.Context) {
	var form walletIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var wallet model.Wallet
	if err := wallet.Delete(form.Name); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(wallet),
	)
	return
}

func walletList(context *gin.Context) {
	var wallets model.Wallets
	length, err := wallets.List()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  wallets,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func walletListTypes(context *gin.Context) {
	types := constant.ListWalletTypes()
	items := itemList{
		Length: len(types),
		Items:  types,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func walletInit(context *gin.Context) {
	wallets := model.Wallets{
		model.Wallet{
			Name:    "Cash",
			Type:    constant.WalletTypes().Cash,
			Balance: decimalFromStringIgnoreError("0.0"),
		},
		model.Wallet{
			Name:    "My Bank",
			Type:    constant.WalletTypes().BankAccount,
			Balance: decimalFromStringIgnoreError("0.0"),
		},
	}

	length, err := wallets.Init()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  wallets,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func walletClear(context *gin.Context) {
	var wallets model.Wallets
	length, err := wallets.Clear()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  wallets,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func walletExpend(context *gin.Context) {
	var form walletTxRxForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tx := form.toExpenseTransaction()

	srcWallet, err := business.Expend(tx)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"src_wallet":  srcWallet,
		"transaction": tx,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(data),
	)
	return
}

func walletReceive(context *gin.Context) {
	var form walletTxRxForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tx := form.toIncomeTransaction()

	dstWallet, err := business.Receive(tx)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"dst_wallet":  dstWallet,
		"transaction": tx,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(data),
	)
	return
}

func walletTransfer(context *gin.Context) {
	var form walletTransferForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	tx := form.toTransferTransaction()

	srcWallet, dstWallet, err := business.Transfer(tx)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"src_wallet":  srcWallet,
		"dst_wallet":  dstWallet,
		"transaction": tx,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(data),
	)
	return
}

func decimalFromStringIgnoreError(num string) decimal.Decimal {
	d, _ := decimal.NewFromString(num)
	return d
}

func (form *walletTransferForm) toTransferTransaction() *model.Transaction {
	tx := form.toTransaction(constant.TransactionTypes().Transfer)
	tx.SrcWallet = form.From
	tx.DstWallet = form.To
	return tx
}

func (form *walletTxRxForm) toExpenseTransaction() *model.Transaction {
	tx := form.toTransaction(constant.TransactionTypes().Expense)
	tx.SrcWallet = form.Name
	return tx
}

func (form *walletTxRxForm) toIncomeTransaction() *model.Transaction {
	tx := form.toTransaction(constant.TransactionTypes().Income)
	tx.DstWallet = form.Name
	return tx
}

func (form *txCreateForm) toTransaction(
	txType string,
) *model.Transaction {
	return &model.Transaction{
		Amount:      form.Amount,
		Type:        txType,
		Category:    form.Category,
		Description: form.Description,
		Date:        form.Date,
	}
}
