package kubernetes

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/tweedproject/tweed/creds"
)

type kubernetesManagerFactory struct{}

func init() {
	creds.Register("kubernetes", NewKubernetesManagerFactory())
}

func NewKubernetesManagerFactory() creds.ManagerFactory {
	return &kubernetesManagerFactory{}
}

func (factory *kubernetesManagerFactory) AddConfig(group *flags.Group) creds.Manager {
	manager := &KubernetesManager{}

	subGroup, err := group.AddGroup("Kubernetes Credential Management", "", manager)
	if err != nil {
		panic(err)
	}

	subGroup.Namespace = "kubernetes"

	return manager
}

func (factory *kubernetesManagerFactory) NewInstance(config interface{}) (creds.Manager, error) {
	return &KubernetesManager{}, nil
}
