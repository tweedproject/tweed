package stencil

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"github.com/opencontainers/runc/libcontainer/specconv"
)

func (f *Factory) Run(stencilRef string, a RunArg) error {
	s, err := f.Get(stencilRef)
	if err != nil {
		return err
	}

	spec := specconv.Example()
	specconv.ToRootless(spec)

	id := uuid.New()

	conf, err := specconv.CreateLibcontainerConfig(&specconv.CreateOpts{
		CgroupName:      id.String(),
		Spec:            spec,
		RootlessEUID:    os.Geteuid() != 0,
		RootlessCgroups: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create libcontainer config: %s", err)
	}

	conf.Rootfs = s.rootfsPath()
	conf.Readonlyfs = true
	conf.Mounts = []*configs.Mount{{
		Source:      a.Store,
		Destination: "/workspace",
	}}

	c, err := f.runc.factory.Create(id.String(), conf)
	if err != nil {
		return fmt.Errorf("failed to create runc container: %s", err)
	}

	err = c.Run(&libcontainer.Process{
		Args:   []string{a.Run},
		Env:    a.Env,
		Stdout: a.Stdout,
		Stderr: a.Stderr,
		Init:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to run script in stencil container: %s", err)
	}
	return nil
}
