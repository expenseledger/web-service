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

var (
	categoryMapper BaseMapper
	walletMapper   BaseMapper
)

// NewCategoryMapper ...
func NewCategoryMapper(model interface{}) Mapper {
	categoryMapper.once.Do(func() {
		categoryMapper.insertStmt = `
			INSERT INTO category (name)
			VALUES (:name)
			RETURNING name;
		`
		categoryMapper.deleteStmt = `
			DELETE FROM category
			WHERE name=:name
			RETURNING name;
		`
		categoryMapper.oneStmt = `
			SELECT name FROM category
			WHERE name=:name;
		`
		categoryMapper.manyStmt = `
			SELECT name FROM category;
		`
		categoryMapper.clearStmt = `
			DELETE FROM category
			RETURNING name;
		`
	})

	categoryMapper.modelType = reflect.TypeOf(model)

	return &categoryMapper
}

// NewWalletMapper ...
func NewWalletMapper(model interface{}) Mapper {
	walletMapper.once.Do(func() {
		walletMapper.insertStmt = `
			INSERT INTO wallet (name, type, balance)
			VALUES (:name, :type, :balance)
			RETURNING name, type, balance;
		`
		walletMapper.deleteStmt = `
			DELETE FROM wallet
			WHERE name=:name
			RETURNING name, type, balance;
		`
		walletMapper.oneStmt = `
			SELECT name, type, balance FROM wallet
			WHERE name=:name;
		`
		walletMapper.manyStmt = `
			SELECT name, type, balance FROM wallet;
		`
		walletMapper.clearStmt = `
			DELETE FROM wallet
			RETURNING name, type, balance;
		`
	})

	walletMapper.modelType = reflect.TypeOf(model)

	return &walletMapper
}
