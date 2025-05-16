package usecases

import (
	"github.com/pkg/errors"
	"stori-card-challenge/process-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/topic"
	"time"
)

type SendMessageUsecase interface {
	Execute(MsgData MsgData) error
}

type sendMessageUsecase struct {
	snsSender topic.SnsSender
}

func NewSendMessageUsecase(snsSender topic.SnsSender) *sendMessageUsecase {
	return &sendMessageUsecase{
		snsSender: snsSender,
	}
}

type MsgData struct {
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Balance   float64              `json:"balance"`
	Monthly   []sns.MonthlySummary `json:"monthly"`
	Email     string               `json:"email"`
}

func (m *MsgData) ToTopicMessage() sns.TopicMessage {
	return sns.TopicMessage{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Balance:   m.Balance,
		Monthly:   m.Monthly,
		CreatedAt: time.Time{},
	}
}

func (s *sendMessageUsecase) Execute(msgData MsgData) error {
	tm := msgData.ToTopicMessage()
	err := s.snsSender.Execute(tm)

	if err != nil {
		return errors.Wrap(err, "cannot send msg to sns")
	}
	return err
}
