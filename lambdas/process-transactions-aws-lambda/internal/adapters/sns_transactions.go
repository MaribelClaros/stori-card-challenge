package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	domain "stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/sns"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/internal/ports"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type snsTransactions struct {
	client   *sns.SNS
	topicARN string
}

// NewSnsTransactions creates a new SNS notifier with default AWS session
func NewSnsTransactions() ports.Notifier {
	sess := session.Must(session.NewSession())
	snsClient := sns.New(sess)

	topicARN := os.Getenv("SNS_TOPIC_ARN")
	if topicARN == "" {
		panic("SNS_TOPIC_ARN environment variable not set")
	}

	return &snsTransactions{
		client:   snsClient,
		topicARN: topicARN,
	}
}

// Execute sends a message to the SNS topic
func (t *snsTransactions) Execute(ctx context.Context, message domain.TopicMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(jsonData)),
		TopicArn: aws.String(t.topicARN),
	}

	_, err = t.client.PublishWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("error publish message to SNS: %w", err)
	}

	return nil
}
