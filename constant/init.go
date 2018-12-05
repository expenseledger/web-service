package constant

import "reflect"

type transactionType struct {
	Income   string
	Expense  string
	Transfer string
}

type walletRole struct {
	SrcWallet string
	DstWallet string
}

type walletType struct {
	Cash        string `db:"cash"`
	BankAccount string `db:"bank_account"`
	Credit      string `db:"credit"`
}

// TransactionType returns the types of an transaction
// @TODO: this should be singleton
func TransactionType() *transactionType {
	return &transactionType{
		Expense:  "EXPENSE",
		Income:   "INCOME",
		Transfer: "TRANSFER",
	}
}

// WalletType ...
func WalletType() *walletType {
	return &walletType{
		Cash:        "CASH",
		BankAccount: "BANK_ACCOUNT",
		Credit:      "CREDIT",
	}
}

// WalletRole ...
func WalletRole() *walletRole {
	return &walletRole{
		SrcWallet: "SRC_WALLET",
		DstWallet: "DST_WALLET",
	}
}

// ListWalletTypes returns all possible types of a wallet
func ListWalletTypes() []string {
	wt := WalletType()
	v := reflect.ValueOf(*wt)
	types := make([]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		types[i] = v.Field(i).String()
	}

	return types
}
