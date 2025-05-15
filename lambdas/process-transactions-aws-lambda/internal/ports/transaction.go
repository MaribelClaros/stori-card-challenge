package ports

import (
	"context"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"
)

type TransactionRepository interface {
	SaveTransactions(ctx context.Context, txs []transaction.Transaction) error
}
