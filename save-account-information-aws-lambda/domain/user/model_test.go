package user

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewUser", func() {
	It("should create a user with correct data and a 6-digit ID", func() {
		firstName := "testFirstName"
		lastName := "testLastName"

		u := NewUser(firstName, lastName)

		Expect(u.FirstName).To(Equal(firstName))
		Expect(u.LastName).To(Equal(lastName))

		idStr := strconv.FormatInt(u.ID, 10)
		Expect(len(idStr)).To(Equal(6))
	})
})
