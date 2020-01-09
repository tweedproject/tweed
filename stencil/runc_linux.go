package stencil

import (
	"fmt"

	"github.com/opencontainers/runc/libcontainer"
)

type runc struct {
	factory libcontainer.Factory
}

func newRunc(dir, binPath string) (*runc, error) {
	f, err := libcontainer.New(dir,
		libcontainer.RootlessCgroupfs,
		libcontainer.InitArgs(binPath, RuncInitCmd),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize RunC: %s", err)
	}

	return &runc{factory: f}, nil
}
