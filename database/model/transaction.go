package model

import (
	"fmt"
	"log"
	"time"

	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/database"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a stored transaction on database
type Transaction struct {
	ID          string          `db:"id"`
	SrcWallet   string          `db:"src_wallet"`
	DstWallet   string          `db:"dst_wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    string          `db:"category"`
	Description string          `db:"description"`
	OccurredAt  time.Time       `db:"occurred_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

type affectedWallet struct {
	TransactionID string    `db:"transaction_id"`
	Wallet        string    `db:"wallet"`
	Role          string    `db:"role"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Insert ...
func (tx *Transaction) Insert() error {
	query := tx.buildInsertSQLStmt()

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	walletRole := constant.WalletRole()
	transactionType := constant.TransactionType()
	switch {
	case !tx.OccurredAt.IsZero() &&
		tx.Type == transactionType.Transfer:
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.OccurredAt,
			tx.SrcWallet,
			walletRole.SrcWallet,
			tx.DstWallet,
			walletRole.DstWallet,
		)

	case !tx.OccurredAt.IsZero():
		wallet, role := tx.SrcWallet, walletRole.SrcWallet
		if tx.Type == transactionType.Income {
			wallet, role = tx.DstWallet, walletRole.DstWallet
		}
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.OccurredAt,
			wallet,
			role,
		)

	case tx.Type == constant.TransactionType().Transfer:
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.SrcWallet,
			walletRole.SrcWallet,
			tx.DstWallet,
			walletRole.DstWallet,
		)

	default:
		wallet, role := tx.SrcWallet, walletRole.SrcWallet
		if tx.Type == transactionType.Income {
			wallet, role = tx.DstWallet, walletRole.DstWallet
		}
		err = stmt.Get(
			tx,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			wallet,
			role,
		)
	}

	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}
	return nil
}

// DeleteAll ...
func (txs *Transactions) DeleteAll() error {
	query := fmt.Sprintf(
		`
		WITH deleted_transfer AS (
			DELETE FROM %s
			WHERE type='TRANSFER'
			RETURNING *
		), deleted_income AS (
			DELETE FROM %s
			WHERE type='INCOME'
			RETURNING *
		), deleted_expense AS (
			DELETE FROM %s
			WHERE type='EXPENSE'
			RETURNING *
		), deleted_tx_wallet AS (
			DELETE FROM %s
			RETURNING *
		)
		SELECT
			t.id AS id,
			w1.wallet AS src_wallet,
			w2.wallet AS dst_wallet,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM deleted_transfer t, deleted_tx_wallet w1, deleted_tx_wallet w2
			WHERE t.id=w1.transaction_id AND t.id=w2.transaction_id AND
			w1.wallet<>w2.wallet AND w1.role='SRC_WALLET'
		UNION ALL
		SELECT
			t.id AS id,
			'' AS src_wallet,
			w.wallet AS dst_wallet,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM deleted_income t, deleted_tx_wallet w
			WHERE t.id=w.transaction_id
		UNION ALL
		SELECT
			t.id AS id,
			w.wallet AS src_wallet,
			'' AS dst_wallet,
			t.amount AS amount,
			t.type AS type,
			t.category AS category,
			t.description AS description,
			t.occurred_at AS occurred_at
			FROM deleted_expense t, deleted_tx_wallet w
			WHERE t.id=w.transaction_id;
		`,
		database.Transaction,
		database.Transaction,
		database.Transaction,
		database.AffectedWallet,
	)

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Println("Error deleting all transactions", err)
		return err
	}

	if err := stmt.Select(txs); err != nil {
		log.Println("Error deleting all transactions", err)
		return err
	}

	return nil
}

func (tx *Transaction) buildInsertSQLStmt() string {
	var query string

	switch {
	case !tx.OccurredAt.IsZero() &&
		tx.Type == constant.TransactionType().Transfer:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description, occurred_at)
				VALUES
				($1, $2, $3, $4, $5)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $6, CAST ($7 AS %s) FROM inserted_tx
				UNION ALL
				SELECT id, $8, CAST ($9 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
			database.WalletRole,
		)

	case !tx.OccurredAt.IsZero():
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description, occurred_at)
				VALUES
				($1, $2, $3, $4, $5)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $6, CAST ($7 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
		)

	case tx.Type == constant.TransactionType().Transfer:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description)
				VALUES
				($1, $2, $3, $4)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $5, CAST ($6 AS %s) FROM inserted_tx
				UNION ALL
				SELECT id, $7, CAST ($8 AS %s) FROM inserted_tx
				RETURNING *
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
			database.WalletRole,
		)

	default:
		query = fmt.Sprintf(
			`
			WITH inserted_tx AS (
				INSERT INTO %s
				(amount, type, category, description)
				VALUES
				($1, $2, $3, $4)
				RETURNING *
			), tx_wallet AS (
				INSERT INTO %s
				(transaction_id, wallet, role)
				SELECT id, $5, CAST ($6 AS %s) FROM inserted_tx
				RETURNING transaction_id, wallet, role
			)
			SELECT id, occurred_at FROM inserted_tx;
			`,
			database.Transaction,
			database.AffectedWallet,
			database.WalletRole,
		)
	}

	return query
}
