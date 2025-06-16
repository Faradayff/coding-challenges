package e2etests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2etests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2etests Suite")
}
