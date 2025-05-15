package transaction

import (
	"context"
	"io"
	"os"
	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProcessor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Handler Suite")
}

var _ = Describe("TransactionsProcessor", func() {
	var (
		mockRepo     *MockTransactionRepository
		mockNotifier *MockNotifier
		mockS3Client *MockS3Client
		processor    *TransactionsProcessor
	)

	BeforeEach(func() {
		mockRepo = &MockTransactionRepository{}
		mockNotifier = &MockNotifier{}
		mockS3Client = &MockS3Client{}
		processor = NewTransactionsProcessorWithS3(mockRepo, mockNotifier, mockS3Client)
	})

	It("should process a valid CSV file and notify", func() {

		file, err := os.Open("testdata/mock_data.csv")
		Expect(err).To(BeNil())
		defer file.Close()

		csvBytes, err := io.ReadAll(file)
		Expect(err).To(BeNil())

		mockS3Client.GetObjectWithContextFunc = func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: io.NopCloser(strings.NewReader(string(csvBytes))),
			}, nil
		}

		mockRepo.SaveTransactionsFunc = func(ctx context.Context, txs []transaction.Transaction) error {
			Expect(len(txs)).To(Equal(4))
			return nil
		}

		mockNotifier.ExecuteFunc = func(ctx context.Context, msg interface{}) error {
			return nil
		}

		err = processor.ProcessCSVRecords(context.TODO(), "bucket", "key")
		Expect(err).To(BeNil())
	})

	It("should return error if CSV is invalid", func() {
		file, err := os.Open("testdata/mock_invalid_data.csv")
		Expect(err).To(BeNil())
		defer file.Close()

		csvBytes, err := io.ReadAll(file)
		Expect(err).To(BeNil())

		mockS3Client.GetObjectWithContextFunc = func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: io.NopCloser(strings.NewReader(string(csvBytes))),
			}, nil
		}

		err = processor.ProcessCSVRecords(context.TODO(), "bucket", "key")
		Expect(err).ToNot(BeNil())
	})
})
