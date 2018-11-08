package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/database/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a transaction in application layer
type Transaction struct {
	ID          string          `db:"id"`
	SrcWallet   string          `db:"src_wallet"`
	DstWallet   string          `db:"dst_wallet"`
	Amount      decimal.Decimal `db:"amount"`
	Type        string          `db:"type"`
	Category    string          `db:"category"`
	Description string          `db:"description"`
	Date        date.Date       `db:"date"`
}

// Transactions is defined just to be used as a receiver
type Transactions []Transaction

// Create ...
func (tx *Transaction) Create() error {
	dbTx := tx.toDBCounterpart()
	if err := dbTx.Insert(); err != nil {
		return err
	}

	tx.fromDBCounterpart(dbTx)
	return nil
}

// Clear ...
func (txs *Transactions) Clear() (int, error) {
	var dbTxs dbmodel.Transactions
	if err := dbTxs.DeleteAll(); err != nil {
		return 0, err
	}

	length := len(dbTxs)
	transactions := make(Transactions, length)
	for i, dbTx := range dbTxs {
		var tx Transaction
		tx.fromDBCounterpart(&dbTx)
		transactions[i] = tx
	}

	*txs = transactions
	return length, nil
}

func (tx *Transaction) toDBCounterpart() *dbmodel.Transaction {

	return &dbmodel.Transaction{
		ID:          tx.ID,
		SrcWallet:   tx.SrcWallet,
		DstWallet:   tx.DstWallet,
		Type:        tx.Type,
		Category:    tx.Category,
		Amount:      tx.Amount,
		Description: tx.Description,
		OccurredAt:  time.Time(tx.Date),
	}
}

func (tx *Transaction) fromDBCounterpart(dbTx *dbmodel.Transaction) {
	tx.ID = dbTx.ID
	tx.SrcWallet = dbTx.SrcWallet
	tx.DstWallet = dbTx.DstWallet
	tx.Amount = dbTx.Amount
	tx.Type = dbTx.Type
	tx.Category = dbTx.Category
	tx.Description = dbTx.Description
	tx.Date = date.Date(dbTx.OccurredAt)
}
