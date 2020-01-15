package vault

import (
	"path"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// A Client reads a vault secret from the given path. It should
// be thread safe!
type Client interface {
	Read(path string) (*vaultapi.Secret, error)
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
}

// Vault converts a vault secret to our completely untyped secret
// data.
type Vault struct {
	Client Client
	Prefix string
}

// Get retrieves the value and expiration of an individual secret
func (v Vault) Get(secretPath string) (interface{}, bool, error) {
	secret, _, found, err := v.findSecret(secretPath)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	val, found := secret.Data["value"]
	if found {
		return val, true, nil
	}

	return secret.Data, true, nil
}

// Get retrieves the value and expiration of an individual secret
func (v Vault) Set(secretPath string, value interface{}) error {
	return v.writeSecret(secretPath, map[string]interface{}{
		"value": value,
	})
}

func (v Vault) findSecret(p string) (*vaultapi.Secret, *time.Time, bool, error) {
	secret, err := v.Client.Read(path.Join(v.Prefix, p))
	if err != nil {
		return nil, nil, false, err
	}

	if secret != nil {
		// The lease duration is TTL: the time in seconds for which the lease is valid
		// A consumer of this secret must renew the lease within that time.
		duration := time.Duration(secret.LeaseDuration) * time.Second / 2
		expiration := time.Now().Add(duration)
		return secret, &expiration, true, nil
	}

	return nil, nil, false, nil
}

func (v Vault) writeSecret(p string, data map[string]interface{}) error {
	_, err := v.Client.Write(path.Join(v.Prefix, p), data)
	return err
}
