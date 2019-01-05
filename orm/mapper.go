package orm

import (
	"reflect"

	"github.com/expenseledger/web-service/constant"
)

type Mapper interface {
	Insert(obj interface{}) (interface{}, error)
	Delete(obj interface{}) (interface{}, error)
	One(obj interface{}) (interface{}, error)
	Update(obj interface{}) (interface{}, error)
	Many(obj interface{}) (interface{}, error)
	Clear() (interface{}, error)
}

var (
	categoryMapper BaseMapper
	walletMapper   BaseMapper
	txMapper       TxMapper
)

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
		walletMapper.updateStmt = `
			UPDATE wallet
			SET balance=:balance
			WHERE name=:name
			RETURNING name, type, balance;
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

func NewTxMapper(model interface{}, txType constant.TransactionType) Mapper {
	txMapper.once.Do(func() {
		txMapper.insertStmt = `
			WITH tx AS (
				INSERT INTO transaction
				(amount, type, category, description, occurred_at)
				VALUES
				(:amount, :type, :category, :description, :occurred_at)
				RETURNING id, amount, type, category, description, occurred_at
			), tx_wallet AS (
				INSERT INTO affected_wallet
				(transaction_id, wallet, role)
				SELECT id, :wallet, :role FROM tx
				RETURNING wallet, role
			)
			SELECT
			tx.id AS id, amount, type, category,
			description, occurred_at, wallet, role
			FROM tx, tx_wallet;
		`
		txMapper.transferStmt = `
			WITH tx AS (
				INSERT INTO transaction
				(amount, type, category, description, occurred_at)
				VALUES
				(:amount, :type, :category, :description, :occurred_at)
				RETURNING id, amount, type, category, description, occurred_at
			), tx_wallet AS (
				INSERT INTO affected_wallet
				(transaction_id, wallet, role)
				SELECT
				id, :src_wallet, CAST ('SRC_WALLET' AS wallet_role) FROM tx
				UNION ALL
				SELECT
				id, :dst_wallet, CAST ('DST_WALLET' AS wallet_role) FROM tx
				RETURNING wallet, role
			)
			SELECT
			id, w1.wallet AS src_wallet, w2.wallet AS dst_wallet,
			amount, type, category, description, occurred_at
			FROM tx, tx_wallet w1, tx_wallet w2
			WHERE w1.role = 'SRC_WALLET' AND w2.role = 'DST_WALLET';
		`
		txMapper.deleteStmt = `
			WITH tx AS (
				DELETE FROM transaction
				WHERE id = :id
				RETURNING id, amount, type, category, description, occurred_at
			), tx_wallet AS (
				DELETE FROM affected_wallet
				WHERE transaction_id = :id
				RETURNING transaction_id, wallet, role
			)
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at
			FROM tx, tx_wallet
			WHERE tx.id = tx_wallet.transaction_id
			ORDER BY role ASC;
		`
		txMapper.oneStmt = `
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at
			FROM transaction t, affected_wallet w
			WHERE t.id = :id AND t.id = w.transaction_id
			ORDER BY role ASC;
		`
		txMapper.manyStmt = `
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at
			FROM transaction t, affected_wallet w
			WHERE w.wallet = :wallet AND t.id = w.transaction_id
			ORDER BY occurred_at ASC, w.created_at ASC, role ASC;
		`
		txMapper.clearStmt = `
			WITH tx AS (
				DELETE FROM transaction
				RETURNING id, amount, type, category, description, occurred_at
			), tx_wallet AS (
				DELETE FROM affected_wallet
				RETURNING transaction_id, wallet, role, created_at
			)
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at
			FROM transaction t, affected_wallet w
			WHERE t.id = w.transaction_id
			ORDER BY occurred_at ASC, w.created_at ASC, role ASC;
		`
	})

	txMapper.modelType = reflect.TypeOf(model)
	txMapper.txType = txType

	return &txMapper
}
