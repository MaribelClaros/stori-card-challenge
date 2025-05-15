package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/app/usecases/transaction"
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/internal/adapters"
)

type App struct {
	useCase transaction.Processor
}

func NewApp() *App {
	repo := adapters.NewTransactionRepository()
	notifier := adapters.NewSnsTransactions()
	useCase := transaction.NewTransactionsProcessor(repo, notifier)
	return &App{useCase: useCase}
}

func (a *App) Handler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3 := record.S3

		log.Printf("Processing file: s3://%s/%s", s3.Bucket.Name, s3.Object.Key)

		err := a.useCase.ProcessCSVRecords(ctx, s3.Bucket.Name, s3.Object.Key)
		if err != nil {
			log.Printf("Error processing file %s: %v", s3.Object.Key, err)
			return err
		}
	}
	return nil
}

func main() {
	app := NewApp()
	lambda.Start(app.Handler)
}
