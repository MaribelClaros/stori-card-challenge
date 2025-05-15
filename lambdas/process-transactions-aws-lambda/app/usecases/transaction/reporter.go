package transaction

import (
	"math"
	"sort"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"
	"time"
)

func CalculateReport(txs []transaction.Transaction) (float64, []sns.MonthlySummary) {
	var balance float64
	monthly := make(map[time.Month]*sns.MonthlySummary)
	debitCount := make(map[time.Month]int)
	creditCount := make(map[time.Month]int)

	for _, tx := range txs {
		balance += tx.Amount

		month := tx.Date.Month()

		if monthly[month] == nil {
			monthly[month] = &sns.MonthlySummary{Month: month}
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
	var reports []sns.MonthlySummary
	for _, report := range monthly {
		reports = append(reports, *report)
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Month < reports[j].Month
	})

	roundBalance := math.Round(balance*100) / 100

	return roundBalance, reports
}
