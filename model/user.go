package model

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/orm"
)

// ListCategories ...
func ListCategories() ([]Category, error) {
	return applyToCategories(list)
}

// ClearCategories ...
func ClearCategories() ([]Category, error) {
	return applyToCategories(clear)
}

// ListWallets ...
func ListWallets() ([]Wallet, error) {
	return applyToWallets(list)
}

// ClearWallets ...
func ClearWallets() ([]Wallet, error) {
	return applyToWallets(clear)
}

func ClearTransactions() ([]Transaction, error) {
	txTypes := constant.TransactionTypes()
	mapper := orm.NewTxMapper(_Transaction{}, txTypes.Expense)

	tmp, err := mapper.Clear()
	if err != nil {
		return nil, err
	}

	_txs := *(tmp.(*[]_Transaction))
	txs := make([]Transaction, 0, len(_txs))
	length := len(_txs)

	for i := 0; i < length; i++ {
		_tx := _txs[i]
		tx := _tx.toTransaction()

		if tx.Type == txTypes.Transfer {
			i++
			_tx = _txs[i]
			tx.To = _tx.Wallet
		}

		txs = append(txs, *tx)
	}

	return txs, nil
}

func applyToCategories(op operation) ([]Category, error) {
	mapper := orm.NewCategoryMapper(Category{})

	var tmp interface{}
	var err error
	switch op {
	case list:
		tmp, err = mapper.Many()
	case clear:
		tmp, err = mapper.Clear()
	}

	if err != nil {
		return nil, err
	}

	c := tmp.(*[]Category)
	categories := *c

	return categories, nil
}

func applyToWallets(op operation) ([]Wallet, error) {
	mapper := orm.NewWalletMapper(Wallet{})

	var tmp interface{}
	var err error
	switch op {
	case list:
		tmp, err = mapper.Many()
	case clear:
		tmp, err = mapper.Clear()
	}

	if err != nil {
		return nil, err
	}

	w := tmp.(*[]Wallet)
	wallets := *w

	return wallets, nil
}
