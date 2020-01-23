package volume

import (
	"context"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/tweedproject/tweed/creds"
)

type Volume struct {
	secrets creds.Secrets
	secret  string
	target  string
	fssrv   *fuse.Server
}

func (v *Volume) Unmount() error {
	return v.fssrv.Unmount()
}

func (v *Volume) mount(ctx context.Context) error {
	root, err := createRoot(ctx, v.secrets, v.secret)
	if err != nil {
		return err
	}
	fssrv, err := fs.Mount(v.target, root, &fs.Options{
		// MountOptions: fuse.MountOptions{
		//   Debug:      true,
		// },
	})
	if err != nil {
		return err
	}
	v.fssrv = fssrv
	go v.fssrv.Serve()
	err = v.fssrv.WaitMount()
	if err != nil {
		return err
	}
	return nil
}
