package account

import (
	"stori-card-challenge/save-account-information-aws-lambda/domain/user"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewAccountForUser", func() {
	It("should return an Account with correct user and balance", func() {
		u := user.User{
			ID:        1,
			FirstName: "testFirstName",
			LastName:  "testLastName",
		}

		totalBalance := 1000.00

		acc := NewAccountForUser(u, totalBalance)

		Expect(acc.User.ID).To(Equal(u.ID))
		Expect(acc.User.FirstName).To(Equal(u.FirstName))
		Expect(acc.User.LastName).To(Equal(u.LastName))
		Expect(acc.TotalBalance).To(Equal(totalBalance))
	})
})
