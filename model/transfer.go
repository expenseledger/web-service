package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/database/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// Transfer the structure represents a transaction in application layer
type Transfer struct {
	ID          string          `json:"id"`
	SrcWallet   string          `json:"src_wallet"`
	DstWallet   string          `json:"dst_wallet"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

// Transfers is defined just to be used as a receiver
type Transfers []Transfer

// Clear ...
func (tfs *Transfers) Clear() (int, error) {
	var dbTfs dbmodel.Transfers

	length, err := dbTfs.DeleteAll()
	if err != nil {
		return 0, err
	}

	var tf Transfer
	transfers := make(Transfers, length)
	for i, dbTf := range dbTfs {
		tf.fromDBCounterpart(&dbTf)
		transfers[i] = tf
	}
	*tfs = transfers

	return length, nil
}

// Create ...
func (tf *Transfer) Create() error {
	dbTf := tf.toDBCounterpart()
	if err := dbTf.Insert(); err != nil {
		return err
	}

	tf.fromDBCounterpart(dbTf)
	return nil
}

// Get ...
func (tf *Transfer) Get() error {
	dbTf := tf.toDBCounterpart()
	if err := dbTf.One(); err != nil {
		return err
	}

	tf.fromDBCounterpart(dbTf)
	return nil
}

func (tf *Transfer) toDBCounterpart() *dbmodel.Transfer {
	var t *time.Time
	if tf.Date != nil {
		_t := time.Time(*tf.Date)
		t = &_t
	}

	return &dbmodel.Transfer{
		ID:          tf.ID,
		SrcWallet:   tf.SrcWallet,
		DstWallet:   tf.DstWallet,
		Amount:      tf.Amount,
		Description: tf.Description,
		OccurredAt:  t,
	}
}

func (tf *Transfer) fromDBCounterpart(dbTf *dbmodel.Transfer) {
	d := date.Date(*dbTf.OccurredAt)
	tf.ID = dbTf.ID
	tf.SrcWallet = dbTf.SrcWallet
	tf.DstWallet = dbTf.DstWallet
	tf.Amount = dbTf.Amount
	tf.Description = dbTf.Description
	tf.Date = &d
}
