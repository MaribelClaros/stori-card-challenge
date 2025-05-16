package user

import (
	"stori-card-challenge/save-account-information-aws-lambda/utils"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUser(firstName string, lastName string) User {
	idGenerator := utils.NewUserIDGenerator()
	return User{
		ID:        idGenerator.GenerateID(),
		FirstName: firstName,
		LastName:  lastName,
	}
}
