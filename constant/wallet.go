package constant

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
