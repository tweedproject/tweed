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

func (e *Exec) Eval() (*ProcessState, error) {
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
		return nil, fmt.Errorf("failed to create libcontainer config: %s", err)
	}

	conf.Rootfs = e.Stencil.rootfsPath()
	conf.Readonlyfs = true

	e.Mounts = append(e.Mounts, "/etc/resolv.conf")
	for _, mount := range e.Mounts {
		conf.Mounts = append(conf.Mounts, &configs.Mount{
			Source:      mount,
			Destination: mount,
		})
	}

	conf.Networks = []*configs.Network{
		{
			Type:    "loopback",
			Address: "127.0.0.1/0",
			Gateway: "localhost",
		},
	}

	c, err := e.Stencil.runc.factory.Create(id.String(), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create runc container: %s", err)
	}
	defer c.Destroy()

	p := libcontainer.Process{
		Args:   []string{e.Run},
		Env:    e.Env,
		Stdout: e.Stdout,
		Stderr: e.Stderr,
		Init:   true,
	}
	err = c.Run(&p)
	if err != nil {
		return nil, fmt.Errorf("failed to run script in stencil container: %s", err)
	}

	state, err := p.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed while waiting on script running in stencil container: %s", err)
	}

	return &ProcessState{
		ExitCode: state.ExitCode(),
		Exited:   state.Exited(),
	}, nil
}
