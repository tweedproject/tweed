package tweed

import (
	"fmt"
)

type Plan struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`

	Metadata       map[string]interface{} `json:"metadata"`
	Free           bool                   `json:"free"`
	Bindable       bool                   `json:"bindable"`
	PlanUpdateable bool                   `json:"plan_updateable"`

	// FIXME: support schemas

	MaintenanceInfo struct {
		Version     string `json:"version"`
		Description string `json:"description"`
	} `json:"maintenance_info"`

	/****** end of OSB API fields ******/

	Service *Service `json:"-"`

	Tweed struct {
		Infrastructure string `json:"infrastructure"`
		Stencil        string `json:"stencil"`

		Credentials map[string]interface{} `json:"credentials"`

		// populated by Tweed, to convey actual instance
		// count information to catalog consumers.
		Provisioned int `json:"provisioned"`
		Limit       int `json:"limit"`

		Config map[string]interface{} `json:"config"`
		Params map[string]Parameter   `json:"params"`

		// populated by Tweed, to make scripting
		// with the API / CLI --json mode easier.
		Reference string `json:"reference"`
	} `json:"tweed"`
}

func (p *Plan) Same(other *Plan) bool {
	return p.ID == other.ID && p.Service.ID == other.Service.ID
}

func (p Plan) ValidateInputs(inputs map[string]string) error {
	var errors []error
	for key, value := range inputs {
		if param, found := p.Tweed.Params[key]; !found {
			errors = append(errors, fmt.Errorf("disallowed parameter '%s'", key))
		} else if err := param.Validate(value); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("errors encountered")
	}
	return nil
}
