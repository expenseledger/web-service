package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/database/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// ExpenseIncome the structure represents a transaction in application layer
type ExpenseIncome struct {
	ID          string          `json:"id"`
	Wallet      string          `json:"wallet"`
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

// ExpenseIncomes is defined just to be used as a receiver
type ExpenseIncomes []ExpenseIncome

// Clear ...
func (eis *ExpenseIncomes) Clear() (int, error) {
	var dbEis dbmodel.ExpenseIncomes

	length, err := dbEis.DeleteAll()
	if err != nil {
		return 0, err
	}

	var ei ExpenseIncome
	expenseincomes := make(ExpenseIncomes, length)
	for i, dbEi := range dbEis {
		ei.fromDBCounterpart(&dbEi)
		expenseincomes[i] = ei
	}
	*eis = expenseincomes

	return length, nil
}

// Create ...
func (ei *ExpenseIncome) Create() error {
	dbEi := ei.toDBCounterpart()
	if err := dbEi.Insert(); err != nil {
		return err
	}

	ei.fromDBCounterpart(dbEi)
	return nil
}

func (ei *ExpenseIncome) toDBCounterpart() *dbmodel.ExpenseIncome {
	var t *time.Time
	if ei.Date != nil {
		_t := time.Time(*ei.Date)
		t = &_t
	}

	return &dbmodel.ExpenseIncome{
		ID:          ei.ID,
		Wallet:      ei.Wallet,
		Amount:      ei.Amount,
		Type:        ei.Type,
		Category:    ei.Category,
		Description: ei.Description,
		OccurredAt:  t,
	}
}

func (ei *ExpenseIncome) fromDBCounterpart(dbEi *dbmodel.ExpenseIncome) {
	d := date.Date(*dbEi.OccurredAt)

	ei.ID = dbEi.ID
	ei.Wallet = dbEi.Wallet
	ei.Amount = dbEi.Amount
	ei.Type = dbEi.Type
	ei.Category = dbEi.Category
	ei.Description = dbEi.Description
	ei.Date = &d
}
