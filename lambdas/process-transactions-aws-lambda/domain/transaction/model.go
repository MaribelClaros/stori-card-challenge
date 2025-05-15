package transaction

import "time"

type Transaction struct {
	ID     int       `json:"Id"`
	Date   time.Time `json:"Date"`
	Amount float64   `json:"Amount"`
}

type TransactionsInformation struct {
	TotalBalance    float64 `json:"total_balance"`
	AvgDebitAmount  float64 `json:"avg_debit_amount"`
	AvgCreditAmount float64 `json:"avg_credit_amount"`
}
