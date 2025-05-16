package usecases

import (
	"errors"
	"math"
	"sort"
	"stori-card-challenge/process-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	infraTransaction "stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/transaction"
	"time"
)

type ProcessTransactionsAndSendEmailUsecase interface {
	ProcessTransactionsAndSendEmail(transactions []transaction.Transaction, email string) (*transaction.TransactionsInformation, error)
}

type processTransactionsAndSendEmailUsecase struct {
	emailSender infraTransaction.EmailSender
}

func NewProcessTransactionsAndSendEmailUsecase(emailSender infraTransaction.EmailSender) *processTransactionsAndSendEmailUsecase {
	return &processTransactionsAndSendEmailUsecase{
		emailSender: emailSender,
	}
}

func (u *processTransactionsAndSendEmailUsecase) ProcessTransactionsAndSendEmail(transactions []transaction.Transaction, email string) (*transaction.TransactionsInformation, error) {

	ts := processDataAndCalculateReport(transactions)

	ti := &transaction.TransactionsInformation{
		Balance: ts.Balance,
	}

	err := u.emailSender.SendEmail(ts, email)

	if err != nil {
		return nil, errors.New("error sending email to user")
	}

	return ti, err

}

func processDataAndCalculateReport(txs []transaction.Transaction) *infraTransaction.TransactionsStatus {
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

	return &infraTransaction.TransactionsStatus{
		Balance: roundBalance,
		Monthly: reports,
	}

}
