package usecases

import (
	"errors"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/transaction/mocks"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGetTransactionUsecase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GetTransactionUsecase Suite")
}

var _ = Describe("GetTransactionUsecase", func() {
	var (
		repo         *mocks.TransactionRepository
		usecase      GetTransactionUsecase
		bucket, key  string
		transactions []transaction.Transaction
	)

	BeforeEach(func() {
		repo = new(mocks.TransactionRepository)
		usecase = NewGetTransactionUsecase(repo)
		bucket = "test-bucket"
		key = "test-key"
	})

	Context("when getting transactions is successful", func() {
		BeforeEach(func() {
			transactions = []transaction.Transaction{
				{
					ID:     1,
					Date:   time.Date(2023, 2, 25, 0, 0, 0, 0, time.UTC),
					Amount: 100.5,
				},
			}
			repo.On("GetTransactionsFromS3", bucket, key).Return(transactions, nil)
		})

		It("should return the transactions", func() {
			result, err := usecase.GetTransactions(bucket, key)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(transactions))
			repo.AssertExpectations(GinkgoT())
		})
	})

	Context("when repository returns an error", func() {
		BeforeEach(func() {
			repo.On("GetTransactionsFromS3", bucket, key).Return(nil, errors.New("s3 error"))
		})

		It("should return an error", func() {
			result, err := usecase.GetTransactions(bucket, key)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error getting transactions from s3"))
			Expect(result).To(BeNil())
			repo.AssertExpectations(GinkgoT())
		})
	})
})
