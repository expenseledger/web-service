package model

import (
	"fmt"
	"log"
	"time"

	"github.com/expenseledger/web-service/database"
	"github.com/shopspring/decimal"
)

// Transfer the structure represents a stored transaction on database
type Transfer struct {
	ID          string          `db:"id"`
	SrcWallet   string          `db:"src_wallet"`
	DstWallet   string          `db:"dst_wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Description string          `db:"description"`
	OccurredAt  *time.Time      `db:"occurred_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// Transfers is defined just to be used as a receiver
type Transfers []Transfer

// Insert ...
func (tf *Transfer) Insert() error {
	fields, names := tf.buildInsertSQLStmt()
	query := fmt.Sprintf(
		`
		INSERT INTO %s
		%s VALUES %s
		RETURNING *;
		`,
		database.Transfer,
		fields,
		names,
	)

	namedStmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	if err := namedStmt.Get(tf, tf); err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	return nil
}

// One ...
func (tf *Transfer) One() error {
	query := fmt.Sprintf(
		`
		SELECT * FROM %s
		WHERE id=:id;
		`,
		database.Transfer,
	)

	namedStmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error selecting a transaction", err)
		return err
	}

	if err := namedStmt.Get(tf, tf); err != nil {
		log.Println("Error selecting a transaction", err)
		return err
	}

	return nil
}

// DeleteAll ...
func (tfs *Transfers) DeleteAll() (int, error) {
	query := fmt.Sprintf(
		`
		DELETE FROM %s
		RETURNING *;
		`,
		database.Transfer,
	)

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	if err := stmt.Select(tfs); err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	return len(*tfs), nil
}

func (tf *Transfer) buildInsertSQLStmt() (fields, names string) {
	fields = "(src_wallet, dst_wallet, amount, description"
	names = "(:src_wallet, :dst_wallet, :amount, :description"

	if tf.OccurredAt != nil {
		fields += ", occurred_at"
		names += ", :occurred_at"
	}

	fields += ")"
	names += ")"

	return
}
