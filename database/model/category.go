package dbmodel

import (
	"log"
	"time"

	"github.com/expenseledger/web-service/database"
)

// Category the structure represents a stored category on database
type Category struct {
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
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
	db := database.GetDB()

	stmt, err := db.PrepareNamed(query)
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
		WHERE name=$1 AND deleted_at IS NULL;
		`
	db := database.GetDB()

	stmt, err := db.Preparex(query)
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
		UPDATE category
		SET deleted_at=now()
		WHERE name=$1 AND deleted_at IS NULL
		RETURNING *;
		`
	db := database.GetDB()

	stmt, err := db.Preparex(query)
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
		SELECT * FROM category
		WHERE deleted_at IS NULL;
		`
	db := database.GetDB()

	stmt, err := db.Preparex(query)
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
	for index, category := range *categories {
		err = category.Insert()
		if err != nil {
			break
		}
		[]Category(*categories)[index] = category
	}

	if err != nil {
		log.Println("Error doing batch insertion categories", err)
		return 0, err
	}

	return len(*categories), nil
}
