package dbmodel

import (
	"log"
	"time"

	"github.com/ExpenseLedger/expense-ledger-web-service/database"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a stored wallet on database
type Wallet struct {
	ID        string          `db:"id"`
	Name      string          `db:"name"`
	Type      WalletType      `db:"type"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
	DeletedAt time.Time       `db:"deleted_at"`
}

// Insert inserts a wallet into the database
func (wallet Wallet) Insert() error {
	query :=
		`
		INSERT INTO wallet (name, type, balance)
		VALUES (:name, :type, :balance)
		`
	db := database.GetDB()

	_, err := db.NamedExec(query, &wallet)
	if err != nil {
		log.Println("Error inserting a wallet", err)
		return err
	}

	return nil
}
