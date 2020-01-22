package stencil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"

	gorunc "github.com/containerd/go-runc"
	"github.com/opencontainers/runc/libcontainer/specconv"
	gospec "github.com/opencontainers/runtime-spec/specs-go"
)

func (e *Exec) Eval() (*ProcessState, error) {
	spec := specconv.Example()
	specconv.ToRootless(spec)

	spec.Root.Path = e.Stencil.rootfsPath()

	id := uuid.New()

	bundlePath := path.Join(e.Stencil.runc.bundles, id.String())
	err := os.MkdirAll(bundlePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create runc bundle dir: %s", err)
	}
	defer os.RemoveAll(bundlePath)

	tmpDir := path.Join(bundlePath, "tmp")
	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create container tmp dir: %s", err)
	}

	e.Mounts = append(e.Mounts, Mount{
		Source:      tmpDir,
		Destination: "/tmp",
		Writable:    true,
	})
	e.Mounts = append(e.Mounts, Mount{
		Source:      "/etc/resolv.conf",
		Destination: "/etc/resolv.conf",
		Writable:    false,
	})
	for _, mount := range e.Mounts {
		opts := []string{"rbind"}
		if mount.Writable {
			opts = append(opts, "rw")
		}
		spec.Mounts = append(spec.Mounts, gospec.Mount{
			Type:        "none",
			Source:      mount.Source,
			Destination: mount.Destination,
			Options:     opts,
		})
	}

	cwd := "/"
	if e.Cwd != "" {
		cwd = e.Cwd
	}
	spec.Process = &gospec.Process{
		Args: e.Args,
		Env:  e.Env,
		Cwd:  cwd,
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
	pipe, err := gorunc.NewPipeIO(
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

	if e.Stdout != nil {
		go io.Copy(e.Stdout, pipe.Stdout())
	}
	if e.Stderr != nil {
		go io.Copy(e.Stderr, pipe.Stderr())
	}

	code, err := e.Stencil.runc.client.Run(ctx, id.String(), bundlePath, &gorunc.CreateOpts{
		IO: pipe,
	})
	if err != nil && !strings.HasSuffix(err.Error(), "did not terminate successfully") {
		return &ProcessState{
			Exited:   true,
			ExitCode: code,
		}, fmt.Errorf("failed to run lifecyle hook with runc: %s", err)
	}

	return &ProcessState{
		ExitCode: code,
		Exited:   true,
	}, nil

}
