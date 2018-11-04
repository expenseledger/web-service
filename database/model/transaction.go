package model

import (
	"log"
	"time"

	"github.com/shopspring/decimal"
)

// Transaction the structure represents a stored transaction on database
type Transaction struct {
	ID          string          `db:"id"`
	SrcWallet   string          `db:"src_wallet"`
	DstWallet   *string         `db:"dst_wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    *string         `db:"category"`
	Description string          `db:"description"`
	OccuredAt   *time.Time      `db:"occured_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Insert ...
func (tx *Transaction) Insert() error {
	query :=
		`
		INSERT INTO transaction
		(src_wallet, amount, type, category, description)
		VALUES
		(:src_wallet, :amount, :type, :category, :description)
		RETURNING *;
		`

	if tx.OccuredAt != nil {
		query =
			`
			INSERT INTO transaction
			(src_wallet, amount, type, category, description, occured_at)
			VALUES
			(:src_wallet, :amount, :type, :category, :description, :occured_at)
			RETURNING *;
			`
	}

	namedStmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	if err := namedStmt.Get(tx, tx); err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	return nil
}

// One ...
func (tx *Transaction) One() error {
	query :=
		`
		SELECT * FROM transaction
		WHERE id=:id;
		`
	namedStmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error selecting a transaction", err)
		return err
	}

	if err := namedStmt.Get(tx, tx); err != nil {
		log.Println("Error selecting a transaction", err)
		return err
	}

	return nil
}

// DeleteAll ...
func (txs *Transactions) DeleteAll() (int, error) {
	query :=
		`
		DELETE FROM transaction
		RETURNING *;
		`

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	if err := stmt.Select(txs); err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	return len(*txs), nil
}
