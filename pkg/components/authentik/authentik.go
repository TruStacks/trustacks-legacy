package authentik

import "fmt"

const (
	componentName = "authentik"
	repository    = "https://charts.goauthentik.io"
	chart         = "authentik/authentik"
)

type Authentik struct {
	skipTLS bool
	domain  string
	port    uint16
	service string
}

func (c *Authentik) Endpoint() string {
	var scheme string
	if c.skipTLS {
		scheme = "http"
	} else {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s.%s:%d/application/o/%s/", scheme, componentName, c.domain, c.port, c.service)
}

func (c *Authentik) WithDomain(domain string) *Authentik {
	c.domain = domain
	return c
}

func (c *Authentik) WithService(service string) *Authentik {
	c.service = service
	return c
}

func (c *Authentik) WithPort(port uint16) *Authentik {
	c.port = port
	return c
}

func (c *Authentik) WithSkipTLS() *Authentik {
	c.skipTLS = true
	return c
}

func (c *Authentik) Values() string {
	return ""
}

func New() *Authentik {
	return &Authentik{port: 443}
}
