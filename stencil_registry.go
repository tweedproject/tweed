package tweed

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/distribution/registry/client"
)

func StencilExists(stencil string) (bool, error) {
	rc, err := client.NewRegistry("localhost:5000", http.DefaultTransport)
	if err != nil {
		return false, fmt.Errorf("failed to connect to local registry: %s", err)
	}

	entries := make([]string, 2)
	ctx := context.Background()
	numFilled, err := rc.Repositories(ctx, entries, stencil)
	if err != io.EOF {
		return false, fmt.Errorf("could stencil: %s in local registry: %s", stencil, err)
	}

	if numFilled != 1 {
		return false, fmt.Errorf("found more then 1 stencil for: %s in local registry: %s", stencil)
	}
	return true, nil
}
