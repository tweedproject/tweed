package volume

import (
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

func (m *Mounter) Mount(target, secret string) (*Volume, error) {
	mount := Volume{
		secrets: m.secrets,
		secret:  secret,
		target:  target,
	}
	err := volume.mount()
	if err != nil {
		return nil, err
	}
	return &volume, nil
}
