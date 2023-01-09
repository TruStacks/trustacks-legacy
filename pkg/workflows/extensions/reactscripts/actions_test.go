package reactscripts

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
)

func TestTestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	defer os.RemoveAll("./artifacts")
	client, err := dagger.Connect(context.TODO(), dagger.WithLogOutput(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}
	project := utils.NewProject("test", filepath.Join("testdata", "reactapp"), client)
	defer client.Close()
	_, err = Test("test", project, client)
	if err != nil {
		t.Fatal(err)
	}
}
