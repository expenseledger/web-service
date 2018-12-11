package orm

import "reflect"

// Mapper defines the relational operations
type Mapper interface {
	Insert(obj interface{}) (interface{}, error)
	Delete(obj interface{}) (interface{}, error)
	One(obj interface{}) (interface{}, error)
	Many() (interface{}, error)
	Clear() (interface{}, error)
}

// NewCategoryMapper ...
func NewCategoryMapper(modelType reflect.Type) Mapper {
	return &BaseMapper{
		modelType: modelType,
		insertStmt: `
			INSERT INTO category (name)
			VALUES (:name)
			RETURNING name;
		`,
		deleteStmt: `
			DELETE FROM category
			WHERE name=:name
			RETURNING name;
		`,
		oneStmt: `
			SELECT name FROM category
			WHERE name=:name;
		`,
		manyStmt: `
			SELECT name FROM category;
		`,
		clearStmt: `
			DELETE FROM category
			RETURNING name;
		`,
	}
}
