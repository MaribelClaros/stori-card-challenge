package sns

import (
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"
	"time"
)

// TopicMessage is the message sent via SNS containing transaction data and summary.
type TopicMessage struct {
	Transactions []transaction.Transaction `json:"transactions"` // List of individual transactions
	Balance      float64                   `json:"balance"`      // Total balance from all transactions
	Monthly      []MonthlySummary          `json:"monthly"`      // Monthly aggregated summaries
}
type MonthlySummary struct {
	Month            time.Month // Calendar month (e.g., time.March)
	TransactionCount int        // Total number of transactions in the month
	AverageDebit     float64    // Average amount of debit transactions (negative values)
	AverageCredit    float64    // Average amount of credit transactions (positive values)
}
