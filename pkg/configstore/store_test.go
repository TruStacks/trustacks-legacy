package configstore

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigReadWrite(t *testing.T) {
	previousDataDir := dataDir
	dataDir = os.TempDir()
	defer func() {
		dataDir = previousDataDir
	}()
	if err := os.MkdirAll(fmt.Sprintf("/%s/test-read-write", dataDir), 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fmt.Sprintf("/%s/test-read-write", dataDir))
	if err := writeConfig("test", map[string]string{"key": "value"}, "test-read-write/config", "system"); err != nil {
		t.Fatal(err)
	}
	config, err := readConfig("test", "test-read-write/config")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "system", config["_aud"])
	assert.Equal(t, "value", config["key"])
}

func TestValueEncryption(t *testing.T) {
	config := map[string]string{
		"secret1": "abc123",
		"secret2": "123xyz",
		"secret3": "abcxyz",
	}
	encrypted, err := encryptValues("password", config)
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := decryptValues("password", encrypted)
	if err != nil {

		t.Fatal(err)
	}
	assert.Equal(t, "abc123", decrypted["secret1"], "got an unexpected secret value")
	assert.Equal(t, "123xyz", decrypted["secret2"], "got an unexpected secret value")
	assert.Equal(t, "abcxyz", decrypted["secret3"], "got an unexpected secret value")
}
