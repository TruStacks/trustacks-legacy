package configstore

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestSync(t *testing.T) {
	previousDataDir := dataDir
	dataDir = os.TempDir()
	defer func() {
		dataDir = previousDataDir
	}()
	if err := os.WriteFile(fmt.Sprintf("/%s/test", dataDir), []byte(""), 0744); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fmt.Sprintf("/%s/test", dataDir))
	client, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	client.MakeBucket(context.Background(), "test", minio.MakeBucketOptions{})
	if err := sync("test", client, "test"); err != nil {
		t.Fatal(err)
	}
}
