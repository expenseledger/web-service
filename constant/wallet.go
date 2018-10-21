package constant

import "reflect"

type walletType struct {
	Cash        string `db:"cash"`
	BankAccount string `db:"bank_account"`
	Credit      string `db:"credit"`
}

// WalletType ...
var WalletType = walletType{
	Cash:        "CASH",
	BankAccount: "BANK_ACCOUNT",
	Credit:      "CREDIT",
}

// ListWalletTypes returns all possible types of a wallet
func ListWalletTypes() []string {
	v := reflect.ValueOf(WalletType)
	types := make([]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		types[i] = v.Field(i).String()
	}

	return types
}
