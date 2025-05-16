package transaction

import (
	"stori-card-challenge/process-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	"time"
)

type TransactionDTO struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

func FromDTOtoTransaction(dto TransactionDTO) transaction.Transaction {
	return transaction.Transaction{
		ID:     dto.ID,
		Date:   dto.Date,
		Amount: dto.Amount,
	}
}

type TransactionsStatus struct {
	Balance float64              `json:"balance"`
	Monthly []sns.MonthlySummary `json:"monthly"`
}
