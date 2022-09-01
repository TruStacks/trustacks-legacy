package toolchain

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetWorkflowCatalog(t *testing.T) {
	config := `workflows:
- name: react
  dependencies:
  - catalog: http://test.com/workflows-catalog
    components:
    - test
`
	mockPlainClone := func(path string, _ bool, _ *git.CloneOptions) (*git.Repository, error) {
		if err := os.WriteFile(filepath.Join(path, "config.yaml"), []byte(config), 0644); err != nil {
			return nil, err
		}
		return nil, nil
	}
	catalog, err := getWorkflowCatalog("https://test.com/workflows-catalog.git", "0.0.0", mockPlainClone)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "react", catalog.Workflows[0].Name, "got an unexpected workflow name")
	assert.Equal(t, "http://test.com/workflows-catalog", catalog.Workflows[0].Dependencies[0].Catalog, "got an unexpected dependency catalog")
	assert.Equal(t, "test", catalog.Workflows[0].Dependencies[0].Components[0], "got an unexpected component name")
}

func TestApplicationCreateChart(t *testing.T) {
	defer patchToolchainRoot()()
	app := &application{toolchain: &toolchain{name: "test"}, name: "test"}
	if err := app.createChart(); err != nil {
		t.Fatal(err)
	}
	assert.FileExists(t, path.Join(app.path(), "Chart.yaml"), "expected the chart file to exist")
	assert.DirExists(t, path.Join(app.path(), "templates"), "expected the templates dir to exist")
}

func TestNewApplication(t *testing.T) {
	defer patchToolchainRoot()()
	config := &applicationConfig{
		Vars:    map[string]string{"var": "test"},
		Secrets: map[string]string{"secret": "test"},
	}
	app, err := newApplication("test", config, &toolchain{name: "test"}, false)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", app.name, "got an unexpected application name")

	// test create existing application without force
	_, err = newApplication("test", config, &toolchain{name: "test"}, false)
	assert.ErrorContains(t, err, "error: application 'test' already exists", "expected an already exists error")

	// test create existing application with force
	_, err = newApplication("test", config, &toolchain{name: "test"}, true)
	assert.Nil(t, err, "expected error to be nil")
}

func TestApplicationAddVars(t *testing.T) {
	defer patchToolchainRoot()()
	vars := map[string]string{
		"name": "test",
		"port": "8888",
	}
	app := application{toolchain: &toolchain{name: "test"}}
	if err := os.MkdirAll(path.Join(app.path(), "templates"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := app.addVars(vars); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path.Join(app.path(), "templates", "application-configmap.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	v := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	varsData := v["data"].(map[string]interface{})
	assert.Equal(t, "test", varsData["name"].(string), "got an unexpected var value")
	assert.Equal(t, "8888", varsData["port"].(string), "got an unexpected var value")
}

func TestApplicationAddSecrets(t *testing.T) {
	defer patchToolchainRoot()()
	secrets := map[string]string{
		"database-password": "password123",
		"registry-password": "passwordXYZ",
	}
	app := application{toolchain: &toolchain{name: "test"}}
	if err := os.MkdirAll(path.Join(app.path(), "templates"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := app.addSecrets(secrets); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path.Join(app.path(), "templates", "application-secret.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	s := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	secretsData := s["data"].(map[string]interface{})
	mockDatabasePassword, err := base64.StdEncoding.DecodeString(secretsData["database-password"].(string))
	if err != nil {
		t.Fatal(err)
	}
	mockRegistryPassword, err := base64.StdEncoding.DecodeString(secretsData["registry-password"].(string))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []byte("password123"), mockDatabasePassword, "got an unexpected secret value")
	assert.Equal(t, []byte("passwordXYZ"), mockRegistryPassword, "got an unexpected secret value")
}

func TestAddCIDriverHooks(t *testing.T) {
	defer patchToolchainRoot()()
	hooksManifest, err := os.ReadFile("testdata/hooks.yaml")
	if err != nil {
		t.Fatal(err)
	}
	catalog := &componentCatalog{
		HookSource: "quay.io/trustacks/test-catalog:latest",
		Components: map[string]component{
			"helloworld": {
				ApplicationHooks: string(hooksManifest),
			},
		},
	}
	tc := &toolchain{name: "test"}
	app := &application{name: "test", toolchain: tc}
	if err := os.MkdirAll(fmt.Sprintf("%s/applications/test/templates", tc.path()), 0755); err != nil {
		t.Fatal(err)
	}
	if err := app.addCIDriverHooks("helloworld", []string{"helloworld"}, catalog, map[string]interface{}{"application": "test"}); err != nil {
		t.Fatal(err)
	}
	assert.FileExists(t, fmt.Sprintf("%s/applications/test/templates/trustacks-application-test-hooks.yaml", tc.path()), "expected hooks manifest to exist")
}
