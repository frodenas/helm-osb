package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/helm-osb"

	"github.com/frodenas/helm-osb/broker"
	"github.com/frodenas/helm-osb/helm"
)

var _ = Describe("Config", func() {
	var (
		config Config

		validConfig = Config{
			LogLevel: "DEBUG",
			BrokerConfig: broker.Config{
				Username: "fake-broker-username",
				Password: "fake-broker-password",
			},
			HelmConfig: helm.Config{
				ReleaseNamePrefix: "fake-release-name-prefix",
				DefaultNamespace:  "fake-default-namespace",
				BinaryLocation:    "helm",
			},
		}
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			config = validConfig
		})

		It("does not return error if all sections are valid", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Log Level is not valid", func() {
			config.LogLevel = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Must provide a non-empty Log Level"))
		})

		It("returns error if Broker configuration is not valid", func() {
			config.BrokerConfig = broker.Config{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Broker configuration"))
		})

		It("returns error if Helm configuration is not valid", func() {
			config.HelmConfig = helm.Config{}

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Helm configuration"))
		})
	})
})
