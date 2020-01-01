package tweed

import (
	"fmt"
	"regexp"
)

func ValidInstanceID(id string) error {
	if ok, _ := regexp.Match(`^[a-zA-Z][a-zA-Z-0-9-]{0,61}[a-zA-Z0-9]$`, []byte(id)); !ok {
		return fmt.Errorf("instance IDs must include only alphanumeric characters or hyphens, be less than 64 characters long, must start with a letter, and not end with a hyphen")
	}
	return nil
}

func ValidInstancePrefix(s string) error {
	if s == "" {
		return nil
	}

	if ok, _ := regexp.Match(`^[a-zA-Z0-9][a-zA-Z0-9-]{0,15}$`, []byte(s)); !ok {
		return fmt.Errorf("instance prefixes must include only alphanumeric characters or hyphens, be less than 16 characters long, and must start with a letter")
	}
	return nil
}

func (c *Core) ValidateCatalog() error {
	errors := make([]error, 0)

	for _, s := range c.Config.Catalog.Services {
		for _, p := range s.Plans {
			if p.Tweed.Infrastructure == "" {
				errors = append(errors, fmt.Errorf("service '%s' / '%s' does not specify an infrastructure", s.Name, p.Name))
			} else if !FileExists(c.path("etc/infrastructures/" + p.Tweed.Infrastructure)) {
				errors = append(errors, fmt.Errorf("service '%s' / '%s' specifies unknown infrastructure: %s", s.Name, p.Name, p.Tweed.Infrastructure))
			}

			if p.Tweed.Stencil == "" {
				errors = append(errors, fmt.Errorf("service '%s' / '%s' does not specify a stencil", s.Name, p.Name))
			} else if exists, err := StencilExists(p.Tweed.Stencil); exists {
				errors = append(errors, fmt.Errorf("service '%s' / '%s' specifies unknown stencil: %s", s.Name, p.Name, p.Tweed.Stencil), err)
			}
		}
	}

	if len(errors) > 0 {
		return Errors{all: errors}
	}
	return nil
}
