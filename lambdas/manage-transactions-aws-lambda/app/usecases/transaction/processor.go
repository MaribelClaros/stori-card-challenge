package transaction

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	topic "stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/internal/ports"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Processor interface {
	ProcessCSVRecords(ctx context.Context, bucket, key string) error
}

// TransactionsProcessor is the use case struct
type TransactionsProcessor struct {
	repo     ports.TransactionRepository
	notifier ports.Notifier
	s3Client ports.S3Client
}

func NewTransactionsProcessorWithS3(transactionRepo ports.TransactionRepository, snsNotifier ports.Notifier, s3Client ports.S3Client) *TransactionsProcessor {
	return &TransactionsProcessor{
		repo:     transactionRepo,
		notifier: snsNotifier,
		s3Client: s3Client,
	}
}

func NewTransactionsProcessor(transactionRepo ports.TransactionRepository, snsNotifier ports.Notifier) *TransactionsProcessor {
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
		return fmt.Errorf("error processing CSV in file: %s", err)
	}

	// Save to DynamoDB
	err = p.repo.SaveTransactions(ctx, txs)
	if err != nil {
		return fmt.Errorf("error saving transactions: %w", err)
	}

	// Generate financial summary
	balance, monthlyReports := CalculateReport(txs)

	// Prepare SNS message
	var summaries []topic.MonthlySummary
	for _, report := range monthlyReports {
		summaries = append(summaries, topic.MonthlySummary{
			Month:            report.Month,
			TransactionCount: report.TransactionCount,
			AverageDebit:     report.AverageDebit,
			AverageCredit:    report.AverageCredit,
		})
	}

	topicMessage := topic.TopicMessage{
		Transactions: txs,
		Balance:      balance,
		Monthly:      summaries,
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

		parseDate, err := time.Parse("01/02", record[1])

		if err != nil {
			return nil, errors.New("invalid Date format")

		}
		date := parseDate

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
