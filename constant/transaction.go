package constant

type transactionType struct {
	Income   string
	Expense  string
	Transfer string
}

// TransactionType ...
var TransactionType = transactionType{
	Income:   "INCOME",
	Expense:  "EXPENSE",
	Transfer: "TRANSFER",
}
