package account

import (
	"stori-card-challenge/save-account-information-aws-lambda/domain/account"
	"time"
)

type AccountDTO struct {
	Id            string    `json:"id" dynamodbav:"id"`
	DateCreated   time.Time `json:"date_created" dynamodbav:"date_created"`
	TotalBalance  float64   `json:"total_balance" dynamodbav:"total_balance"`
	UserID        int64     `json:"user_id" dynamodbav:"user_id"`
	UserFirstName string    `json:"user_first_name" dynamodbav:"user_first_name"`
	UserLastName  string    `json:"user_last_name" dynamodbav:"user_last_name"`
}

func FromAccountToDTO(a *account.Account) *AccountDTO {
	return &AccountDTO{
		Id:            a.Id,
		DateCreated:   a.DateCreated,
		TotalBalance:  a.TotalBalance,
		UserID:        a.User.ID,
		UserFirstName: a.User.FirstName,
		UserLastName:  a.User.LastName,
	}
}
