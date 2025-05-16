package usecases

import (
	"errors"
	"stori-card-challenge/process-transactions-aws-lambda/domain/transaction"
	"stori-card-challenge/process-transactions-aws-lambda/internal/infrastructure/transaction/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProcessTransactionsAndSendEmailUsecase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProcessTransactionsAndSendEmailUsecase Suite")
}

var _ = Describe("ProcessTransactionsAndSendEmailUsecase", func() {
	var (
		emailSender *mocks.EmailSender
		usecase     ProcessTransactionsAndSendEmailUsecase
		email       string
	)

	BeforeEach(func() {
		emailSender = new(mocks.EmailSender)
		usecase = NewProcessTransactionsAndSendEmailUsecase(emailSender)
		email = "test@example.com"
	})

	Context("when processing and sending email is successful", func() {
		It("should process transactions and send email", func() {
			transactions := []transaction.Transaction{
				{ID: 1, Date: time.Date(2023, 7, 10, 0, 0, 0, 0, time.UTC), Amount: 100},
				{ID: 2, Date: time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC), Amount: -50},
			}

			// Use mock.Anything for TransactionsStatus
			emailSender.On("SendEmail", mock.Anything, email).Return(nil)

			info, err := usecase.ProcessTransactionsAndSendEmail(transactions, email)

			Expect(err).To(BeNil())
			Expect(info).NotTo(BeNil())
			Expect(info.Balance).To(Equal(50.0))
			emailSender.AssertCalled(GinkgoT(), "SendEmail", mock.Anything, email)
		})
	})

	Context("when SendEmail fails", func() {
		It("should return an error", func() {
			transactions := []transaction.Transaction{
				{ID: 1, Date: time.Date(2023, 7, 10, 0, 0, 0, 0, time.UTC), Amount: 100},
			}

			emailSender.On("SendEmail", mock.Anything, email).Return(errors.New("email error"))

			info, err := usecase.ProcessTransactionsAndSendEmail(transactions, email)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error sending email to user"))
			Expect(info).To(BeNil())
			emailSender.AssertNumberOfCalls(GinkgoT(), "SendEmail", 1)
		})
	})
})
