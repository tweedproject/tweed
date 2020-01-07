package stencil

import (
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer"
)

type runc struct {
	factory libcontainer.Factory
}

func newRunc(dir string) (*runc, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get current executable path: %s", err)
	}

	f, err := libcontainer.New(dir,
		libcontainer.RootlessCgroupfs,
		libcontainer.InitArgs(ex, RuncInitCmd),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize RunC: %s", err)
	}

	return &runc{factory: f}, nil
}
