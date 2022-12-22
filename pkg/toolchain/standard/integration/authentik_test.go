//go:build integration

package integration

import (
	"log"
	"testing"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

func TestAuthentikLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}
	t.Parallel()
	namespace := "authentik-integration-test"
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		log.Fatal(err)
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		log.Fatal(err)
	}
	storageConfig, err := NewTestStorageConfig("authentik")
	if err != nil {
		log.Fatal(err)
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		log.Fatal(err)
	}
	defer PurgeBucket("authentik")
	c := authentik.New(profile.Profile{Domain: "authentik-integration-test.local.gd", Port: 8081, Insecure: true})
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
		log.Fatal(err)
	}
}
