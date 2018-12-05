package business

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/model"
)

// Expend inserts and rebalances wallet and transaction, respectively
func Expend(tx *model.Transaction) (*model.Wallet, error) {
	srcWallet := model.Wallet{
		Name: tx.SrcWallet,
	}
	if err := srcWallet.Get(); err != nil {
		return nil, err
	}

	if err := tx.Create(); err != nil {
		return nil, err
	}

	srcWallet.Balance = srcWallet.Balance.Sub(tx.Amount)
	if err := srcWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, err
	}

	return &srcWallet, nil
}

// Receive inserts and rebalances wallet and transaction, respectively
func Receive(tx *model.Transaction) (*model.Wallet, error) {
	dstWallet := model.Wallet{
		Name: tx.DstWallet,
	}
	if err := dstWallet.Get(); err != nil {
		return nil, err
	}

	if err := tx.Create(); err != nil {
		return nil, err
	}

	dstWallet.Balance = dstWallet.Balance.Add(tx.Amount)
	if err := dstWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, err
	}

	return &dstWallet, nil
}

// Transfer transfers an amount from wallet to another wallet
func Transfer(tx *model.Transaction) (*model.Wallet, *model.Wallet, error) {
	srcWallet := model.Wallet{
		Name: tx.SrcWallet,
	}
	if err := srcWallet.Get(); err != nil {
		return nil, nil, err
	}

	dstWallet := model.Wallet{
		Name: tx.DstWallet,
	}
	if err := dstWallet.Get(); err != nil {
		return nil, nil, err
	}

	if err := tx.Create(); err != nil {
		return nil, nil, err
	}

	srcWallet.Balance = srcWallet.Balance.Sub(tx.Amount)
	if err := srcWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	dstWallet.Balance = dstWallet.Balance.Add(tx.Amount)
	if err := dstWallet.Update(true); err != nil {
		// @TODO: Roll back
		return nil, nil, err
	}

	return &srcWallet, &dstWallet, nil
}

// DeleteTransaction deletes a transaction and rebalances wallets
func DeleteTransaction(
	tx *model.Transaction,
) (*model.Wallet, *model.Wallet, error) {
	if err := tx.Delete(); err != nil {
		return nil, nil, err
	}

	transactionType := constant.TransactionTypes()
	srcWallet := &model.Wallet{
		Name: tx.SrcWallet,
	}
	dstWallet := &model.Wallet{
		Name: tx.DstWallet,
	}

	switch tx.Type {
	case transactionType.Transfer:
		if err := srcWallet.Get(); err != nil {
			return nil, nil, err
		}
		if err := dstWallet.Get(); err != nil {
			return nil, nil, err
		}
		srcWallet.Balance = srcWallet.Balance.Add(tx.Amount)
		if err := srcWallet.Update(true); err != nil {
			// @TODO: Roll back
			return nil, nil, err
		}
		dstWallet.Balance = dstWallet.Balance.Sub(tx.Amount)
		if err := dstWallet.Update(true); err != nil {
			// @TODO: Roll back
			return nil, nil, err
		}
	case transactionType.Expense:
		if err := srcWallet.Get(); err != nil {
			return nil, nil, err
		}
		srcWallet.Balance = srcWallet.Balance.Add(tx.Amount)
		if err := srcWallet.Update(true); err != nil {
			// @TODO: Roll back
			return nil, nil, err
		}
		dstWallet = nil
	case transactionType.Income:
		if err := dstWallet.Get(); err != nil {
			return nil, nil, err
		}
		dstWallet.Balance = dstWallet.Balance.Sub(tx.Amount)
		if err := dstWallet.Update(true); err != nil {
			// @TODO: Roll back
			return nil, nil, err
		}
		srcWallet = nil
	}

	return srcWallet, dstWallet, nil
}
