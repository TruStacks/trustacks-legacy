package components

type Component interface {
	SetValue(string, interface{})
}

type OIDCProviderComponent interface {
	Component
	Endpoint() string
	WithDomain(string) *Component
	WithPort(string) *Component
	WithService(string) *Component
}
