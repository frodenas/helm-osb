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
	contextLogKey                = "context"
	instanceIDLogKey             = "instance-id"
	bindingIDLogKey              = "binding-id"
	detailsLogKey                = "details"
	asyncAllowedLogKey           = "async-allowed"
	operationDataLogKey          = "operation-data"
	provisionedServiceSpecLogKey = "provisioned-service-spec"
	deprovisionServiceSpecLogKey = "deprovision-service-spec"
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
	b.logger.Debug("services", lager.Data{
		contextLogKey: ctx,
	})

	services := []brokerapi.Service{}

	brokerCatalog, err := json.Marshal(b.config.Catalog)
	if err != nil {
		b.logger.Error("marshal-error", err)
		return services
	}

	if err = json.Unmarshal(brokerCatalog, services); err != nil {
		b.logger.Error("unmarshal-error", err)
		return services
	}

	return services
}

func (b *Broker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	b.logger.Debug("provision", lager.Data{
		contextLogKey:      ctx,
		instanceIDLogKey:   instanceID,
		detailsLogKey:      details,
		asyncAllowedLogKey: asyncAllowed,
	})

	provisionedServiceSpec := brokerapi.ProvisionedServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return provisionedServiceSpec, brokerapi.ErrAsyncRequired
	}

	provisionParameters := ProvisionParameters{}
	if b.config.AllowUserProvisionParameters {
		if err := mapstructure.Decode(details.RawParameters, &provisionParameters); err != nil {
			return provisionedServiceSpec, fmt.Errorf("Error parsing provision parameters: %s", err)
		}
	}

	servicePlan, ok := b.config.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if !ok {
		return provisionedServiceSpec, fmt.Errorf("Plan `%s` for Service `%s` not found in Catalog", details.PlanID, details.ServiceID)
	}

	if err := b.helmClient.Install(instanceID, servicePlan.Metadata.Helm.Chart, servicePlan.Metadata.Helm.Repository, servicePlan.Metadata.Helm.Version); err != nil {
		return provisionedServiceSpec, err
	}

	b.logger.Debug("provision", lager.Data{
		provisionedServiceSpecLogKey: provisionedServiceSpec,
	})

	return provisionedServiceSpec, nil
}

func (b *Broker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	b.logger.Debug("update", lager.Data{
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

	return updateServiceSpec, nil
}

func (b *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	b.logger.Debug("deprovision", lager.Data{
		contextLogKey:      ctx,
		instanceIDLogKey:   instanceID,
		detailsLogKey:      details,
		asyncAllowedLogKey: asyncAllowed,
	})

	deprovisionServiceSpec := brokerapi.DeprovisionServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return deprovisionServiceSpec, brokerapi.ErrAsyncRequired
	}

	if err := b.helmClient.Delete(instanceID); err != nil {
		return deprovisionServiceSpec, err
	}

	b.logger.Debug("deprovision", lager.Data{
		deprovisionServiceSpecLogKey: deprovisionServiceSpec,
	})

	return deprovisionServiceSpec, nil
}

func (b *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	b.logger.Debug("bind", lager.Data{
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

	return binding, nil
}

func (b *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	b.logger.Debug("unbind", lager.Data{
		contextLogKey:    ctx,
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
		detailsLogKey:    details,
	})

	// TODO

	return nil
}

func (b *Broker) LastOperation(ctx context.Context, instanceID string, operationData string) (brokerapi.LastOperation, error) {
	b.logger.Debug("last-operation", lager.Data{
		contextLogKey:       ctx,
		instanceIDLogKey:    instanceID,
		operationDataLogKey: operationData,
	})

	lastOperation := brokerapi.LastOperation{State: brokerapi.Failed}

	// TODO

	return lastOperation, nil
}
