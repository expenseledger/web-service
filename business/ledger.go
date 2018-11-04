package business

import (
	"github.com/expenseledger/web-service/model"
	"github.com/expenseledger/web-service/view"
)

// InsertExpense inserts an expense transaction into a wallet
func InsertExpense(
	expense *model.ExpenseIncome,
) (*view.Transaction, *model.Wallet, error) {
	wallet := model.Wallet{
		Name: expense.Wallet,
	}
	if err := wallet.Get(); err != nil {
		return nil, nil, err
	}

	if err := expense.Create(); err != nil {
		return nil, nil, err
	}

	wallet.Balance = wallet.Balance.Sub(expense.Amount)
	if err := wallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	var tx view.Transaction
	tx.FromExpenseIncome(expense)

	return &tx, &wallet, nil
}

// InsertIncome inserts an income transaction into a wallet
func InsertIncome(
	income *model.ExpenseIncome,
) (*view.Transaction, *model.Wallet, error) {
	wallet := model.Wallet{
		Name: income.Wallet,
	}
	if err := wallet.Get(); err != nil {
		return nil, nil, err
	}

	if err := income.Create(); err != nil {
		return nil, nil, err
	}

	wallet.Balance = wallet.Balance.Add(income.Amount)
	if err := wallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	var tx view.Transaction
	tx.FromExpenseIncome(income)

	return &tx, &wallet, nil
}

// Transfer transfers an amount from wallet to another wallet
func Transfer(
	tf *model.Transfer,
) (*view.Transaction, *model.Wallet, *model.Wallet, error) {
	srcWallet := model.Wallet{
		Name: tf.SrcWallet,
	}
	if err := srcWallet.Get(); err != nil {
		return nil, nil, nil, err
	}

	dstWallet := model.Wallet{
		Name: tf.DstWallet,
	}
	if err := dstWallet.Get(); err != nil {
		return nil, nil, nil, err
	}

	if err := tf.Create(); err != nil {
		return nil, nil, nil, err
	}

	srcWallet.Balance = srcWallet.Balance.Sub(tf.Amount)
	if err := srcWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, nil, err
	}

	dstWallet.Balance = dstWallet.Balance.Add(tf.Amount)
	if err := dstWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, nil, err
	}

	var tx view.Transaction
	tx.FromTransfer(tf)

	return &tx, &srcWallet, &dstWallet, nil
}

// ClearTransactions deletes all transactions
func ClearTransactions() (int, []view.Transaction, error) {
	var (
		eis model.ExpenseIncomes
		tfs model.Transfers
	)
	eisLength, err := eis.Clear()
	if err != nil {
		return 0, nil, err
	}

	tfsLength, err := tfs.Clear()
	if err != nil {
		return 0, nil, err
	}

	length := eisLength + tfsLength
	txs := make([]view.Transaction, length)
	for _, ei := range eis {
		var tx view.Transaction
		tx.FromExpenseIncome(&ei)
		txs = append(txs, tx)
	}
	for _, tf := range tfs {
		var tx view.Transaction
		tx.FromTransfer(&tf)
		txs = append(txs, tx)
	}

	return length, txs, nil
}
