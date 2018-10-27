package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/database/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a transaction in presentation layer
type Transaction struct {
	ID          string          `json:"id"`
	SrcWallet   string          `json:"src_wallet"`
	DstWallet   *string         `json:"dst_wallet"`
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Create ...
func (tx *Transaction) Create() error {
	dbTx := tx.toDBCounterpart()
	if err := dbTx.Insert(); err != nil {
		return err
	}

	tx.newFromDBCounterpart(dbTx)
	return nil
}

func (tx *Transaction) toDBCounterpart() *dbmodel.Transaction {
	var t *time.Time
	if tx.Date != nil {
		_t := time.Time(*tx.Date)
		t = &_t
	}

	return &dbmodel.Transaction{
		ID:          tx.ID,
		SrcWallet:   tx.SrcWallet,
		DstWallet:   tx.DstWallet,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Category:    tx.Category,
		Description: tx.Description,
		OccuredAt:   t,
	}
}

func (tx *Transaction) newFromDBCounterpart(dbTx *dbmodel.Transaction) {
	var d *date.Date
	if dbTx.OccuredAt != nil {
		_d := date.Date(*dbTx.OccuredAt)
		d = &_d
	}

	tx.ID = dbTx.ID
	tx.SrcWallet = dbTx.SrcWallet
	tx.DstWallet = dbTx.DstWallet
	tx.Amount = dbTx.Amount
	tx.Type = dbTx.Type
	tx.Category = dbTx.Category
	tx.Description = dbTx.Description
	tx.Date = d
}
