package stencil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"

	gorunc "github.com/containerd/go-runc"
	"github.com/opencontainers/runc/libcontainer/specconv"
	gospec "github.com/opencontainers/runtime-spec/specs-go"
)

func (e *Exec) Eval() (*ProcessState, error) {
	spec := specconv.Example()
	specconv.ToRootless(spec)

	spec.Root.Path = e.Stencil.rootfsPath()

	err := copyResolvConf(spec.Root.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to copy resolv.conf: %s", err)
	}

	for _, mount := range e.Mounts {
		spec.Mounts = append(spec.Mounts, gospec.Mount{
			Type:        "none",
			Source:      mount,
			Destination: mount,
			Options:     []string{"rbind", "rw"},
		})
	}

	id := uuid.New()

	bundlePath := path.Join(e.Stencil.runc.bundles, id.String())
	err = os.MkdirAll(bundlePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create runc bundle dir: %s", err)
	}

	specConfig := path.Join(bundlePath, "config.json")
	data, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal runc spec: %s", err)
	}
	err = ioutil.WriteFile(specConfig, data, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to write runc spec config: %s", err)
	}

	ctx := context.Background()
	cpipe, err := gorunc.NewPipeIO(
		os.Geteuid(),
		os.Getgid(),
		func(opt *gorunc.IOOption) {
			opt.OpenStdin = false
			opt.OpenStdin = true
			opt.OpenStderr = true
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup runc io: %s", err)
	}

	csocket, err := gorunc.NewTempConsoleSocket()
	if err != nil {
		return nil, fmt.Errorf("failed to setup runc socket: %s", err)
	}

	resolv, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to open resolv.conf: %s", err)
	}

	err = e.Stencil.runc.client.Create(ctx, id.String(), bundlePath, &gorunc.CreateOpts{
		IO:            cpipe,
		ConsoleSocket: csocket,
		ExtraFiles:    []*os.File{resolv},
	})
	if err != nil {
		stderr, _ := ioutil.ReadAll(cpipe.Stderr())
		stdout, _ := ioutil.ReadAll(cpipe.Stderr())
		return nil, fmt.Errorf(
			"failed to create runc container: %s\nSTDOUT: %s\nSTDERR: %s",
			err, string(stdout), string(stderr))
	}

	epipe, err := gorunc.NewPipeIO(
		os.Geteuid(),
		os.Getgid(),
		func(opt *gorunc.IOOption) {
			opt.OpenStdin = false
			opt.OpenStdin = true
			opt.OpenStderr = true
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup runc io: %s", err)
	}

	go io.Copy(e.Stdout, epipe.Stdout())
	go io.Copy(e.Stderr, epipe.Stderr())

	err = e.Stencil.runc.client.Exec(ctx, id.String(), gospec.Process{
		Terminal: false,
		Args:     e.Args,
		Env:      e.Env,
		Cwd:      "/",
	}, &gorunc.ExecOpts{
		IO: epipe,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exec process in runc container: %s", err)
	}

	return &ProcessState{
		ExitCode: 0,
		Exited:   true,
	}, nil

}

func copyResolvConf(rootfs string) error {
	in, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(path.Join(rootfs, "/etc/resolv.conf"))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
