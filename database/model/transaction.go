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
	Category    string          `db:"category"`
	Description string          `db:"description"`
	OccuredAt   *time.Time      `db:"occured_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Insert ...
func (tx *Transaction) Insert() error {
	var query string

	switch {
	case tx.OccuredAt != nil && tx.DstWallet != nil:
		query =
			`
			INSERT INTO transaction
			(src_wallet, dst_wallet, amount, type, category, description, occured_at)
			SELECT
			w1.name, w2.name, :amount, :type, c.name, :description, :occured_at
			FROM wallet w1, wallet w2, category c
			WHERE
			w1.name=:src_wallet AND	w2.name=:dst_wallet AND	c.name=:category
			RETURNING *;
			`
	case tx.OccuredAt == nil && tx.DstWallet != nil:
		query =
			`
			INSERT INTO transaction
			(src_wallet, dst_wallet, amount, type, category, description)
			SELECT
			w1.name, w2.name, :amount, :type, c.name, :description
			FROM wallet w1, wallet w2, category c
			WHERE
			w1.name=:src_wallet AND	w2.name=:dst_wallet AND	c.name=:category
			RETURNING *;
		`
	case tx.OccuredAt != nil && tx.DstWallet == nil:
		query =
			`
			INSERT INTO transaction
			(src_wallet, amount, type, category, description, occured_at)
			SELECT
			w.name, :amount, :type, c.name, :description, :occured_at
			FROM wallet w, category c
			WHERE
			w.name=:src_wallet AND c.name=:category
			RETURNING *;
			`
	case tx.OccuredAt == nil && tx.DstWallet == nil:
		query =
			`
			INSERT INTO transaction
			(src_wallet, amount, type, category, description)
			SELECT
			w.name, :amount, :type, c.name, :description
			FROM wallet w, category c
			WHERE
			w.name=:src_wallet AND c.name=:category
			RETURNING *;
			`
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	if err := stmt.Get(tx, tx); err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	return nil
}

// DeleteAll ...
func (transactions *Transactions) DeleteAll() (int, error) {
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

	if err := stmt.Select(transactions); err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	return len(*transactions), nil
}
