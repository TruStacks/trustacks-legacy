package backend

import (
	"context"
	"strings"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StorageConfig struct {
	URL             string
	AccessKeyID     string
	SecretAccessKey string
}

// NewStorageConfig creates a storage configuration for an external
// s3 datastore.
func NewStorageConfig(config StorageConfig, namespace string, clientset kubernetes.Interface) error {
	secretName := "ts-storage-config"
	pwd, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: map[string][]byte{
			"url":               []byte(config.URL),
			"password":          []byte(pwd),
			"access-key-id":     []byte(config.AccessKeyID),
			"secret-access-key": []byte(config.SecretAccessKey),
		},
	}
	if _, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{}); err != nil {
		if strings.Contains(err.Error(), "not found") {
			if _, err := clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
				return err
			}
		}
		if !strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}
