package sns

import (
	"stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/transaction"
)

type TopicMessage struct {
	Transactions []transaction.Transaction `json:"transactions"`
}
