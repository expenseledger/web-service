package model

import (
	dbmodel "github.com/expenseledger/web-service/database/model"
	"github.com/jinzhu/copier"
)

// Category the structure represents a category in presentation layer
type Category struct {
	Name string `json:"name"`
}

// Categories is defined just to be used as a receiver
type Categories []Category

// Create ...
func (category *Category) Create() error {
	var dbCategory dbmodel.Category

	copier.Copy(&dbCategory, &category)

	if err := dbCategory.Insert(); err != nil {
		return err
	}

	copier.Copy(category, &dbCategory)
	return nil
}

// Get ...
func (category *Category) Get(name string) error {
	var dbCategory dbmodel.Category
	if err := dbCategory.One(name); err != nil {
		return err
	}

	copier.Copy(category, &dbCategory)
	return nil
}

// Delete ...
func (category *Category) Delete(name string) error {
	var dbCategory dbmodel.Category
	if err := dbCategory.Delete(name); err != nil {
		return err
	}

	copier.Copy(category, &dbCategory)
	return nil
}

// List ...
func (categories *Categories) List() (int, error) {
	var dbCategories dbmodel.Categories

	length, err := dbCategories.All()
	if err != nil {
		return 0, err
	}

	copier.Copy(categories, &dbCategories)
	return length, nil
}

// Init insert default categories
func (categories *Categories) Init() (int, error) {
	var dbCategories dbmodel.Categories
	copier.Copy(&dbCategories, categories)

	length, err := dbCategories.BatchInsert()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Clear ...
func (categories *Categories) Clear() (int, error) {
	var dbCategories dbmodel.Categories

	length, err := dbCategories.DeleteAll()
	if err != nil {
		return 0, err
	}

	copier.Copy(categories, &dbCategories)
	return length, nil
}
