package transaction

import (
	"context"
	"io"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/sns"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"
)

// ---------------- TransactionRepository ----------------

type MockTransactionRepository struct {
	SaveTransactionsFunc func(ctx context.Context, txs []transaction.Transaction) error
}

func (m *MockTransactionRepository) SaveTransactions(ctx context.Context, txs []transaction.Transaction) error {
	if m.SaveTransactionsFunc != nil {
		return m.SaveTransactionsFunc(ctx, txs)
	}
	return nil
}

// ---------------- Notifier ----------------

type MockNotifier struct {
	ExecuteFunc func(ctx context.Context, message interface{}) error
}

func (m *MockNotifier) Execute(ctx context.Context, message sns.TopicMessage) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, message)
	}
	return nil
}

// ---------------- S3Client ----------------

type MockS3Client struct {
	GetObjectWithContextFunc func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
}

func (m *MockS3Client) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	if m.GetObjectWithContextFunc != nil {
		return m.GetObjectWithContextFunc(ctx, input, opts...)
	}
	return &s3.GetObjectOutput{
		Body: io.NopCloser(nil),
	}, nil
}
