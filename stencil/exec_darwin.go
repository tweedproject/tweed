package stencil

import (
	"errors"
)

func (e *Exec) Eval() (*ProcessState, error) {
	return nil, errors.New("runc is not supported on macOS")
}
