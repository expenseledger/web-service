package model

import (
	"reflect"

	"github.com/expenseledger/web-service/orm"
)

// Category the structure represents a category in presentation layer
type Category struct {
	Name string `json:"name"`
}

// Categories is defined just to be used as a receiver
type Categories []Category

// Create ...
func (category *Category) Create() error {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(*category))

	tmp, err := mapper.Insert(category)
	if err != nil {
		return err
	}

	c := tmp.(*Category)
	*category = *c

	return nil
}

// Get ...
func (category *Category) Get() error {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(*category))

	tmp, err := mapper.One(category)
	if err != nil {
		return err
	}

	c := tmp.(*Category)
	*category = *c

	return nil
}

// Delete ...
func (category *Category) Delete() error {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(*category))

	tmp, err := mapper.Delete(category)
	if err != nil {
		return err
	}

	c := tmp.(*Category)
	*category = *c

	return nil
}

// List ...
func (categories *Categories) List() (int, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(Category{}))

	tmp, err := mapper.Many()
	if err != nil {
		return 0, err
	}

	c := tmp.(*[]Category)
	*categories = *c

	return len(*categories), nil
}

// Init insert default categories
func (categories *Categories) Init() (int, error) {
	for _, c := range *categories {
		err := c.Create()
		if err != nil {
			return 0, err
		}
	}
	return len(*categories), nil
}

// Clear ...
func (categories *Categories) Clear() (int, error) {
	mapper := orm.NewCategoryMapper(reflect.TypeOf(Category{}))

	tmp, err := mapper.Clear()
	if err != nil {
		return 0, err
	}

	c := tmp.(*[]Category)
	*categories = *c

	return len(*categories), nil
}
