package model

import (
	"github.com/expenseledger/web-service/constant"
	"github.com/expenseledger/web-service/orm"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a wallet in presentation layer
type Wallet struct {
	Name    string              `json:"name" db:"name"`
	Type    constant.WalletType `json:"type" db:"type"`
	Balance decimal.Decimal     `json:"balance" db:"balance"`
	UserId  string              `json:"userId" db:"user_id"`
}

// CreateWallet inserts wallet to DB
func CreateWallet(
	name string,
	t constant.WalletType,
	balance decimal.Decimal,
	userId string,
) (*Wallet, error) {
	w := Wallet{Name: name, Type: t, Balance: balance, UserId: userId}
	mapper := orm.NewWalletMapper(w)

	tmp, err := mapper.Insert(&w)
	if err != nil {
		return nil, err
	}

	return tmp.(*Wallet), nil
}

// GetWallet returns matching wallet from DB
func GetWallet(name string, userId string) (*Wallet, error) {
	return applyToWallet(name, one, userId)
}

// DeleteWallet removes wallet from DB
func DeleteWallet(name string, userId string) (*Wallet, error) {
	return applyToWallet(name, delete, userId)
}

// ListWallets ...
func ListWallets(userId string) ([]Wallet, error) {
	return applyToWallets(list, userId)
}

// ClearWallets ...
func ClearWallets(userId string) ([]Wallet, error) {
	return applyToWallets(clear, userId)
}

func (wallet *Wallet) Expend(tx *Transaction) error {
	wallet.Balance = wallet.Balance.Sub(tx.Amount)
	return wallet.update()
}

func (wallet *Wallet) Receive(tx *Transaction) error {
	wallet.Balance = wallet.Balance.Add(tx.Amount)
	return wallet.update()
}

func (wallet *Wallet) update() error {
	mapper := orm.NewWalletMapper(*wallet)
	if _, err := mapper.Update(wallet); err != nil {
		return err
	}
	return nil
}

func applyToWallet(name string, op operation, userId string) (*Wallet, error) {
	w := Wallet{Name: name, UserId: userId}
	mapper := orm.NewWalletMapper(w)

	var tmp interface{}
	var err error
	switch op {
	case delete:
		tmp, err = mapper.Delete(&w)
	case one:
		tmp, err = mapper.One(&w)
	}

	if err != nil {
		return nil, err
	}

	return tmp.(*Wallet), nil
}

func applyToWallets(op operation, userId string) ([]Wallet, error) {
	wallet := Wallet{UserId: userId}
	mapper := orm.NewWalletMapper(wallet)

	var tmp interface{}
	var err error
	switch op {
	case list:
		tmp, err = mapper.Many(&wallet)
	case clear:
		tmp, err = mapper.Clear()
	}

	if err != nil {
		return nil, err
	}

	wallets := *(tmp.(*[]Wallet))
	return wallets, nil
}
