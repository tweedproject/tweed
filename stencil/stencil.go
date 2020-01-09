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
	go s.bundle.once.Do(func() {
		fmt.Println("loading bundle")
		bundle, err := s.registry.loadStencilBundle(s.Reference)
		if err != nil {
			panic(fmt.Errorf("failed to create stencil bundle :%s", err))
		}
		fmt.Println("bundle loaded")
		s.bundle.path = bundle
		close(s.bundle.ready)
	})
}

func (s *Stencil) bundlePath() string {
	if s.bundle.path != "" {
		return s.bundle.path
	}
	<-s.bundle.ready
	return s.bundle.path
}

func (s *Stencil) bundleConfig() string {
	return path.Join(s.bundlePath(), "config.json")
}

func (s *Stencil) rootfsPath() string {
	return path.Join(s.bundlePath(), "rootfs")
}
