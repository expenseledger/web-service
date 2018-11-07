package model

import (
	"errors"
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

	var results []affectedWallet
	switch {
	case !tx.OccurredAt.IsZero() &&
		tx.Type == constant.TransactionType().Transfer:
		err = stmt.Select(
			&results,
			tx.Amount,
			tx.Type,
			tx.Category,
			tx.Description,
			tx.OccurredAt,
			tx.SrcWallet,
			constant.WalletRole().SrcWallet,
			tx.DstWallet,
			constant.WalletRole().DstWallet,
		)
	default:
		err = errors.New("Nahnah")
	}

	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	tx.ID = results[0].TransactionID
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
			)
			INSERT INTO %s
			(transaction_id, wallet, role)
			SELECT id, $6, CAST ($7 AS wallet_role) FROM inserted_tx
			UNION ALL
			SELECT id, $8, CAST ($9 AS wallet_role) FROM inserted_tx
			RETURNING *;
			`,
			database.Transaction,
			database.AffectedWallet,
		)
	}

	return query
}
