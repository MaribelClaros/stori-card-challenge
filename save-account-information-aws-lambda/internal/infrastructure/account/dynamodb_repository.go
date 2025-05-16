package account

import (
	"context"
	"log"
	"stori-card-challenge/save-account-information-aws-lambda/domain/account"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type AccountDBRepository interface {
	SaveUserAccount(ctx context.Context, a *account.Account) error
}

type accountDBRepository struct {
	db    *dynamodb.DynamoDB
	table string
}

func NewAccountDBRepository(db *dynamodb.DynamoDB, table string) *accountDBRepository {
	return &accountDBRepository{
		db:    db,
		table: table,
	}
}

func (r *accountDBRepository) SaveUserAccount(ctx context.Context, a *account.Account) error {

	aDto := FromAccountToDTO(a)

	attributeValues, err := dynamodbattribute.MarshalMap(aDto)

	if err != nil {
		return errors.Wrapf(err, "error mapping account DTO")
	}

	input := dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(r.table),
	}

	log.Print("input received:", input)

	_, err = r.db.PutItem(&input)
	if err != nil {
		log.Print("Got error calling PutItem:", err)
		return errors.Wrapf(err, "error putting item in dynamo db")

	}

	return err

}
