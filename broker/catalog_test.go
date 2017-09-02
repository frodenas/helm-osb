package broker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/helm-osb/broker"
)

var _ = Describe("Catalog", func() {
	var (
		catalog      Catalog
		service1     Service
		service2     Service
		servicePlan1 ServicePlan
		servicePlan2 ServicePlan
	)

	BeforeEach(func() {
		servicePlan1 = ServicePlan{
			ID:          "fake-service-plan-1",
			Name:        "Fake Service Plan 1 Name",
			Description: "Fake Service Plan 1 Description",
		}

		servicePlan2 = ServicePlan{
			ID:          "fake-service-plan-2",
			Name:        "Fake Service Plan 2 Name",
			Description: "Fake Service Plan 2 Description",
		}

		service1 = Service{
			ID:          "fake-service-1",
			Name:        "Fake Service 1 Name",
			Description: "Fake Service 1 Description",
			Plans:       []ServicePlan{servicePlan1},
		}

		service2 = Service{
			ID:          "fake-service-2",
			Name:        "Fake Service 2 Name",
			Description: "Fake Service 2 Description",
			Plans:       []ServicePlan{servicePlan2},
		}

		catalog = Catalog{
			Services: []Service{service1, service2},
		}
	})

	Describe("Validate", func() {
		It("does not return error if all fields are valid", func() {
			err := catalog.Validate()

			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if Services are not valid", func() {
			catalog.Services = []Service{
				Service{},
			}

			err := catalog.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Services configuration"))
		})
	})

	Describe("FindService", func() {
		It("returns true and the Service if it is found", func() {
			service, found := catalog.FindService("fake-service-1")
			Expect(found).To(BeTrue())
			Expect(service).To(Equal(service1))
		})

		It("returns false if it is not found", func() {
			_, found := catalog.FindService("fake-service-3")
			Expect(found).To(BeFalse())
		})
	})

	Describe("FindServicePlan", func() {
		It("returns true and the Service Plan if it is found", func() {
			plan, found := catalog.FindServicePlan("fake-service-1", "fake-service-plan-1")
			Expect(found).To(BeTrue())
			Expect(plan).To(Equal(servicePlan1))
		})

		It("returns false if it is not found", func() {
			_, found := catalog.FindServicePlan("fake-service-1", "ffake-service-plan-2")
			Expect(found).To(BeFalse())
		})
	})
})

var _ = Describe("Service", func() {
	var (
		service Service
	)

	BeforeEach(func() {
		service = Service{
			ID:              "fake-service",
			Name:            "Fake Service Name",
			Description:     "Fake Service Description",
			Tags:            []string{"service"},
			Requires:        []string{"syslog"},
			Bindable:        true,
			Metadata:        &ServiceMetadata{},
			DashboardClient: &ServiceDashboardClient{},
			PlanUpdateable:  true,
			Plans: []ServicePlan{
				ServicePlan{
					ID:          "fake-service-plan",
					Name:        "Fake Service Plan Name",
					Description: "Fake Service Plan Description",
					Metadata:    &ServicePlanMetadata{},
					Free:        true,
					Bindable:    true,
				},
			},
		}
	})

	Describe("Validate", func() {
		It("does not return error if all fields are valid", func() {
			err := service.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if ID is empty", func() {
			service.ID = ""

			err := service.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty ID"))
		})

		It("returns error if Name is empty", func() {
			service.Name = ""

			err := service.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Name"))
		})

		It("returns error if Description is empty", func() {
			service.Description = ""

			err := service.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Description"))
		})

		It("returns error if it does not have plans", func() {
			service.Plans = []ServicePlan{}

			err := service.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must contain at least one plan"))
		})

		It("returns error if Plans are not valid", func() {
			service.Plans = []ServicePlan{
				ServicePlan{},
			}

			err := service.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Validating Plans configuration for Service"))
		})
	})
})

var _ = Describe("ServicePlan", func() {
	var (
		servicePlan ServicePlan
	)

	BeforeEach(func() {
		servicePlan = ServicePlan{
			ID:          "fake-service-plan",
			Name:        "Fake Service Plan Name",
			Description: "Fake Service Plan Description",
			Metadata:    &ServicePlanMetadata{},
			Free:        true,
			Bindable:    true,
		}
	})

	Describe("Validate", func() {
		It("does not return error if all fields are valid", func() {
			err := servicePlan.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if ID is empty", func() {
			servicePlan.ID = ""

			err := servicePlan.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty ID"))
		})

		It("returns error if Name is empty", func() {
			servicePlan.Name = ""

			err := servicePlan.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Name"))
		})

		It("returns error if Description is empty", func() {
			servicePlan.Description = ""

			err := servicePlan.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Must provide a non-empty Description"))
		})
	})
})
