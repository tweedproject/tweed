package tweed

import "errors"

func (c *Core) GetOrMakeStencilRootfs(stencil string) (string, error) {
	return "", errors.New("Running the broker not supported on macOS because lac of RunC support")
}
