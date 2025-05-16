package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"stori-card-challenge/save-account-information-aws-lambda/domain/account"
	"stori-card-challenge/save-account-information-aws-lambda/domain/user"
	infraAccount "stori-card-challenge/save-account-information-aws-lambda/internal/infrastructure/account"
	usecases "stori-card-challenge/save-account-information-aws-lambda/internal/usecases/account"
	"stori-card-challenge/save-account-information-aws-lambda/utils"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	aws_config_path = "/var/task/aws_config.json"
)

type SNSMsg struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	TotalBalance float64 `json:"balance"`
	Email        string  `json:"email"`
}

func HandleSNS(ctx context.Context, SNSEvent events.SNSEvent) error {
	// Read AWS configuration from JSON file
	config, err := utils.ReadAWSConfig(aws_config_path)
	if err != nil {
		fmt.Println("Error reading AWS config:", err)
		return err
	}

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
	})

	if err != nil {
		log.Print("Error creating session:", err)
		return err
	}

	dynamoClient := utils.CreateDBConnection(sess)

	saveAccountRepository := infraAccount.NewAccountDBRepository(dynamoClient, config.DynamoTable)
	saveAccountUsecase := usecases.NewSaveAccountUsecase(saveAccountRepository)

	for _, message := range SNSEvent.Records {

		var snsMsg SNSMsg
		if err := json.Unmarshal([]byte(message.SNS.Message), &snsMsg); err != nil {
			log.Printf("Error parsing SNS message: %v", err)
			continue
		}
		log.Printf("Received SNS message: %+v", snsMsg)

		saveAccountModel := createUserAccountFromSNS(snsMsg)
		log.Printf("creating account for: %s", snsMsg.Email)

		err = saveAccountUsecase.Execute(ctx, saveAccountModel)

		if err != nil {
			log.Printf("error account for: %s", snsMsg.Email)
			return err
		}
		log.Print("account created successfully")

	}

	return nil
}

func createUserAccountFromSNS(msg SNSMsg) *account.Account {
	usr := user.NewUser(msg.FirstName, msg.LastName)

	return account.NewAccountForUser(usr, msg.TotalBalance)

}

func main() {
	lambda.Start(HandleSNS)

}
