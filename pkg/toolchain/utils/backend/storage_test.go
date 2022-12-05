package backend

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestNewStorageConfig(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	config := StorageConfig{
		URL:             "test.storage.local",
		AccessKeyID:     "access123",
		SecretAccessKey: "secret123",
	}
	if err := NewStorageConfig(config, "default", clientset); err != nil {
		t.Fatal(err)
	}
}
