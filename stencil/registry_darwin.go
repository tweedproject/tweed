package stencil

import "errors"

func (r *registry) loadStencilBundle(stencil string) (string, error) {
	return "", errors.New("Running the broker not supported on macOS because lac of RunC support")
}
