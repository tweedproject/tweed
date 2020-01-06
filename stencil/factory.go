package stencil

import (
	"path"
	"sync"
)

const (
	stencilsSubDir   = "run/stencils"
	stencilsDepotDir = "run/depot"
)

type Factory struct {
	stencils []*Stencil
	registry *registry
	mux      sync.Mutex
}

func NewFactory(root string) *Factory {
	registry := registry{
		stencilsDir: path.Join(root, stencilsSubDir),
	}
	stencils := make([]*Stencil, 0)
	return &Factory{
		registry: &registry,
		stencils: stencils,
	}
}

// Load adds a Stencil to the in memory list of stencils
// it also starts an asynchronous load of the oci bundle onto the file system
// it is safe to use this method from multiple goroutines
func (f *Factory) Load(ref string) {
	f.Get(ref).loadBundle()
}

// Get retrieves a pointer to a Stencil by reference from the internal stencils list
// if it does not exist it will add a new Stencil to the list
// it is safe to use this method from multiple goroutines
func (f *Factory) Get(ref string) *Stencil {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, s := range f.stencils {
		if s.Reference == ref {
			return s
		}
	}
	s := NewStencil(ref, f.registry)
	f.stencils = append(f.stencils, s)
	return s
}
