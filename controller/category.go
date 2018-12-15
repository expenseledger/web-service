package controller

import (
	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

func createCategory(context *gin.Context) {
	var form CategoryIDForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	category, err := MakeCategory(&form, Create)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
	return
}

func getCategory(context *gin.Context) {
	var form CategoryIDForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	category, err := MakeCategory(&form, Get)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
	return
}

func deleteCategory(context *gin.Context) {
	var form CategoryIDForm
	if err := bindJSON(context, &form); err != nil {
		buildFailedContext(context, err)
		return
	}

	iCategory, err := MakeCategory(&form, Get)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	category := iCategory.(*model.Category)
	if _, err = model.DeleteCategory(category); err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
	return
}

func listCategories(context *gin.Context) {
	categories, err := model.ListCategories()
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(categories),
		Items:  categories,
	}

	buildSuccessContext(context, items)
	return
}

func initCategories(context *gin.Context) {
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
			buildFailedContext(context, err)
			return
		}
		categories[i] = *category.(*model.Category)
	}

	items := itemList{
		Length: length,
		Items:  categories,
	}

	buildSuccessContext(context, items)
	return
}

func categoryClear(context *gin.Context) {
	categories, err := model.ClearCategories()
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	items := itemList{
		Length: len(categories),
		Items:  categories,
	}

	buildSuccessContext(context, items)
	return
}
