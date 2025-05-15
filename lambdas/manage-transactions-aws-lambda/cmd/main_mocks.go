package main

import (
	"context"
)

type MockProcessor struct {
	ProcessCSVRecordsFunc func(ctx context.Context, bucket, key string) error
	Called                bool
}

func (m *MockProcessor) ProcessCSVRecords(ctx context.Context, bucket, key string) error {
	m.Called = true
	if m.ProcessCSVRecordsFunc != nil {
		return m.ProcessCSVRecordsFunc(ctx, bucket, key)
	}
	return nil
}
