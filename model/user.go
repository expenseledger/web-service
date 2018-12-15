package model

import (
	"github.com/expenseledger/web-service/orm"
)

// ListCategories ...
func ListCategories() ([]Category, error) {
	return applyToCategories(list)
}

// ClearCategories ...
func ClearCategories() ([]Category, error) {
	return applyToCategories(clear)
}

func applyToCategories(op operation) ([]Category, error) {
	mapper := orm.NewCategoryMapper(Category{})

	var tmp interface{}
	var err error
	switch op {
	case list:
		tmp, err = mapper.Many()
	case clear:
		tmp, err = mapper.Clear()
	}

	if err != nil {
		return nil, err
	}

	c := tmp.(*[]Category)
	categories := *c

	return categories, nil
}
