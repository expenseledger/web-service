package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/db/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// Transaction the structure represents a transaction in application layer
type Transaction struct {
	ID          string          `json:"id"`
	SrcWallet   string          `json:"src_wallet"`
	DstWallet   string          `json:"dst_wallet"`
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        date.Date       `json:"date"`
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

// Get ...
func (tx *Transaction) Get() error {
	dbTx := tx.toDBCounterpart()
	if err := dbTx.One(); err != nil {
		return err
	}

	tx.fromDBCounterpart(dbTx)
	return nil
}

// Delete ...
func (tx *Transaction) Delete() error {
	dbTx := tx.toDBCounterpart()
	if err := dbTx.Delete(); err != nil {
		return err
	}

	tx.fromDBCounterpart(dbTx)
	return nil
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
