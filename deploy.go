package tweed

import (
	"fmt"
)

func (c *Core) Provision(inst *Instance) (string, error) {
	if inst.ID == "" {
		return "", fmt.Errorf("no service instance id specified")
	}

	if c.Count(inst.Plan) >= inst.Plan.Tweed.Limit {
		return "", fmt.Errorf("too many instances of this service have been provisioned")
	}

	t, err := inst.Provision()
	if err != nil {
		return "", err
	}

	if c.instances == nil {
		c.instances = make(map[string]*Instance)
	}
	c.instances[inst.ID] = inst

	return t.id, nil
}

func (c *Core) Bind(id, bid string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("no service instance id specified")
	}
	if bid == "" {
		return "", fmt.Errorf("no service instance binding id specified")
	}

	if inst, ok := c.instances[id]; ok {
		t, err := inst.Bind(bid)
		if err != nil {
			return "", err
		}

		return t.id, nil

	} else {
		return "", fmt.Errorf("instance '%s' not found", id)
	}
}

func (c *Core) Unbind(id, bid string) (string, error) {
	if inst, ok := c.instances[id]; ok {
		t, err := inst.Unbind(bid)
		if err != nil {
			return "", err
		}

		return t.id, nil

	} else {
		return "", fmt.Errorf("instance '%s' not found", id)
	}
}

func (c *Core) Deprovision(id string) (string, bool, error) {
	if inst, ok := c.instances[id]; ok {
		if inst.IsGone() {
			return "", true, nil
		}
		t, err := inst.Deprovision()
		if err != nil {
			return "", false, err
		}

		return t.id, false, nil

	} else {
		return "", false, fmt.Errorf("instance '%s' not found", id)
	}
}
