package controller

import (
	"net/http"

	"github.com/expenseledger/web-service/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type categoryIdentifyForm struct {
	Name string `json:"name" binding:"required"`
}

func categoryCreate(context *gin.Context) {
	var form categoryIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var category model.Category
	copier.Copy(&category, &form)

	if err := category.Create(); err != nil {
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
	var form categoryIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var category model.Category
	if err := category.Get(form.Name); err != nil {
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
	var form categoryIdentifyForm
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(
			http.StatusBadRequest,
			buildNonsuccessResponse(err, nil),
		)
		return
	}

	var category model.Category
	if err := category.Delete(form.Name); err != nil {
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
