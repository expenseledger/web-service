package model

import (
	"fmt"
	"log"
	"time"

	"github.com/expenseledger/web-service/database"
	"github.com/shopspring/decimal"
)

// ExpenseIncome the structure represents a stored transaction on database
type ExpenseIncome struct {
	ID          string          `db:"id"`
	Wallet      string          `db:"wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    string          `db:"category"`
	Description string          `db:"description"`
	OccurredAt  *time.Time      `db:"occurred_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

// ExpenseIncomes is defined just to be used as a receiver
type ExpenseIncomes []ExpenseIncome

// Insert ...
func (ei *ExpenseIncome) Insert() error {
	fields, names := ei.buildInsertSQLStmt()
	query := fmt.Sprintf(
		`
		INSERT INTO %s
		%s VALUES %s
		RETURNING *;
		`,
		database.ExpenseIncome,
		fields,
		names,
	)

	namedStmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	if err := namedStmt.Get(ei, ei); err != nil {
		log.Println("Error inserting a transaction", err)
		return err
	}

	return nil
}

// DeleteAll ...
func (eis *ExpenseIncomes) DeleteAll() (int, error) {
	query := fmt.Sprintf(
		`
		DELETE FROM %s
		RETURNING *;
		`,
		database.ExpenseIncome,
	)

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	if err := stmt.Select(eis); err != nil {
		log.Println("Error deleting all transactions", err)
		return 0, err
	}

	return len(*eis), nil
}

func (ei *ExpenseIncome) buildInsertSQLStmt() (fields, names string) {
	fields = "(wallet, amount, type, category, description"
	names = "(:wallet, :amount, :type, :category, :description"

	if ei.OccurredAt != nil {
		fields += ", occurred_at"
		names += ", :occurred_at"
	}

	fields += ")"
	names += ")"

	return
}
