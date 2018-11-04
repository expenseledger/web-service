package view

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/type/date"
	"github.com/shopspring/decimal"
)

// Transaction the structure holding a transaction in presentation layer
type Transaction struct {
	ID          string          `json:"id"`
	SrcWallet   *string         `json:"src_wallet"`
	DstWallet   *string         `json:"dst_wallet"`
	Type        string          `json:"type"`
	Category    *string         `json:"category"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Date        *date.Date      `json:"date"`
}

// FromExpenseIncome converts expense/income to transaction
func (tx *Transaction) FromExpenseIncome(ei *model.ExpenseIncome) {
	tx.ID = ei.ID
	tx.Category = &ei.Category
	tx.Amount = ei.Amount
	tx.Description = ei.Description
	tx.Date = ei.Date

	switch ei.Type {
	case constant.TransactionType.Expense:
		tx.SrcWallet = &ei.Wallet
		tx.Type = constant.TransactionType.Expense
	case constant.TransactionType.Income:
		tx.DstWallet = &ei.Wallet
		tx.Type = constant.TransactionType.Income
	}
}

// FromTransfer converts transfer to transaction
func (tx *Transaction) FromTransfer(tf *model.Transfer) {
	tx.ID = tf.ID
	tx.SrcWallet = &tf.SrcWallet
	tx.DstWallet = &tf.DstWallet
	tx.Type = constant.TransactionType.Transfer
	tx.Amount = tf.Amount
	tx.Description = tf.Description
	tx.Date = tf.Date
}
