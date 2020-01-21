package tweed

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tweedproject/tweed/creds"
	"github.com/tweedproject/tweed/stencil"
)

type Core struct {
	Root             string
	HTTPAuthUsername string
	HTTPAuthPassword string
	HTTPAuthRealm    string

	Config         Config
	StencilFactory *stencil.Factory
	SecretManager  creds.Secrets

	// FIXME track corrupted instances read at startup
	instances map[string]*Instance

	oops struct {
		max    int
		index  []string
		faults map[string]oops
	}
}

func (c *Core) path(rel string) string {
	return fmt.Sprintf("%s/%s", c.Root, rel)
}

func (c *Core) Scan() error {
	idirs, err := ioutil.ReadDir(c.path("data/instances"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, ifi := range idirs {
		if ifi.IsDir() {
			id := ifi.Name()
			b, err := ioutil.ReadFile(c.path("data/instances/" + id + "/instance.mf"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "found corrupted service instance '%s' -- unable to read %s/instance.mf: %s\n", id, id, err)
				continue
			}

			inst, err := ParseInstance(c.Config.Catalog, c.StencilFactory, c.SecretManager, c.Root, b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "found corrupted service instance '%s' -- unable to parse %s/instance.mf: %s\n", id, id, err)
				continue
			}
			inst.Prefix = c.Config.Prefix
			inst.VaultPrefix = c.Config.Vault.Prefix
			if err := inst.RefreshBindings(); err != nil {
				fmt.Fprintf(os.Stderr, "found corrupted service instance '%s' -- unable to read bindings from the vault: %s\n", id, err)
				continue
			}

			if c.instances == nil {
				c.instances = make(map[string]*Instance)
			}
			c.instances[inst.ID] = &inst
		}
	}

	return nil
}

func (c *Core) Count(plan *Plan) int {
	n := 0
	for _, inst := range c.instances {
		if inst.State != "gone" && plan.Same(inst.Plan) {
			n++
		}
	}
	return n
}

func (c *Core) UpdateCatalogNumbers() {
	for _, s := range c.Config.Catalog.Services {
		for _, p := range s.Plans {
			p.Tweed.Provisioned = 0
		}
	}
	for _, inst := range c.instances {
		if inst.State != "gone" {
			inst.Plan.Tweed.Provisioned++
		}
	}
}
