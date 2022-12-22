package standard

import (
	"fmt"
	"os"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

func Install(name string, prof profile.Profile) error {
	namespace := fmt.Sprintf("trustacks-toolchain-%s", name)
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		return err
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		return err
	}
	storageConfig := backend.StorageConfig{
		URL:             os.Getenv("STORAGE_URL"),
		AccessKeyID:     os.Getenv("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		return err
	}
	return nil
}
