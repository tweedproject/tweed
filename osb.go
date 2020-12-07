package tweed

import (
	"context"
	"encoding/json"
	"fmt"
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
		var servicePlans []brokerapi.ServicePlan
		for _, p := range s.Plans {
			sp := brokerapi.ServicePlan{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Free:        &p.Free,
				Bindable:    &p.Bindable,
			}
			servicePlans = append(servicePlans, sp)
		}
		dc := brokerapi.Service{
			Name:                 s.Name,
			ID:                   s.ID,
			Description:          s.Description,
			Tags:                 s.Tags,
			Bindable:             s.Bindable,
			InstancesRetrievable: s.InstancesRetrievable,
			BindingsRetrievable:  s.BindingsRetrievable,
			PlanUpdatable:        s.PlanUpdateable,
			Plans:                servicePlans,
		}
		services = append(services, dc)
	}
	return services, nil
}

func (b broker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	// FIXME: Uncomment after prefixing is done
	// if err := ValidInstanceID(instanceID); err != nil {
	// 	return brokerapi.ProvisionedServiceSpec{}, err
	// }
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
	_, gone, err := b.c.Deprovision(instanceID)
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	if gone {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("service instance '%s' already deprovisioned", instanceID)
	}

	return brokerapi.DeprovisionServiceSpec{
		IsAsync: true,
	}, nil
}

func (b broker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	inst, ok := b.c.instances[instanceID]
	if !ok {
		return brokerapi.GetInstanceDetailsSpec{}, fmt.Errorf("service instance '%s' not found", instanceID)
	}

	return brokerapi.GetInstanceDetailsSpec{
		PlanID:     inst.Plan.ID,
		Parameters: inst.UserParameters,
	}, nil
}

func (b broker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	panic("not implemented")
}

func (b broker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	inst, ok := b.c.instances[instanceID]
	if !ok {
		return brokerapi.LastOperation{State: brokerapi.Failed}, fmt.Errorf("service instance '%s' not found", instanceID)
	}

	if inst.State == "provisioning" || inst.State == "deprovisioning" || inst.State == "binding" || inst.State == "unbinding" {
		return brokerapi.LastOperation{State: brokerapi.InProgress}, nil
	}

	if inst.State == "quiet" || inst.State == "gone" {
		return brokerapi.LastOperation{State: brokerapi.Succeeded}, nil
	}

	return brokerapi.LastOperation{State: brokerapi.Failed}, fmt.Errorf("operation failed: %s", inst.State)
}

func (b broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	_, err := b.c.Bind(instanceID, bindingID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("unable to bind service instance '%s': %s", instanceID, err)
	}

	return brokerapi.Binding{
		IsAsync:       true,
		AlreadyExists: false,
	}, nil
}

func (b broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	_, err := b.c.Unbind(instanceID, bindingID)
	if err != nil {
		return brokerapi.UnbindSpec{}, fmt.Errorf("unable to unbind service instance '%s': %s", instanceID, err)
	}

	return brokerapi.UnbindSpec{
		IsAsync: true,
	}, nil
}

func (b broker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	inst, ok := b.c.instances[instanceID]
	if !ok {
		return brokerapi.GetBindingSpec{}, fmt.Errorf("service instance '%s' not found", instanceID)
	}

	_, ok = inst.Bindings[bindingID]
	if !ok {
		return brokerapi.GetBindingSpec{}, fmt.Errorf("binding '%s' not found for service instance '%s'", bindingID, instanceID)
	}

	return brokerapi.GetBindingSpec{Credentials: inst.Plan.Tweed.Credentials}, nil
}

func (b broker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	inst, ok := b.c.instances[instanceID]
	if !ok {
		return brokerapi.LastOperation{State: brokerapi.Failed}, fmt.Errorf("service instance '%s' not found", instanceID)
	}

	_, ok = inst.Bindings[bindingID]
	if !ok {
		return brokerapi.LastOperation{State: brokerapi.InProgress}, nil
	}
	return brokerapi.LastOperation{State: brokerapi.Succeeded}, nil
}
