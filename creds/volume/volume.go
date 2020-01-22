package volume

import (
	"context"

	"github.com/blang/vfs/memfs"
	"github.com/tweedproject/tweed/creds"
	"golang.org/x/tools/godoc/vfs"
)

type Volume struct {
	secrets creds.Secrets
	secret  string
	target  string
	root    *credRoot
}

func (v *Volume) Close() error {
	return nil
}

func (v *Volume) mount(ctx context.Context) error {
	root, err := createRoot(ctx, v.secrets, v.secret)
	if err != nil {
		return err
	}
	v.root = root
	fs := mountfs.Create(vfs.OS())
	return fs.Mount(memfs.Create(), v.target)
}
