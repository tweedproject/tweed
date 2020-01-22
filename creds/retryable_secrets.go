package creds

import (
	"fmt"
	"time"

	"github.com/concourse/retryhttp"
)

type SecretRetryConfig struct {
	Attempts int           `long:"secret-retry-attempts" default:"5"  description:"The number of attempts secret will be retried to be fetched, in case a retryable error happens."`
	Interval time.Duration `long:"secret-retry-interval" default:"1s" description:"The interval between secret retry retrieval attempts."`
}

type RetryableSecrets struct {
	secrets     Secrets
	retryConfig SecretRetryConfig
}

func NewRetryableSecrets(secrets Secrets, retryConfig SecretRetryConfig) Secrets {
	return &RetryableSecrets{secrets: secrets, retryConfig: retryConfig}
}

// Get retrieves the value and expiration of an individual secret
func (rs RetryableSecrets) Get(secretPath string) (interface{}, bool, error) {
	r := &retryhttp.DefaultRetryer{}
	for i := 0; i < rs.retryConfig.Attempts-1; i++ {
		result, exists, err := rs.secrets.Get(secretPath)
		if err != nil && r.IsRetryable(err) {
			time.Sleep(rs.retryConfig.Interval)
			continue
		}
		return result, exists, err
	}
	result, exists, err := rs.secrets.Get(secretPath)
	if err != nil {
		err = fmt.Errorf("%s (after %d retries)", err, rs.retryConfig.Attempts)
	}
	return result, exists, err
}

// Set sets the value and expiration of an individual secret
func (rs RetryableSecrets) Set(secretPath string, value interface{}) error {
	r := &retryhttp.DefaultRetryer{}
	for i := 0; i < rs.retryConfig.Attempts-1; i++ {
		err := rs.secrets.Set(secretPath, value)
		if err != nil && r.IsRetryable(err) {
			time.Sleep(rs.retryConfig.Interval)
			continue
		}
		return err
	}
	err := rs.secrets.Set(secretPath, value)
	if err != nil {
		err = fmt.Errorf("%s (after %d retries)", err, rs.retryConfig.Attempts)
	}
	return err
}
