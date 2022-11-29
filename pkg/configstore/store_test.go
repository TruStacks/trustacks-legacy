package configstore

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type mockWriteConfigCallback struct {
	called bool
}

func (m *mockWriteConfigCallback) call() error {
	m.called = true
	return nil
}

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
	mockCallback := &mockWriteConfigCallback{}
	if err := writeConfig("test", map[string]string{"key": "value"}, "test-read-write/config", "system", mockCallback.call); err != nil {
		t.Fatal(err)
	}
	config, err := readConfig("test", "test-read-write/config")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "system", config["_aud"])
	assert.Equal(t, "value", config["key"])
	assert.True(t, mockCallback.called)
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

func TestExportValuesToFile(t *testing.T) {
	check := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}
	testVars := map[string]string{
		"first":  "1",
		"second": "2",
	}
	testSecrets := map[string]string{
		"secret1": "asdf1",
		"secret2": "qwerty2",
	}
	dir := os.TempDir()
	mockCallback := mockWriteConfigCallback{}
	err := writeConfig("vars", testVars, "test-read-write/config", "system", mockCallback.call)
	check(err)
	err = writeConfig("secrets", testSecrets, "test-read-write/config", "system", mockCallback.call)
	check(err)
	path := fmt.Sprintf("%s%s", dir, "/config.yaml")
	filePath, err := exportValuesToFile(path)
	check(err)
	defer os.RemoveAll(filePath)
	file, err := os.ReadFile(filePath)
	check(err)
	var config string
	err = yaml.Unmarshal(file, &config)
	check(err)
	assert.Contains(t, config, "qwerty2")
	assert.NotContains(t, config, "random string of text")
}
