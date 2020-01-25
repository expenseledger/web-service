package controller

import (
	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
)

type categoryIDForm struct {
	Name string `json:"name" binding:"required"`
}

func createCategory(context *gin.Context) {
	var form categoryIDForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	category, err := model.CreateCategory(form.Name)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
}

func getCategory(context *gin.Context) {
	var form categoryIDForm
	if err := bindJSON(context, &form); err != nil {
		return
	}

	category, err := model.GetCategory(form.Name)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
}

func deleteCategory(context *gin.Context) {
	var form categoryIDForm
	if err := bindJSON(context, &form); err != nil {
		buildFailedContext(context, err)
		return
	}

	category, err := model.DeleteCategory(form.Name)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
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
}

func initCategories(context *gin.Context) {
	names := []string{
		"Food And Drink",
		"Transportation",
		"Shopping",
		"Bill",
		"Withdraw",
	}

	length := len(names)
	categories := make([]*model.Category, length)
	for i, name := range names {
		category, err := model.CreateCategory(name)
		if err != nil {
			buildFailedContext(context, err)
			return
		}
		categories[i] = category
	}

	items := itemList{
		Length: length,
		Items:  categories,
	}

	buildSuccessContext(context, items)
}

func clearCategories(context *gin.Context) {
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
}
