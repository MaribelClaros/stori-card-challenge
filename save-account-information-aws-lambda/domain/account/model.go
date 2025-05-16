package account

import (
	"stori-card-challenge/save-account-information-aws-lambda/domain/user"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id           string    `json:"id"`
	DateCreated  time.Time `json:"date_created"`
	TotalBalance float64   `json:"total_balance"`
	User         user.User `json:"user"`
}

func NewAccountForUser(user user.User, totalBalance float64) *Account {
	return &Account{
		Id:           uuid.New().String(),
		DateCreated:  time.Now().UTC(),
		TotalBalance: totalBalance,
		User:         user,
	}
}
