package creds

//go:generate counterfeiter . SecretsFactory

type SecretsFactory interface {
	// NewSecrets returns an instance of a secret manager, capable of retrieving individual secrets
	NewSecrets() Secrets
}

//go:generate counterfeiter . Secrets

type Secrets interface {
	// Every credential manager needs to be able to return (secret, exists, error) based on the secret path
	Get(string) (interface{}, bool, error)

	// Every credential manager needs to be able to set (secret) based on the secret path
	Set(string, interface{}) error
}
