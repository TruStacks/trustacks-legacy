package authentik

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthentikEndpoint(t *testing.T) {
	t.Run("local-endpoint", func(t *testing.T) {
		c := New().
			WithDomain("local.gd").
			WithPort(8081).
			WithService("test").
			WithSkipTLS()
		assert.Equal(t, "http://authentik.local.gd:8081/application/o/test/", c.Endpoint())
	})

	t.Run("secure-endpoint", func(t *testing.T) {
		c := New().
			WithDomain("test.trustacks.io").
			WithService("test")
		assert.Equal(t, "https://authentik.test.trustacks.io:443/application/o/test/", c.Endpoint())
	})
}
