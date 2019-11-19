package tweed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Vault struct {
		Prefix string `json:"prefix"`
		URL    string `json:"url"`
		Token  string `json:"token"`
	} `json:"vault"`

	Prefix          string                    `json:"prefix"`
	Infrastructures map[string]Infrastructure `json:"infrastructures"`
	Catalog         Catalog                   `json:"catalog"`
}

func ParseConfig(b []byte) (Config, error) {
	c := Config{}
	c.Vault.Prefix = "secret/tweed"

	err := json.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}

	for i, s := range c.Catalog.Services {
		for j, p := range s.Plans {
			s.Plans[j].Service = &c.Catalog.Services[i]
			s.Plans[j].Tweed.Reference = fmt.Sprintf("%s/%s", s.ID, p.ID)
		}
	}

	return c, nil
}

func ReadConfig(path string) (Config, error) {
	c := Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	return ParseConfig(b)
}
