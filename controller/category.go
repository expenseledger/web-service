package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

func categoryCreate(context *gin.Context) {
	var form CategoryIDForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	category, err := MakeCategory(&form, Create)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(category),
	)
	return
}

func categoryGet(context *gin.Context) {
	var form CategoryIDForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	category, err := MakeCategory(&form, Get)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(category),
	)
	return
}

func categoryDelete(context *gin.Context) {
	var form CategoryIDForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	iCategory, err := MakeCategory(&form, Get)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	category := iCategory.(*model.Category)
	if err = category.Delete(); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(category),
	)
	return
}

func categoryList(context *gin.Context) {
	var categories model.Categories
	length, err := categories.List()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  categories,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func categoryInit(context *gin.Context) {
	categories := model.Categories{
		model.Category{
			Name: "Food And Drink",
		},
		model.Category{
			Name: "Transportation",
		},
		model.Category{
			Name: "Shopping",
		},
		model.Category{
			Name: "Bill",
		},
		model.Category{
			Name: "Withdraw",
		},
	}

	length, err := categories.Init()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  categories,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func categoryClear(context *gin.Context) {
	var categories model.Categories
	length, err := categories.Clear()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: length,
		Items:  categories,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}
