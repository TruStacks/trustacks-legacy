package server

import (
	"encoding/json"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateToolchainConfig(t *testing.T) {
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
	path, err := generateToolchainConfig("test", "https://local.gd/test/toolchain", tf.Name(), map[string]interface{}{"domain": "local.gd"})
	if err != nil {
		t.Fatal(err)
	}
	configRaw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal(configRaw, &config); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", config["name"].(string), "got an unexpected config name")
	assert.Equal(t, "https://local.gd/test/toolchain", config["source"].(string), "got an unexpected config source")
	assert.Equal(t, "443", config["parameters"].(map[string]interface{})["ingressPort"], "got an unexpected config ingressPort parameter")
	assert.Equal(t, "public", config["parameters"].(map[string]interface{})["network"], "got an unexpected config network parameter")
	assert.Equal(t, "true", config["parameters"].(map[string]interface{})["tls"], "got an unexpected config tls parameter")
	assert.Equal(t, "test.local.gd", config["parameters"].(map[string]interface{})["domain"], "got an unexpected config domain parameter")
}
