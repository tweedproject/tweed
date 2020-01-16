package stencil

import (
	"path"

	gorunc "github.com/containerd/go-runc"
)

type runc struct {
	client  *gorunc.Runc
	bundles string
}

func newRunc(root string) *runc {
	return &runc{
		bundles: path.Join(root, stencilsBundlesDir),
		client: &gorunc.Runc{
			Root: path.Join(root, stencilsRuncDir),
		},
	}
}
