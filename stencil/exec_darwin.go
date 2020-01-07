package stencil

import (
	"errors"
)

func (f *Factory) Run(stencilRef string, a RunArg) error {
	return errors.New("runc is not supported on macOS")
}
