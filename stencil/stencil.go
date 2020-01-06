package stencil

import (
	"fmt"
	"sync"
)

type Stencil struct {
	Reference string
	valid     *bool
	bundle    bundle
	registry  *registry
}

func NewStencil(ref string, registry *registry) *Stencil {
	var once sync.Once
	var m sync.Mutex
	cond := sync.NewCond(&m)
	return &Stencil{
		Reference: ref,
		valid:     nil,
		bundle: bundle{
			path: "",
			once: &once,
			cond: cond,
		},
		registry: registry,
	}
}

type bundle struct {
	path    string
	once    *sync.Once
	cond    *sync.Cond
	getPath chan string
}

func (s *Stencil) Valid() (bool, error) {
	return false, nil
}

func (s *Stencil) loadBundle() {
	go s.bundle.once.Do(func() {
		s.bundle.cond.L.Lock()
		rootfs, err := s.registry.loadStencilBundle(s.Reference)
		if err != nil {
			panic(fmt.Errorf("failed to create stencil bundle :%s", err))
		}
		s.bundle.path = rootfs
		s.bundle.cond.Broadcast()
		s.bundle.cond.L.Unlock()
	})
}

func (s *Stencil) bundlePath() string {
	if s.bundle.path != "" {
		return s.bundle.path
	}
	s.loadBundle()
	s.bundle.cond.L.Lock()
	s.bundle.cond.Wait()
	s.bundle.cond.L.Unlock()
	return s.bundle.path
}
