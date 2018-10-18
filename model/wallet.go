package model

import (
	"time"

	"github.com/ExpenseLedger/expense-ledger-web-service/database/model"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

// Wallet the structure represents a wallet in presentation layer
type Wallet struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Type      dbmodel.WalletType `json:"type"`
	Balance   decimal.Decimal    `json:"balance"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// Create ...
func (wallet Wallet) Create() error {
	var dbWallet dbmodel.Wallet

	copier.Copy(&dbWallet, &wallet)

	return dbWallet.Insert()
}
