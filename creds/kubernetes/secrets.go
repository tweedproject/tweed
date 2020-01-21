package kubernetes

import (
	"fmt"

	"code.cloudfoundry.org/lager"

	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Secrets struct {
	logger lager.Logger

	client    kubernetes.Interface
	namespace string
}

// Get retrieves the value and expiration of an individual secret
func (secrets Secrets) Get(secretPath string) (interface{}, bool, error) {
	secret, found, err := secrets.findSecret(secrets.namespace, secretPath)
	if err != nil {
		secrets.logger.Error("failed-to-fetch-secret", err, lager.Data{
			"namespace":   secrets.namespace,
			"secret-name": secretPath,
		})
		return nil, false, err
	}

	if found {
		return secrets.getValueFromSecret(secret)
	}

	secrets.logger.Info("secret-not-found", lager.Data{
		"namespace":   secrets.namespace,
		"secret-name": secretPath,
	})

	return nil, false, nil
}

// Get retrieves the value and expiration of an individual secret
func (secrets Secrets) Set(secretPath string, value interface{}) (err error) {
	_, found, _ := secrets.Get(secretPath)

	data := make(map[string][]byte)
	switch v := value.(type) {
	case map[string]interface{}:
		for k, j := range v {
			data[k] = []byte(fmt.Sprintf("%v", j))
		}
	case map[string]string:
		for k, j := range v {
			data[k] = []byte(j)
		}
	default:
		data = map[string][]byte{"value": []byte(fmt.Sprintf("%v", value))}
	}

	if found {
		secrets.logger.Info("update-secret", lager.Data{
			"namespace":   secrets.namespace,
			"secret-name": secretPath,
		})

		err = secrets.updateSecret(secrets.namespace, secretPath, data)
	} else {
		secrets.logger.Info("create-secret", lager.Data{
			"namespace":   secrets.namespace,
			"secret-name": secretPath,
		})

		err = secrets.createSecret(secrets.namespace, secretPath, data)
	}

	return err
}

func (secrets Secrets) getValueFromSecret(secret *v1.Secret) (interface{}, bool, error) {
	val, found := secret.Data["value"]
	if found {
		return string(val), true, nil
	}

	stringified := map[string]interface{}{}
	for k, v := range secret.Data {
		stringified[k] = string(v)
	}

	return stringified, true, nil
}

func (secrets Secrets) findSecret(namespace, name string) (*v1.Secret, bool, error) {
	var secret *v1.Secret
	var err error

	secret, err = secrets.client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})

	if err != nil && k8serr.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	} else {
		return secret, true, err
	}
}

func (secrets Secrets) createSecret(namespace, name string, data map[string][]byte) error {
	_, err := secrets.client.CoreV1().Secrets(namespace).Create(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Data:       data,
	})

	return err
}

func (secrets Secrets) updateSecret(namespace, name string, data map[string][]byte) error {
	_, err := secrets.client.CoreV1().Secrets(namespace).Update(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Data:       data,
	})

	return err
}
