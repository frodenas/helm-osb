package helm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/helm-osb/helm"
)

var _ = Describe("Config", func() {
	var (
		config Config

		validConfig = Config{
			ReleaseNamePrefix: "fake-release-name-prefix",
			DefaultNamespace:  "fake-default-namespace",
			BinaryLocation:    "helm",
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

		It("returns error if Release Name Prefix is not valid", func() {
			config.ReleaseNamePrefix = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Must provide a non-empty Release Name Prefix"))
		})

		It("returns error if Default Namespace is not valid", func() {
			config.DefaultNamespace = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Must provide a non-empty Default Namespace"))
		})

		It("returns error if Binary Location is not valid", func() {
			config.BinaryLocation = ""

			err := config.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Must provide a non-empty Binary Location"))
		})
	})
})
