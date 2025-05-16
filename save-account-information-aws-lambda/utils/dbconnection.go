package utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateDBConnection(sess *session.Session) *dynamodb.DynamoDB {
	dbClient := dynamodb.New(sess)

	return dbClient
}
