package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIV1GenerateToolchainConfig(t *testing.T) {
	spec := `ingressPort: "443"
network: "public"
tls: "true"
`
	tf, err := os.CreateTemp("", "toolchain-parameters")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tf.Name())
	if _, err := tf.Write([]byte(spec)); err != nil {
		t.Fatal(err)
	}
	path, config, err := (&apiV1{}).generateToolchainConfig("test", "https://local.gd/test/toolchain", tf.Name(), map[string]interface{}{"domain": "local.gd"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected the config file to exist")
	}
	params := config["parameters"].(map[string]interface{})
	assert.Equal(t, "test", config["name"].(string), "got an unexpected config name")
	assert.Equal(t, "https://local.gd/test/toolchain", config["source"].(string), "got an unexpected config source")
	assert.Equal(t, "443", params["ingressPort"], "got an unexpected config ingressPort parameter")
	assert.Equal(t, "public", params["network"], "got an unexpected config network parameter")
	assert.Equal(t, "true", params["tls"], "got an unexpected config tls parameter")
	assert.Equal(t, "test.local.gd", params["domain"], "got an unexpected config domain parameter")
}
