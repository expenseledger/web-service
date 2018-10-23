package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
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

	var wallet model.Wallet
	if err := wallet.Get(form.Name); err != nil {
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

func decimalFromStringIgnoreError(num string) decimal.Decimal {
	d, _ := decimal.NewFromString(num)
	return d
}
