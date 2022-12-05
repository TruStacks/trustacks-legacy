package integration

import (
	"log"
	"os"
	"testing"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

func TestAuthentikLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}
	namespace := "authentik-integration-test"
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		log.Fatal(err)
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		log.Fatal(err)
	}
	storageConfig := backend.StorageConfig{
		URL:             os.Getenv("STORAGE_URL"),
		AccessKeyID:     os.Getenv("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		log.Fatal(err)
	}
	sso := authentik.New("test.local.gd", 80, true)
	if err := sso.Install(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := sso.Upgrade(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := sso.Rollback(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := sso.Uninstall(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.DeleteNamespace(); err != nil {
		log.Fatal(err)
	}
}
