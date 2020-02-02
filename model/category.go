package model

import (
	"github.com/expenseledger/web-service/orm"
)

// Category the structure represents a category in presentation layer
type Category struct {
	Name   string `json:"name" db:"name"`
	UserId string `json:"userId" db:"user_id"`
}

// CreateCategory inserts category to DB
func CreateCategory(name string, userId string) (*Category, error) {
	return applyToCategory(name, insert, userId)
}

// GetCategory returns matching category from DB
func GetCategory(name string, userId string) (*Category, error) {
	return applyToCategory(name, one, userId)
}

// DeleteCategory removes category from DB
func DeleteCategory(name string, userId string) (*Category, error) {
	return applyToCategory(name, delete, userId)
}

// ListCategories ...
func ListCategories(userId string) ([]Category, error) {
	return applyToCategories(list, userId)
}

// ClearCategories ...
func ClearCategories(userId string) ([]Category, error) {
	return applyToCategories(clear, userId)
}

func applyToCategory(name string, op operation, userId string) (*Category, error) {
	c := Category{Name: name, UserId: userId}
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

func applyToCategories(op operation, userId string) ([]Category, error) {
	mapper := orm.NewCategoryMapper(Category{UserId: userId})

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
