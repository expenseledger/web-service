package model

import (
	"time"

	"github.com/expenseledger/web-service/database/model"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a wallet in presentation layer
type Wallet struct {
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Balance   decimal.Decimal `json:"balance"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Wallets is defined just to be used as a receiver
type Wallets []Wallet

// Create ...
func (wallet *Wallet) Create() error {
	var dbWallet dbmodel.Wallet

	copier.Copy(&dbWallet, &wallet)

	if err := dbWallet.Insert(); err != nil {
		return err
	}

	copier.Copy(wallet, &dbWallet)
	return nil
}

// Get ...
func (wallet *Wallet) Get(name string) error {
	var dbWallet dbmodel.Wallet
	if err := dbWallet.One(name); err != nil {
		return err
	}

	copier.Copy(wallet, &dbWallet)
	return nil
}

// Delete ...
func (wallet *Wallet) Delete(name string) error {
	var dbWallet dbmodel.Wallet
	if err := dbWallet.Delete(name); err != nil {
		return err
	}

	copier.Copy(wallet, &dbWallet)
	return nil
}

// List ...
func (wallets *Wallets) List() (int, error) {
	var dbWallets dbmodel.Wallets

	length, err := dbWallets.All()
	if err != nil {
		return 0, err
	}

	copier.Copy(wallets, &dbWallets)
	return length, nil
}
