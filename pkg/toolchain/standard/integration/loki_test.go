package integration

import (
	"log"
	"testing"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/loki"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

func TestLokiLifecycleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
	}
	t.Parallel()
	namespace := "loki-integration-test"
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		log.Fatal(err)
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		log.Fatal(err)
	}
	storageConfig, err := NewTestStorageConfig("loki")
	if err != nil {
		log.Fatal(err)
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		log.Fatal(err)
	}
	defer PurgeBucket("loki")
	pfl := profile.Profile{Domain: "loki-integration-test.local.gd", Port: 8081, Insecure: true}
	sso := authentik.New(pfl)
	if err := sso.Install(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	c := loki.New(pfl)
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
	if err := sso.Uninstall(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.DeleteNamespace(); err != nil {
		log.Fatal(err)
	}
}
