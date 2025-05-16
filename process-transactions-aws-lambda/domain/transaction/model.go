package transaction

import "time"

type Transaction struct {
	ID     int       `json:"Id"`
	Date   time.Time `json:"Date"`
	Amount float64   `json:"Amount"`
}

type TransactionsInformation struct {
	Balance float64 `json:"balance"`
}
