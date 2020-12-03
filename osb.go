package tweed

import (
	"context"
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (c *Core) OSB() http.Handler {
	return brokerapi.New(c.OSBBroker(), lager.NewLogger("tweed"), brokerapi.BrokerCredentials{
		Username: c.HTTPAuthUsername,
		Password: c.HTTPAuthPassword,
	})
}

func (c *Core) OSBBroker() brokerapi.ServiceBroker {
	return broker{c: c}
}

type broker struct {
	c *Core
}

func (b broker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	var services []brokerapi.Service
	for _, s := range b.c.Config.Catalog.Services {
		var dc brokerapi.Service
		dc.Name = s.Name
		dc.ID = s.ID
		dc.Description = s.Description
		dc.Tags = s.Tags
		dc.Bindable = s.Bindable
		dc.InstancesRetrievable = s.InstancesRetrievable
		dc.BindingsRetrievable = s.BindingsRetrievable
		dc.PlanUpdatable = s.PlanUpdateable
		services = append(services, dc)
	}
	return services, nil
}

func (b broker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	if err := ValidInstanceID(instanceID); err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}
	plan, err := b.c.Config.Catalog.FindPlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	var params map[string]interface{}
	err = json.Unmarshal(details.RawParameters, &params)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	_, err = b.c.Provision(&Instance{
		ID:             instanceID,
		Plan:           plan,
		Root:           b.c.Root,
		Prefix:         b.c.Config.Prefix,
		VaultPrefix:    b.c.Config.Vault.Prefix,
		UserParameters: params,
	})

	return brokerapi.ProvisionedServiceSpec{
		IsAsync:       true,
		AlreadyExists: false,
	}, err
}

func (b broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	panic("not implemented")
}

func (b broker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	panic("not implemented")
}

func (b broker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	panic("not implemented")
}

func (b broker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}

func (b broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	panic("not implemented")
}

func (b broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	panic("not implemented")
}

func (b broker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	panic("not implemented")
}

func (b broker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}
