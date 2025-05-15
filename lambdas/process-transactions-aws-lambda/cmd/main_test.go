package main

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAppHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Handler Suite")
}

var _ = Describe("App Handler", func() {
	var (
		mockProcessor *MockProcessor
		app           *App
	)

	BeforeEach(func() {
		mockProcessor = &MockProcessor{}
		app = &App{useCase: mockProcessor}
	})

	It("should process an S3 event successfully", func() {
		mockProcessor.ProcessCSVRecordsFunc = func(ctx context.Context, bucket, key string) error {
			Expect(bucket).To(Equal("test-bucket"))
			Expect(key).To(Equal("file.csv"))
			return nil
		}

		s3Event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{Name: "test-bucket"},
						Object: events.S3Object{Key: "file.csv"},
					},
				},
			},
		}

		err := app.Handler(context.Background(), s3Event)
		Expect(err).To(BeNil())
		Expect(mockProcessor.Called).To(BeTrue())
	})

	It("should return error if useCase fails", func() {
		mockProcessor.ProcessCSVRecordsFunc = func(ctx context.Context, bucket, key string) error {
			return errors.New("fail")
		}

		s3Event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{Name: "test-bucket"},
						Object: events.S3Object{Key: "file.csv"},
					},
				},
			},
		}

		err := app.Handler(context.Background(), s3Event)
		Expect(err).To(MatchError("fail"))
	})
})
