package controller

import (
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/pkg"
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

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	category, err := model.CreateCategory(form.Name, userId)
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

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	category, err := model.GetCategory(form.Name, userId)
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

	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	category, err := model.DeleteCategory(form.Name, userId)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	buildSuccessContext(context, category)
}

func listCategories(context *gin.Context) {
	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	categories, err := model.ListCategories(userId)
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
	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

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
		category, err := model.CreateCategory(name, userId)
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
	userId, err := pkg.GetUserId(context)
	if err != nil {
		buildFailedContext(context, err)
		return
	}

	categories, err := model.ClearCategories(userId)
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
