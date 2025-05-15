package transaction

import (
	"time"

	"stori-card-challenge/lambdas/process-transactions-aws-lambda/domain/transaction"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter", func() {

	Context("when calculating a report from transactions", func() {
		It("should return total balance and correct monthly summaries", func() {
			layout := "01/02"

			t1, _ := time.Parse(layout, "07/15")
			t2, _ := time.Parse(layout, "07/28")
			t3, _ := time.Parse(layout, "08/02")
			t4, _ := time.Parse(layout, "08/13")

			transactions := []transaction.Transaction{
				{ID: 1, Date: t1, Amount: 60.50},
				{ID: 2, Date: t2, Amount: 10.00},
				{ID: 3, Date: t3, Amount: -20.46},
				{ID: 4, Date: t4, Amount: -10.30},
			}

			total, summaries := CalculateReport(transactions)

			Expect(total).To(Equal(39.74))
			Expect(summaries).To(HaveLen(2))

			Expect(summaries[0].Month).To(Equal(time.July))
			Expect(summaries[0].TransactionCount).To(Equal(2))
			Expect(summaries[0].AverageCredit).To(Equal(35.25))
			Expect(summaries[0].AverageDebit).To(Equal(0.0))

			Expect(summaries[1].Month).To(Equal(time.August))
			Expect(summaries[1].TransactionCount).To(Equal(2))
			Expect(summaries[1].AverageCredit).To(Equal(0.0))
			Expect(summaries[1].AverageDebit).To(Equal(-15.38))
		})
	})
})
