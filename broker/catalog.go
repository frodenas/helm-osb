package broker

import (
	"fmt"
)

type Catalog struct {
	Services []Service `json:"services,omitempty"`
}

type Service struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Tags            []string                `json:"tags,omitempty"`
	Requires        []string                `json:"requires,omitempty"`
	Bindable        bool                    `json:"bindable"`
	Metadata        *ServiceMetadata        `json:"metadata,omitempty"`
	DashboardClient *ServiceDashboardClient `json:"dashboard_client,omitempty"`
	PlanUpdateable  bool                    `json:"plan_updateable,omitempty"`
	Plans           []ServicePlan           `json:"plans"`
}

type ServiceMetadata struct {
	DisplayName         string `json:"displayName,omitempty"`
	ImageURL            string `json:"imageUrl,omitempty"`
	LongDescription     string `json:"longDescription,omitempty"`
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
	DocumentationURL    string `json:"documentationUrl,omitempty"`
	SupportURL          string `json:"supportUrl,omitempty"`
}

type ServiceDashboardClient struct {
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type ServicePlan struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Metadata    *ServicePlanMetadata `json:"metadata,omitempty"`
	Free        bool                 `json:"free,omitempty"`
	Bindable    bool                 `json:"bindable,omitempty"`
}

type ServicePlanMetadata struct {
	DisplayName string            `json:"displayName,omitempty"`
	Bullets     []string          `json:"bullets,omitempty"`
	Costs       []ServicePlanCost `json:"costs,omitempty"`
	Helm        HelmConfig        `json:"helm"`
}

type HelmConfig struct {
	Chart      string `json:"chart"`
	Repository string `json:"repository,omitempty"`
	Version    string `jsob:"version,omitempty"`
}

type ServicePlanCost struct {
	Amount map[string]float64 `json:"amount,omitempty"`
	Unit   string             `json:"unit,omitempty"`
}

func (c Catalog) Validate() error {
	for _, service := range c.Services {
		if err := service.Validate(); err != nil {
			return fmt.Errorf("Validating Services configuration: %s", err)
		}
	}

	return nil
}

func (c Catalog) FindService(serviceID string) (service Service, found bool) {
	for _, service := range c.Services {
		if service.ID == serviceID {
			return service, true
		}
	}

	return service, false
}

func (c Catalog) FindServicePlan(serviceID string, planID string) (plan ServicePlan, found bool) {
	service, ok := c.FindService(serviceID)
	if !ok {
		return plan, false
	}

	for _, plan := range service.Plans {
		if plan.ID == planID {
			return plan, true
		}
	}

	return plan, false
}

func (s Service) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("Must provide a non-empty ID (%+v)", s)
	}

	if s.Name == "" {
		return fmt.Errorf("Must provide a non-empty Name (%+v)", s)
	}

	if s.Description == "" {
		return fmt.Errorf("Must provide a non-empty Description (%+v)", s)
	}

	if len(s.Plans) == 0 {
		return fmt.Errorf("Must contain at least one plan (%+v)", s)
	}

	for _, servicePlan := range s.Plans {
		if err := servicePlan.Validate(); err != nil {
			return fmt.Errorf("Validating Plans configuration for Service `%s`: %s", s.Name, err)
		}
	}

	return nil
}

func (sp ServicePlan) Validate() error {
	if sp.ID == "" {
		return fmt.Errorf("Must provide a non-empty ID (%+v)", sp)
	}

	if sp.Name == "" {
		return fmt.Errorf("Must provide a non-empty Name (%+v)", sp)
	}

	if sp.Description == "" {
		return fmt.Errorf("Must provide a non-empty Description (%+v)", sp)
	}

	if sp.Metadata == nil {
		return fmt.Errorf("Must provide a non-empty Helm configuration (%+v)", sp)
	}

	if err := sp.Metadata.Helm.Validate(); err != nil {
		return fmt.Errorf("Validating Helm configuration for Service Plan `%s`: %s", sp.Name, err)
	}

	return nil
}

func (hc HelmConfig) Validate() error {
	if hc.Chart == "" {
		return fmt.Errorf("Must provide a non-empty Chart (%+v)", hc)
	}

	return nil
}
