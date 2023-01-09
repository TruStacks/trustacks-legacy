package docker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend/storagetest"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	"gopkg.in/yaml.v3"
)

func TestGetChart(t *testing.T) {
	path, err := (&Docker{}).GetChart()
	if err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	expectedFiles := []string{"Chart.yaml", "charts"}
	for _, f := range files {
		assert.Contains(t, expectedFiles, f.Name())
	}
	charts, err := os.ReadDir(filepath.Join(path, "charts"))
	if err != nil {
		t.Fatal(err)
	}
	// check that the chart version matches
	fd, err := os.Open(fmt.Sprintf("%s/Chart.yaml", path))
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(fd)
	if err != nil {
		t.Fatal(err)
	}
	chartYaml := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &chartYaml); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, chartVersion, chartYaml["version"].(string))
	assert.Equal(t, "docker.tgz", charts[0].Name())
}

func TestDockerLifecycleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}
	t.Parallel()
	namespace := "docker-integration-test"
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		t.Fatal(err)
	}
	storageConfig, err := storagetest.NewTestStorageConfig("docker")
	if err != nil {
		t.Fatal(err)
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		t.Fatal(err)
	}
	defer storagetest.PurgeBucket("docker")
	pfl := profile.Profile{Domain: "docker-integration-test.local.gd", Port: 8081, Insecure: true}
	c := New(pfl)
	if err := c.Install(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Upgrade(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Rollback(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Uninstall(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.DeleteNamespace(); err != nil {
		t.Fatal(err)
	}
}
