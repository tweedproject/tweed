package tweed

import (
	"fmt"
)

type Catalog struct {
	Services []Service `json:"services"`
}

func (c Catalog) FindPlan(service, plan string) (*Plan, error) {
	for _, s := range c.Services {
		if s.ID == service {
			for _, p := range s.Plans {
				if p.ID == plan {
					return p, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("service/plan %s/%s not found", service, plan)
}

func (c *Core) LoadCatalogStencils() {
	for _, s := range c.Config.Catalog.Services {
		for _, p := range s.Plans {
			c.StencilFactory.Load(p.Tweed.Stencil)
		}
	}
}
