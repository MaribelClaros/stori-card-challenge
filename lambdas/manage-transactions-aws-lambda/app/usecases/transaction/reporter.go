package transaction

import (
	"math"
	"sort"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"
	"time"
)

type MonthlyReport struct {
	Month            time.Month // Calendar month (e.g., time.March)
	TransactionCount int        // Total number of transactions in the month
	AverageDebit     float64    // Average amount of debit transactions (negative values)
	AverageCredit    float64    // Average amount of credit transactions (positive values)
}

func CalculateReport(txs []transaction.Transaction) (float64, []MonthlyReport) {
	var balance float64
	monthly := make(map[time.Month]*MonthlyReport)
	debitCount := make(map[time.Month]int)
	creditCount := make(map[time.Month]int)

	for _, tx := range txs {
		balance += tx.Amount

		month := tx.Date.Month()

		if monthly[month] == nil {
			monthly[month] = &MonthlyReport{Month: month}
		}

		monthly[month].TransactionCount++
		if tx.Amount >= 0 {
			monthly[month].AverageCredit += tx.Amount
			creditCount[month]++
		} else {
			monthly[month].AverageDebit += tx.Amount
			debitCount[month]++
		}
	}

	// Calculate averages
	for month, report := range monthly {
		if creditCount[month] > 0 {
			report.AverageCredit /= float64(creditCount[month])
		}
		if debitCount[month] > 0 {
			report.AverageDebit /= float64(debitCount[month])
		}
	}

	// Convert map to sorted slice
	var reports []MonthlyReport
	for _, report := range monthly {
		reports = append(reports, *report)
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Month < reports[j].Month
	})

	roundBalance := math.Round(balance*100) / 100

	return roundBalance, reports
}
