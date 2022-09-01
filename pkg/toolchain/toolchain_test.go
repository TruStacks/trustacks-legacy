package toolchain

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func patchToolchainRoot() func() {
	previousToolchainRoot := toolchainRoot
	d, err := os.MkdirTemp("", "toolchain-root")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(d)
	toolchainRoot = d
	return func() {
		toolchainRoot = previousToolchainRoot
	}
}

func TestNewToolchain(t *testing.T) {
	defer patchToolchainRoot()()
	mockPlainClone := func(basePath string, _ bool, _ *git.CloneOptions) (*git.Repository, error) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path.Join(basePath, "config.yaml"), []byte(""), 0644); err != nil {
			return nil, err
		}
		return nil, nil
	}
	tc, err := newToolchain("test", "http://test.com/toolchain-catalog.git", "0.0.0", false, mockPlainClone)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", tc.name, "got an unexpected toolchain name")
}

func TestLoadToolchainConfig(t *testing.T) {
	config, err := loadToolchainConfig(filepath.Join("testdata", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if config.Parameters["test"].(string) != "value" {
		t.Fatal("got an unexpected config parameter value")
	}
}

func TestGetToolchainCatalog(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{
  "hookSource":"quay.io/trustacks/test:latest",
  "components":{
    "test":{
      "repository":"https://test-charts.trustacks.io",
      "chart":"test/test",
      "version":"1.1.1"
    }
  }
}`)); err != nil {
			t.Fatal(err)
		}
	}))
	catalog, err := getToolchainCatalog(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if catalog.HookSource != "quay.io/trustacks/test:latest" {
		t.Fatal("got an unexpected hook source")
	}
	if catalog.Components["test"].Repo != "https://test-charts.trustacks.io" {
		t.Fatal("got an unexpected repo")
	}
	if catalog.Components["test"].Chart != "test/test" {
		t.Fatal("got an unexpected chart")
	}
	if catalog.Components["test"].Version != "1.1.1" {
		t.Fatal("got an unexpected chart")
	}
}

func TestAddComponents(t *testing.T) {
	defer patchToolchainRoot()()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/charts/helloworld-1.0.0.tgz":
			data, err := os.ReadFile("testdata/helloworld-1.0.0.tgz")
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.Write(data)
			if err != nil {
				t.Fatal(err)
			}
		}
	}))
	catalog := &componentCatalog{
		Components: map[string]component{
			"helloworld": {
				Repo:    fmt.Sprintf("%s/charts", ts.URL),
				Chart:   "helloworld",
				Version: "1.0.0",
				Hooks:   "",
			},
		},
	}
	tc := &toolchain{name: "test"}
	if err := tc.addComponents([]string{"helloworld"}, catalog); err != nil {
		t.Fatal(err)
	}
	assert.DirExists(t, fmt.Sprintf("%s/components/helloworld", tc.path()), "expected chart to exist")
}

func TestAddHooks(t *testing.T) {
	defer patchToolchainRoot()()
	hooksManifest, err := os.ReadFile("testdata/hooks.yaml")
	if err != nil {
		t.Fatal(err)
	}
	catalog := &componentCatalog{
		HookSource: "quay.io/trustacks/test-catalog:latest",
		Components: map[string]component{
			"helloworld": {
				Hooks: string(hooksManifest),
			},
		},
	}
	tc := &toolchain{}
	if err := os.MkdirAll(fmt.Sprintf("%s/components", tc.path()), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("cp", "-R", "testdata/helloworld", fmt.Sprintf("%s/components/helloworld", tc.path()))
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if err := tc.addHooks([]string{"helloworld"}, catalog, map[string]interface{}{"testParam": "test"}); err != nil {
		t.Fatal(err)
	}
	assert.FileExists(t, fmt.Sprintf("%s/components/helloworld/templates/trustacks-hooks.yaml", tc.path()), "expected hooks manifest to exist")
}

func TestToolchainAddSubchartValues(t *testing.T) {
	defer patchToolchainRoot()()
	catalog := &componentCatalog{
		HookSource: "quay.io/trustacks/test-catalog:latest",
		Components: map[string]component{
			"helloworld": {
				Values: `username: username
password: password`,
			},
		},
	}
	parameters := map[string]interface{}{
		"username": "username",
		"password": "password",
	}
	tc := &toolchain{}
	if err := os.MkdirAll(path.Join(tc.componentsPath(), "helloworld"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := tc.addSubChartValues([]string{"helloworld"}, catalog, parameters); err != nil {
		t.Fatal(err)
	}
	values, err := os.ReadFile(path.Join(tc.componentsPath(), "helloworld", "override-values.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	expectedValues := `username: username
password: password`

	if string(values) != expectedValues {
		t.Fatal("got an unexpected values output")
	}
}

func TestConfigJoinParameters(t *testing.T) {
	defer patchToolchainRoot()()
	catalogConfig := componentCatalogConfig{
		Parameters: []componentCatalogConfigParameters{
			{Name: "test", Default: ""},
			{Name: "port", Default: "8080"},
		},
	}
	config := &toolchainConfig{
		Parameters: map[string]interface{}{
			"test": "value",
		},
	}
	joined := (&toolchain{}).join(config.Parameters, catalogConfig.Parameters)
	if joined["test"].(string) != "value" {
		t.Fatal("expected test value to be set")
	}
	if joined["port"].(string) != "8080" {
		t.Fatal("expected default port value to be set")
	}
}

func TestCreateAgeKeySecret(t *testing.T) {
	defer patchToolchainRoot()()
	tc := &toolchain{}
	if err := tc.createAgeKeySecret(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path.Join(tc.path(), "chart", "templates", "sops-age-secret.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var secret map[string]interface{}
	if err := yaml.Unmarshal(data, &secret); err != nil {
		t.Fatal(err)
	}
	stringData := secret["stringData"].(map[string]interface{})
	assert.Contains(t, stringData["age.agekey"].(string), "AGE-SECRET-KEY", "got an unexpected age private key")
	assert.Contains(t, stringData["age.agepub"].(string), "age", "got an unexpected age public key")
}

func TestFromConfig(t *testing.T) {
	defer patchToolchainRoot()()
	config := `dependencies:
  - catalog: http://test-catalog.local
`
	if err := os.MkdirAll(filepath.Join(toolchainRoot, "test"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(toolchainRoot, "test", "config.yaml"), []byte(config), 0666); err != nil {
		t.Fatal(err)
	}
	tc, err := newToolchainFromConfig("test")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "http://test-catalog.local", tc.Dependencies[0].Catalog, "got an unexpected dependency catalog")
}
