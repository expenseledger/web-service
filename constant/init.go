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

// TransactionTypes returns the types of an transaction
// @TODO: this should be singleton
func TransactionTypes() *transactionType {
	return &transactionType{
		Expense:  "EXPENSE",
		Income:   "INCOME",
		Transfer: "TRANSFER",
	}
}

// WalletTypes ...
func WalletTypes() *walletType {
	return &walletType{
		Cash:        "CASH",
		BankAccount: "BANK_ACCOUNT",
		Credit:      "CREDIT",
	}
}

// WalletRoles ...
func WalletRoles() *walletRole {
	return &walletRole{
		SrcWallet: "SRC_WALLET",
		DstWallet: "DST_WALLET",
	}
}

// ListWalletTypes returns all possible types of a wallet
func ListWalletTypes() []string {
	wt := WalletTypes()
	v := reflect.ValueOf(*wt)
	types := make([]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		types[i] = v.Field(i).String()
	}

	return types
}
