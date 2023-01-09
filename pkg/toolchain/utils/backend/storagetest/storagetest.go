package storagetest

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const namespace = "integration-test-utils"

// NewTestStorageConfig creates a storage configuration for a minio
// bucket storage server.
func NewTestStorageConfig(component string) (backend.StorageConfig, error) {
	config := backend.StorageConfig{
		URL:             fmt.Sprintf("s3:http://minio.%s.svc.cluster.local:9000/%s", namespace, component),
		AccessKeyID:     "root",
		SecretAccessKey: "rootpass123",
	}
	minioClient, err := minio.New(
		"minio.integration-test-utils.svc.cluster.local:9000",
		&minio.Options{
			Creds:  credentials.NewStaticV4("root", "rootpass123", ""),
			Secure: false,
		},
	)
	if err != nil {
		return config, err
	}
	ctx := context.Background()
	if err := minioClient.MakeBucket(ctx, component, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, component)
		if errBucketExists != nil || !exists {
			return config, err
		}
	}
	return config, nil
}

// PurgeBucket deletes the bucket from the minio server.
func PurgeBucket(component string) error {
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		log.Fatal(err)
	}
	cmd := []string{
		"mc alias set s3 http://localhost:9000 root rootpass123",
		fmt.Sprintf(`mc rb s3/%s --force`, component),
	}
	clientset := dispatcher.Clientset()
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{LabelSelector: "app.kubernetes.io/name=minio"})
	if err != nil {
		return err
	}
	if err := dispatcher.ExecCommand(pods.Items[0].Name, "minio", strings.Join(cmd, "\n"), namespace); err != nil {
		return err
	}
	return nil
}
