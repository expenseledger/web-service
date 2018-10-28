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

// Clear ...
func (txs *Transactions) Clear() (int, error) {
	var dbTxs dbmodel.Transactions

	length, err := dbTxs.DeleteAll()
	if err != nil {
		return 0, err
	}

	transactions := make(Transactions, 0, length)
	for _, dbTx := range dbTxs {
		tx := newFromDBCounterpart(&dbTx)
		transactions = append(transactions, tx)
	}
	*txs = transactions

	return length, nil
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

func newFromDBCounterpart(dbTx *dbmodel.Transaction) Transaction {
	var d *date.Date
	if dbTx.OccuredAt != nil {
		_d := date.Date(*dbTx.OccuredAt)
		d = &_d
	}

	return Transaction{
		ID:          dbTx.ID,
		SrcWallet:   dbTx.SrcWallet,
		DstWallet:   dbTx.DstWallet,
		Amount:      dbTx.Amount,
		Type:        dbTx.Type,
		Category:    dbTx.Category,
		Description: dbTx.Description,
		Date:        d,
	}
}
