package sns

import (
	"time"
)

type TopicMessage struct {
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Balance   float64          `json:"balance"`
	Monthly   []MonthlySummary `json:"monthly"`
	CreatedAt time.Time        `json:"created_at"`
}

type MonthlySummary struct {
	Month            time.Month // Calendar month (e.g., time.March)
	TransactionCount int        // Total number of transactions in the month
	AverageDebit     float64    // Average amount of debit transactions (negative values)
	AverageCredit    float64    // Average amount of credit transactions (positive values)
}
