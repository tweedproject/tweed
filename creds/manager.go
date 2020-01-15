package creds

import "log"

type Manager interface {
	IsConfigured() bool
	Validate() error
	Health() (*HealthResponse, error)
	Init(log.Logger) error
	Close(log.Logger)

	NewSecretsFactory(log.Logger) (SecretsFactory, error)
}

type HealthResponse struct {
	Response interface{} `json:"response,omitempty"`
	Error    string      `json:"error,omitempty"`
	Method   string      `json:"method,omitempty"`
}
