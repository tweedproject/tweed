package stencil

import (
	"fmt"
	"path"
	"sync"
)

type Stencil struct {
	Reference string
	valid     *bool
	bundle    bundle
	registry  *registry
	runc      *runc
}

func NewStencil(ref string, registry *registry, runc *runc) *Stencil {
	var once sync.Once
	return &Stencil{
		Reference: ref,
		valid:     nil,
		bundle: bundle{
			path:  "",
			once:  &once,
			ready: make(chan interface{}),
		},
		registry: registry,
		runc:     runc,
	}
}

type bundle struct {
	path  string
	once  *sync.Once
	ready chan interface{}
}

func (s *Stencil) loadBundle() {
	// function below will only be called once per sencil
	go s.bundle.once.Do(func() {
		bundle, err := s.registry.loadStencilBundle(s.Reference)
		if err != nil {
			panic(fmt.Errorf("failed to create stencil bundle :%s", err))
		}
		s.bundle.path = bundle
		// close the ready channel so all channel subscribers can go ahead
		// and read the bundle path
		close(s.bundle.ready)
	})
}

func (s *Stencil) bundlePath() string {
	if s.bundle.path != "" {
		return s.bundle.path
	}
	// this function is safe to call multiple times due to the internal usages of sync.Once
	s.loadBundle()
	// wait for the channel to be closed (signals bundle.path has been written)
	<-s.bundle.ready
	return s.bundle.path
}

func (s *Stencil) bundleConfig() string {
	return path.Join(s.bundlePath(), "config.json")
}

func (s *Stencil) rootfsPath() string {
	return path.Join(s.bundlePath(), "rootfs")
}
