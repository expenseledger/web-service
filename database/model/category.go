package model

import (
	"log"
	"time"

	"github.com/expenseledger/web-service/database"
)

// Category the structure represents a stored category on database
type Category struct {
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Categories is defined just to be used as a receiver
type Categories []Category

// Insert ...
func (category *Category) Insert() error {
	query :=
		`
		INSERT INTO category (name)
		VALUES (:name)
		RETURNING *;
		`

	stmt, err := database.DB().PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a category", err)
		return err
	}

	if err := stmt.Get(category, category); err != nil {
		log.Println("Error inserting a category", err)
		return err
	}

	return nil
}

// One ...
func (category *Category) One(name string) error {
	query :=
		`
		SELECT * FROM category
		WHERE name=$1;
		`

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error selecting a category", err)
		return err
	}

	if err := stmt.Get(category, name); err != nil {
		log.Println("Error selecting a category", err)
		return err
	}

	return nil
}

// Delete ...
func (category *Category) Delete(name string) error {
	query :=
		`
		DELETE FROM category
		WHERE name=$1
		RETURNING *;
		`

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error deleting a category", err)
		return err
	}

	if err := stmt.Get(category, name); err != nil {
		log.Println("Error deleting a category", err)
		return err
	}

	return nil
}

// All gets all categories
func (categories *Categories) All() (int, error) {
	query :=
		`
		SELECT * FROM category;
		`

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error selecting all categories", err)
		return 0, err
	}

	if err := stmt.Select(categories); err != nil {
		log.Println("Error selecting all categories", err)
		return 0, err
	}

	return len(*categories), nil
}

// BatchInsert ...
func (categories *Categories) BatchInsert() (int, error) {
	var err error
	length := len(*categories)
	insertedCategories := make(Categories, length)
	for i, category := range *categories {
		err = category.Insert()
		if err != nil {
			break
		}
		insertedCategories[i] = category
	}

	if err != nil {
		log.Println("Error doing batch insertion categories", err)
		return 0, err
	}

	*categories = insertedCategories
	return length, nil
}

// DeleteAll ...
func (categories *Categories) DeleteAll() (int, error) {
	query :=
		`
		DELETE FROM category
		RETURNING *;
		`

	stmt, err := database.DB().Preparex(query)
	if err != nil {
		log.Println("Error deleting all categories", err)
		return 0, err
	}

	if err := stmt.Select(categories); err != nil {
		log.Println("Error deleting all categories", err)
		return 0, err
	}

	return len(*categories), nil
}
