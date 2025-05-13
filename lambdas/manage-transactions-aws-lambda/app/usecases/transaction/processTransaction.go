armame un processTransactionTest.go en base a este código, usando ginko y gin gonic:

package transaction

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	topic "stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/internal/infraestructure/dynamodb"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/internal/infraestructure/sns"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TransactionsProcessor is the use case struct
type TransactionsProcessor struct {
	repo     dynamodb.TransactionRepository
	notifier sns.Notifier
	s3Client *s3.S3
}

func NewTransactionsProcessor(transactionRepo dynamodb.TransactionRepository, snsNotifier sns.Notifier) *TransactionsProcessor {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	return &TransactionsProcessor{
		repo:     transactionRepo,
		notifier: snsNotifier,
		s3Client: s3Client,
	}
}

// ProcessCSVRecords reads a CSV file from S3, parses transactions,
// saves them to DynamoDB and sends a summary via SNS
func (p *TransactionsProcessor) ProcessCSVRecords(ctx context.Context, bucket, key string) error {

	// Get object from S3
	goo, err := p.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("cannot get S3 object: %w", err)
	}
	defer goo.Body.Close()

	records, err := csv.NewReader(goo.Body).ReadAll()
	if err != nil {
		return fmt.Errorf("cannot read CSV: %w", err)
	}

	txs, err := validateAndProcessCSVRecords(records)

	if err != nil {
		return fmt.Errorf("error processing CSV in file %s", key)
	}

	// Save to DynamoDB
	err = p.repo.SaveTransactions(ctx, txs)
	if err != nil {
		return fmt.Errorf("error saving transactions: %w", err)
	}

	topicMessage := topic.TopicMessage{
		Transactions: txs,
	}

	err = p.notifier.Execute(ctx, topicMessage)
	if err != nil {
		return fmt.Errorf("error sending SNS notification: %w", err)
	}

	log.Printf("File processed successfully!")
	return nil
}

func validateAndProcessCSVRecords(records [][]string) ([]transaction.Transaction, error) {

	var txs []transaction.Transaction

	for i, record := range records {
		//ignore first row (title)
		if i == 0 {
			continue
		}

		if len(record) != 3 {
			return nil, errors.New("invalid number of fields in CSV record")
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, errors.New("invalid ID format")
		}

		_, err = time.Parse("01/02", record[1])

		if err != nil {
			return nil, errors.New("invalid Date format")

		}
		date := record[1]

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, errors.New("invalid Amount format")
		}

		tx := transaction.Transaction{
			ID:     id,
			Date:   date,
			Amount: amount,
		}
		txs = append(txs, tx)

	}
	return txs, nil

}
