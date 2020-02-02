package controller

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/pkg"
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

func createWallet(context *gin.Context) {
	var form walletCreateForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	userId, err := pkg.GetUserId(context)

	if err != nil {
		return
	}

	wallet, err := model.CreateWallet(form.Name, form.Type, form.Balance, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, wallet)
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
}

func listWalletTypes(context *gin.Context) {
	types := constant.ListWalletTypes()
	items := itemList{
		Length: len(types),
		Items:  types,
	}

	buildSuccessContext(context, items)
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
}

func decimalFromStringIgnoreError(num string) decimal.Decimal {
	d, _ := decimal.NewFromString(num)
	return d
}
