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

var (
	tt transactionType
	wt walletType
	wr walletRole
)

func init() {
	tt.Expense = "EXPENSE"
	tt.Income = "INCOME"
	tt.Transfer = "TRANSFER"

	wt.Cash = "CASH"
	wt.BankAccount = "BANK_ACCOUNT"
	wt.Credit = "CREDIT"

	wr.SrcWallet = "SRC_WALLET"
	wr.DstWallet = "DST_WALLET"
}

// TransactionType ...
func TransactionType() *transactionType {
	return &tt
}

// WalletType ...
func WalletType() *walletType {
	return &wt
}

// WalletRole ...
func WalletRole() *walletRole {
	return &wr
}

// ListWalletTypes returns all possible types of a wallet
func ListWalletTypes() []string {
	v := reflect.ValueOf(wt)
	types := make([]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		types[i] = v.Field(i).String()
	}

	return types
}
