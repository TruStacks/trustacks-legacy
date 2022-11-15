package standard

import (
	"fmt"
	"reflect"

	"github.com/trustacks/trustacks/pkg/components/authentik"
)

type DefaultToolchainConfig struct {
	Domain           string `json:"domain"`
	Local            bool   `json:"local"`
	OIDCClientSecret string `json:"oidcClientSecret" trustacks:"encrypt"`
}

type DefaultToolchain struct {
}

func NewDefaultToolchain(config DefaultToolchainConfig) *DefaultToolchain {
	t := reflect.TypeOf(config)
	f, ok := t.FieldByName("OIDCClientSecret")
	if ok {
		fmt.Println(f)
	}
	tc := &DefaultToolchain{}
	oidcProvider := authentik.New().WithDomain(config.Domain)
	if config.Local {
		oidcProvider.WithSkipTLS()
	}
	return tc
}
