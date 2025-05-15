package adapters

import (
	"context"
	"fmt"
	"os"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/internal/ports"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// repository implements ports.Repository
type transactionRepository struct {
	client    *dynamodb.DynamoDB
	tableName string
}

// NewTransactionRepository creates a new DynamoDB repository
func NewTransactionRepository() ports.TransactionRepository {
	sess := session.Must(session.NewSession())
	client := dynamodb.New(sess)

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		panic("DYNAMODB_TABLE_NAME environment variable not set")
	}

	return &transactionRepository{
		client:    client,
		tableName: tableName,
	}
}

// SaveTransactions stores a list of transactions in DynamoDB
func (t *transactionRepository) SaveTransactions(ctx context.Context, txs []transaction.Transaction) error {
	for _, tx := range txs {
		item, err := dynamodbattribute.MarshalMap(tx)
		if err != nil {
			return fmt.Errorf("error marshaling transaction: %w", err)
		}

		_, err = t.client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(t.tableName),
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("error put item in DynamoDB: %w", err)
		}
	}

	return nil
}
