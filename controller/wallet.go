package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/pkg/type/date"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type walletCreateForm struct {
	Name    string              `json:"name" binding:"required"`
	Type    constant.WalletType `json:"type" binding:"required"`
	Balance decimal.Decimal     `json:"balance"`
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

func createWallet(context *gin.Context) {
	var form walletCreateForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	wallet, err := model.CreateWallet(form.Name, form.Type, form.Balance)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, wallet)
	return
}

func getWallet(context *gin.Context) {
	var form walletIdentifyForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	wallet, err := model.GetWallet(form.Name)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, wallet)
	return
}

func deleteWallet(context *gin.Context) {
	var form walletIdentifyForm
	if err := bindJSON(context, &form); err != nil {
		buildFailedContext(context, err)
		return
	}

	wallet, err := model.DeleteWallet(form.Name)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, wallet)
	return
}

func listWallets(context *gin.Context) {
	wallets, err := model.ListWallets()
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(wallets),
		Items:  wallets,
	}

	buildSuccessContext(context, items)
	return
}

func listWalletTypes(context *gin.Context) {
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

func initWallets(context *gin.Context) {
	recipes := []walletCreateForm{
		walletCreateForm{
			Name:    "Cash",
			Type:    constant.WalletTypes().Cash,
			Balance: decimalFromStringIgnoreError("0.0"),
		},
		walletCreateForm{
			Name:    "My Bank",
			Type:    constant.WalletTypes().BankAccount,
			Balance: decimalFromStringIgnoreError("0.0"),
		},
	}

	length := len(recipes)
	wallets := make([]*model.Wallet, length)
	for i, recipe := range recipes {
		wallet, err := model.CreateWallet(
			recipe.Name,
			recipe.Type,
			recipe.Balance,
		)
		if err != nil {
			buildFailedContext(context, err)
			return
		}
		wallets[i] = wallet
	}

	items := itemList{
		Length: length,
		Items:  wallets,
	}

	buildSuccessContext(context, items)
	return
}

func clearWallets(context *gin.Context) {
	wallets, err := model.ClearWallets()
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(wallets),
		Items:  wallets,
	}

	buildSuccessContext(context, items)
	return
}

// func walletExpend(context *gin.Context) {
// 	var form walletTxRxForm
// 	if err := context.ShouldBindJSON(&form); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	tx := form.toExpenseTransaction()

// 	srcWallet, err := business.Expend(tx)
// 	if err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"src_wallet":  srcWallet,
// 		"transaction": tx,
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(data),
// 	)
// 	return
// }

// func walletReceive(context *gin.Context) {
// 	var form walletTxRxForm
// 	if err := context.ShouldBindJSON(&form); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	tx := form.toIncomeTransaction()

// 	dstWallet, err := business.Receive(tx)
// 	if err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"dst_wallet":  dstWallet,
// 		"transaction": tx,
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(data),
// 	)
// 	return
// }

// func walletTransfer(context *gin.Context) {
// 	var form walletTransferForm
// 	if err := context.ShouldBindJSON(&form); err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	tx := form.toTransferTransaction()

// 	srcWallet, dstWallet, err := business.Transfer(tx)
// 	if err != nil {
// 		context.JSON(
// 			http.StatusBadRequest,
// 			buildNonsuccessResponse(err, nil),
// 		)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"src_wallet":  srcWallet,
// 		"dst_wallet":  dstWallet,
// 		"transaction": tx,
// 	}

// 	context.JSON(
// 		http.StatusOK,
// 		buildSuccessResponse(data),
// 	)
// 	return
// }

func decimalFromStringIgnoreError(num string) decimal.Decimal {
	d, _ := decimal.NewFromString(num)
	return d
}

// func (form *walletTransferForm) toTransferTransaction() *model.Transaction {
// 	tx := form.toTransaction(constant.TransactionTypes().Transfer)
// 	tx.SrcWallet = form.From
// 	tx.DstWallet = form.To
// 	return tx
// }

// func (form *walletTxRxForm) toExpenseTransaction() *model.Transaction {
// 	tx := form.toTransaction(constant.TransactionTypes().Expense)
// 	tx.SrcWallet = form.Name
// 	return tx
// }

// func (form *walletTxRxForm) toIncomeTransaction() *model.Transaction {
// 	tx := form.toTransaction(constant.TransactionTypes().Income)
// 	tx.DstWallet = form.Name
// 	return tx
// }

// func (form *txCreateForm) toTransaction(
// 	txType string,
// ) *model.Transaction {
// 	return &model.Transaction{
// 		Amount:      form.Amount,
// 		Type:        txType,
// 		Category:    form.Category,
// 		Description: form.Description,
// 		Date:        form.Date,
// 	}
// }
