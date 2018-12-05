package model

import (
	"time"

	dbmodel "github.com/expenseledger/web-service/db/model"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a wallet in presentation layer
type Wallet struct {
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Balance   decimal.Decimal `json:"balance"`
	UpdatedAt time.Time       `json:"last_update"`
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
func (wallet *Wallet) Get() error {
	dbWallet := dbmodel.Wallet{
		Name: wallet.Name,
	}

	if err := dbWallet.One(); err != nil {
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

// Init inserts default wallets
func (wallets *Wallets) Init() (int, error) {
	var dbWallets dbmodel.Wallets
	copier.Copy(&dbWallets, wallets)

	length, err := dbWallets.BatchInsert()
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Clear ...
func (wallets *Wallets) Clear() (int, error) {
	var dbWallets dbmodel.Wallets

	length, err := dbWallets.DeleteAll()
	if err != nil {
		return 0, err
	}

	copier.Copy(wallets, &dbWallets)
	return length, nil
}

// Update ...
func (wallet *Wallet) Update(replacing bool) error {
	var dbWallet dbmodel.Wallet
	copier.Copy(&dbWallet, wallet)

	var err error
	if replacing {
		err = dbWallet.Save()
	} else {
		err = dbWallet.Update()
	}

	if err != nil {
		return err
	}

	copier.Copy(wallet, &dbWallet)
	return nil
}
