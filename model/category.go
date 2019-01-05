package model

import (
	"github.com/expenseledger/web-service/orm"
)

// Category the structure represents a category in presentation layer
type Category struct {
	Name string `json:"name" db:"name"`
}

// CreateCategory inserts category to DB
func CreateCategory(name string) (*Category, error) {
	return applyToCategory(name, insert)
}

// GetCategory returns matching category from DB
func GetCategory(name string) (*Category, error) {
	return applyToCategory(name, one)
}

// DeleteCategory removes category from DB
func DeleteCategory(name string) (*Category, error) {
	return applyToCategory(name, delete)
}

// ListCategories ...
func ListCategories() ([]Category, error) {
	return applyToCategories(list)
}

// ClearCategories ...
func ClearCategories() ([]Category, error) {
	return applyToCategories(clear)
}

func applyToCategory(name string, op operation) (*Category, error) {
	c := Category{Name: name}
	mapper := orm.NewCategoryMapper(c)

	var tmp interface{}
	var err error
	switch op {
	case insert:
		tmp, err = mapper.Insert(&c)
	case delete:
		tmp, err = mapper.Delete(&c)
	case one:
		tmp, err = mapper.One(&c)
	}

	if err != nil {
		return nil, err
	}

	return tmp.(*Category), nil
}

func applyToCategories(op operation) ([]Category, error) {
	mapper := orm.NewCategoryMapper(Category{})

	var tmp interface{}
	var err error
	switch op {
	case list:
		tmp, err = mapper.Many(&struct{}{})
	case clear:
		tmp, err = mapper.Clear()
	}

	if err != nil {
		return nil, err
	}

	categories := *(tmp.(*[]Category))
	return categories, nil
}
