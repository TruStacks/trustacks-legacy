package api

import (
	"context"
	"fmt"
	"os"

	"github.com/teris-io/shortid"
	"github.com/trustacks/trustacks/pkg/toolchain/standard"
	"github.com/trustacks/trustacks/pkg/toolchain/state"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	"github.com/trustacks/trustacks/pkg/workflows/react"
	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

type APIV1Handler struct{}

// InstallToolchain installs the toolchain components.
func (h *APIV1Handler) InstallToolchain(kind, name string, profile map[string]interface{}) error {
	namespace := fmt.Sprintf("ts-toolchain-%s", name)
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
	sm := state.NewStateManager(namespace, dispatcher.Clientset())
	switch kind {
	case "standard":
		desiredConfig, err := standard.NewDesiredConfig(profile)
		if err != nil {
			return err
		}
		if err := sm.Save(desiredConfig); err != nil {
			return err
		}
		go standard.Install(namespace, desiredConfig, sm, dispatcher)
	}
	return nil
}

// UpgradeToolchain upgrades the toolchain components.
func (h *APIV1Handler) UpgradeToolchain(kind, name string) error {
	namespace := fmt.Sprintf("ts-toolchain-%s", name)
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		return err
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		return err
	}
	sm := state.NewStateManager(namespace, dispatcher.Clientset())
	switch kind {
	case "standard":
		desiredConfig := &standard.DesiredConfig{}
		if err := sm.Load(desiredConfig); err != nil {
			return err
		}
		go standard.Upgrade(namespace, desiredConfig, sm, dispatcher)
	}
	return nil
}

// UninstallToolchain uninstalls the toolchain components.
func (h *APIV1Handler) UninstallToolchain(kind, name string) error {
	namespace := fmt.Sprintf("ts-toolchain-%s", name)
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		return err
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		return err
	}
	sm := state.NewStateManager(namespace, dispatcher.Clientset())
	switch kind {
	case "standard":
		desiredConfig := &standard.DesiredConfig{}
		if err := sm.Load(desiredConfig); err != nil {
			return err
		}
		go standard.Uninstall(namespace, desiredConfig, sm, dispatcher)
	}
	return nil
}

type TestHandler struct{}

// ExecuteWorkflow .
func (h *TestHandler) ExecuteWorkflow(application, repo string) error {
	c, err := temporalClient.Dial(temporalClient.Options{
		HostPort:  "temporal-frontend-headless.ts-toolchain-dev.svc.cluster.local:7233",
		Namespace: application,
	})
	if err != nil {
		return err
	}
	defer c.Close()
	options := temporalClient.StartWorkflowOptions{
		TaskQueue:   application,
		RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 1},
	}
	runID, err := shortid.Generate()
	if err != nil {
		return err
	}
	params := react.ReactWorkflowParams{
		Name:  application,
		Repo:  repo,
		RunID: runID,
	}
	_, err = c.ExecuteWorkflow(context.Background(), options, react.ReactWorkflow, params)
	if err != nil {
		return err
	}
	return nil
}
