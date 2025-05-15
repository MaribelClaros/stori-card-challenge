package ports

import (
	"context"
	domain "stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/sns"
)

type Notifier interface {
	Execute(ctx context.Context, message domain.TopicMessage) error
}
