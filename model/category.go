package model

import (
	"github.com/expenseledger/web-service/orm"
)

// Category the structure represents a category in presentation layer
type Category struct {
	Name string `json:"name" db:"name"`
}

type operation int

const (
	insert operation = iota
	delete
	one
	list
	clear
)

// CreateCategory inserts category to DB
func CreateCategory(name string) (*Category, error) {
	return applyToCategory(name, insert)
}

// GetCategory inserts category to DB
func GetCategory(name string) (*Category, error) {
	return applyToCategory(name, one)
}

// DeleteCategory remove category from DB
func DeleteCategory(name string) (*Category, error) {
	return applyToCategory(name, delete)
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
