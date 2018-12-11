package controller

import (
	"reflect"

	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/orm"
)

type CategoryIDForm struct {
	Name string `json:"name" binding:"required"`
}

func (form *CategoryIDForm) Model() interface{} {
	return &model.Category{
		Name: form.Name,
	}
}

func MakeCategory(recipe Modeler, intent Intent) (interface{}, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(model.Category{}))
	category := recipe.Model()
	return makeModel(mapper, category, intent)
}
