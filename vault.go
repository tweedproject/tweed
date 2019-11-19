package tweed

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (c *Core) SetupVault() error {
	err := ioutil.WriteFile(c.path(".saferc"), []byte(`
version: 1
current: the-vault
vaults:
  the-vault:
    url:          `+c.Config.Vault.URL+`
    token:        `+c.Config.Vault.Token+`
    no_strongbox: true
    skip_verify:  true
`), 0600)
	if err != nil {
		return fmt.Errorf("failed to write %s: %s\n", c.path(".saferc"), err)
	}

	err = ioutil.WriteFile(c.path(".svtoken"), []byte(`
vault: `+c.Config.Vault.URL+`
token: `+c.Config.Vault.Token+`
skip_verify: true
`), 0600)
	if err != nil {
		return fmt.Errorf("failed to write %s: %s\n", c.path(".svtoken"), err)
	}

	return nil
}

func (c *Core) WaitForVault() {
	_, err := run1(Exec{
		Run: c.path("bin/await-vault"),
		Env: []string{
			"HOME=" + c.path(""),
			"PATH=" + os.Getenv("PATH"),
			"LANG=" + os.Getenv("LANG"),
			"VAULT=" + c.Config.Vault.Prefix,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to wait for vault: %s\n", err)
		os.Exit(1)
	}
}
