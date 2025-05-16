package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/topic"
	infraTransaction "stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/transaction"
	usecases "stori-card-challenge/process-transactions-aws-lambda/internal/usecases/transaction"
	"stori-card-challenge/process-transactions-aws-lambda/utils"
)

const (
	aws_config_path = "/var/task/aws_config.json"
)

type RequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// HandleAPIGatewayProxyRequest is the Lambda handler function.
func HandleAPIGatewayProxyRequest(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	config, err := utils.ReadAWSConfig(aws_config_path)
	if err != nil {
		log.Println("Error reading AWS config:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
	})

	if err != nil {
		log.Printf("Error creating session: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Could not create aws session!",
		}, nil
	}

	var requestBody RequestBody
	if err := json.Unmarshal([]byte(r.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Wrong type of request",
		}, nil
	}

	transactionsRepository := infraTransaction.NewGetTransactionRepository(session)

	getTransactionsUsecase := usecases.NewGetTransactionUsecase(transactionsRepository)

	transactions, err := getTransactionsUsecase.GetTransactions(config.S3Bucket, config.ObjectKey)

	fmt.Printf("Config bucket: %s, %s\n", config.S3Bucket, config.ObjectKey)

	if err != nil {
		fmt.Println("Could not retrieve transactions:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "could not retrieve transactions!",
		}, nil
	}

	emailSender := infraTransaction.NewGetEmailSender(session)
	processAndSendEmailUsecase := usecases.NewProcessTransactionsAndSendEmailUsecase(emailSender)

	transactionInfo, err := processAndSendEmailUsecase.ProcessTransactionsAndSendEmail(transactions, requestBody.Email)

	if err != nil {
		log.Println("error processing and sending email:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	msgData := ToMsgData(requestBody, *transactionInfo)

	snsSender := topic.NewSnsSender(session, config.TopicArn)

	sendSnsMessageUsecase := usecases.NewSendMessageUsecase(snsSender)

	err = sendSnsMessageUsecase.Execute(msgData)

	if err != nil {
		log.Print("Warning: could not send message", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Email has been sent to user!",
	}, nil
}

func ToMsgData(r RequestBody, tInfo transaction.TransactionsInformation) usecases.MsgData {
	return usecases.MsgData{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Balance:   tInfo.Balance,
	}
}

func main() {
	// Register Lambda handlers
	lambda.Start(HandleAPIGatewayProxyRequest)
}
