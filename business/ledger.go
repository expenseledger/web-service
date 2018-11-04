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

// Transfer transfers an amount from wallet to another wallet
func Transfer(
	tf *model.Transfer,
) (*model.Wallet, *model.Wallet, error) {
	srcWallet := model.Wallet{
		Name: tf.SrcWallet,
	}
	if err := srcWallet.Get(); err != nil {
		return nil, nil, err
	}

	dstWallet := model.Wallet{
		Name: tf.DstWallet,
	}
	if err := dstWallet.Get(); err != nil {
		return nil, nil, err
	}

	if err := tf.Create(); err != nil {
		return nil, nil, err
	}

	srcWallet.Balance = srcWallet.Balance.Sub(tf.Amount)
	if err := srcWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	dstWallet.Balance = dstWallet.Balance.Add(tf.Amount)
	if err := dstWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	return &srcWallet, &dstWallet, nil
}
