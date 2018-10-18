package dbmodel

// WalletType ...
type WalletType string

const (
	Cash        WalletType = "CASH"
	BankAccount WalletType = "BANK_ACCOUNT"
	Credit      WalletType = "CREDIT"
)
