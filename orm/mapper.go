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
			INSERT INTO category (name, user_id)
			VALUES (:name, :user_id)
			RETURNING name, user_id;
		`
		categoryMapper.deleteStmt = `
			DELETE FROM category
			WHERE name=:name
			AND user_id=:user_id
			RETURNING name, user_id;
		`
		categoryMapper.oneStmt = `
			SELECT name, user_id
			FROM category
			WHERE name=:name
			AND user_id=:user_id;
		`
		categoryMapper.manyStmt = `
			SELECT name, user_id
			FROM category
			WHERE user_id=:user_id;
		`
		categoryMapper.clearStmt = `
			DELETE FROM category
			WHERE user_id=:user_id
			RETURNING name, user_id;
		`
	})

	categoryMapper.modelType = reflect.TypeOf(model)

	return &categoryMapper
}

func NewWalletMapper(model interface{}) Mapper {
	walletMapper.once.Do(func() {
		walletMapper.insertStmt = `
			INSERT INTO wallet (name, type, balance, user_id)
			VALUES (:name, :type, :balance, :user_id)
			RETURNING name, type, balance, user_id;
		`
		walletMapper.deleteStmt = `
			DELETE FROM wallet
			WHERE name=:name
			AND user_id=:user_id
			RETURNING name, type, balance, user_id;
		`
		walletMapper.oneStmt = `
			SELECT name, type, balance, user_id
			FROM wallet
			WHERE name=:name
			AND user_id=:user_id;
		`
		walletMapper.updateStmt = `
			UPDATE wallet
			SET balance=:balance
			WHERE name=:name
			AND user_id=:user_id
			RETURNING name, type, balance, user_id;
		`
		walletMapper.manyStmt = `
			SELECT name, type, balance, user_id 
			FROM wallet
			WHERE user_id=:user_id;
		`
		walletMapper.clearStmt = `
			DELETE FROM wallet
			WHERE user_id=:user_id
			RETURNING name, type, balance, user_id;
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
				(amount, type, category, description, occurred_at, user_id)
				VALUES
				(:amount, :type, :category, :description, :occurred_at, :user_id)
				RETURNING id, amount, type, category, description, occurred_at
			), tx_wallet AS (
				INSERT INTO affected_wallet
				(transaction_id, wallet, role, user_id)
				SELECT id, :wallet, :role, :user_id FROM tx
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
				(amount, type, category, description, occurred_at, user_id)
				VALUES
				(:amount, :type, :category, :description, :occurred_at, :user_id)
				RETURNING id, amount, type, category, description, occurred_at, user_id
			), tx_wallet AS (
				INSERT INTO affected_wallet
				(transaction_id, wallet, role, user_id)
				SELECT
				id, :src_wallet, CAST ('SRC_WALLET' AS wallet_role) FROM tx
				UNION ALL
				SELECT
				id, :dst_wallet, CAST ('DST_WALLET' AS wallet_role) FROM tx
				RETURNING wallet, role
			)
			SELECT
			id, w1.wallet AS src_wallet, w2.wallet AS dst_wallet,
			amount, type, category, description, occurred_at, user_id
			FROM tx, tx_wallet w1, tx_wallet w2
			WHERE w1.role = 'SRC_WALLET' AND w2.role = 'DST_WALLET';
		`
		txMapper.deleteStmt = `
			WITH tx AS (
				DELETE FROM transaction
				WHERE id = :id
				AND user_id = :user_id
				RETURNING id, amount, type, category, description, occurred_at, user_id
			), tx_wallet AS (
				DELETE FROM affected_wallet
				WHERE transaction_id = :id
				AND user_id = :user_id
				RETURNING transaction_id, wallet, role
			)
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at, user_id
			FROM tx, tx_wallet
			WHERE tx.id = tx_wallet.transaction_id
			ORDER BY role ASC;
		`
		txMapper.oneStmt = `
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at, t.user_id
			FROM transaction t, affected_wallet w
			WHERE t.id = :id AND t.id = w.transaction_id
			AND t.user_id = :user_id AND t.user_id = w.user_id
			ORDER BY role ASC;
		`
		txMapper.manyStmt = `
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at, user_id
			FROM transaction t, affected_wallet w
			WHERE t.id IN (
				SELECT transaction_id FROM affected_wallet
				WHERE wallet = :wallet
				AND user_id = :user_id
			) AND t.id = w.transaction_id 
			AND user_id = :user_id
			ORDER BY occurred_at ASC, w.created_at ASC, role ASC;
		`
		txMapper.clearStmt = `
			WITH tx AS (
				DELETE FROM transaction
				WHERE user_id = :user_id
				RETURNING id, amount, type, category, description, occurred_at, user_id
			), tx_wallet AS (
				DELETE FROM affected_wallet
				WHERE user_id = :user_id
				RETURNING transaction_id, wallet, role, created_at
			)
			SELECT
			id, wallet, role, amount, type, category, description, occurred_at, user_id
			FROM transaction t, affected_wallet w
			WHERE t.id = w.transaction_id
			ORDER BY occurred_at ASC, w.created_at ASC, role ASC;
		`
	})

	txMapper.modelType = reflect.TypeOf(model)
	txMapper.txType = txType

	return &txMapper
}
