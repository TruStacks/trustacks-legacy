package configstore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

func sync(path string, client *minio.Client, bucket string) error {
	file, err := os.Open(filepath.Join(dataDir, path))
	if err != nil {
		return err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = client.PutObject(
		context.Background(),
		bucket,
		fmt.Sprintf("%s/%s", path, time.Now().Format(time.RFC3339)),
		file,
		fileStat.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	return err
}
