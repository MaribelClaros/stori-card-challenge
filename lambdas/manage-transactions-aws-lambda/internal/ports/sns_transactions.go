package ports

import (
	"context"
	domain "stori-card-challenge/lambdas/manage-transactions-aws-lambda/domain/sns"
)

type Notifier interface {
	Execute(ctx context.Context, message domain.TopicMessage) error
}
