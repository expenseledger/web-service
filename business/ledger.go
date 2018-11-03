package business

import (
	"github.com/expenseledger/web-service/model"
)

// InsertExpense inserts an expense transaction into a wallet
func InsertExpense(tx *model.Transaction) (*model.Wallet, error) {
	if err := tx.Create(); err != nil {
		return nil, err
	}

	wallet := model.Wallet{
		Name: tx.SrcWallet,
	}
	if err := wallet.Get(); err != nil {
		// @TODO: Roll back
		return nil, err
	}

	wallet.Balance = wallet.Balance.Sub(tx.Amount)
	if err := wallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, err
	}
	return &wallet, nil
}
