package kubernetes

import (
	"code.cloudfoundry.org/lager"
	"k8s.io/client-go/kubernetes"

	"github.com/tweedproject/tweed/creds"
)

type kubernetesFactory struct {
	logger lager.Logger

	client    kubernetes.Interface
	namespace string
}

func NewKubernetesFactory(logger lager.Logger, client kubernetes.Interface, namespace string) *kubernetesFactory {
	factory := &kubernetesFactory{
		logger:    logger,
		client:    client,
		namespace: namespace,
	}

	return factory
}

func (factory *kubernetesFactory) NewSecrets() creds.Secrets {
	return &Secrets{
		logger:    factory.logger,
		client:    factory.client,
		namespace: factory.namespace,
	}
}
