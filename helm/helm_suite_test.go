package helm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCloudFormationBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helm Suite")
}
