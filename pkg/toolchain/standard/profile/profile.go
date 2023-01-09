package profile

import (
	"errors"
)

type Profile struct {
	Domain   string `json:"domain"`
	Port     uint16 `json:"port"`
	Insecure bool   `json:"insecure"`
}

func (p *Profile) Validate() error {
	if p.Domain == "" {
		return errors.New("profile validation error")
	}
	return nil
}
