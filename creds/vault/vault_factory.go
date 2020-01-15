package vault

import (
	"time"

	"github.com/tweedproject/tweed/creds"
)

// The vaultFactory will return a vault implementation of vars.Variables.
type vaultFactory struct {
	client   Client
	prefix   string
	loggedIn <-chan struct{}
}

func NewVaultFactory(client Client, loggedIn <-chan struct{}, prefix string, sharedPath string) *vaultFactory {
	factory := &vaultFactory{
		client:   client,
		prefix:   prefix,
		loggedIn: loggedIn,
	}

	return factory
}

// NewSecrets will block until the loggedIn channel passed to the constructor signals a successful login.
func (factory *vaultFactory) NewSecrets() creds.Secrets {
	select {
	case <-factory.loggedIn:
	case <-time.After(5 * time.Second):
	}

	return &Vault{
		Client: factory.client,
		Prefix: factory.prefix,
	}
}
