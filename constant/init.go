package constant

import (
	"reflect"
	"sync"
)

type transactionType struct {
	once     sync.Once
	Income   string
	Expense  string
	Transfer string
}

type walletRole struct {
	once      sync.Once
	SrcWallet string
	DstWallet string
}

type walletType struct {
	once        sync.Once
	Cash        string
	BankAccount string
	Credit      string
}

var (
	wt walletType
	tt transactionType
	wr walletRole
)

// TransactionTypes returns the types of a transaction
func TransactionTypes() *transactionType {
	tt.once.Do(func() {
		tt.Expense = "EXPENSE"
		tt.Income = "INCOME"
		tt.Transfer = "TRANSFER"
	})
	return &tt
}

// WalletTypes returns the types of a wallet
func WalletTypes() *walletType {
	wt.once.Do(func() {
		wt.Cash = "CASH"
		wt.BankAccount = "BANK_ACCOUNT"
		wt.Credit = "CREDIT"
	})
	return &wt
}

// WalletRoles returns roles of a wallet
func WalletRoles() *walletRole {
	wr.once.Do(func() {
		wr.SrcWallet = "SRC_WALLET"
		wr.DstWallet = "DST_WALLET"
	})
	return &wr
}

// ListWalletTypes returns the types of a wallet as a slice of strings
func ListWalletTypes() []string {
	wt := WalletTypes()
	v := reflect.ValueOf(wt)
	types := make([]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		types[i] = v.Field(i).String()
	}

	return types
}
