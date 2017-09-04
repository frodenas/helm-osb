package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/mitchellh/mapstructure"
	"github.com/pivotal-cf/brokerapi"

	"github.com/frodenas/helm-osb/helm"
)

const (
	contextLogKey       = "context"
	instanceIDLogKey    = "instance-id"
	bindingIDLogKey     = "binding-id"
	detailsLogKey       = "details"
	asyncAllowedLogKey  = "async-allowed"
	operationDataLogKey = "operation-data"
	responseLogKey      = "response"
)

type Broker struct {
	config     Config
	helmClient *helm.Client
	logger     lager.Logger
}

func New(config Config, helmClient *helm.Client, logger lager.Logger) *Broker {
	return &Broker{
		config:     config,
		helmClient: helmClient,
		logger:     logger.Session("broker"),
	}
}

func (b *Broker) Services(ctx context.Context) []brokerapi.Service {
	b.logger.Debug("services-parameters", lager.Data{
		contextLogKey: ctx,
	})

	services := []brokerapi.Service{}

	brokerCatalog, err := json.Marshal(b.config.Catalog.Services)
	if err != nil {
		b.logger.Error("marshal-error", err)
		return services
	}

	if err = json.Unmarshal(brokerCatalog, &services); err != nil {
		b.logger.Error("unmarshal-error", err)
		return services
	}

	b.logger.Debug("services-response", lager.Data{
		responseLogKey: services,
	})

	return services
}

func (b *Broker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	b.logger.Debug("provision-parameters", lager.Data{
		contextLogKey:      ctx,
		instanceIDLogKey:   instanceID,
		detailsLogKey:      details,
		asyncAllowedLogKey: asyncAllowed,
	})

	provisionedServiceSpec := brokerapi.ProvisionedServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return provisionedServiceSpec, brokerapi.ErrAsyncRequired
	}

	servicePlan, ok := b.config.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if !ok {
		return provisionedServiceSpec, fmt.Errorf("Plan `%s` for Service `%s` not found in Catalog", details.PlanID, details.ServiceID)
	}

	provisionParameters := ProvisionParameters{}
	if servicePlan.Metadata.Helm.Values != nil {
		for k, v := range *servicePlan.Metadata.Helm.Values {
			provisionParameters[k] = v
		}
	}

	if b.config.AllowUserProvisionParameters {
		if err := mapstructure.Decode(details.RawParameters, &provisionParameters); err != nil {
			return provisionedServiceSpec, fmt.Errorf("Error parsing provision parameters: %s", err)
		}
	}

	err := b.helmClient.InstallRelease(
		instanceID,
		servicePlan.Metadata.Helm.Chart,
		servicePlan.Metadata.Helm.Repository,
		servicePlan.Metadata.Helm.Version,
		provisionParameters)
	if err != nil {
		return provisionedServiceSpec, err
	}

	b.logger.Debug("provision-response", lager.Data{
		responseLogKey: provisionedServiceSpec,
	})

	return provisionedServiceSpec, nil
}

func (b *Broker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	b.logger.Debug("update-parameters", lager.Data{
		contextLogKey:      ctx,
		instanceIDLogKey:   instanceID,
		detailsLogKey:      details,
		asyncAllowedLogKey: asyncAllowed,
	})

	updateServiceSpec := brokerapi.UpdateServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return updateServiceSpec, brokerapi.ErrAsyncRequired
	}

	updateParameters := UpdateParameters{}
	if b.config.AllowUserUpdateParameters {
		if err := mapstructure.Decode(details.RawParameters, &updateParameters); err != nil {
			return updateServiceSpec, fmt.Errorf("Error parsing update parameters: %s", err)
		}
	}

	// TODO

	b.logger.Debug("update-response", lager.Data{
		responseLogKey: updateServiceSpec,
	})

	return updateServiceSpec, nil
}

func (b *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	b.logger.Debug("deprovision-parameters", lager.Data{
		contextLogKey:      ctx,
		instanceIDLogKey:   instanceID,
		detailsLogKey:      details,
		asyncAllowedLogKey: asyncAllowed,
	})

	deprovisionServiceSpec := brokerapi.DeprovisionServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return deprovisionServiceSpec, brokerapi.ErrAsyncRequired
	}

	if err := b.helmClient.DeleteRelease(instanceID); err != nil {
		return deprovisionServiceSpec, err
	}

	b.logger.Debug("deprovision-response", lager.Data{
		responseLogKey: deprovisionServiceSpec,
	})

	return deprovisionServiceSpec, nil
}

func (b *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	b.logger.Debug("bind-parameters", lager.Data{
		contextLogKey:    ctx,
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})

	binding := brokerapi.Binding{}

	bindParameters := BindParameters{}
	if b.config.AllowUserBindParameters {
		if err := mapstructure.Decode(details.RawParameters, &bindParameters); err != nil {
			return binding, fmt.Errorf("Error parsing bind parameters: %s", err)
		}
	}

	// TODO

	b.logger.Debug("bind-response", lager.Data{
		responseLogKey: binding,
	})

	return binding, nil
}

func (b *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	b.logger.Debug("unbind-parameters", lager.Data{
		contextLogKey:    ctx,
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})

	// TODO

	return nil
}

func (b *Broker) LastOperation(ctx context.Context, instanceID string, operationData string) (brokerapi.LastOperation, error) {
	b.logger.Debug("last-operation-parameters", lager.Data{
		contextLogKey:       ctx,
		instanceIDLogKey:    instanceID,
		operationDataLogKey: operationData,
	})

	lastOperation := brokerapi.LastOperation{State: brokerapi.Failed}

	status, description, err := b.helmClient.ReleaseStatus(instanceID)
	if err != nil {
		return lastOperation, err
	}

	switch status {
	case "SUCCEEDED":
		lastOperation.State = brokerapi.Succeeded
		lastOperation.Description = description
	case "INPROGRESS":
		lastOperation.State = brokerapi.InProgress
	}

	b.logger.Debug("last-operation-response", lager.Data{
		responseLogKey: lastOperation,
	})

	return lastOperation, nil
}
