package tweed

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Infrastructure struct {
	Type string `json:"type"`

	URL           string `json:"url"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	CACertificate string `json:"ca_certificate"`

	KubeConfig string `json:"kubeconfig"`
}

func (infra Infrastructure) Valid() error {
	switch infra.Type {
	case "bosh":
		if infra.URL == "" {
			return fmt.Errorf("bosh infrastructure missing BOSH Director URL")
		}
		if infra.ClientID == "" {
			return fmt.Errorf("bosh infrastructure missing BOSH Director Client ID")
		}
		if infra.ClientSecret == "" {
			return fmt.Errorf("bosh infrastructure missing BOSH Director Client Secret")
		}
		if infra.CACertificate == "" {
			return fmt.Errorf("bosh infrastructure missing BOSH Director CA Certificate")
		}
		return nil

	case "kubernetes":
		if infra.KubeConfig == "" {
			return fmt.Errorf("k8s infrastructure missing KubeConfig")
		}
		return nil

	default:
		return fmt.Errorf("unknown infrastructure type")
	}
}

func (infra Infrastructure) Render() (string, error) {
	switch infra.Type {
	case "bosh":
		return `# sourced
export BOSH_ENVIRONMENT=` + infra.URL + `
export BOSH_CLIENT=` + infra.ClientID + `
export BOSH_CLIENT_SECRET=` + infra.ClientSecret + `
export BOSH_CA_CERT='` + infra.CACertificate + `'
`, nil

	case "kubernetes":
		return infra.KubeConfig, nil

	default:
		return "", fmt.Errorf("unknown infrastructure type")
	}
}

func (c *Core) SetupInfrastructures() error {
	os.MkdirAll(c.path("etc/infrastructures"), 0777)
	for name, infra := range c.Config.Infrastructures {
		s, err := infra.Render()
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(c.path("etc/infrastructures/"+name), []byte(s), 0666); err != nil {
			return err
		}
		if err := ioutil.WriteFile(c.path("etc/infrastructures/"+name+".type"), []byte(infra.Type), 0666); err != nil {
			return err
		}
	}
	return nil
}
