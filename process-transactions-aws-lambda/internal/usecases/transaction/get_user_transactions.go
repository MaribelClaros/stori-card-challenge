package usecases

import (
	"errors"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	infraTransaction "stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/transaction"
)

type GetTransactionUsecase interface {
	GetTransactions(bucket, key string) ([]transaction.Transaction, error)
}

type getTransactionUsecase struct {
	transactionRepository infraTransaction.TransactionRepository
}

func NewGetTransactionUsecase(transactionRepository infraTransaction.TransactionRepository) *getTransactionUsecase {
	return &getTransactionUsecase{
		transactionRepository: transactionRepository,
	}
}

func (u *getTransactionUsecase) GetTransactions(bucket, key string) ([]transaction.Transaction, error) {

	transactions, err := u.transactionRepository.GetTransactionsFromS3(bucket, key)

	if err != nil {
		return nil, errors.New("error getting transactions from s3")

	}

	return transactions, nil

}
