package stencil

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/distribution/registry/client"
)

const (
	localRegistry = "localhost:5000"
)

type registry struct {
	stencilsDir string
}

func ValidateStencilReference(stencil string) (bool, error) {
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
