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

type txCreateCommonForm struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

type walletTxRxForm struct {
	Name     string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
	txCreateCommonForm
}

type walletTransferForm struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	txCreateCommonForm
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
			Type:    constant.WalletType.Cash,
			Balance: decimalFromStringIgnoreError("0.0"),
		},
		model.Wallet{
			Name:    "My Bank",
			Type:    constant.WalletType.BankAccount,
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

	expense := form.toExpense()

	wallet, err := business.InsertExpense(expense)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"wallet":  wallet,
		"expense": expense,
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

	income := form.toIncome()

	wallet, err := business.InsertIncome(income)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"wallet": wallet,
		"income": income,
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

	transfer := form.toTransfer()

	srcWallet, dstWallet, err := business.Transfer(transfer)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	data := map[string]interface{}{
		"src_wallet": srcWallet,
		"dst_wallet": dstWallet,
		"transfer":   transfer,
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

func (form *walletTxRxForm) toExpense() *model.ExpenseIncome {
	ei := form.toExpenseIncome(constant.TransactionType.Expense)
	ei.Wallet = form.Name
	ei.Category = form.Category
	return ei
}

func (form *walletTxRxForm) toIncome() *model.ExpenseIncome {
	ei := form.toExpenseIncome(constant.TransactionType.Income)
	ei.Wallet = form.Name
	ei.Category = form.Category
	return ei
}

func (form *walletTransferForm) toTransfer() *model.Transfer {
	return &model.Transfer{
		SrcWallet:   form.From,
		DstWallet:   form.To,
		Amount:      form.Amount,
		Description: form.Description,
		Date:        form.Date,
	}
}

func (form *txCreateCommonForm) toExpenseIncome(
	txType string,
) *model.ExpenseIncome {
	return &model.ExpenseIncome{
		Amount:      form.Amount,
		Type:        txType,
		Description: form.Description,
		Date:        form.Date,
	}
}
