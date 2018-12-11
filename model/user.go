package model

import (
	"reflect"

	"github.com/expenseledger/web-service/orm"
)

// DeleteCategory ...
func DeleteCategory(category *Category) (*Category, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(*category))

	tmp, err := mapper.Delete(category)
	if err != nil {
		return nil, err
	}

	return tmp.(*Category), nil
}

// ListCategories ...
func ListCategories() ([]Category, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(Category{}))

	tmp, err := mapper.Many()
	if err != nil {
		return nil, err
	}

	c := tmp.(*[]Category)
	categories := *c

	return categories, nil
}

// ClearCategories ...
func ClearCategories() ([]Category, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(Category{}))

	tmp, err := mapper.Clear()
	if err != nil {
		return nil, err
	}

	c := tmp.(*[]Category)
	categories := *c

	return categories, nil
}
