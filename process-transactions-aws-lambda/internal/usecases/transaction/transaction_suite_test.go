package usecases

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUsecases(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Usercases Transaction Suite")
}
