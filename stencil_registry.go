package tweed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/distribution/registry/client"
)

const (
	localRegistry  = "localhost:5000"
	stencilsSubDir = "run/stencils"
)

func StencilExists(stencil string) (bool, error) {
	rc, err := client.NewRegistry(localRegistry, http.DefaultTransport)
	if err != nil {
		return false, fmt.Errorf("failed to connect to local registry: %s", err)
	}

	entries := make([]string, 2)
	ctx := context.Background()
	numFilled, err := rc.Repositories(ctx, entries, stencil)
	if err != io.EOF {
		return false, fmt.Errorf("unable to find stencil: %s in local registry: %s", stencil, err)
	}

	if numFilled != 1 {
		return false, fmt.Errorf("found more then 1 stencil for: %s in local registry: %s", stencil)
	}
	return true, nil
}

func (c *Core) LoadCatalogStencils() error {
	errors := make([]error, 0)
	for _, s := range c.Config.Catalog.Services {
		for _, p := range s.Plans {
			_, err := c.GetOrMakeStencilRootfs(p.Tweed.Stencil)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return Errors{all: errors}
	}
	return nil
}

func isNonEmptyDir(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return false
	}
	return true
}
