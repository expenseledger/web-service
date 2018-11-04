package business

import (
	"github.com/expenseledger/web-service/model"
)

// InsertExpense inserts an expense transaction into a wallet
func InsertExpense(expense *model.ExpenseIncome) (*model.Wallet, error) {
	wallet := model.Wallet{
		Name: expense.Wallet,
	}
	if err := wallet.Get(); err != nil {
		return nil, err
	}

	if err := expense.Create(); err != nil {
		return nil, err
	}

	wallet.Balance = wallet.Balance.Sub(expense.Amount)
	if err := wallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, err
	}
	return &wallet, nil
}

// InsertIncome inserts an income transaction into a wallet
func InsertIncome(income *model.ExpenseIncome) (*model.Wallet, error) {
	wallet := model.Wallet{
		Name: income.Wallet,
	}
	if err := wallet.Get(); err != nil {
		return nil, err
	}

	if err := income.Create(); err != nil {
		return nil, err
	}

	wallet.Balance = wallet.Balance.Add(income.Amount)
	if err := wallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, err
	}
	return &wallet, nil
}
