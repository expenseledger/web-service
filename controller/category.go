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
	if _, err = model.DeleteCategory(category); err != nil {
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
	categories, err := model.ListCategories()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: len(categories),
		Items:  categories,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}

func categoryInit(context *gin.Context) {
	recipes := []CategoryIDForm{
		CategoryIDForm{
			Name: "Food And Drink",
		},
		CategoryIDForm{
			Name: "Transportation",
		},
		CategoryIDForm{
			Name: "Shopping",
		},
		CategoryIDForm{
			Name: "Bill",
		},
		CategoryIDForm{
			Name: "Withdraw",
		},
	}

	length := len(recipes)
	categories := make([]model.Category, length)
	for i := 0; i < length; i++ {
		category, err := MakeCategory(&recipes[i], Create)
		if err != nil {
			context.JSON(
				http.StatusBadRequest,
				buildNonsuccessResponse(err, nil),
			)
			return
		}
		categories[i] = *category.(*model.Category)
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
	categories, err := model.ClearCategories()
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	items := itemList{
		Length: len(categories),
		Items:  categories,
	}

	context.JSON(
		http.StatusOK,
		buildSuccessResponse(items),
	)
	return
}
