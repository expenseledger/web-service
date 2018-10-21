package controller

import (
	"net/http"

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
