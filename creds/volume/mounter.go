package volume

import (
	"context"

	"github.com/tweedproject/tweed/creds"
)

type Mounter struct {
	secrets creds.Secrets
}

func NewMounter(s creds.Secrets) *Mounter {
	return &Mounter{
		secrets: s,
	}
}

func (m *Mounter) Mount(ctx context.Context, target, secret string) (*Volume, error) {
	volume := Volume{
		secrets: m.secrets,
		secret:  secret,
		target:  target,
	}
	err := volume.mount(ctx)
	if err != nil {
		return nil, err
	}
	return &volume, nil
}
