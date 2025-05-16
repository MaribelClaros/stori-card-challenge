package account

import (
	"stori-card-challenge/save-account-information-aws-lambda/domain/account"
	"stori-card-challenge/save-account-information-aws-lambda/domain/user"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromAccountToDTO", func() {
	It("should convert Account to DTO correctly", func() {
		now := time.Now()
		u := user.User{
			ID:        1,
			FirstName: "testFirstName",
			LastName:  "testLastName",
		}
		a := &account.Account{
			Id:           "testID",
			DateCreated:  now,
			TotalBalance: 1000.5,
			User:         u,
		}

		dto := FromAccountToDTO(a)

		Expect(dto.Id).To(Equal(a.Id))
		Expect(dto.DateCreated).To(Equal(a.DateCreated))
		Expect(dto.TotalBalance).To(Equal(a.TotalBalance))
		Expect(dto.UserFirstName).To(Equal(a.User.FirstName))
		Expect(dto.UserLastName).To(Equal(a.User.LastName))
		Expect(dto.UserID).To(Equal(a.User.ID))
	})
})
