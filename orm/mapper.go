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

var baseMapper BaseMapper

// NewCategoryMapper ...
func NewCategoryMapper(modelType reflect.Type) Mapper {
	baseMapper.once.Do(func() {
		baseMapper.insertStmt = `
			INSERT INTO category (name)
			VALUES (:name)
			RETURNING name;
		`
		baseMapper.deleteStmt = `
			DELETE FROM category
			WHERE name=:name
			RETURNING name;
		`
		baseMapper.oneStmt = `
			SELECT name FROM category
			WHERE name=:name;
		`
		baseMapper.manyStmt = `
			SELECT name FROM category;
		`
		baseMapper.clearStmt = `
			DELETE FROM category
			RETURNING name;
		`
	})

	baseMapper.modelType = modelType

	return &baseMapper
}
