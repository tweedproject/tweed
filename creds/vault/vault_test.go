package vault_test

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tweedproject/tweed/creds/vault"
)

type MockSecret struct {
	path   string
	secret *vaultapi.Secret
}

type MockClient struct {
	secrets []MockSecret
}

func (msr *MockClient) Read(lookupPath string) (*vaultapi.Secret, error) {
	Expect(lookupPath).ToNot(BeNil())

	for _, secret := range msr.secrets {
		if lookupPath == secret.path {
			return secret.secret, nil
		}
	}

	return nil, nil
}

func (c *MockClient) Write(lookupPath string, value map[string]interface{}) (*vaultapi.Secret, error) {
	Expect(lookupPath).ToNot(BeNil())
	s := vaultapi.Secret{Data: value}
	c.secrets = append(c.secrets, MockSecret{path: lookupPath, secret: &s})
	return &s, nil
}

var _ = Describe("Vault", func() {
	var v *vault.Vault
	var c *MockClient

	BeforeEach(func() {
		c = &MockClient{[]MockSecret{}}
		v = &vault.Vault{Client: c}
	})

	Describe("Get()", func() {
		It("should get secret from path with prefix", func() {
			v.Client = &MockClient{[]MockSecret{
				{
					path: "/creds/foo",
					secret: &vaultapi.Secret{
						Data: map[string]interface{}{"value": "bar"},
					},
				}},
			}
			v.Prefix = "/creds"
			value, found, err := v.Get("foo")
			Expect(value).To(BeEquivalentTo("bar"))
			Expect(found).To(BeTrue())
			Expect(err).To(BeNil())
		})

		It("should get secret with full path without Prefix", func() {
			v.Prefix = "/"
			v.Client = &MockClient{[]MockSecret{
				{
					path: "/creds/bar/foo",
					secret: &vaultapi.Secret{
						Data: map[string]interface{}{"value": "bar"},
					},
				}},
			}
			value, found, err := v.Get("/creds/bar/foo")
			Expect(value).To(BeEquivalentTo("bar"))
			Expect(found).To(BeTrue())
			Expect(err).To(BeNil())
		})

	})

	Describe("Set()", func() {
		It("should set secret with path in prefix", func() {
			v.Prefix = "/creds"
			err := v.Set("foo", "bar")
			Expect(err).To(BeNil())
			Expect(c.secrets[0].path).To(Equal("/creds/foo"))
			Expect(c.secrets[0].secret.Data).To(BeEquivalentTo(map[string]interface{}{"value": "bar"}))
		})

		It("should set secret with full path without Prefix", func() {
			err := v.Set("/creds/bar/foo", "bar")
			Expect(err).To(BeNil())
			Expect(c.secrets[0].path).To(Equal("/creds/bar/foo"))
			Expect(c.secrets[0].secret.Data).To(BeEquivalentTo(map[string]interface{}{"value": "bar"}))
		})

	})
})
