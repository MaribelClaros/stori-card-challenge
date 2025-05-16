package topic

import (
	"encoding/json"
	"log"
	topic "stori-card-challenge/process-transactions-aws-lambda/domain/sns"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
)

type SnsSender interface {
	Execute(topicMessage topic.TopicMessage) error
}

type snsSender struct {
	session  *session.Session
	topicArn string
}

func NewSnsSender(session *session.Session, topicArn string) *snsSender {
	return &snsSender{
		session:  session,
		topicArn: topicArn,
	}
}

func (s *snsSender) Execute(msg topic.TopicMessage) error {

	svc := sns.New(s.session)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	_, err = svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(jsonData)),
		TopicArn: aws.String(s.topicArn),
	})

	if err != nil {
		return errors.Wrap(err, "cannot send msg to sns")
	}

	return err

}
