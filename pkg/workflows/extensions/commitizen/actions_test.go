package commitizen

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
	"gotest.tools/v3/assert"
)

func TestVersionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	defer os.RemoveAll("./artifacts")
	client, err := dagger.Connect(context.TODO(), dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}
	project := utils.NewProject("test", filepath.Join("testdata", "hello"), client)
	defer client.Close()
	container, err := Version("version", project, client)
	if err != nil {
		t.Fatal(err)
	}
	f := container.File(utils.ArtifactPath("version", "version"))
	version, err := f.Contents(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0.1.0", version)
}
