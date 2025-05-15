package sns

import "stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"

// TopicMessage is the message sent via SNS containing transaction data and summary.
type TopicMessage struct {
	Transactions []transaction.Transaction `json:"transactions"` // List of individual transactions
	Balance      float64                   `json:"balance"`      // Total balance from all transactions
	Monthly      []MonthlySummary          `json:"monthly"`      // Monthly aggregated summaries
}
type MonthlySummary struct {
	Month            string  `json:"month"`             // Month name (e.g., "January")
	TransactionCount int     `json:"transaction_count"` // Number of transactions in the month
	AverageDebit     float64 `json:"average_debit"`     // Average debit amount (negative values)
	AverageCredit    float64 `json:"average_credit"`    // Average credit amount (positive values)
}
