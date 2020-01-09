package stencil

import (
	"log"
	"path"
	"sync"
)

const (
	stencilsSubDir     = "run/stencils"
	stencilsRuncDir    = "run/runc"
	stencilsBundlesDir = "run/bundles"
)

type Factory struct {
	stencils []*Stencil
	registry *registry
	mux      sync.Mutex
	runc     *runc
	logger   *log.Logger
}

// NewFactory setup a new Factory for managing Stencil bundles and executing
// life-cycle hooks using RunC
func NewFactory(root string, logger *log.Logger) *Factory {
	return &Factory{
		registry: &registry{
			stencilsDir: path.Join(root, stencilsSubDir),
		},
		stencils: make([]*Stencil, 0),
		runc:     newRunc(root),
		logger:   logger,
	}
}

// Load adds a Stencil to the in memory list of stencils
// it also starts an asynchronous load of the oci bundle onto the file system
// it is safe to use this method from multiple goroutines
func (f *Factory) Load(ref string) error {
	s, err := f.Get(ref)
	if err != nil {
		return err
	}
	s.loadBundle()
	return nil
}

// Get retrieves a pointer to a Stencil by reference from the internal stencils list
// if it does not exist it will add a new Stencil to the list
// it is safe to use this method from multiple goroutines
func (f *Factory) Get(ref string) (*Stencil, error) {
	if valid, err := ValidateStencilReference(ref); !valid {
		return nil, err
	}

	f.mux.Lock()
	defer f.mux.Unlock()
	for _, s := range f.stencils {
		if s.Reference == ref {
			return s, nil
		}
	}
	s := NewStencil(ref, f.registry, f.runc)
	f.stencils = append(f.stencils, s)
	return s, nil
}
