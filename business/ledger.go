package business

import (
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
